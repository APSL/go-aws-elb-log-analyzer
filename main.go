package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// SVC is a global session to S3
var SVC *s3.S3

// IPRecord test
type IPRecord struct {
	Hits    int64
	Elapsed []float64
}

// TopIP is a counter of IP clients
var TopIP map[string]*IPRecord

func main() {

	var awsProfile string
	flag.StringVar(&awsProfile, "profile", "", "a string var")

	var awsRegion string
	flag.StringVar(&awsRegion, "region", "eu-west-1", "a string var")

	var awsBucket string
	flag.StringVar(&awsBucket, "bucket", "", "a string var")

	var awsPrefix string
	flag.StringVar(&awsPrefix, "prefix", "", "a string var")

	var strStart string
	flag.StringVar(&strStart, "start", "", "a string var")

	var strEnd string
	flag.StringVar(&strEnd, "end", "", "a string var")

	var limit int
	flag.IntVar(&limit, "limit", 10, "a int var")

	flag.Parse()

	TopIP = make(map[string]*IPRecord)

	// Specify profile to load for the session's config
	sess, _ := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
	})

	// Select region
	SVC = s3.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	start, _ := time.Parse("2006-01-02 15:04:05 -0700", strStart)
	end, _ := time.Parse("2006-01-02 15:04:05 -0700", strEnd)

	log.Printf("Time Range: %s - %s", start.String(), end.String())
	awsPrefix = fmt.Sprintf("%s/%d/%02d/%02d", awsPrefix, start.Year(), start.Month(), start.Day())

	log.Printf("Bucket: %s/%s", awsBucket, awsPrefix)

	// Start S3 file reading
	s3page(SVC, awsBucket, awsPrefix, start, end, nil)

	// Print results
	log.Println("Top clients by hits")
	IPbyHits(limit)

	// Print results
	log.Println("Top of slowest clients")
	IPbyElapsedMedian(limit)

}

func s3page(SVC *s3.S3, bucket string, prefix string, start time.Time, end time.Time, NextToken *string) {
	params := &s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket), // Required
		MaxKeys:           aws.Int64(1000),
		Prefix:            aws.String(prefix),
		ContinuationToken: NextToken,
	}
	resp, err := SVC.ListObjectsV2(params)

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
		//fmt.Printf("    File %v\n", file.Date)
		if InTimeSpan(start, end, file.Date) {
			log.Printf("Reading %s\n", file.Key)
			file.Download(start, end)
		}
	}

	// If NextContinuationToken it's empty, no more content to load
	if resp.NextContinuationToken == nil {
		return
	}

	// Continue with the next "page". It's like click in "more" or scroll down
	s3page(SVC, bucket, prefix, start, end, resp.NextContinuationToken)
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
