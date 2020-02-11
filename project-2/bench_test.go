package project_2

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"runtime/pprof"
	"testing"

	"github.com/pingcap/talentplan/tidb/util"
)

var version = flag.String("version", "", "sort code version")

const Example string = "Example"

func benchmarkDataScale() (DataSize, int) {
	dataSize := DataSize(100 * MB)
	nMapFiles := 20
	return dataSize, nMapFiles
}

func BenchmarkExampleURLTop(b *testing.B) {
	flag.Parse()
	if *version != "" {
		f, err := util.CreateProfile(Example, util.Cpu, *version)
		if err == nil {
			if err := pprof.StartCPUProfile(&f); err != nil {
				log.Fatalf("could not start %s mr CPU profile: %s", Example, err.Error())
			}
			defer pprof.StopCPUProfile()
		}
	}
	rounds := ExampleURLTop10(GetMRCluster().NWorkers())
	benchmarkURLTop10(b, rounds)
	if *version != "" {
		f, err := util.CreateProfile(Example, util.Mem, *version)
		if err == nil {
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(&f); err != nil {
				log.Fatalf("could not write %s mr memory profile: %s", Example, err.Error())
			}
			f.Close()
		}
	}
}

func BenchmarkURLTop(b *testing.B) {
	flag.Parse()
	if *version != "" {
		f, err := util.CreateProfile(Example, util.Cpu, *version)
		if err == nil {
			if err := pprof.StartCPUProfile(&f); err != nil {
				log.Fatalf("could not start %s mr CPU profile: %s", Example, err.Error())
			}
			defer pprof.StopCPUProfile()
		}
	}
	rounds := URLTop10(GetMRCluster().NWorkers())
	benchmarkURLTop10(b, rounds)
	if *version != "" {
		f, err := util.CreateProfile(Example, util.Mem, *version)
		if err == nil {
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(&f); err != nil {
				log.Fatalf("could not write %s mr memory profile: %s", Example, err.Error())
			}
			f.Close()
		}
	}
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
