package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// IPRecord test
type IPRecord struct {
	sync.RWMutex
	Hits    int64
	Elapsed []float64
}

// TopMutex is..
var TopMutex = &sync.RWMutex{}

// TopIP is a counter of IP clients
var TopIP map[string]*IPRecord

// AnalyzerQueue Queue of files to analyce
var AnalyzerQueue chan []byte

// AnalyzerDone Signal on complete files download
var AnalyzerDone chan bool

var start time.Time
var end time.Time

var wg sync.WaitGroup

// AnalyzerDispatch starts the go routine to analyce each record
func AnalyzerDispatch(s_start time.Time, s_end time.Time) {

	start = s_start
	end = s_end

	TopIP = make(map[string]*IPRecord)

	AnalyzerQueue = make(chan []byte)

	wg.Add(1)
	go analyzerReader()

}

// AnalyzerFinished control the end of the analyce process
func AnalyzerFinished() {
	defer log.Printf("Finish the analycis")

	log.Println("Waiting for finish")
	wg.Wait()
}

func analyzerReader() {
	defer wg.Done()

	for {
		select {
		case filelog := <-AnalyzerQueue:
			if filelog == nil {
				return
			}
			wg.Add(1)
			go analyzerFile(string(filelog))
		}
	}

	fmt.Printf("Waiting %v", wg)

}

func analyzerFile(filelog string) {
	defer wg.Done()

	log.Printf("Analyzing %s", filelog)

	count := 0

	file, err := os.Open(filelog)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		file.Close()
		log.Printf("%d lines where processed", count)
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := NewLineLog(scanner.Text())

		if InTimeSpan(start, end, line.Time) {
			RecordIP(line.IPClient, 1, line.Elapsed)

		}

		count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
