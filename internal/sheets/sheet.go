package sheets

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gopkg.in/Iwark/spreadsheet.v2"
	"time"
)

const (
	clientSecretPath = "client_secret.json"
)

type (
	GoogleCli struct {
		Svc           *spreadsheet.Service
		SpreadSheetID string
		SheetName     string
		Sheet         *spreadsheet.Sheet
		V4Svc         *sheets.Service
	}

	SpreadsheetPushRequest struct {
		SpreadsheetId string        `json:"spreadsheet_id"`
		Range         string        `json:"range"`
		Values        []interface{} `json:"values"`
	}

	Report struct {
		UserID        string  `json:"userID,omitempty"`
		Username      string  `json:"username,omitempty"`
		Fees          float64 `json:"fees,omitempty"`
		ReportCreated string  `json:"reportCreated,omitempty"`
		ReportPaid    string  `json:"reportPaid,omitempty"`
		UUID          string  `json:"uuid,omitempty"`
	}
)

func NewGoogleClient() (*GoogleCli, error) {
	service, err := spreadsheet.NewService()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	v4Svc, err := sheets.NewService(ctx, option.WithCredentialsFile(clientSecretPath), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return nil, err
	}

	svc := &GoogleCli{
		Svc:           service,
		SpreadSheetID: "1wTk6czcyFBOlO5ak1HNWZEjW4cCPkBsJRPVjhCmbNyY",
		SheetName:     "reports",
		V4Svc:         v4Svc,
	}

	sheet, err := svc.getSheet()
	if err != nil {
		return nil, err
	}

	if sheet != nil {
		svc.Sheet = sheet
	} else {
		return nil, errors.New("Google cli failed to be created ")
	}

	return svc, nil
}

func (s *GoogleCli) UpdateRow(rowID int, columnID int, value string) error {
	s.Sheet.Update(rowID, columnID, value)
	err := s.Sheet.Synchronize()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *GoogleCli) getSheet() (*spreadsheet.Sheet, error) {
	spreadSheet, err := s.Svc.FetchSpreadsheet(s.SpreadSheetID)
	if err != nil {
		return nil, err
	}

	// get a sheet by the title.
	sheet, err := spreadSheet.SheetByTitle(s.SheetName)
	if err != nil {
		return nil, err
	}

	return sheet, nil
}

func (s *GoogleCli) InsertNewRows(reports []*Report) error {
	for _, report := range reports {
		reportCreated := time.Now().Format(time.RFC3339)
		req := &SpreadsheetPushRequest{
			SpreadsheetId: s.SpreadSheetID,
			Range:         "reports",
			Values:        []interface{}{report.UUID, report.Username, report.UserID, report.Fees, reportCreated, "false"},
		}

		err := s.WriteToSpreadsheet(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (s *GoogleCli) WriteToSpreadsheet(object *SpreadsheetPushRequest) error {
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, object.Values)

	res, err := s.V4Svc.Spreadsheets.Values.Append(object.SpreadsheetId, object.Range, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		fmt.Println("Unable to update data to sheet  ", err)
	}
	fmt.Println("spreadsheet push ", res)

	return nil
}

func (s *GoogleCli) UpdateRowWithPaymentDate(uuid string) error {
	reportPaidCol := 5
	uuidCol := 0
	recordRow := 0

	for _, row := range s.Sheet.Rows {
		if recordRow != 0 {
			break
		}
		for c, data := range row {
			if c > uuidCol {
				break
			}
			if data.Value == uuid {
				recordRow = int(data.Row)
			}
		}
	}

	err := s.UpdateRow(recordRow, reportPaidCol, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}
