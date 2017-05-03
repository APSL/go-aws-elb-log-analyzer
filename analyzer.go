package main

import (
	"bufio"
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
var AnalyzerQueue chan *FileLog

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

	AnalyzerQueue = make(chan *FileLog, 10)

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
		case f, more := <-AnalyzerQueue:
			if !more {
				return
			}
			wg.Add(1)
			go analyzerFile(f)
		}
	}

}

func analyzerFile(f *FileLog) {
	defer wg.Done()

	log.Printf("Analyzing %s", f.Filename)
	count := 0

	file, err := os.Open(f.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		file.Close()
		log.Printf("%s - %d lines where processed", f.Filename, count)
	}()

	f.Pointer = NewFilePointer(f.Filename)
	reader := bufio.NewReader(file)
	pos := int64(0)

	for {
		l, err := reader.ReadString(byte('\n'))
		if err != nil {
			return
		}

		line := NewLineLog(l, f, pos)
		pos += int64(len(l))

		if InTimeSpan(start, end, line.Time) {
			RecordAdd(line)

			if analyze {
				RecordIP(line.IPClient, 1, line.Elapsed)
			}

			if err != nil {
				log.Fatalf("Error: %v", err)
			}
		}
		count++
	}
}
