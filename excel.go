package toold

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
)

//Excel 结构
type Excel struct {
	File  *xlsx.File
	Sheet *xlsx.Sheet
	Row   *xlsx.Row
	Cell  *xlsx.Cell
	Err   error
}

var file *xlsx.File
var sheet *xlsx.Sheet
var row *xlsx.Row
var cell *xlsx.Cell
var err error

//GetExcelExt GetImageExt
func GetExcelExt(ext string) string {
	extL := strings.ToLower(ext)
	switch extL {
	case ".xlsx", ".xls":
		break
	case "xlsx", "xls":
		ext = fmt.Sprintf(".%v", ext)
		break
	default:
		return ""
	}
	return ext
}

//NewExcels 初始化
func NewExcels() *Excel {
	return &Excel{}
}

//CreateXLSX 创建
func (e *Excel) CreateXLSX() {
	e.File, e.Sheet, e.Err = createXLSX()
}

/*
SetTitle 设置标题
*/
func (e *Excel) SetTitle(arrayTitle []string) (err error) {
	style := xlsx.NewStyle()
	fill := *xlsx.NewFill("solid", "FFFFFFFF", "00000000")
	font := *xlsx.NewFont(16, "Verdana")
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")

	style.Fill = fill
	style.Font = font
	style.Border = border
	style.Alignment.Horizontal = "center"
	style.Alignment.Vertical = "center"
	style.Alignment.WrapText = false
	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = false
	row = e.Sheet.AddRow()
	for i := 0; i < len(arrayTitle); i++ {
		cell = row.AddCell()
		row.SetHeight(18)
		cell.SetStyle(style)
		cell.Value = arrayTitle[i]
	}
	if err != nil {
		return err
	}

	e.Sheet.SetColWidth(0, len(arrayTitle), 25.0)
	return nil
}

/*
AddContent 添加内容
*/
func (e *Excel) AddContent(arrayContent []string) (err error) {

	styleC := xlsx.NewStyle()
	styleC.Alignment.Horizontal = "center"
	styleC.Alignment.Vertical = "center"
	styleC.Alignment.WrapText = true

	row = e.Sheet.AddRow()
	for i := 0; i < len(arrayContent); i++ {
		cell = row.AddCell()
		//row.SetHeight(15)
		cell.SetStyle(styleC)

		val := arrayContent[i]
		if len(val) != 0 && (IsNumberFromString(val) || IsMinusNumberFromString(val)) {
			cell.SetInt64(ConversionToInt64(val))
		} else {
			cell.Value = arrayContent[i]
		}
	}
	return nil
}

type ExcelCellType int

const (
	ExcelCellTypeText ExcelCellType = iota
	ExcelCellTypeSelect
	ExcelCellTypeURI
)

type ExcelCell struct {
	Select []string
	Val    string
	Type   ExcelCellType
}

func (e *Excel) AddCell(cls []*ExcelCell) {
	styleC := xlsx.NewStyle()
	styleC.Alignment.Horizontal = "center"
	styleC.Alignment.Vertical = "center"
	styleC.Alignment.WrapText = true
	row = e.Sheet.AddRow()
	for _, info := range cls {
		cell := row.AddCell()
		cell.SetStyle(styleC)
		if info.Type == ExcelCellTypeSelect {
			val := xlsx.NewXlsxCellDataValidation(true)
			val.SetDropList(info.Select)
			cell.SetDataValidation(val)
			cell.SetString(info.Val)
		} else if info.Type == ExcelCellTypeURI {
			styleC := xlsx.NewStyle()
			styleC.Alignment.Horizontal = "center"
			styleC.Alignment.Vertical = "center"
			styleC.Alignment.WrapText = true
			styleC.ApplyAlignment = true
			styleC.Font.Underline = true
			styleC.Font.Color = "FF0000FF"
			cell.SetStyle(styleC)
			cell.SetFormula(fmt.Sprintf("=HYPERLINK(\"%v\")", info.Val))
		} else {
			cell.SetString(info.Val)
		}
	}
}

/*
AddContent 添加内容
*/
func (e *Excel) AddContentLink(arrayContent []string, link map[string]string) (err error) {

	styleC := xlsx.NewStyle()
	styleC.Alignment.Horizontal = "center"
	styleC.Alignment.Vertical = "center"
	styleC.Alignment.WrapText = true
	styleC.Font.Underline = false

	row = e.Sheet.AddRow()
	for i := 0; i < len(arrayContent); i++ {
		cell = row.AddCell()
		//row.SetHeight(15)
		cell.SetStyle(styleC)

		val := arrayContent[i]
		if link != nil {
			con, is := link[val]
			if is {
				styleC := xlsx.NewStyle()
				styleC.Alignment.Horizontal = "center"
				styleC.Alignment.Vertical = "center"
				styleC.Alignment.WrapText = true
				styleC.ApplyAlignment = true
				styleC.Font.Underline = true
				styleC.Font.Color = "FF0000FF"
				cell.SetStyle(styleC)
				cell.SetFormula(fmt.Sprintf("=HYPERLINK(\"%v\")", con))
				continue
			}
		}
		if len(val) != 0 && (IsNumberFromString(val) || IsMinusNumberFromString(val)) {
			cell.SetInt64(ConversionToInt64(val))
		} else {
			cell.Value = arrayContent[i]
		}
	}
	return nil
}

