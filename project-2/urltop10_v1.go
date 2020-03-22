package project_2

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

var countV1 int64

// URLTop10 generates RoundsArgs for getting the 10 most frequent URLs.
// There are two rounds in this approach.
// The first round will do url count.
// The second will sort results generated in the first round and
// get the 10 most frequent URLs.
func URLTop10V1(nWorkers int) RoundsArgs {
	var args RoundsArgs
	// round 1: do url count.
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs.
	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

// URLCountMap is the map function in the first round.
func URLCountMapV1(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	kvs := make([]KeyValue, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		kvs = append(kvs, KeyValue{Key: l})
	}
	return kvs
}

// URLCountReduce is the reduce function in the first round.
func URLCountReduceV1(key string, values []string) string {
	return fmt.Sprintf("%s %s\n", strconv.Itoa(len(values)), key)
}

// URLTop10Map is the map function in the second round.
func URLTop10MapV1(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	kvs := make([]KeyValue, 0, len(lines))
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		words := strings.Split(line, " ")
		kvs = append(kvs, KeyValue{words[0], words[1]})
	}
	atomic.StoreInt64(&countV1, 0)
	return kvs
}

// URLTop10Reduce is the reduce function in the second round.
func URLTop10ReduceV1(key string, values []string) string {
	sort.Strings(values)
	var result string
	for _, val := range values {
		if atomic.LoadInt64(&countV1) != 10 {
			result += fmt.Sprintf("%s: %s\n", val, key)
			atomic.AddInt64(&countV1, 1)
		}
	}
	return result
}
