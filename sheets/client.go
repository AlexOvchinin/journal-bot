package sheets

import (
	"context"
	"encoding/base64"
	"os"

	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetClient interface {
	OpenSpreadsheet(spreadsheetId string) (Spreadsheet, error)
}

type Client struct {
	Service *sheets.Service
}

func NewClient() SpreadsheetClient {
	keyJsonBase64 := os.Getenv("GOOGLE_SA_KEY")
	credBytes, err := base64.StdEncoding.DecodeString(keyJsonBase64)
	if err != nil {
		logger.Error(err)
		return nil
	}

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsJSON(credBytes))
	if err != nil {
		logger.Error(zap.Error(err))
		return nil
	}

	return &Client{
		Service: sheetsService,
	}
}

func (client *Client) OpenSpreadsheet(spreadsheetId string) (Spreadsheet, error) {
	spreadsheet, err := client.Service.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil || spreadsheet.HTTPStatusCode != 200 {
		logger.Errorf("Error while opening spreasheet with id %v", spreadsheetId, zap.Error(err))
		return nil, err
	}

	return &spreadsheetImpl{
		spreadsheet: spreadsheet,
		service:     client.Service,
	}, nil
}
