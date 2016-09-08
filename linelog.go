package main

import (
	"net"
	"regexp"
	"strconv"
	"time"
)

// LineLog is the struct to analyce and store a line
type LineLog struct {
	raw      string
	URL      string
	Time     time.Time
	Elapsed  float64
	Method   string
	IPClient net.IP
}

// NewLineLog export
func NewLineLog(raw string) *LineLog {

	line := &LineLog{
		raw: raw,
	}

	line.parse()

	return line
}

func (l *LineLog) parse() {
	re := regexp.MustCompile(`(?P<date>[^Z]+Z) (?P<elb>[^\s]+) (?P<ipclient>[^:]+?):[0-9]+ (?P<ipnode>[^:]+?):[0-9]+ (?P<reqtime>[0-9\.]+) (?P<backtime>[0-9\.]+) (?P<restime>[0-9\.]+) (?P<elbcode>[0-9]{3}) (?P<backcode>[0-9]{3}) (?P<lenght1>[0-9]+) (?P<lenght2>[0-9]+) "(?P<Method>[A-Z]+) (?P<URL>[^"]+) HTTP/[0-9\.]+".*`)
	n1 := re.SubexpNames()
	r2 := re.FindAllStringSubmatch(l.raw, -1)

	if r2 == nil {
		return
	}

	for i, n := range r2[0] {
		switch n1[i] {
		case "date":
			l.Time, _ = time.Parse(time.RFC3339Nano, n)
			break
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
		}
		//fmt.Printf("%d. match='%s'\tname='%s'\n", i, n, n1[i])
	}

	/**

	l.Time, _ = time.Parse(time.RFC3339Nano, re.ReplaceAllString(l.raw, "${date}"))
	l.IPClient = net.ParseIP(re.ReplaceAllString(l.raw, "${ipclient}"))
	l.URL = re.ReplaceAllString(l.raw, "${URL}")
	l.Method = re.ReplaceAllString(l.raw, "${Method}")
	reqtime, _ := strconv.ParseFloat(re.ReplaceAllString(l.raw, "${reqtime}"), 64)
	backtime, _ := strconv.ParseFloat(re.ReplaceAllString(l.raw, "${backtime}"), 64)
	restime, _ := strconv.ParseFloat(re.ReplaceAllString(l.raw, "${restime}"), 64)
	l.Elapsed = reqtime + backtime + restime
	**/

	/*
		if l.Elapsed > 10 {
			fmt.Printf("----------------\n")
			fmt.Printf("RAW: %s\n", l.raw)
			fmt.Printf("    reqtime: %.5f | %s\n", reqtime, re.ReplaceAllString(l.raw, "${reqtime}"))
			fmt.Printf("   backtime: %.5f | %s\n", backtime, re.ReplaceAllString(l.raw, "${backtime}"))
			fmt.Printf("    restime: %.5f | %s\n", restime, re.ReplaceAllString(l.raw, "${restime}"))
			fmt.Printf("    Elapsed: %.5f\n\n", l.Elapsed)
		}
	*/
}
