package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/montanaflynn/stats"
)

// AddHit to a record IP
func (r *IPRecord) AddHit() {
	r.Hits++
}

// AddElapsed to a record IP
func (r *IPRecord) AddElapsed(elapsed float64) {
	r.Elapsed = append(r.Elapsed, elapsed)
}

// RecordIP to top
func RecordIP(ip net.IP, hit int, elapsed float64) {

	if _, ok := TopIP[ip.String()]; ok {
		TopIP[ip.String()].AddHit()
		TopIP[ip.String()].AddElapsed(elapsed)
	} else {
		TopIP[ip.String()] = &IPRecord{Hits: 1}
		TopIP[ip.String()].AddElapsed(elapsed)
	}

}

// IPSortedInt export
type IPSortedInt struct {
	Name  string
	Value int64
}

// IPSortedFloat export
type IPSortedFloat struct {
	Name    string
	Value   float64
	Perc80  float64
	Average float64
}

func (p IPSortedInt) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Value)
}

// ByHits amount of times found at the logs
type ByHits []IPSortedInt

func (a ByHits) Len() int           { return len(a) }
func (a ByHits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHits) Less(i, j int) bool { return a[i].Value > a[j].Value }

// IPbyHits order the repited IP
func IPbyHits(limit int) {
	ips := make([]IPSortedInt, len(TopIP))

	i := 0
	for k, v := range TopIP {
		ips[i] = IPSortedInt{k, v.Hits}
		i++
	}

	sort.Sort(ByHits(ips))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

	fmt.Fprintln(w, fmt.Sprintf("%s \t %s", "IP", "Hits"))

	c := 0
	for _, ip := range ips {

		fmt.Fprintln(w, fmt.Sprintf("%s \t %v", ip.Name, ip.Value))

		c++
		if c >= limit {
			break
		}
	}
	w.Flush()

}

// ByElapsedMedian amount of times found at the logs
type ByElapsedMedian []IPSortedFloat

func (a ByElapsedMedian) Len() int           { return len(a) }
func (a ByElapsedMedian) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByElapsedMedian) Less(i, j int) bool { return a[i].Value > a[j].Value }

// IPbyElapsedMedian order the repited IP
func IPbyElapsedMedian(limit int) {
	ips := make([]IPSortedFloat, len(TopIP))

	i := 0
	for k, v := range TopIP {
		median, _ := stats.Median(v.Elapsed)
		perc80, _ := stats.Percentile(v.Elapsed, 80)
		sum, _ := stats.Sum(v.Elapsed)
		average := sum / float64(len(v.Elapsed))

		ips[i] = IPSortedFloat{k, median, perc80, average}
		i++
	}

	sort.Sort(ByElapsedMedian(ips))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	c := 0
	fmt.Fprintln(w, fmt.Sprintf("%s \t %s \t %s \t %s", "IP", "Median", "Percentile 80", "Average"))
	for _, ip := range ips {

		fmt.Fprintln(w, fmt.Sprintf("%s \t %.10f \t %.10f \t %.10f", ip.Name, ip.Value, ip.Perc80, ip.Average))

		c++
		if c >= limit {
			break
		}
	}
	w.Flush()

}
