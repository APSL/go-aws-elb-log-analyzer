//
package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"time"

	"strings"

	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	// AWSLogs/88888888888/elasticloadbalancing/eu-west-1/2016/08/28/888888888888_elasticloadbalancing_eu-west-1_MY-ELB-NAME_20160828T1615Z_192.168.2.1_58yr1tfa.log
	// AWSLogs/88888888888/elasticloadbalancing/eu-west-1/2017/05/03/888888888888_elasticloadbalancing_eu-west-1_app.MY-ALB-NAME.4b47fffb401dca2f_20170503T0000Z_52.51.226.102_30bo67wj.log.gz
	reFileDate = regexp.MustCompile(".*_(?P<year>[0-9]{4})(?P<month>[0-9]{2})(?P<day>[0-9]{2})T(?P<h>[0-9]{2})(?P<m>[0-9]{2})Z_.*")
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
	Pointer  *FilePointer
	ElbType  string
}

// NewFileLog start the file log process
func NewFileLog(bucket *string, key *string) *FileLog {
	b := *bucket
	name := *key

	file := &FileLog{
		Bucket:   b,
		Key:      name,
		Filename: fmt.Sprintf("%s", path.Base(name)),
		ElbType:  "classic",
	}

	if strings.HasSuffix(file.Filename, ".gz") {
		file.ElbType = "app"
	}

	file.parseDate()
	return file
}

func (f *FileLog) parseDate() {
	reversed := fmt.Sprintf(
		"${%s}-${%s}-${%s}T${%s}:${%s}:00Z", // 2006-01-02T15:04:05Z
		reFileDate.SubexpNames()[1], reFileDate.SubexpNames()[2], reFileDate.SubexpNames()[3], reFileDate.SubexpNames()[4], reFileDate.SubexpNames()[5])
	fomated := reFileDate.ReplaceAllString(f.Key, reversed)
	t, _ := time.Parse(time.RFC3339, fomated)
	f.Date = t
}

// Download and proccess
func (f *FileLog) Download(start time.Time, end time.Time) {

	tmpfilename := fmt.Sprintf(".downloading__%s", path.Base(f.Filename))

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

		fmt.Println(f.Filename)

		if strings.HasSuffix(f.Filename, ".gz") {

			f.Filename = strings.Replace(f.Filename, ".gz", "", -1)
			fmt.Println("Is gz", f.Filename)
			//file.Seek(0, 0)

			zw, _ := gzip.NewReader(file)

			finalFile, err := os.Create(f.Filename)
			if err != nil {
				log.Panic(err)
				return
			}

			io.Copy(finalFile, zw)
			zw.Close()
			finalFile.Close()

			os.Remove(tmpfilename)

		} else {
			file.Close()
			os.Rename(tmpfilename, f.Filename)
		}

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

}
