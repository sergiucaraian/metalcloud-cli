package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

//ConsoleIOChannel represents an IO channel, typically stdin and stdout but could be anything
type ConsoleIOChannel struct {
	Stdin  io.Reader
	Stdout io.Writer
}

var consoleIOChannelInstance ConsoleIOChannel

var once sync.Once

//GetConsoleIOChannel returns the console channel singleton
func GetConsoleIOChannel() *ConsoleIOChannel {
	once.Do(func() {

		consoleIOChannelInstance = ConsoleIOChannel{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
		}
	})

	return &consoleIOChannelInstance
}

//GetStdout returns the configured output channel
func GetStdout() io.Writer {
	return GetConsoleIOChannel().Stdout
}

//GetStdin returns the configured input channel
func GetStdin() io.Reader {
	return GetConsoleIOChannel().Stdin
}

//SetConsoleIOChannel configures the stdin and stdout to be used by all io with
func SetConsoleIOChannel(in io.Reader, out io.Writer) {
	channel := GetConsoleIOChannel()
	channel.Stdin = in
	channel.Stdout = out
}

//GetTableHeader returns the row for header (all cells strings but of the length specified in the schema)
func GetTableHeader(schema []SchemaField) string {
	var alteredSchema []SchemaField
	var header []interface{}

	for _, field := range schema {
		alteredSchema = append(alteredSchema, SchemaField{
			FieldType: TypeString,
			FieldSize: field.FieldSize,
		})
		header = append(header, field.FieldName)
	}
	return GetTableRow(header, alteredSchema)
}

//GetTableRow returns the string for a row with the | delimiter
func GetTableRow(row []interface{}, schema []SchemaField) string {
	var rowStr []string
	for i, field := range schema {
		switch field.FieldType {
		case TypeInt:
			rowStr = append(rowStr, fmt.Sprintf(fmt.Sprintf(" %%-%dd", field.FieldSize), row[i].(int)))
		case TypeString:
			rowStr = append(rowStr, fmt.Sprintf(fmt.Sprintf(" %%-%ds", field.FieldSize), row[i].(string)))
		case TypeFloat:
			rowStr = append(rowStr, fmt.Sprintf(fmt.Sprintf(" %%-%d.%df", field.FieldSize, field.FieldPrecision), row[i].(float64)))
		default:
			rowStr = append(rowStr, fmt.Sprintf(fmt.Sprintf(" %%-+%dv", field.FieldSize), row[i]))
		}
	}
	return "|" + strings.Join(rowStr, "|") + "|"
}

// GetCellSize calculates how wide a cell is by converting it to string and measuring it's size
func GetCellSize(d interface{}, field *SchemaField) int {
	var s string
	switch field.FieldType {
	case TypeInt:
		s = fmt.Sprintf("%d", d.(int))
	case TypeString:
		s = d.(string)
	case TypeFloat:
		s = fmt.Sprintf(fmt.Sprintf("%%.%df", field.FieldPrecision), d.(float64))
	default:
		s = fmt.Sprintf("%v", d)

	}
	return len(s)
}

//AdjustFieldSizes expands field sizes to match the widest cell
func AdjustFieldSizes(data [][]interface{}, schema *[]SchemaField) {
	rowSize := len(*schema)
	for i := 0; i < rowSize; i++ {
		f := (*schema)[i]

		//iterate over the entire column
		rowCount := len(data)

		maxLen := f.FieldSize

		if len(f.FieldName) > maxLen {
			maxLen = len(f.FieldName)
		}

		for k := 0; k < rowCount; k++ {
			cellSize := GetCellSize(data[k][i], &f)
			if cellSize > maxLen {
				maxLen = cellSize
			}
		}
		if maxLen > f.FieldSize {
			(*schema)[i].FieldSize = maxLen + 1 //we leave a little room to the right
		}
	}
}

// GetTableDelimiter returns a delimiter row for the schema
func GetTableDelimiter(schema []SchemaField) string {
	row := "+"
	for _, field := range schema {
		for i := 0; i < field.FieldSize+1; i++ {
			row += "-"
		}
		row += "+"
	}
	return row
}

//GetTableAsString returns the string representation of a table.
func GetTableAsString(data [][]interface{}, schema []SchemaField) string {
	var rows []string

	rows = append(rows, GetTableDelimiter(schema))
	rows = append(rows, GetTableHeader(schema))
	rows = append(rows, GetTableDelimiter(schema))
	for _, row := range data {
		rows = append(rows, GetTableRow(row, schema))
	}
	rows = append(rows, GetTableDelimiter(schema))

	return strings.Join(rows, "\n") + "\n"
}

func printTableHeader(schema []SchemaField) {
	fmt.Println(GetTableHeader(schema))
}

func printTableRow(row []interface{}, schema []SchemaField) {
	fmt.Println(GetTableRow(row, schema))
}