/*
Save 保存
*/
func (e *Excel) Save(title string) (err error) {
	var name = title + ".xlsx"
	err = e.File.Save(name)
	return err
}

func createXLSX() (*xlsx.File, *xlsx.Sheet, error) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	return file, sheet, err
}

/*
ExcelX 多级参加excel sheet
*/
type ExcelX struct {
	File  *xlsx.File
	Sheet map[string]*xlsx.Sheet
}

//ExcelCreatelFile ExcelCreatelFile
func ExcelCreatelFile() ExcelX {
	return ExcelX{
		File:  xlsx.NewFile(),
		Sheet: make(map[string]*xlsx.Sheet),
	}
}

//AddSheet AddSheet
func (e ExcelX) AddSheet(name string) (*xlsx.Sheet, error) {
	sheet, err := e.File.AddSheet(name)
	if err != nil {
		return nil, err
	}
	e.Sheet[name] = sheet
	return sheet, nil
}

//SetTitle SetTitle
func (e ExcelX) SetTitle(sheetName string, arrayTitle []string) error {
	style := xlsx.NewStyle()
	fill := *xlsx.NewFill("solid", "FFFFFFFF", "00000000")
	font := *xlsx.NewFont(14, "Verdana")
	border := *xlsx.NewBorder("thin", "thin", "thin", "thin")

	style.Fill = fill
	style.Font = font
	style.Border = border
	style.Alignment.Horizontal = "center"
	style.Alignment.Vertical = "center"
	style.Alignment.WrapText = false
	style.ApplyFill = true
	style.ApplyFont = true
	style.ApplyBorder = false

	row = e.Sheet[sheetName].AddRow()
	for i := 0; i < len(arrayTitle); i++ {
		cell = row.AddCell()
		row.SetHeight(18)
		cell.SetStyle(style)
		cell.Value = arrayTitle[i]
	}
	if err != nil {
		return err
	}

	e.Sheet[sheetName].SetColWidth(0, len(arrayTitle), 25.0)
	return nil
}

//AddCellContent AddCellContent
func (e ExcelX) AddCellContent(sheetName string, arrayContent []string) error {
	styleC := xlsx.NewStyle()
	styleC.Alignment.Horizontal = "center"
	styleC.Alignment.Vertical = "center"
	styleC.Alignment.WrapText = true
	row = e.Sheet[sheetName].AddRow()
	for i := 0; i < len(arrayContent); i++ {
		cell = row.AddCell()
		//row.SetHeight(15)
		cell.SetStyle(styleC)
		val := arrayContent[i]
		if len(val) != 0 && (IsNumberFromString(val) || IsMinusNumberFromString(val)) {
			cell.SetInt64(ConversionToInt64(val))
		} else {
			cell.Value = arrayContent[i]
		}
	}
	return nil
}

//AddCellLinkContent AddCellLinkContent
func (e ExcelX) AddCellLinkContent(sheetName string, arrayContent []string, index map[int]bool) error {

	row = e.Sheet[sheetName].AddRow()
	for i := 0; i < len(arrayContent); i++ {
		cell = row.AddCell()
		styleC := xlsx.NewStyle()
		styleC.Alignment.Horizontal = "center"
		styleC.Alignment.Vertical = "center"
		styleC.Alignment.WrapText = true
		//row.SetHeight(15)
		if index[i] {
			styleC.Font.Underline = true
			styleC.Font.Color = "FF0000FF"
			cell.SetStyle(styleC)
			cell.SetFormula(fmt.Sprintf("=HYPERLINK(\"%v\")", arrayContent[i]))
		} else {
			cell.SetStyle(styleC)
			val := arrayContent[i]
			if len(val) != 0 && (IsNumberFromString(val) || IsMinusNumberFromString(val)) {
				cell.SetInt64(ConversionToInt64(val))
			} else {
				cell.Value = arrayContent[i]
			}
		}
	}
	return nil
}

//AddCellContentHeight AddCellContentHeight
func (e ExcelX) AddCellContentHeight(sheetName string, height float64, arrayContent []string) error {
	styleC := xlsx.NewStyle()
	styleC.Alignment.Horizontal = "center"
	styleC.Alignment.Vertical = "center"
	styleC.Alignment.WrapText = true
	// styleC. = width
	row = e.Sheet[sheetName].AddRow()

	for i := 0; i < len(arrayContent); i++ {
		cell = row.AddCell()
		row.SetHeight(height)
		cell.SetStyle(styleC)
		val := arrayContent[i]
		if len(val) != 0 && (IsNumberFromString(val) || IsMinusNumberFromString(val)) {
			cell.SetInt64(ConversionToInt64(val))
		} else {
			cell.Value = arrayContent[i]
		}
	}
	return nil
}

//SetWidth SetWidth
func (e ExcelX) SetWidth(sheetName string, startcol, endcol int, width float64) {
	e.Sheet[sheetName].SetColWidth(startcol, endcol, width)
}

/*
Save 保存
*/
func (e ExcelX) Save(title string) (err error) {
	var name = title + ".xlsx"
	err = e.File.Save(name)
	return err
}
