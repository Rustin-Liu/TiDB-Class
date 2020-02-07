package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
)

// KeyValue is a type used to hold the key/value pairs passed to the map and reduce functions.
type KeyValue struct {
	Key   string
	Value string
}

// ReduceF function from MIT 6.824 LAB1.
type ReduceF func(key string, values []string) string

// MapF function from MIT 6.824 LAB1.
type MapF func(filename string, contents string) []KeyValue

// jobPhase indicates whether a task is scheduled as a map or reduce task.
type jobPhase string

const (
	mapPhase    jobPhase = "mapPhase"
	reducePhase          = "reducePhase"
)

type task struct {
	dataDir    string
	jobName    string
	mapFile    string   // only for map, the input file.
	phase      jobPhase // are we in mapPhase or reducePhase?
	taskNumber int      // this task's index in the current phase.
	nMap       int      // number of map tasks.
	nReduce    int      // number of reduce tasks.
	mapF       MapF     // map function used in this job.
	reduceF    ReduceF  // reduce function used in this job.
	wg         sync.WaitGroup
}

// MRCluster represents a map-reduce cluster.
type MRCluster struct {
	nWorkers int
	wg       sync.WaitGroup
	taskCh   chan *task
	exit     chan struct{}
}

var singleton = &MRCluster{
	nWorkers: runtime.NumCPU(),
	taskCh:   make(chan *task),
	exit:     make(chan struct{}),
}

func init() {
	singleton.Start()
}

// GetMRCluster returns a reference to a MRCluster.
func GetMRCluster() *MRCluster {
	return singleton
}

// NWorkers returns how many workers there are in this cluster.
func (c *MRCluster) NWorkers() int { return c.nWorkers }

// Start starts this cluster.
func (c *MRCluster) Start() {
	for i := 0; i < c.nWorkers; i++ {
		c.wg.Add(1)
		go c.worker()
	}
}

func (c *MRCluster) worker() {
	defer c.wg.Done()
	for {
		select {
		case t := <-c.taskCh:
			if t.phase == mapPhase {
				content, err := ioutil.ReadFile(t.mapFile)
				if err != nil {
					panic(err)
				}

				fs := make([]*os.File, t.nReduce)
				bs := make([]*bufio.Writer, t.nReduce)
				for i := range fs {
					rpath := reduceName(t.dataDir, t.jobName, t.taskNumber, i)
					fs[i], bs[i] = CreateFileAndBuf(rpath)
				}
				results := t.mapF(t.mapFile, string(content))
				for _, kv := range results {
					enc := json.NewEncoder(bs[ihash(kv.Key)%t.nReduce])
					if err := enc.Encode(&kv); err != nil {
						log.Fatalln(err)
					}
				}
				for i := range fs {
					SafeClose(fs[i], bs[i])
				}
			} else {
				var keys []string                   // store all keys for sort.
				var kvs = make(map[string][]string) // store all key-value pairs from nMap imm files.

				// read nMap imm files from map workers
				for i := 0; i < t.nMap; i++ {
					fileName := reduceName(t.dataDir, t.jobName, i, t.taskNumber)
					imm, err := os.Open(fileName)
					if err != nil {
						log.Fatalln(err)
					}
					var kv KeyValue
					dec := json.NewDecoder(imm)
					err = dec.Decode(&kv)
					for err == nil {
						// is this key seen?
						if _, ok := kvs[kv.Key]; !ok {
							keys = append(keys, kv.Key)
						}
						kvs[kv.Key] = append(kvs[kv.Key], kv.Value)
						// decode repeatedly until an error
						err = dec.Decode(&kv)
					}
					err = imm.Close()
					if err != nil {
						log.Fatalln(err)
					}
				}
				sort.Strings(keys)
				outFileName := mergeName(t.dataDir, t.jobName, t.taskNumber)
				out, err := os.Create(outFileName)
				if err != nil {
					log.Fatalln(err)
				}
				for _, key := range keys {
					if _, err = fmt.Fprintf(out, "%v", t.reduceF(key, kvs[key])); err != nil {
						log.Printf("write [key: %s] to file %s failed", key, outFileName)
					}
				}
				err = out.Close()
				if err != nil {
					log.Fatalln(err)
				}
			}
			t.wg.Done()
		case <-c.exit:
			return
		}
	}
}

// Shutdown shutdowns this cluster.
func (c *MRCluster) Shutdown() {
	close(c.exit)
	c.wg.Wait()
}

// Submit submits a job to this cluster.
func (c *MRCluster) Submit(jobName, dataDir string, mapF MapF, reduceF ReduceF, mapFiles []string, nReduce int) <-chan []string {
	notify := make(chan []string)
	go c.run(jobName, dataDir, mapF, reduceF, mapFiles, nReduce, notify)
	return notify
}

func (c *MRCluster) run(jobName, dataDir string, mapF MapF, reduceF ReduceF, mapFiles []string, nReduce int, notify chan<- []string) {
	// map phase
	nMap := len(mapFiles)
	tasks := make([]*task, 0, nMap)
	for i := 0; i < nMap; i++ {
		t := &task{
			dataDir:    dataDir,
			jobName:    jobName,
			mapFile:    mapFiles[i],
			phase:      mapPhase,
			taskNumber: i,
			nReduce:    nReduce,
			nMap:       nMap,
			mapF:       mapF,
		}
		t.wg.Add(1)
		tasks = append(tasks, t)
		go func() { c.taskCh <- t }()
	}
	for _, t := range tasks {
		t.wg.Wait()
	}
	tasks = make([]*task, 0, nReduce)
	for i := 0; i < nReduce; i++ {
		t := &task{
			dataDir:    dataDir,
			jobName:    jobName,
			phase:      reducePhase,
			taskNumber: i,
			nReduce:    nReduce,
			reduceF:    reduceF,
			nMap:       nMap,
		}
		t.wg.Add(1)
		tasks = append(tasks, t)
		go func() { c.taskCh <- t }()
	}
	for _, t := range tasks {
		t.wg.Wait()
	}
	notify <- []string{merge(jobName, dataDir, nReduce)}
}

func merge(jobName, dataDir string, nReduce int) string {
	resultFileName := jobName + ".txt"
	file, err := os.Create(resultFileName)
	if err != nil {
		log.Fatalln(err)
	}
	w := bufio.NewWriter(file)
	for i := 0; i < nReduce; i++ {
		fileName := mergeName(dataDir, jobName, i)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatalln(err)
		}
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(content)
		if err != nil {
			log.Fatalln(err)
		}
		err = file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}
	SafeClose(file, w)
	return resultFileName
}
