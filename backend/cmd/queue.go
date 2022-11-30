package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/model"
)

var GenerateReport = func() {
	go func() {
		for range time.Tick(time.Second * 5) {
			reportUnprocessed, err := container.ReportRepository.GetReportByStatus(model.StatusUnprocessedReport)
			if err != nil {
				continue
			}

			for i := 0; i < len(reportUnprocessed); i++ {
				err := container.ReportRepository.UpdateStatusReport(reportUnprocessed[i].ID.Hex(), model.StatusProcessingReport)
				if err != nil {
					RemoveErrorReport(reportUnprocessed, i)
				}
			}

			for i := 0; i < len(reportUnprocessed); i++ {
				err := HandleReport(reportUnprocessed[i])
				if err != nil {
					container.Logger().Errorf(err.Error())
					container.ReportRepository.UpdateStatusReport(reportUnprocessed[i].ID.Hex(), model.StatusUnprocessedReport)
					continue
				}

				container.Logger().Infof("handle report: %s", reportUnprocessed[i].Name)
				container.ReportRepository.UpdateStatusReport(reportUnprocessed[i].ID.Hex(), model.StatusProcessedReport)
			}
		}
	}()
}

func HandleReport(report model.Report) error {
	buf := &bytes.Buffer{}
	wr := csv.NewWriter(buf)

	e := wr.Write([]string{"Date", "Longitude", "Latitude", "Address", "TimeStamp", "Speed"})
	if e != nil {
		return e
	}
	wr.Flush()
	if err := wr.Error(); err != nil {
		return err
	}

	dataChan, errChan := container.ActivityLogRepository.FilterReport(report.Filter)
	running := true
	for {
		if !running {
			break
		}

		select {
		case data, ok := <-dataChan:
			if !ok {
				running = false
				continue
			}

			oneRow := []string{
				data.Date.String(),
				fmt.Sprintf("%f", data.Longitude),
				fmt.Sprintf("%f", data.Latitude),
				data.Address,
				fmt.Sprintf("%d", data.TimeStamp),
				data.Speed,
			}
			if err := wr.Write(oneRow); err != nil {
				buf = &bytes.Buffer{}
				running = false
			}

			wr.Flush()
			if err := wr.Error(); err != nil {
				buf = &bytes.Buffer{}
				running = false
			}
		case err, ok := <-errChan:
			if err != nil || !ok {
				buf = &bytes.Buffer{}
			}
			running = false
		}
	}

	if buf.Len() < 1 {
		return fmt.Errorf("empty file")
	}

	return os.WriteFile(report.Filename, buf.Bytes(), 0o644)
}

func RemoveErrorReport(data []model.Report, index int) []model.Report {
	if data == nil || len(data) == 0 {
		return make([]model.Report, 0)
	}

	result := make([]model.Report, 0)
	for i := 0; i < len(data); i++ {
		if i == index {
			continue
		}
		result = append(result, data[i])
	}
	return result
}
