# go-aws-elb-log-analyzer #

This project was create to analyse the logs of ELB stored in S3. I parse and process the logs in memory. Don't need to download the log files to the local hard drive or send by a pipe.

Why develop new tool in go? we need something simple and fast, easy like execute a binary and the posibility to deploy this binary anywhere without  dependency issues. Just copy and execute.

It's compatible with credentials of AWS-CLI http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-multiple-profiles

## Example report ##

```
2016/08/30 16:23:41 Time Range: 2016-08-17 01:00:00 +0000 +0000 - 2016-08-17 07:08:00 +0000 +0000
2016/08/30 16:23:41 Bucket: ?.....?
2016/08/30 16:23:43 Loading..
2016/08/30 16:34:10 81464 lines where processed
2016/08/30 16:34:10 Top clients by hits
          IP | Hits
  332.18.?.? | 529914
  337.89.?.? | 137409
 303.183.?.? | 119717
  332.68.?.? | 40220
  332.79.?.? | 23825
 322.122.?.? | 23262
  332.79.?.? | 23185
 322.122.?.? | 22862
 322.122.?.? | 19040
  354.64.?.? | 17315
2016/08/30 16:34:10 Top of slowest clients
              IP |         Median |  Percentile 80 | Average
    303.82.??.?? |  22.1508900000 |  25.6659920000 | 17.9154478571
    375.99.??.?? |  16.2751465000 |  21.2413640000 | 15.8590785000
    374.37.??.?? |   9.6639520000 |   9.6639520000 | 10.6361656667
    352.25.??.?? |   8.6189680000 |  11.1341940000 | 6.0397899639
    369.57.??.?? |   8.2205860000 |  16.9552370000 | 10.5844959695
   321.138.??.?? |   7.8579020000 |   9.6017140000 | 5.4473516061
   310.177.??.?? |   6.7928660000 |   9.0229420000 | 6.5690146818
    350.97.??.?? |   6.6750390000 |   7.2816815000 | 6.7309778000
    369.57.??.?? |   6.5800730000 |  15.0966750000 | 9.9973729459
   322.122.??.?? |   6.5737150000 |  16.2252460000 | 9.6619743944
```

## How to use ##

```bash
./goelbanalyzer \
  --profile="my_project" \
  --region="eu-west-1" \
  --bucket="my-elb-bucket-logs" \
  --prefix="eng/AWSLogs/9999999999999/elasticloadbalancing/eu-west-1" \
  --start="2016-08-17 01:00:00 +0000" \
  --end="2016-08-17 07:08:00 +0000"
```

# How to compile #
The final binary name will be "goelbanalyzer"
```bash
go build -o goelbanalyzer *.go
```