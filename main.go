package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// SVC is a global session to S3
var SVC *s3.S3

// Dirname where save the downloaded logs
var Dirname string

func main() {

	flag.StringVar(&Dirname, "dir", "/tmp", "a string var")

	var awsProfile string
	flag.StringVar(&awsProfile, "profile", "", "a string var")

	var awsRegion string
	flag.StringVar(&awsRegion, "region", "eu-west-1", "a string var")

	var awsBucket string
	flag.StringVar(&awsBucket, "bucket", "", "a string var")

	var strStart string
	flag.StringVar(&strStart, "start", "", "a string var")

	var strEnd string
	flag.StringVar(&strEnd, "end", "1h", "a string var")

	var strPrefix string
	flag.StringVar(&strPrefix, "prefix", "", "a string var")

	var strSave string
	flag.StringVar(&strSave, "save", "complete.log", "a string var")

	var top int
	flag.IntVar(&top, "top", 10, "a int var")

	var analyze bool
	flag.BoolVar(&analyze, "analyze", false, "a bool var")

	flag.Parse()

	// Specify profile to load for the session's config
	sess, _ := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
	})

	InitRecords()

	// Select region
	SVC = s3.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	start, _ := time.Parse("2006-01-02 15:04:05 -0700", strStart)
	margin, _ := time.ParseDuration(strEnd)
	end := start.Add(margin)

	if analyze {
		AnalyzerDispatch(start, end)
	}

	log.Printf("Time Range: %s - %s", start.String(), end.String())

	// awsPrefix = fmt.Sprintf("%s/%d/%02d/%02d", awsPrefix, start.Year(), start.Month(), start.Day())

	// log.Printf("Bucket: %s/%s", awsBucket, awsPrefix)

	// Start S3 file reading
	s3page(SVC, awsBucket, 100, "", start, end, analyze, nil)

	if analyze {
		AnalyzerQueue <- []byte(nil)
		AnalyzerFinished()

		PrintSortLog(strSave)
	}

}

func s3page(SVC *s3.S3, bucket string, maxkeys int64, prefix string, start time.Time, end time.Time, analyze bool, NextToken *string) {
	params := &s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket), // Required
		MaxKeys:           aws.Int64(maxkeys),
		Prefix:            aws.String(prefix),
		ContinuationToken: NextToken,
	}
	resp, err := SVC.ListObjectsV2(params)

	region := *SVC.Config.Region

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	log.Println("Loading..")

	// Pretty-print the response data.
	for f := range resp.Contents {
		file := NewFileLog(params.Bucket, resp.Contents[f].Key)

		if prefix == "" {
			searchkey := "/elasticloadbalancing/" + region + "/"
			index := strings.Index(file.Key, searchkey)
			if index <= 0 {
				// Continue with the next "page". It's like click in "more" or scroll down
				s3page(SVC, bucket, 100, prefix, start, end, analyze, resp.NextContinuationToken)
				return
			}
			index = index + len(searchkey)
			prefix = fmt.Sprintf("%s%d/%02d/%02d", file.Key[0:index], start.Year(), start.Month(), start.Day())

			log.Printf("SET Prefix: %s\n", prefix)
		}

		if InTimeSpan(start, end, file.Date) {
			log.Printf("Reading %s\n", file.Key)
			file.Download(start, end)
			if analyze {
				AnalyzerQueue <- []byte(file.Filename)
			}
		}
	}

	// If NextContinuationToken it's empty, no more content to load
	if resp.NextContinuationToken == nil {
		return
	}

	// Continue with the next "page". It's like click in "more" or scroll down
	s3page(SVC, bucket, 100, prefix, start, end, analyze, resp.NextContinuationToken)
}

// InTimeSpan if the record it's in the time range
func InTimeSpan(start, end, check time.Time) bool {
	if check.Equal(start) {
		return true
	}
	if check.Equal(end) {
		return true
	}
	return check.After(start) && check.Before(end)
}
