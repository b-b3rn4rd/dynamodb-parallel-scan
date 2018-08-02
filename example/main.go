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

	tableName := "me_table_name"
	parallelScanWorkers := 1
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
