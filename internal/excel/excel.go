package excel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Excel interface {
	SetCellValue(cell string, data any, fillingRequired bool) error
	SetHeader(month time.Month) error
	AddData(data *Rows) error
	addRow(rowNum int, data Row) error
	Save() error
}

type excelImpl struct {
	filename    string
	file        *excelize.File
	sheetName   string
	baseStyle   int
	filledStyle int
	columnNames []string
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
			Size:   11,
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
			Size:   11,
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
	columnNames := []string{
		"Клиенты", "Минуты", "Сумма за минуты без НДС", "Номера, кол-во/шт", "Описание дополнительных услуг",
		"Сумма за доп услуги, без НДС", "Сумма за доп услуги, с НДС", "Итоговая сумма, без НДС", "Итоговая сумма, с НДС",
		"Компания", "Номер УПД", "Сумма за VPBX, с НДС", "Сумма от БЛ в БИ, руб с НДС", "Кол-во исходящих вызовов",
	}

	return &excelImpl{
		log:         logger,
		filename:    filename,
		file:        file,
		baseStyle:   baseStyle,
		filledStyle: filledStyle,
		sheetName:   sheetName,
		columnNames: columnNames,
	}
}

// SetHeader - Создает шапку в excel файле с названиями столбцов
func (ex *excelImpl) SetHeader(month time.Month) error {
	ex.log.Info("Set header")
	topLeftCell := "B1"
	bottomRightCell := "N1"
	err := ex.SetCellValue("A1", "Месяц", false)
	if err != nil {
		return err
	}
	err = ex.file.MergeCell(ex.sheetName, topLeftCell, bottomRightCell)
	if err != nil {
		return err
	}

	err = ex.SetCellValue("B1", month, false)
	if err != nil {
		return err
	}

	err = ex.file.SetCellStyle(ex.sheetName, topLeftCell, bottomRightCell, ex.baseStyle)
	if err != nil {
		return err
	}

	// Установка высоты строки
	err = ex.file.SetRowHeight(ex.sheetName, 1, 30)
	if err != nil {
		return err
	}

	for colNum, v := range ex.columnNames {
		err := ex.SetCellValue(cell(numberToExcelCol(colNum+1), 2), v, false)
		if err != nil {
			return err
		}
	}
	ex.log.Info("Set header successfully")
	return nil
}

// AddData - Добавляет данные в Excel файл
func (ex *excelImpl) AddData(data *Rows) error {
	ex.log.Info("Adding data")
	err := ex.SetHeader(data.Month)
	if err != nil {
		ex.log.Error("Failed to set header", zap.Error(err))
		return err
	}

	ex.log.Info("Adding BL data ")
	start := 3
	blData := data.BL
	for rowNum := start; rowNum < len(blData.Data); rowNum++ {
		err = ex.addRow(rowNum, blData.Data[rowNum-start])
		if err != nil {
			ex.log.Error("Failed to add row", zap.Error(err))
		}
	}

	ex.log.Info("Adding total value BL data ")
	err = ex.addTotalValue(start+len(blData.Data), blData)
	if err != nil {
		ex.log.Error("Failed to add row", zap.Error(err))
	}

	ex.log.Info("Adding BI data ")
	start = start + len(blData.Data) + 1
	biData := data.BI
	for rowNum := start; rowNum < len(biData.Data); rowNum++ {
		err = ex.addRow(rowNum, biData.Data[rowNum-start])
		if err != nil {
			ex.log.Error("Failed to add row", zap.Error(err))
		}
	}

	ex.log.Info("Adding total value BI data ")
	err = ex.addTotalValue(start+len(biData.Data), biData)
	if err != nil {
		ex.log.Error("Failed to add row", zap.Error(err))
	}

	ex.log.Info("Adding data successfully")
	return nil
}

// SetCellValue - Вносит данные (data) в указанную ячейку (cell).
// Если необходима заливка и жирный шрифт, используется флаг fillingRequired
func (ex *excelImpl) SetCellValue(cell string, data any, fillingRequired bool) error {
	err := ex.file.SetCellValue(ex.sheetName, cell, data)
	if err != nil {
		return err
	}

	if fillingRequired {
		err = ex.file.SetCellStyle(ex.sheetName, cell, cell, ex.filledStyle)
		if err != nil {
			return err
		}
	} else {
		err = ex.file.SetCellStyle(ex.sheetName, cell, cell, ex.baseStyle)
		if err != nil {
			return err
		}
	}

	return nil
}

// Save - сохранение файла
func (ex *excelImpl) Save() error {
	// TODO: Возможно стоит рассмотреть вариант, когда после добавление
	// 	вызывается данный метод и сюда передается месяц и мы формируем
	// 	имя файла тут, а не при инициализации объекта Excel
	return ex.file.SaveAs(ex.filename)
}

// AddRow - Добавляет строку в Excel файл.
//
// Параметры:
//
//	rowNum - номер строки в которую кладутся данные
//	offset - число, необходимое, чтобы разница rowNum - offset = 0. Необходимо для правильного определения столбца
//	data - данные о клиенте
//	fillReq - флаг определяющий надо ли выделять ячейки или нет
func (ex *excelImpl) addRow(rowNum int, data Row) error {
	for colNum, v := range data.Flatten() {
		err := ex.SetCellValue(cell(numberToExcelCol(colNum+1), rowNum), v, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// addTotalValue - добавляет итоговые значения после основных данных
func (ex *excelImpl) addTotalValue(rowNum int, data CompanyData) error {
	err := ex.SetCellValue(cell("A", rowNum), "Итого", true)
	if err != nil {
		return err
	}

	err = ex.SetCellValue(cell("B", rowNum), data.SumMinutesCount, true)
	if err != nil {
		return err
	}

	err = ex.SetCellValue(cell("D", rowNum), data.SumNumbersCount, true)
	if err != nil {
		return err
	}

	err = ex.SetCellValue(cell("G", rowNum), data.SumAdditionalServices, true)
	if err != nil {
		return err
	}

	err = ex.SetCellValue(cell("I", rowNum), data.SumTotalAmountWithTax, true)
	if err != nil {
		return err
	}

	err = ex.SetCellValue(cell("N", rowNum), data.SumCallsCount, true)
	if err != nil {
		return err
	}
	return nil
}

// cell - возвращает название ячейки в нужном формате
// Пример:
//
// cell("A", 2) -> "A2"
func cell(cell string, num int) string {
	return strings.ToUpper(fmt.Sprintf("%s%d", cell, num))
}
func calculatingColumnWidth(columnName string) int {
	panic("implement me")
	return 0
}

// numberToExcelCol - Возвращает букву(-ы) для указания ячейки в Excel.
//
//	n - параметр смещения по алфавиту от буквы А
//
// Пример:
//
//		n = 0 -> "A"
//		n = 1 -> "B"
//		n = 26 -> "Z"
//	 n = 27 -> "AA"
//		n = 703 -> "AAA"
func numberToExcelCol(n int) string {
	var res string
	letterA := rune(65)
	if n == 0 {
		return string(letterA)
	}
	for n > 0 {
		n--
		res = string(letterA+rune(n%26)) + res
		n /= 26
	}
	return res
}
