package scan

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/sirupsen/logrus"
)

type DynamodDBItems []map[string]*dynamodb.AttributeValue

type TableName string
type ParallelScanWorkers int
type ScanItemLimit int

type ParallelScan struct {
	logger              *logrus.Logger
	svc                 dynamodbiface.DynamoDBAPI
	scanItemLimit       int
	parallelScanWorkers int
	tableName           string
}

func New(svc dynamodbiface.DynamoDBAPI, tableName TableName, parallelScanWorkers ParallelScanWorkers, scanItemLimit ScanItemLimit, logger *logrus.Logger) *ParallelScan {
	return &ParallelScan{
		logger:              logger,
		svc:                 svc,
		tableName:           string(tableName),
		parallelScanWorkers: int(parallelScanWorkers),
		scanItemLimit:       int(scanItemLimit),
	}
}

func (p *ParallelScan) Scan() (DynamodDBItems, error) {
	var wg sync.WaitGroup

	scannedSegmentItems := make(chan DynamodDBItems)
	scannedItems := make(DynamodDBItems, 0)

	fields := logrus.Fields{
		"tableName":       p.tableName,
		"numberOfThreads": p.parallelScanWorkers,
		"itemLimit":       p.scanItemLimit,
	}

	p.logger.WithFields(fields).Info("Scanning items")

	totalSegments := p.parallelScanWorkers

	for segment := 0; segment < totalSegments; segment++ {
		wg.Add(1)
		go func(segment int) {
			defer wg.Done()

			totalScannedItemCount := 0

			fields := logrus.Fields{
				"tableName":     p.tableName,
				"segment":       segment,
				"itemLimit":     p.scanItemLimit,
				"totalSegments": totalSegments,
			}

			p.logger.WithFields(fields).Info("Scanning items from segment")

			err := p.svc.ScanPages(&dynamodb.ScanInput{
				TableName:     aws.String(p.tableName),
				Limit:         aws.Int64(int64(p.scanItemLimit)),
				TotalSegments: aws.Int64(int64(totalSegments)),
				Segment:       aws.Int64(int64(segment)),
			}, func(page *dynamodb.ScanOutput, lastPage bool) bool {
				scannedSegmentItems <- page.Items
				totalScannedItemCount += int(aws.Int64Value(page.Count))
				return !lastPage
			})

			if err != nil {
				logrus.WithError(err).Error("error while scanning pages")
				return
			}

			fields["totalScannedItemCount"] = totalScannedItemCount
			p.logger.WithFields(fields).Info("Scanned scannedItems from segment")

		}(segment)
	}

	go func() {
		wg.Wait()
		close(scannedSegmentItems)
	}()

	for itemsInSegment := range scannedSegmentItems {
		scannedItems = append(scannedItems, itemsInSegment...)
	}

	return scannedItems, nil
}

func Session() *session.Session {
	return session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		},
	))
}

func Dynamodb(session *session.Session) dynamodbiface.DynamoDBAPI {
	return dynamodb.New(session)
}

func Logger() *logrus.Logger {
	logger := logrus.New()
	logrus.SetLevel(logrus.DebugLevel)
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.DebugLevel

	return logger
}