func printTableDelimiter(schema []SchemaField) {
	fmt.Println(GetTableDelimiter(schema))
}

func printTable(data [][]interface{}, schema []SchemaField) {
	fmt.Print(GetTableAsString(data, schema))
}

//GetTableAsJSONString returns a MarshalIndent string for the given data
func GetTableAsJSONString(data [][]interface{}, schema []SchemaField) (string, error) {
	dataAsMap := make([]interface{}, len(data))

	for k, row := range data {
		rowAsMap := make(map[string]interface{}, len(schema))
		for i, field := range schema {
			rowAsMap[field.FieldName] = row[i]
		}
		dataAsMap[k] = rowAsMap
	}

	ret, err := json.MarshalIndent(dataAsMap, "", "\t")
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

//GetTableAsCSVString returns a csv
func GetTableAsCSVString(data [][]interface{}, schema []SchemaField) (string, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	csvWriter := csv.NewWriter(writer)

	rowStr := make([]string, len(schema))
	for i, field := range schema {
		rowStr[i] = field.FieldName
	}

	csvWriter.Write(rowStr)

	for _, row := range data {
		for i, field := range schema {
			switch field.FieldType {
			case TypeInt:
				rowStr[i] = fmt.Sprintf("%d", row[i].(int))
			case TypeString:
				rowStr[i] = row[i].(string)
			case TypeFloat:
				rowStr[i] = fmt.Sprintf("%f", row[i].(float64))
			case TypeInterface:
				rowStr[i] = fmt.Sprintf("%v", row[i])
			case TypeDateTime:
				rowStr[i] = row[i].(string)
			default:
				rowStr[i] = fmt.Sprintf("%v", row[i])
			}
		}
		csvWriter.Write(rowStr)
	}

	writer.Flush()
	csvWriter.Flush()

	return buf.String(), nil
}

func truncateString(s string, length int) string {
	str := s
	if len(str) > 0 {
		return str[:length] + "..."
	}
	return ""
}

func renderTable(tableName string, topLine string, format string, data [][]interface{}, schema []SchemaField) (string, error) {
	var sb strings.Builder

	switch format {
	case "json", "JSON":
		ret, err := GetTableAsJSONString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)
	case "csv", "CSV":
		ret, err := GetTableAsCSVString(data, schema)
		if err != nil {
			return "", err
		}
		sb.WriteString(ret)

	default:
		if topLine != "" {
			sb.WriteString(fmt.Sprintf("%s\n", topLine))
		} else {
			user := GetUserEmail()
			sb.WriteString(fmt.Sprintf("%s I have access to as user %s:\n", tableName, user))
		}

		AdjustFieldSizes(data, &schema)

		sb.WriteString(GetTableAsString(data, schema))

		sb.WriteString(fmt.Sprintf("Total: %d %s\n\n", len(data), tableName))
	}

	return sb.String(), nil
}

//transposeTable turns columns into rows. It assumes an uniform length table
func transposeTable(data [][]interface{}) [][]interface{} {

	dataT := [][]interface{}{}

	if len(data) == 0 {
		return dataT
	}

	tableLength := len(data)
	rowLength := len(data[0])

	for j := 0; j < rowLength; j++ {

		newRow := []interface{}{}

		for i := 0; i < tableLength; i++ {

			newRow = append(newRow, data[i][j])
		}

		dataT = append(dataT, newRow)
	}

	return dataT
}

func convertToStringTable(data [][]interface{}) [][]interface{} {
	dataS := [][]interface{}{}

	for _, row := range data {
		newRow := []interface{}{}
		for _, v := range row {
			if v == nil {
				v = " "
			}
			newRow = append(newRow, fmt.Sprintf("%v", v))
		}
		dataS = append(dataS, newRow)
	}
	return dataS
}

//renderTransposedTable renders the text format as a key-value table. json and csv formats remain the same as render table
func renderTransposedTable(tableName string, topLine string, format string, data [][]interface{}, schema []SchemaField) (string, error) {

	if format != "" {
		return renderTable(tableName, topLine, format, data, schema)
	}

	headerRow := []interface{}{}
	for _, s := range schema {
		headerRow = append(headerRow, s.FieldName)
	}

	dataAsStrings := convertToStringTable(data)
	newData := [][]interface{}{}
	newData = append(newData, headerRow)
	for _, row := range dataAsStrings {
		newData = append(newData, row)
	}

	dataTransposed := transposeTable(newData)

	newSchema := []SchemaField{
		{
			FieldName: "KEY",
			FieldType: TypeString,
			FieldSize: 5,
		},
		{
			FieldName: "VALUE",
			FieldType: TypeString,
			FieldSize: 5,
		},
	}

	return renderTable(tableName, topLine, format, dataTransposed, newSchema)

}
