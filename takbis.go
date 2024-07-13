package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func parseHeaderChangePattern(pattern string) map[string]string {
	mapping := make(map[string]string)
	pairs := strings.Split(pattern, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, ":")
		if len(kv) == 2 {
			mapping[kv[0]] = kv[1]
		}
	}
	return mapping
}

func parceCellChangePattern(pattern string) map[string]string {
	mapping := make(map[string]string)
	pairs := strings.Split(pattern, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "->")
		if len(kv) == 2 {
			mapping[kv[0]] = kv[1]
		}
	}
	return mapping
}

func (app *App) ModifyExcelWithTakbis(excelPath string, takbisPaths []string, headerMatchPattern string, cellChangeRule string) string {
	excelHeaders, excelRows, excel, err := ReadExcel(excelPath)
	if err != nil {
		return err.Error()
	}
	sheetName := excel.GetSheetList()[0]

	headerMatchMap := parseHeaderChangePattern(headerMatchPattern)
	cellChangeMap := parceCellChangePattern(cellChangeRule)

	runtime.LogInfo(app.ctx, "Modifying Excel with Takbis")

	runtime.LogInfo(app.ctx, "Header match map: "+fmt.Sprint(headerMatchMap))
	runtime.LogInfo(app.ctx, "Cell change map: "+fmt.Sprint(cellChangeMap))

	// Create reverse maps for quick lookup
	excelHeaderIdx := make(map[string]int)

	for i, header := range excelHeaders {
		excelHeaderIdx[header] = i
	}

	runtime.LogInfo(app.ctx, "Excel header index: "+fmt.Sprint(excelHeaderIdx))

	// Put takbis rows and headers into a map
	takbisHeadersList := make([][]string, 0)
	takbisRowsList := make([][][]string, 0)

	runtime.LogInfo(app.ctx, "Reading takbis rows and headers")

	for i, takbisPath := range takbisPaths {
		runtime.WindowExecJS(appContext, `window.setTakbisMessage("Takbis dosyaları okunuyor: `+fmt.Sprintf("%d/%d", i+1, len(takbisPaths))+`");`)

		currentTakbisHeaders, currentTakbisRows, _, err := ReadExcel(takbisPath)

		if err != nil {
			return err.Error()
		}

		takbisHeadersList = append(takbisHeadersList, currentTakbisHeaders)
		takbisRowsList = append(takbisRowsList, currentTakbisRows)
	}

	runtime.LogDebug(app.ctx, "Takbis headers and rows are read")

	for excelRowNumber := 0; excelRowNumber < len(excelRows); excelRowNumber++ {
		runtime.WindowExecJS(appContext, `window.setTakbisMessage("`+fmt.Sprintf("%d/%d", excelRowNumber, len(excelRows))+`");`)
		excelRow := excelRows[excelRowNumber]

		for i := 0; i < 2; i++ {

			for takbisNumber := 0; takbisNumber < len(takbisHeadersList); takbisNumber++ {
				breakOut := false
				takbisHeaders := takbisHeadersList[takbisNumber]
				takbisRows := takbisRowsList[takbisNumber]

				takbisHeaderIdx := make(map[string]int)
				for i, header := range takbisHeaders {
					takbisHeaderIdx[header] = i
				}

				for takbisRowNumber := 0; takbisRowNumber < len(takbisRows); takbisRowNumber++ {
					takbisRow := takbisRows[takbisRowNumber]
					rowMatch := true

					for targetHeader, takbisHeader := range headerMatchMap {
						targetIdx, ok1 := excelHeaderIdx[targetHeader]
						takbisIdx, ok2 := takbisHeaderIdx[takbisHeader]

						if strings.TrimSpace(excelRow[targetIdx]) == "" {
							rowMatch = false
							break
						}

						if ok1 && ok2 && targetIdx < len(excelRow) && takbisIdx < len(takbisRow) {
							if i == 1 {
								if !LooseEqualWithoutLastWord(takbisRow[takbisIdx], excelRow[targetIdx]) {
									rowMatch = false
									break
								}
							} else {
								if !LooseEqual(excelRow[targetIdx], takbisRow[takbisIdx]) {
									rowMatch = false
									break
								}
							}
						} else {
							rowMatch = false
							break
						}
					}

					if rowMatch {
						i++
						runtime.LogDebugf(app.ctx, "Matched row in target Excel: %s\nwith row in Takbis Excel: %s", excelRow, takbisRow)

						// Print comparison values
						for targetHeader, takbisHeader := range headerMatchMap {
							targetIdx, ok1 := excelHeaderIdx[targetHeader]
							takbisIdx, ok2 := takbisHeaderIdx[takbisHeader]

							if ok1 && ok2 && targetIdx < len(excelRow) && takbisIdx < len(takbisRow) {
								runtime.LogDebugf(app.ctx, "%s:%s", excelRow[targetIdx], takbisRow[takbisIdx])
							}
						}

						// Update the cells in the target Excel row based on cellChangeRule
						for takbisHeader, targetHeader := range cellChangeMap {
							takbisIdx, ok1 := takbisHeaderIdx[takbisHeader]
							targetIdx, ok2 := excelHeaderIdx[targetHeader]

							if ok1 && ok2 && takbisIdx < len(takbisRow) && targetIdx < len(excelRow) {
								isFloat := false

								temp := takbisRow[takbisIdx]
								temp = strings.ReplaceAll(temp, ",", ".")

								floatValue, err := strconv.ParseFloat(temp, 64)
								if err == nil {
									isFloat = true
								}

								if isFloat {
									excel.SetCellFloat(sheetName, alphabet[targetIdx]+fmt.Sprintf("%v", excelRowNumber+2), floatValue, 2, 64)
								} else {
									excel.SetCellDefault(sheetName, alphabet[targetIdx]+fmt.Sprintf("%v", excelRowNumber+2), takbisRow[takbisIdx])
								}
							}
						}

						breakOut = true
						break
					}
				}

				if breakOut {
					break
				}
			}

		}

	}

	runtime.LogInfo(app.ctx, "Attempting to save Excel file")

	err = excel.SaveAs(excelPath)

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return err.Error()
	}

	runtime.WindowExecJS(appContext, `window.setTakbisMessage("Excel dosyası güncellendi");`)
	return ""
}

func LooseEqual(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	aTitleCase := toTitleCaseWord(a)
	bTitleCase := toTitleCaseWord(b)

	return aTitleCase == bTitleCase
}

func LooseEqualWithoutLastWord(a, b string) bool {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)

	aSplit := strings.Split(a, " ")

	if len(aSplit) > 1 {
		// remove last word
		a = strings.Join(aSplit[:len(aSplit)-1], " ")
	}

	return LooseEqual(a, b)
}
