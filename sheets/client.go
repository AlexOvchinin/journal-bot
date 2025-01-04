package sheets

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient interface {
	OpenSpreadsheet(spreadsheetId string) (Spreadsheet, error)
}

type Client struct {
	Service *sheets.Service
}

func CreateClient() SpreadsheetClient {
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

	return &Client{
		Service: sheetsService,
	}
}

func (client *Client) OpenSpreadsheet(spreadsheetId string) (Spreadsheet, error) {
	spreadsheet, err := client.Service.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil || spreadsheet.HTTPStatusCode != 200 {
		fmt.Println(err)
		return nil, err
	}

	return &spreadsheetImpl{
		spreadsheet: spreadsheet,
		service:     client.Service,
	}, nil
}
