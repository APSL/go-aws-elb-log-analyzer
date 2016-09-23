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
	TopMutex.Lock()
	defer TopMutex.Unlock()

	if _, ok := TopIP[ip.String()]; ok {
		TopIP[ip.String()].AddHit()
		TopIP[ip.String()].AddElapsed(elapsed)
	} else {
		TopIP[ip.String()] = &IPRecord{Hits: 1}
		TopIP[ip.String()].AddElapsed(elapsed)
	}
}

// IPSorted export
type IPSorted struct {
	Name       string
	ByColumn   string
	Hits       int64
	Median     float64
	Percentile float64
	Average    float64
}

func (p IPSorted) String() string {
	return fmt.Sprintf("%s", p.Name)
}

// BuildIPSorter amount of times found at the logs
type BuildIPSorter []IPSorted

func (a BuildIPSorter) Len() int      { return len(a) }
func (a BuildIPSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BuildIPSorter) Less(i, j int) bool {

	switch a[i].ByColumn {
	case "median":
		return a[i].Median > a[j].Median
	case "percentile":
		return a[i].Percentile > a[j].Percentile
	default:
		return a[i].Hits > a[j].Hits
	}
}

// IPbyHits order by hits
func createTable(byColumn string) []IPSorted {
	ips := make([]IPSorted, len(TopIP))

	i := 0
	for k, v := range TopIP {
		hits := v.Hits
		median, _ := stats.Median(v.Elapsed)
		percentile, _ := stats.Percentile(v.Elapsed, 95)
		sum, _ := stats.Sum(v.Elapsed)
		average := sum / float64(len(v.Elapsed))

		ips[i] = IPSorted{
			Name:       k,
			Hits:       hits,
			Median:     median,
			Percentile: percentile,
			Average:    average,
			ByColumn:   byColumn,
		}
		i++
	}

	return ips
}

func printTable(ips []IPSorted, limit int) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	c := 0
	fmt.Fprintln(w, fmt.Sprintf("%s \t %s \t %s \t %s \t %s", "IP", "Hits", "Median latency", "Percentile 90 latency", "Average latency"))
	for _, ip := range ips {

		fmt.Fprintln(w, fmt.Sprintf("%s \t %d \t %.10f \t %.10f \t %.10f", ip.Name, ip.Hits, ip.Median, ip.Percentile, ip.Average))

		c++
		if c >= limit {
			break
		}
	}
	w.Flush()
}

// PrintBy order by byColumn (median, percentile..)
func PrintBy(limit int, byColumn string) {

	ips := createTable(byColumn)

	sort.Sort(BuildIPSorter(ips))

	printTable(ips, limit)

}
