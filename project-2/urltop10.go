package project_2

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

var count int64

// URLTop10 generates RoundsArgs for getting the 10 most frequent URLs.
// There are two rounds in this approach.
// The first round will do url count.
// The second will sort results generated in the first round and
// get the 10 most frequent URLs.
func URLTop10(nWorkers int) RoundsArgs {
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
func URLCountMap(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	result := make([]KeyValue, 0, len(lines))
	kvs := make(map[string]int)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		count, ok := kvs[l]
		if ok {
			kvs[l] = count + 1
		} else {
			kvs[l] = 1
		}
	}
	for k, v := range kvs {
		result = append(result, KeyValue{Key: k, Value: strconv.Itoa(v)})
	}
	return result
}

// URLCountReduce is the reduce function in the first round.
func URLCountReduce(key string, values []string) string {
	total := 0
	for _, count := range values {
		temp, err := strconv.Atoi(count)
		if err == nil {
			total += temp
		}
	}
	return fmt.Sprintf("%s %s\n", strconv.Itoa(total), key)
}

// URLTop10Map is the map function in the second round.
func URLTop10Map(filename string, contents string) []KeyValue {
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
	atomic.StoreInt64(&count, 0)
	return kvs
}

// URLTop10Reduce is the reduce function in the second round.
func URLTop10Reduce(key string, values []string) string {
	sort.Strings(values)
	var result string
	for _, val := range values {
		if atomic.LoadInt64(&count) != 10 {
			result += fmt.Sprintf("%s: %s\n", val, key)
			atomic.AddInt64(&count, 1)
		}
	}
	return result
}
