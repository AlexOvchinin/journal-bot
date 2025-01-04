package updater

import (
	"errors"
	"fmn/journalbot/sheets"
	"time"
)

const (
	layoutISO = "2006-01-02"
	location  = "Asia/Yekaterinburg"
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

func (updater *Updater) ProcessUpdate(spreadsheetId string, timestamp int64, content string) error {
	loc, _ := time.LoadLocation(location)
	currentDate := time.Unix(timestamp, 0).In(loc).Format(layoutISO)

	spreadsheet, err := updater.Client.OpenSpreadsheet(spreadsheetId)
	if err != nil {
		return errors.New(errorUpdate)
	}

	sheet := spreadsheet.GetSheet(0)
	emptyCell, _ := sheet.GetNextEmptyCell()
	row, col, err := sheets.GetCellCoords(emptyCell)
	if err != nil {
		return errors.New(errorUpdate)
	}

	if isUpdatable(sheet, currentDate, row, col) {
		err := sheet.Update(sheets.GetCoordsCell(row-1, col+1), content)
		if err != nil {
			return errors.New(errorUpdate)
		}
	} else {
		err := sheet.Append([][]interface{}{{currentDate, content}})
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
