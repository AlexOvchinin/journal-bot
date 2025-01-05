package sheets

import (
	"errors"

	"go.uber.org/zap"
	"google.golang.org/api/sheets/v4"
)

type Spreadsheet interface {
	GetFirstSheet() Sheet
}

type spreadsheetImpl struct {
	spreadsheet *sheets.Spreadsheet
	service     *sheets.Service
}

func (spreadsheet *spreadsheetImpl) GetFirstSheet() Sheet {
	if len(spreadsheet.spreadsheet.Sheets) == 0 {
		logger.Error("Spreadsheet contains no sheets", zap.Error(errors.New("")))
		return nil
	}

	return &sheetImpl{
		spreadsheet: spreadsheet.spreadsheet,
		service:     spreadsheet.service,
		title:       spreadsheet.spreadsheet.Sheets[0].Properties.Title,
	}
}
