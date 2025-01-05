package updater

import (
	"errors"
	"fmn/journalbot/sheets"
	"fmt"
	"strings"
	"time"
)

const (
	// errors
	errorUpdate = "update-error"
)

type Updater struct {
	Client sheets.SpreadsheetClient
}

func NewUpdater() *Updater {
	return &Updater{
		Client: sheets.NewClient(),
	}
}

var loc, _ = time.LoadLocation("Asia/Yekaterinburg")

func (updater *Updater) ProcessUpdate(spreadsheetId string, timestamp int64, content string) error {
	currentDate := time.Unix(timestamp, 0).In(loc).Format(time.DateOnly)
	formattedMessage := formatMessage(timestamp, content)

	spreadsheet, err := updater.Client.OpenSpreadsheet(spreadsheetId)
	if err != nil {
		return errors.New(errorUpdate)
	}

	sheet := spreadsheet.GetFirstSheet()
	emptyCell, err := sheet.GetNextEmptyCell()
	if err != nil {
		return errors.New(errorUpdate)
	}

	row, col, err := sheets.GetCellCoords(emptyCell)
	if err != nil {
		return errors.New(errorUpdate)
	}

	if isUpdatable(sheet, currentDate, row, col) {
		err := sheet.Update(sheets.GetCoordsCell(row-1, col+1), formattedMessage)
		if err != nil {
			return errors.New(errorUpdate)
		}
	} else {
		err := sheet.Append([][]interface{}{{currentDate, formattedMessage}})
		if err != nil {
			return errors.New(errorUpdate)
		}
	}

	return nil
}

func isUpdatable(sheet sheets.Sheet, currentDate string, row int32, col int32) bool {
	if row == 0 {
		return false
	}

	sheetDate, err := sheet.GetContent(sheets.GetCoordsCell(row-1, col))
	if err != nil {
		return false
	}

	return currentDate == sheetDate
}

func formatMessage(timestamp int64, content string) string {
	time := time.Unix(timestamp, 0).In(loc).Format(time.TimeOnly)
	processedContent := strings.ReplaceAll(content, "\n", "\n    ")
	return fmt.Sprintf("[%v]:\n    %v", time, processedContent)
}
