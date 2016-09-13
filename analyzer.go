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
	Hits    int64
	Elapsed []float64
}

// TopMutex is..
var TopMutex = &sync.Mutex{}

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
func AnalyzerDispatch(vStart time.Time, vEnd time.Time) {

	start = vStart
	end = vEnd

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
		log.Printf("%s - %d lines where processed", filelog, count)
	}()

	filepointer := NewFilePointer(filelog)
	reader := bufio.NewReader(file)
	pos := int64(0)

	for {
		l, err := reader.ReadString(byte('\n'))
		if err != nil {
			break
		}

		line := NewLineLog(l, *filepointer, pos)
		pos += int64(len(l))

		if InTimeSpan(start, end, line.Time) {
			RecordAdd(*line)

			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
		count++
	}
}
