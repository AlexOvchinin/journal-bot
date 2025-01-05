package sheets

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/api/sheets/v4"
)

type Sheet interface {
	GetNextEmptyCell() (string, error)
	GetContent(range_ string) (string, error)
	Append(rows [][]interface{}) error
	Update(cell string, content string) error
}

type sheetImpl struct {
	spreadsheet *sheets.Spreadsheet
	service     *sheets.Service
	title       string
}

func (sheet *sheetImpl) GetNextEmptyCell() (string, error) {
	row := &sheets.ValueRange{
		Values: [][]interface{}{{}},
	}

	response, err := sheet.service.Spreadsheets.Values.
		Append(sheet.spreadsheet.SpreadsheetId, sheet.title, row).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		Context(context.Background()).
		Do()
	if err != nil || response.HTTPStatusCode != 200 {
		logger.Error("Error while appending values", zap.Error(err))
		return "", errors.New(errorSpreadsheet)
	}

	return tryRemoveSheetTitle(response.Updates.UpdatedRange, sheet.title), nil
}

func tryRemoveSheetTitle(range_ string, sheetTitle string) string {
	escapedTitle := fmt.Sprintf("'%v'!", sheetTitle)

	if strings.Contains(range_, escapedTitle) {
		_, after, _ := strings.Cut(range_, escapedTitle)
		return after
	}

	title := sheetTitle + "!"
	if strings.Contains(range_, title) {
		_, after, _ := strings.Cut(range_, title)
		return after
	}

	return range_
}

func (sheet *sheetImpl) GetContent(range_ string) (string, error) {
	result, err := sheet.service.Spreadsheets.Values.
		Get(sheet.spreadsheet.SpreadsheetId, fmt.Sprintf("%s!%s", sheet.title, range_)).
		Do()

	if err != nil || result.HTTPStatusCode != 200 {
		logger.Error("Error while getting cells content", zap.Error(err))
		return "", errors.New(errorSpreadsheet)
	}

	return result.Values[0][0].(string), nil
}

func (sheet *sheetImpl) Append(rows [][]interface{}) error {
	row := &sheets.ValueRange{
		Values: rows,
	}

	response2, err := sheet.service.Spreadsheets.Values.
		Append(sheet.spreadsheet.SpreadsheetId, sheet.title, row).
		ValueInputOption("RAW").
		InsertDataOption("INSERT_ROWS").
		Context(context.Background()).
		Do()
	if err != nil || response2.HTTPStatusCode != 200 {
		logger.Error("Error while appending values", zap.Error(err))
		return errors.New(errorSpreadsheet)
	}

	return nil
}

func formatRangeEscaped(sheetTitle string, range_ string) string {
	return fmt.Sprintf("'%s'!%s", sheetTitle, range_)
}

func (sheet *sheetImpl) Update(cell string, content string) error {
	currentContent, err := sheet.GetContent(cell)
	if err != nil {
		logger.Error("Error while getting content", zap.Error(err))
		return errors.New(errorSpreadsheet)
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{currentContent + "\n" + content}},
	}

	response, err := sheet.service.Spreadsheets.Values.
		Update(sheet.spreadsheet.SpreadsheetId, formatRangeEscaped(sheet.title, cell), valueRange).
		ValueInputOption("RAW").
		Do()

	if err != nil || response.HTTPStatusCode != 200 {
		logger.Error("Error while updating values", zap.Error(err))
		return errors.New(errorSpreadsheet)
	}

	return nil
}
