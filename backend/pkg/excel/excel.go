package excel

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelWriter Excel 写入器
type ExcelWriter struct {
	file *excelize.File
}

// NewExcelWriter 创建 Excel 写入器
func NewExcelWriter() *ExcelWriter {
	return &ExcelWriter{
		file: excelize.NewFile(),
	}
}

// WriteHeaders 写入表头
func (w *ExcelWriter) WriteHeaders(sheetName string, headers []string) error {
	// 创建工作表
	w.file.NewSheet(sheetName)

	// 写入表头
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), 1)
		w.file.SetCellValue(sheetName, cell, header)
	}

	return nil
}

// WriteRow 写入一行数据
func (w *ExcelWriter) WriteRow(sheetName string, rowNum int, data []interface{}) error {
	for i, value := range data {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), rowNum)
		w.file.SetCellValue(sheetName, cell, value)
	}
	return nil
}

// WriteRows 写入多行数据
func (w *ExcelWriter) WriteRows(sheetName string, startRow int, data [][]interface{}) error {
	for i, row := range data {
		if err := w.WriteRow(sheetName, startRow+i, row); err != nil {
			return err
		}
	}
	return nil
}

// WriteTo 写入到 io.Writer
func (w *ExcelWriter) WriteTo(writer io.Writer) error {
	return w.file.Write(writer)
}

// Close 关闭 Excel 文件
func (w *ExcelWriter) Close() error {
	return w.file.Close()
}

// ExcelReader Excel 读取器
type ExcelReader struct {
	file *excelize.File
}

// NewExcelReader 从 io.Reader 创建 Excel 读取器
func NewExcelReader(reader io.Reader) (*ExcelReader, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	return &ExcelReader{file: file}, nil
}

// GetSheetNames 获取所有工作表名称
func (r *ExcelReader) GetSheetNames() []string {
	return r.file.GetSheetList()
}

// ReadRows 读取工作表的所有行
func (r *ExcelReader) ReadRows(sheetName string) ([][]interface{}, error) {
	rows, err := r.file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	// 转换为 [][]interface{}
	result := make([][]interface{}, len(rows))
	for i, row := range rows {
		result[i] = make([]interface{}, len(row))
		for j, cell := range row {
			result[i][j] = cell
		}
	}

	return result, nil
}

// Close 关闭 Excel 文件
func (r *ExcelReader) Close() error {
	return r.file.Close()
}

// StringToInt 字符串转整数
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// StringToInt64 字符串转 int64
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// StringToFloat64 字符串转 float64
func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// StringToTime 字符串转时间
func StringToTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// InterfaceToString 接口转字符串
func InterfaceToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
