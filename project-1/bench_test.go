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

var quickSortCpuProfile = flag.String("quickCpu", "", "write quick sort cpu profile `file`")
var quickSortMemProfile = flag.String("quickMem", "", "write quick sort memory profile to `file`")

func BenchmarkMergeSort(b *testing.B) {
	flag.Parse()
	if *quickSortCpuProfile != "" {
		f, err := os.Create(*quickSortCpuProfile)
		if err != nil {
			log.Fatal("could not create quick sort CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start quick sort CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
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

	if *quickSortMemProfile != "" {
		f, err := os.Create(*quickSortMemProfile)
		if err != nil {
			log.Fatal("could not create quick sort memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write quick sort memory profile: ", err)
		}
		f.Close()
	}
}

var mergeSortCpuProfile = flag.String("mergeCpu", "", "write merge sort cpu profile `file`")
var mergeSortMemProfile = flag.String("mergeMem", "", "write merge sort memory profile to `file`")

func BenchmarkNormalSort(b *testing.B) {
	flag.Parse()
	if *mergeSortCpuProfile != "" {
		f, err := os.Create(*mergeSortCpuProfile)
		if err != nil {
			log.Fatal("could not create merge sort CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start merge sort CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
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
		sort.Slice(src, func(i, j int) bool { return src[i] < src[j] })
	}

	if *mergeSortMemProfile != "" {
		f, err := os.Create(*mergeSortMemProfile)
		if err != nil {
			log.Fatal("could not create merge sort memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write merge sort memory profile: ", err)
		}
		f.Close()
	}
}
