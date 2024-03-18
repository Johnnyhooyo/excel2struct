// Package main 读取本地excel 根据不同的字段生成struct
package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"unicode"
)

func main() {
	file, _ := xlsx.OpenFile("./demo.xlsx")

	var structCode strings.Builder

	for _, sheet := range file.Sheets {
		structCode.WriteString(fmt.Sprintf("type %s struct {\n", sheet.Name))

		columnInfo := make(map[string]int)

		title := sheet.Rows[0]
		for i, cell := range title.Cells {
			columnName := cell.String()
			columnInfo[columnName] = i
		}

		for _, row := range sheet.Rows[1:] {
			var fieldName string
			if index, exists := columnInfo["name"]; exists {
				fieldName = row.Cells[index].String()
			}
			var fieldJson string
			if index, exists := columnInfo["json"]; exists {
				fieldJson = row.Cells[index].String()
			}
			var fieldType string
			if index, exists := columnInfo["type"]; exists {
				fieldType = row.Cells[index].String()
				switch fieldType {
				case "bigint":
					fieldType = "int64"
				case "int":
					fieldType = "int"
				case "time":
					fieldType = "time.Time"
				default:
					fieldType = "string"
				}
			}
			var fieldDesc string
			if index, exists := columnInfo["desc"]; exists {
				fieldDesc = row.Cells[index].String()
				fieldDesc = strings.ReplaceAll(fieldDesc, "\n", "\t")
			}
			var fieldMark string
			if index, exists := columnInfo["mark"]; exists {
				fieldMark = row.Cells[index].String()
				fieldMark = strings.ReplaceAll(fieldMark, "\n", "\t")
			}

			if fieldName == "" && fieldJson == "" {
				fmt.Printf("sheet %s not contain fieldName or FieldJson\n", sheet.Name)
				continue
			}

			if fieldName == "" && fieldJson != "" {
				fieldName = convertCamelCase(fieldJson)
			}

			if fieldJson == "" && fieldName != "" {
				fieldJson = convertUnderline(fieldName)
			}

			structCode.WriteString(fmt.Sprintf("    %s %s", fieldName, fieldType))
			structCode.WriteString(fmt.Sprintf(" `json:\"%s\"`", fieldJson))

			if len(fieldDesc) > 0 {
				structCode.WriteString(fmt.Sprintf(" // %s %s", fieldDesc, fieldMark))
			}

			structCode.WriteString("\n")
		}

		structCode.WriteString("}\n\n")
	}

	fmt.Println(structCode.String())
}

func convertCamelCase(input string) string {
	// 将输入字符串按下划线分割成切片
	parts := strings.Split(input, "_")

	// 遍历切片并将每个单词的首字母大写
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	// 将切片中的单词连接起来
	output := strings.Join(parts, "")

	return output
}

func convertUnderline(input string) string {
	var output []rune

	for i, char := range input {
		if unicode.IsUpper(char) {
			if i > 0 && !unicode.IsUpper(rune(input[i-1])) {
				output = append(output, '_')
			}
			output = append(output, unicode.ToLower(char))
		} else {
			output = append(output, char)
		}
	}

	return string(output)
}
