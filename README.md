Golang DynamoDB parallel scan
=================================
The following repository provides DynamoDB scan function that performs a parallel scan exampled [here](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Scan.html#Scan.ParallelScan)

Usage
==============

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/b-b3rn4rd/dynamodb-parallel-scan/scan"
)

func main() {
	start := time.Now()

	ctx := context.Background()

	tableName := "me_table_name````"
	parallelScanWorkers := 5
	scanItemLimit := 600

	parallelScanner := scan.SetupParallelScanner(
		ctx,
		scan.TableName(tableName),
		scan.ParallelScanWorkers(parallelScanWorkers),
		scan.ScanItemLimit(scanItemLimit),
	)

	items, err := parallelScanner.Scan()

	if err != nil {
		panic(fmt.Sprintf("Sad times, an epic error has occured, %s", err))
	}

	elapsed := time.Since(start)

	fmt.Println(fmt.Sprintf("Scanning %d took %s", len(items), elapsed))
}

```

Running example locally

```bash

$ cd example
$ go run main.go 
$ {"itemLimit":600,"level":"info","msg":"Scanning items","numberOfThreads":5,"tableName":"me_table_name","time":"2018-08-02T16:08:10+10:00"}
  {"itemLimit":600,"level":"info","msg":"Scanning items from segment","segment":3,"tableName":"me_table_name","time":"2018-08-02T16:08:10+10:00","totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanning items from segment","segment":2,"tableName":"me_table_name","time":"2018-08-02T16:08:10+10:00","totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanning items from segment","segment":0,"tableName":"me_table_name","time":"2018-08-02T16:08:10+10:00","totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanning items from segment","segment":1,"tableName":"me_table_name","time":"2018-08-02T16:08:10+10:00","totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanning items from segment","segment":4,"tableName":"me_table_name","time":"2018-08-02T16:08:10+10:00","totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanned scannedItems from segment","segment":3,"tableName":"me_table_name","time":"2018-08-02T16:08:11+10:00","totalScannedItemCount":493,"totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanned scannedItems from segment","segment":2,"tableName":"me_table_name","time":"2018-08-02T16:08:11+10:00","totalScannedItemCount":442,"totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanned scannedItems from segment","segment":1,"tableName":"me_table_name","time":"2018-08-02T16:08:11+10:00","totalScannedItemCount":486,"totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanned scannedItems from segment","segment":0,"tableName":"me_table_name","time":"2018-08-02T16:08:11+10:00","totalScannedItemCount":514,"totalSegments":5}
  {"itemLimit":600,"level":"info","msg":"Scanned scannedItems from segment","segment":4,"tableName":"me_table_name","time":"2018-08-02T16:08:11+10:00","totalScannedItemCount":525,"totalSegments":5}
  Scanning 2460 took 925.077649ms
```