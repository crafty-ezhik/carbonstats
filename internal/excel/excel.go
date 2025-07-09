package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type Excel interface {
	SetCellValue(cell string, data any) error
	Save() error
}

type excelImpl struct {
	filename    string
	file        *excelize.File
	sheetName   string
	baseStyle   int
	filledStyle int
	log         *zap.Logger
}

func New(logger *zap.Logger, filename string) Excel {
	// Создание файла
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error("Failed to close file", zap.Error(err))
		}
	}()

	// Создание рабочего листа
	sheetName := "Sheet1"
	sheet, err := file.NewSheet(sheetName)
	if err != nil {
		logger.Error("Failed to create new sheet", zap.Error(err))
		return nil
	}

	// Делаем лист активным
	file.SetActiveSheet(sheet)

	// Добавляем к имени файла расширение .xlsx
	filename = fmt.Sprintf("%s.xlsx", filename)

	// Создание стиля
	border := []excelize.Border{
		{Type: BorderLeft, Color: ColorBlack, Style: 1},
		{Type: BorderRight, Color: ColorBlack, Style: 1},
		{Type: BorderTop, Color: ColorBlack, Style: 1},
		{Type: BorderBottom, Color: ColorBlack, Style: 1},
	}

	alignment := &excelize.Alignment{
		Horizontal: AlignmentCenter,
		Vertical:   AlignmentCenter,
		WrapText:   true,
	}

	baseStyle, err := file.NewStyle(&excelize.Style{
		Border:    border,
		Alignment: alignment,
		Font: &excelize.Font{
			Bold:   false,
			Italic: false,
			Family: "Arial",
			Size:   14,
			Color:  ColorBlack,
		},
	})
	if err != nil {
		logger.Error("Failed to create baseStyle style", zap.Error(err))
	}

	// Создание стиля с заливкой
	filledStyle, err := file.NewStyle(&excelize.Style{
		Border:    border,
		Alignment: alignment,
		Font: &excelize.Font{
			Bold:   true,
			Italic: false,
			Family: "Arial",
			Size:   14,
			Color:  ColorBlack,
		},
		Fill: excelize.Fill{
			Type:  "pattern",
			Color: []string{ColorLightGray},
		},
	})
	if err != nil {
		logger.Error("Failed to create filledStyle style", zap.Error(err))
	}

	return &excelImpl{
		log:         logger,
		filename:    filename,
		file:        file,
		baseStyle:   baseStyle,
		filledStyle: filledStyle,
		sheetName:   sheetName,
	}
}

func (ex *excelImpl) SetCellValue(cell string, data any) error {
	return ex.file.SetCellValue(ex.sheetName, cell, data)
}

func (ex *excelImpl) Save() error {
	return ex.file.SaveAs(ex.filename)
}
