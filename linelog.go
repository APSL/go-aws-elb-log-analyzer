package main

import (
	"net"
	"os"
	"regexp"
	"time"
)

// LineLog is the struct to analyce and store a line
type LineLog struct {
	URL      string
	Time     time.Time
	Elapsed  float64
	Method   string
	IPClient net.IP
	Seek     int64
	Filelog  FilePointer
	Len      int
}

// NewLineLog export
func NewLineLog(raw string, filelog FilePointer, seek int64) *LineLog {

	line := &LineLog{
		Filelog: filelog,
		Seek:    seek,
		Len:     len([]byte(raw)),
	}

	line.parse(raw)

	return line
}

func (l *LineLog) parse(raw string) {

	re := regexp.MustCompile(`(?P<date>[^Z]+Z) (?P<elb>[^\s]+) (?P<ipclient>[^:]+?):[0-9]+ (?P<ipnode>[^:]+?):[0-9]+ (?P<reqtime>[0-9\.]+) (?P<backtime>[0-9\.]+) (?P<restime>[0-9\.]+) (?P<elbcode>[0-9]{3}) (?P<backcode>[0-9]{3}) (?P<lenght1>[0-9]+) (?P<lenght2>[0-9]+) "(?P<Method>[A-Z]+) (?P<URL>[^"]+) HTTP/[0-9\.]+".*`)
	n1 := re.SubexpNames()
	r2 := re.FindAllStringSubmatch(raw, -1)

	if r2 == nil {
		return
	}

	for i, n := range r2[0] {
		switch n1[i] {
		case "date":
			l.Time, _ = time.Parse(time.RFC3339Nano, n)
			break
			/*
				case "ipclient":
					l.IPClient = net.ParseIP(n)
					break
				case "Method":
					l.Method = n
					break
				case "URL":
					l.URL = n
					break
				case "reqtime":
					_time, _ := strconv.ParseFloat(n, 64)
					l.Elapsed = l.Elapsed + _time
					break
				case "backtime":
					_time, _ := strconv.ParseFloat(n, 64)
					l.Elapsed = l.Elapsed + _time
					break
				case "restime":
					_time, _ := strconv.ParseFloat(n, 64)
					l.Elapsed = l.Elapsed + _time
					break
			*/
		}

	}

}

// FilePointer method to move in log files
type FilePointer struct {
	Filename string
	File     *os.File
}

// NewFilePointer create new struc to locate line in file
func NewFilePointer(filename string) *FilePointer {
	f := &FilePointer{
		Filename: filename,
	}

	f.File, _ = os.Open(f.Filename)

	return f
}

// GetLine of the file
func (f *FilePointer) GetLine(seek int64, len int) []byte {

	// f.File.Seek(seek, len)
	buf := make([]byte, len)
	_, _ = f.File.ReadAt(buf, seek)

	return buf

}
