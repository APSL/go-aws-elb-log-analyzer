package main

import (
	"bufio"
	"log"
	"os"
	"sort"
)

// LogRecords stuct to sort all records
type LogRecords []*LineLog

// RawRecords is full list of records
var rawRecords LogRecords

// InitRecords start the records
func InitRecords() {
	rawRecords = make([]*LineLog, 0)
}

// RecordAdd include record in big slice
func RecordAdd(line *LineLog) {
	rawRecords = append(rawRecords, line)
}

func (a LogRecords) Len() int           { return len(a) }
func (a LogRecords) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a LogRecords) Less(i, j int) bool { return a[j].Time.After(a[i].Time) }

// saveSortedLog order by date
func saveSortedLog(save string) {
	file, e := os.OpenFile(save, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_APPEND, 0660)
	if e != nil {
		log.Panicln("Coudn't open file for writeing", save)
	}
	defer log.Printf("Sorted log file saved in %s", file.Name())
	defer file.Close()

	ob := bufio.NewWriter(file)
	defer ob.Flush()

	sort.Sort(rawRecords)

	for _, r := range rawRecords {
		ob.Write(r.Filelog.Pointer.GetLine(r.Seek, r.Len))
	}

}
