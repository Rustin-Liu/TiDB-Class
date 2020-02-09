package project_1

import (
	"flag"
	"log"
	"runtime"
	"runtime/pprof"
	"sort"
	"testing"

	"github.com/pingcap/talentplan/tidb/util"
)

const Merge string = "Merge"

var version = flag.String("version", "", "sort code version")

func BenchmarkMergeSort(b *testing.B) {
	flag.Parse()
	if *version != "" {
		f, err := util.CreateProfile(Merge, util.Cpu, *version)
		if err == nil {
			if err := pprof.StartCPUProfile(&f); err != nil {
				log.Fatalf("could not start %s sort CPU profile: %s", Merge, err.Error())
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
		f, err := util.CreateProfile(Merge, util.Mem, *version)
		if err == nil {
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(&f); err != nil {
				log.Fatalf("could not write %s sort memory profile: %s", Merge, err.Error())
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
