//
package main

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// FileLog general structure
type FileLog struct {
	Bucket string
	Key    string
	Date   time.Time
}

// NewFileLog start the file log process
func NewFileLog(bucket *string, key *string) *FileLog {
	b := *bucket
	name := *key

	file := &FileLog{
		Bucket: b,
		Key:    name,
	}
	file.parseDate()
	return file
}

func (f *FileLog) parseDate() {

	//fmt.Println(resp.Contents[f])
	// AWSLogs/88888888888/elasticloadbalancing/eu-west-1/2016/08/28/888888888888_elasticloadbalancing_eu-west-1_MY-ELB-NAME_20160828T1615Z_192.168.2.1_58yr1tfa.log
	re := regexp.MustCompile(".*_(?P<year>[0-9]{4})(?P<month>[0-9]{2})(?P<day>[0-9]{2})T(?P<h>[0-9]{2})(?P<m>[0-9]{2})Z_.*.log")
	reversed := fmt.Sprintf(
		"${%s}-${%s}-${%s}T${%s}:${%s}:00Z", // 2006-01-02T15:04:05Z
		re.SubexpNames()[1], re.SubexpNames()[2], re.SubexpNames()[3], re.SubexpNames()[4], re.SubexpNames()[5])
	fomated := re.ReplaceAllString(f.Key, reversed)
	t, _ := time.Parse(time.RFC3339, fomated)
	f.Date = t
}

// Download and proccess
func (f *FileLog) Download(start time.Time, end time.Time) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(f.Bucket), // Required
		Key:    aws.String(f.Key),    // Required
	}
	resp, _ := SVC.GetObject(params)

	count := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {

		/**
		  2016-08-27T23:56:06.983879Z
		  MY-ELB-NAME
		  8.8.8.8:47769
		  192.168.1.1:80
		  0.00002
		  0.213758
		  0.000023
		  200
		  200
		  928
		  832
		  "POST http://www.example:80/somepath/index.php?mobile=true HTTP/1.1"
		  "-"
		  -
		  -
		  **/

		line := NewLineLog(scanner.Text())

		if InTimeSpan(start, end, line.Time) {
			RecordIP(line.IPClient, 1, line.Elapsed)
			count++
		}
	}

	log.Printf("%d lines where processed", count)

}
