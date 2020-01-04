package project_1

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"testing"
)

type FileType string

const (
	CPU FileType = "Cpu"
	MEM FileType = "Mem"
)

const MERGE string = "Merge"

var version = flag.String("version", "", "sort code version")

func BenchmarkMergeSort(b *testing.B) {
	flag.Parse()
	if *version != "" {
		f, err := createProfile(MERGE, CPU, *version)
		if err == nil {
			if err := pprof.StartCPUProfile(&f); err != nil {
				log.Fatalf("could not start %s sort CPU profile: %s", MERGE, err.Error())
			}
			defer pprof.StopCPUProfile()
		}
	}

	numElements := 16 << 20
	src := make([]int64, numElements)
	original := make([]int64, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		MergeSort(src)
	}

	if *version != "" {
		f, err := createProfile(MERGE, MEM, *version)
		if err == nil {
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(&f); err != nil {
				log.Fatalf("could not write %s sort memory profile: %s", MERGE, err.Error())
			}
			f.Close()
		}
	}
}

func BenchmarkNormalSort(b *testing.B) {
	numElements := 16 << 20
	src := make([]int64, numElements)
	original := make([]int64, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		sort.Slice(src, func(i, j int) bool { return src[i] < src[j] })
	}
}

func createProfile(sortType string, fileType FileType, version string) (os.File, error) {
	f, err := os.Create("prof/" + version + "/" + sortType + string(fileType) + ".prof")
	if err != nil {
		log.Fatalf("could not create %s sort %s profile: %s", sortType, fileType, err.Error())
		return *f, err
	}
	return *f, nil
}
