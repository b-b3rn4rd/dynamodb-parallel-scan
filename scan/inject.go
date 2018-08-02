//+build wireinject

package scan

import (
	"context"

	"github.com/google/go-cloud/wire"
)

var parallelScannerSet = wire.NewSet(New, Dynamodb, Session, Logger)

func SetupParallelScanner(ctx context.Context, tableName TableName, parallelScanWorkers ParallelScanWorkers, scanItemLimit ScanItemLimit) *ParallelScan {
	panic(wire.Build(
		parallelScannerSet,
	))
}
