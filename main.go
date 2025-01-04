package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	// "golang.org/x/oauth2/google"
	fmnsheets "fmn/journal/sheets"

	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
)

type Service struct {
	Sheets        map[string]string
	SheetsService *sheets.Service
	SheetsClient  fmnsheets.SpreadsheetClient
}

func loadSheetsMapping() map[string]string {
	result := make(map[string]string)
	result[os.Getenv("FIRST_USER_USERNAME")] = os.Getenv("FIRST_USER_SPREADSHEET_ID")
	return result
}

func initSheetsService() *sheets.Service {
	keyJsonBase64 := os.Getenv("GOOGLE_SA_KEY")
	credBytes, err := base64.StdEncoding.DecodeString(keyJsonBase64)
	if err != nil {
		// log.Error(err)
		fmt.Println(err)
		return nil
	}

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsJSON(credBytes))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return sheetsService
}

const (
	layoutISO = "2006-01-02"
)

func main() {
	service := &Service{
		Sheets:        loadSheetsMapping(),
		SheetsService: initSheetsService(),
		SheetsClient:  fmnsheets.CreateClient(),
	}

	// config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	// if err != nil {
	// 	// log.Error(err)
	// 	fmt.Println(err)
	// 	return
	// }

	// client := config.Client(ctx)

	spreadsheetId := os.Getenv("FIRST_USER_SPREADSHEET_ID")

	// response, err := service.SheetsService.Spreadsheets.Get(spreadsheetId).Do()
	// if err != nil || response.HTTPStatusCode != 200 {
	// 	fmt.Println(err)
	// 	return
	// }

	// row := &sheets.ValueRange{
	// 	Values: [][]interface{}{{"1", "ABC", "abc@gmail.com"}},
	// }

	// response2, err := service.SheetsService.Spreadsheets.Values.
	// 	Append(spreadsheetId, response.Sheets[0].Properties.Title, row).
	// 	ValueInputOption("USER_ENTERED").
	// 	InsertDataOption("INSERT_ROWS").
	// 	Context(context.Background()).
	// 	Do()
	// if err != nil || response2.HTTPStatusCode != 200 {
	// 	fmt.Println(err)
	// 	return
	// }

	date := time.Now().Format(layoutISO)
	content := "I ate some veggies"

	spreadsheet, err := service.SheetsClient.OpenSpreadsheet(spreadsheetId)
	if err != nil {
		fmt.Println(err)
	}

	sheet := spreadsheet.GetSheet(0)
	emptyCell, _ := sheet.GetNextEmptyCell()
	row, col, err := fmnsheets.GetCellCoords(emptyCell)
	if err != nil {
		fmt.Println(err)
	}
	currentContent, _ := sheet.GetContent(fmnsheets.GetCoordsCell(row-1, col))
	if row > 0 && currentContent == date {
		sheet.Update(fmnsheets.GetCoordsCell(row-1, col+1), content)
	} else {
		sheet.Append([][]interface{}{{date, content}})
	}

	// bottomRightCellContent, _ := sheet.GetContent(bottomRightCell)
	// if bottomRightCellContent != date {
	// 	err := sheet.Append(bottomRightCell, content)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// } else {
	// 	err := sheet.Append(bottomRightCell, content)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
	fmt.Println("I'm main")
}
