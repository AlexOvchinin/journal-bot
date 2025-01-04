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

	filledDateContent, _ := sheet.GetContent(sheets.GetCoordsCell(row-1, col))
	if row > 0 && filledDateContent == currentDate {
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
