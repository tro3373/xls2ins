package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func Gen(config Config, files []string) error {
	for _, file := range files {
		bookConfig := config.FindBookConfig(file)
		if bookConfig == nil {
			continue
		}
		err := generateSqlsForBook(config, file, bookConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateSqlsForBook(config Config, file string, bookConfig *BookConfig) (resErr error) {
	log.Infof("==> File:%s", file)

	f, err := excelize.OpenFile(file)
	if err != nil {
		return err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			resErr = err
		}
	}()
	for _, sheetConfig := range bookConfig.SheetConfigs {
		if err := sheetConfig.validate(); err != nil {
			return err
		}
		err := generateSql(sheetConfig, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateSql(sheetConfig SheetConfig, f *excelize.File) error {

	row := sheetConfig.StartRow
	sheet := sheetConfig.SheetName
	format := sheetConfig.SqlFormat

	for {
		noEmptyValueExist := false
		var args []any
		for _, col := range sheetConfig.SqlArgCols {
			cell := fmt.Sprintf("%s%d", col, row)
			val, err := f.GetCellValue(sheet, cell)
			if err != nil {
				return err
			}
			if len(val) > 0 {
				noEmptyValueExist = true
			}
			args = append(args, val)
		}
		if !noEmptyValueExist {
			break
		}
		fmt.Println(fmt.Sprintf(format, args...))
		row += 1
	}
	return nil
}
