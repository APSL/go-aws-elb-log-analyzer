# Download logs from ELB and concat all of the in one file in order #

This project was created to analyze ELB logs stored in S3. Logs are pared and processed in memory. Don't need to download the log files to the local hard drive or send by a pipe.

Why develop new tool in go? We need something simple and fast, easy like execute a binary with the possibility to deploy this binary anywhere without any dependency issues. Just copy and execute.

It's compatible with credentials of AWS-CLI http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-multiple-profiles

## How to use ##

In case your ELB has a dedicated bucket

```bash
goelbanalyzer \
  --profile="my_project" \
  --bucket="my-elb-bucket-logs" \
  --start="2016-08-17 01:00:00 +0000" \
  --end="2h"
```

## Concat and sort file ##

After locate and download all the files whitch are in the time frame, this software will join in one file and sort it by date. So you will have a uniq file with all the records in order.

Finale filename **complete.log**

To change the name of the file you can use "-filename", example:

```bash
goelbanalyzer \
  --profile="my_project" \
  --bucket="my-elb-bucket-logs" \
  --start="2016-08-17 01:00:00 +0000" \
  --end="2h" \
  --filename="myfull.log"
```

To disable the contact, use '--join 0' in the parameters.


## Example report ##

### Command ###

In case your ELB has a dedicated bucket

```bash
goelbanalyzer \
  --profile="my_project" \
  --bucket="my-elb-bucket-logs" \
  --start="2016-08-17 01:00:00 +0000" \
  --end="2h"
  --analyse
```

```
2016/08/30 16:23:41 Time Range: 2016-08-17 01:00:00 +0000 +0000 - 2016-08-17 07:08:00 +0000 +0000
2016/08/30 16:23:41 Bucket: ?.....?
2016/08/30 16:23:43 Loading..
2016/08/30 16:34:10 81464 lines where processed
2016/08/30 16:34:10 Top clients by hits
***** TOP by hists
              IP |   Hits |  Median latency |  Percentile 90 latency | Average latency
    52.18.999.11 |  24438 |    0.1749155000 |          79.1077260000 | 13.1908529439
    47.89.999.11 |   4542 |    1.1352130000 |           8.6127860000 | 2.5885444800
  203.183.999.11 |   2157 |    0.1801270000 |          78.2074270000 | 13.3125795345
    52.68.999.11 |   1848 |    0.1780620000 |           1.4558380000 | 0.5151559892
  120.132.999.11 |    954 |    0.1222715000 |           1.1246040000 | 0.2924876447
    13.95.999.11 |    936 |    4.7280725000 |           9.3240620000 | 4.1887639006
  120.132.999.11 |    897 |    0.1253540000 |           1.1191700000 | 0.3108006778
  120.132.999.11 |    889 |    0.1248950000 |           1.1063860000 | 0.2765490641
  120.132.999.11 |    782 |    0.1207545000 |           1.1104930000 | 0.2850222621
   169.57.999.11 |    704 |   68.1784910000 |         104.9732110000 | 65.1308670170
  218.244.999.11 |    628 |    5.6638660000 |          10.0451020000 | 5.6309517452
   169.57.999.11 |    438 |   70.5174620000 |         102.7895760000 | 66.3744711644
  222.122.999.11 |    382 |    0.1758160000 |          74.2409350000 | 15.1374491204
  222.122.999.11 |    303 |    2.0829630000 |           8.6622070000 | 2.7980553102
  213.182.999.11 |    300 |   79.8319910000 |         107.8641890000 | 72.3946263033
  222.122.999.11 |    266 |    0.5016970000 |          75.3529030000 | 16.5369312406
  133.242.999.11 |    255 |    5.9971280000 |          10.6584430000 | 6.3336784039
  153.120.999.11 |    247 |    5.1963570000 |           8.8504320000 | 4.9291020162
   175.99.999.11 |    198 |    2.5931665000 |           8.9435810000 | 3.8328212879
  222.122.999.11 |    177 |    0.0503840000 |          51.6966140000 | 6.0209316102

***** TOP by median latency
              IP  |  Hits |  Median latency |  Percentile 90 latency | Average latency
   182.213.999.11 |   300 |   79.8319910000 |         107.8641890000 | 72.3946263033
   213.182.999.11 |    82 |   71.9786335000 |         103.3641420000 | 65.9088648171
    169.57.999.11 |   438 |   70.5174620000 |         102.7895760000 | 66.3744711644
    169.57.999.11 |   704 |   68.1784910000 |         104.9732110000 | 65.1308670170
   213.182.999.11 |    71 |   57.0327860000 |         109.0471830000 | 55.9872595070
   213.182.999.11 |    31 |   54.0556680000 |         101.3834090000 | 52.8765958710
    121.78.999.11 |     2 |   38.3978040000 |          69.1167380000 | 38.3978040000
   213.182.999.11 |    53 |   37.4184290000 |          94.6283570000 | 42.5396707547
   222.122.999.11 |     3 |   36.7208290000 |          47.4126790000 | 28.1037163333
   125.141.999.11 |    42 |   15.6014080000 |          81.4945850000 | 26.3813173810
   125.141.999.11 |    76 |   13.3390590000 |         103.4595620000 | 28.4127705921
   221.139.999.11 |     1 |   11.0201320000 |          11.0201320000 | 11.0201320000
     65.52.999.11 |     1 |    8.5374680000 |           8.5374680000 | 8.5374680000
   218.213.999.11 |     1 |    8.1529430000 |           8.1529430000 | 8.1529430000
    119.81.999.11 |     3 |    7.5322690000 |           8.5617910000 | 7.1656163333
   182.162.999.11 |    22 |    6.2606350000 |           9.7995030000 | 6.6234874091
   133.242.999.11 |   255 |    5.9971280000 |          10.6584430000 | 6.3336784039
   218.153.999.11 |     4 |    5.6652730000 |           8.5056560000 | 5.0062797500
   218.244.999.11 |   628 |    5.6638660000 |          10.0451020000 | 5.6309517452
   153.120.999.11 |   247 |    5.1963570000 |           8.8504320000 | 4.9291020162
```

# How to install #

Download the bin/goelbanalyzer and copy in /usr/local/bin
```
chmod +x /usr/local/bin/goelbanalyzer
```

# How to compile #
The final binary name will be "goelbanalyzer"
```bash
go build -o goelbanalyzer *.go
```

## Help args ##

```
Usage of goelbanalyzer:
  -analyze
    	Analyze the logs to find top requests and top slow requests
  -bucket string
    	Name of the S3 bucket
  -end string
    	Time after start, example of 30 minutes: 30m (default "1h")
  -filename string
    	Name of the final log file (default "complete.log")
  -join
    	Contact and sort all logs in one file (default true)
  -prefix string
    	Prefix or folder used in S3 bucket
  -profile string
    	Profile credentials used by aws cli
  -region string
    	Name of the region in AWS (default "eu-west-1")
  -start string
    	Date and time to start the download. Example: 2016-09-16 05:00:00 +0000
```