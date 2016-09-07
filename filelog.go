//
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// FileLog general structure
type FileLog struct {
	Bucket   string
	Key      string
	Filename string
	Date     time.Time
	Cache    string
	Start    time.Time
	End      time.Time
}

// NewFileLog start the file log process
func NewFileLog(bucket *string, key *string) *FileLog {
	b := *bucket
	name := *key

	file := &FileLog{
		Bucket:   b,
		Key:      name,
		Filename: fmt.Sprintf("/tmp/%s", path.Base(name)),
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

	tmpfilename := fmt.Sprintf("/tmp/.downloading__%s", path.Base(f.Filename))

	// Return in case the file exists
	if _, err := os.Stat(f.Filename); err == nil {
		return
	}

	// Delete partial file in case
	if _, err := os.Stat(tmpfilename); err == nil {
		os.Remove(tmpfilename)
		log.Printf("Deleted partial file %s", tmpfilename)
	}

	file, err := os.Create(tmpfilename)
	if err != nil {
		log.Fatal("Failed to create file", err)
	}
	defer func() {
		file.Close()
		os.Rename(tmpfilename, f.Filename)
		log.Printf("Downloaded %s", f.Filename)
	}()

	downloader := s3manager.NewDownloaderWithClient(SVC, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
	})

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(f.Bucket), // Required
			Key:    aws.String(f.Key),    // Required
		})
	if err != nil {
		fmt.Println("Failed to download file", err)
		return
	}

	/***
	params := &s3.GetObjectInput{
		Bucket: aws.String(f.Bucket), // Required
		Key:    aws.String(f.Key),    // Required
	}
	resp, _ := SVC.GetObject(params)

	fmt.Printf("Body %v", resp.Body)

	count := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {


		line := NewLineLog(scanner.Text())

		if InTimeSpan(start, end, line.Time) {
			RecordIP(line.IPClient, 1, line.Elapsed)
			count++
		}
	}

	log.Printf("%d lines where processed", count)
	**/

}
