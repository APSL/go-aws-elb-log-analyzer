package main

import (
	"fmt"
	"os"
	"sort"
)

// RawRecords is full list of records
var RawRecords []LineLog

// InitRecords start the records
func InitRecords() {
	RawRecords = make([]LineLog, 0)
}

// RecordAdd include record in big slice
func RecordAdd(line LineLog) {
	RawRecords = append(RawRecords, line)
}

// LogRecords stuct to sort all records
type LogRecords []LineLog

func (a LogRecords) Len() int           { return len(a) }
func (a LogRecords) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LogRecords) Less(i, j int) bool { return a[j].Time.After(a[i].Time) }

// PrintSortLog order by date
func PrintSortLog(save string) {

	file, _ := os.OpenFile(save, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0660)

	sort.Sort(LogRecords(RawRecords))

	for l := range RawRecords {
		// fmt.Printf("%s - %s\n", RawRecords[l].Time, RawRecords[l].Filelog)
		buf := RawRecords[l].Filelog.GetLine(RawRecords[l].Seek, RawRecords[l].Len)
		fmt.Printf("%s", string(buf))
		file.Write(buf)
	}

	file.Close()

}
