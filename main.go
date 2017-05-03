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

// Analyze flag
var analyze bool
var joinFiles bool
var filename string

func main() {

	var awsProfile string
	flag.StringVar(&awsProfile, "profile", "", "Profile credentials used by aws cli")

	var awsRegion string
	flag.StringVar(&awsRegion, "region", "eu-west-1", "Name of the region in AWS")

	var awsBucket string
	flag.StringVar(&awsBucket, "bucket", "", "Name of the S3 bucket")

	var strStart string
	flag.StringVar(&strStart, "start", "", "Date and time to start the download. Example: 2016-09-16 05:00:00 +0000")

	var strEnd string
	flag.StringVar(&strEnd, "end", "1h", "Time after start, example of 30 minutes: 30m")

	var strPrefix string
	flag.StringVar(&strPrefix, "prefix", "", "Prefix or folder used in S3 bucket")

	flag.StringVar(&filename, "filename", "complete.log", "Name of the final log file")

	flag.BoolVar(&joinFiles, "join", true, "Contact and sort all logs in one file")

	flag.BoolVar(&analyze, "analyze", false, "Analyze the logs to find top requests and top slow requests")

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

	if joinFiles != false || analyze != false {
		AnalyzerDispatch(start, end)
	}

	log.Printf("Time Range: %s - %s", start.String(), end.String())

	// Start S3 file reading
	s3page(SVC, awsBucket, 100, strPrefix, start, end, nil)

	if joinFiles != false || analyze != false {
		close(AnalyzerQueue)
		AnalyzerFinished()

		saveSortedLog(filename)
	}

	if analyze {
		fmt.Println("")
		fmt.Println("***** TOP by hists")
		PrintBy(20, "hits")
		fmt.Println("")
		fmt.Println("***** TOP by median latency")
		PrintBy(20, "median")
	}

}

func s3page(SVC *s3.S3, bucket string, maxkeys int64, prefix string, start time.Time, end time.Time, NextToken *string) {
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
				s3page(SVC, bucket, 100, prefix, start, end, resp.NextContinuationToken)
				return
			}
			index = index + len(searchkey)
			prefix = fmt.Sprintf("%s%d/%02d/%02d", file.Key[0:index], start.Year(), start.Month(), start.Day())

			log.Printf("SET Prefix: %s\n", prefix)
		}

		if InTimeSpan(start, end, file.Date) {
			log.Printf("Reading %s\n", file.Key)
			file.Download(start, end)
			if joinFiles != false || analyze != false {
				AnalyzerQueue <- file
			}
		}
	}

	// If NextContinuationToken it's empty, no more content to load
	if resp.NextContinuationToken == nil {
		return
	}

	// Continue with the next "page". It's like click in "more" or scroll down
	s3page(SVC, bucket, 100, prefix, start, end, resp.NextContinuationToken)
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
