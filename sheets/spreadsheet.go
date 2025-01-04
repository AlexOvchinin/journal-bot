package sheets

import (
	"fmt"

	"google.golang.org/api/sheets/v4"
)

type Spreadsheet interface {
	GetSheet(index int) Sheet
}

type spreadsheetImpl struct {
	spreadsheet *sheets.Spreadsheet
	service     *sheets.Service
}

func (spreadsheet *spreadsheetImpl) GetSheet(index int) Sheet {
	if len(spreadsheet.spreadsheet.Sheets) == 0 {
		fmt.Println("Spreadsheet contains no sheets")
		return nil
	}

	return &sheetImpl{
		spreadsheet: spreadsheet.spreadsheet,
		service:     spreadsheet.service,
		title:       spreadsheet.spreadsheet.Sheets[0].Properties.Title,
	}
}
