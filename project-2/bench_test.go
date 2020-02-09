package main

import (
	"fmt"
	"runtime"
	"testing"
)

func benchmarkDataScale() (DataSize, int) {
	dataSize := DataSize(100 * MB)
	nMapFiles := 20
	return dataSize, nMapFiles
}

func BenchmarkExampleURLTop(b *testing.B) {
	rounds := ExampleURLTop10Args(GetMRCluster().NWorkers())
	benchmarkURLTop10(b, rounds)
}

func benchmarkURLTop10(b *testing.B, rounds RoundsArgs) {
	if len(rounds) == 0 {
		b.Fatalf("no rounds arguments, please finish your code")
	}
	mr := GetMRCluster()
	dataSize, nMapFiles := benchmarkDataScale()
	// run cases.
	gens := AllCaseGenFs()
	b.ResetTimer()
	for i, gen := range gens {
		// generate data.
		prefix := dataPrefix(i, dataSize, nMapFiles)
		c := gen(prefix, int(dataSize), nMapFiles)

		runtime.GC()

		// run map-reduce rounds
		inputFiles := c.MapFiles
		for idx, r := range rounds {
			jobName := fmt.Sprintf("Case%d-Round%d", i, idx)
			ch := mr.Submit(jobName, prefix, r.MapFunc, r.ReduceFunc, inputFiles, r.NReduce)
			inputFiles = <-ch
		}

		// check result
		if len(inputFiles) != 1 {
			panic("the length of result file list should be 1")
		}
		result := inputFiles[0]

		if errMsg, ok := CheckFile(c.ResultFile, result); !ok {
			b.Fatalf("Case%d FAIL, dataSize=%v, nMapFiles=%v, %v\n", i, dataSize, nMapFiles, errMsg)
		} else {
			fmt.Printf("Case%d PASS, dataSize=%v, nMapFiles=%v\n", i, dataSize, nMapFiles)
		}
	}
}
