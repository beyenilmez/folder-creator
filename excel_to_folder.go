package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenthenguyen/docx"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (a *App) GetExcelFileDialog() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Excel dosyası seçin",
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel",
				Pattern:     "*.xlsx;*.xls",
			},
		},
	})

	if err != nil {
		runtime.LogWarning(a.ctx, err.Error())
		return ""
	}

	return path
}

func (a *App) GetWordFileDialog() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Word dosyası seçin",
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Word",
				Pattern:     "*.docx",
			},
		},
	})

	if err != nil {
		runtime.LogWarning(a.ctx, err.Error())
		return ""
	}

	return path
}

func (a *App) GetCopyFolderDialog() string {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Kopyalanacak klasörü seçin",
		CanCreateDirectories: true,
	})

	if err != nil {
		runtime.LogWarning(a.ctx, err.Error())
		return ""
	}

	return path
}

func (a *App) GetTargetFolderDialog() string {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Hedef klasörü seçin",
		CanCreateDirectories: true,
	})

	if err != nil {
		runtime.LogWarning(a.ctx, err.Error())
		return ""
	}

	return path
}

func (a *App) CreateFolders(excelPath string, wordPath string, copyFolderPath string, targetPath string, folderNamePattern string, wordFileNamePattern string) string {
	excelFile, err := excelize.OpenFile(excelPath)

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err.Error()
	}

	sheetName := excelFile.GetSheetList()[0]

	// Get all the rows in the Sheet1.
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err.Error()
	}

	var folderNames []string

	var headers []string

	for i, row := range rows {
		if i == 0 {
			headers = append(headers, row...)

			runtime.LogInfo(a.ctx, "Headers: "+strings.Join(headers, ", "))
			continue
		}

		folderName := folderNamePattern
		for k, header := range headers {
			colCell := row[k]

			// remove whitespace
			colCell = strings.TrimSpace(colCell)

			// Convert new lines to _
			colCell = strings.ReplaceAll(colCell, "\r\n", "_")
			colCell = strings.ReplaceAll(colCell, "\n", "_")

			// Convert tabs to _
			colCell = strings.ReplaceAll(colCell, "\t", "_")

			// replace slashes
			colCell = strings.ReplaceAll(colCell, "/", "_")

			// Replace placeholders in the pattern
			placeholder := "{" + header + "}"
			titlePlaceholder := "{{" + header + "}}"

			// Title case conversion for specific placeholders
			if strings.Contains(folderName, titlePlaceholder) {
				words := strings.Split(colCell, " ")
				for i, word := range words {
					words[i] = cases.Title(language.Turkish).String(word)
				}
				colCell = strings.Join(words, "_")
				folderName = strings.ReplaceAll(folderName, titlePlaceholder, colCell)
			}

			// Simple replacement
			folderName = strings.ReplaceAll(folderName, placeholder, colCell)
		}

		folderName = strings.TrimSpace(folderName)

		folderNames = append(folderNames, folderName)
	}

	//remove first row
	rows = rows[1:]

	for i := 0; i < len(folderNames); i++ {
		folderName := folderNames[i]
		targetFolderPath := filepath.Join(targetPath, folderName)

		if _, err := os.Stat(targetFolderPath); os.IsNotExist(err) {
			err = os.MkdirAll(targetFolderPath, 0o755)
			if err != nil {
				runtime.LogError(a.ctx, "Failed to create folder: "+err.Error())
				continue
			}
		}

		if copyFolderPath != "" {
			// Recursively copy contents of copyFolderPath to targetFolderPath
			err = filepath.Walk(copyFolderPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				relativePath, err := filepath.Rel(copyFolderPath, path)
				if err != nil {
					return err
				}
				targetPath := filepath.Join(targetFolderPath, relativePath)
				if info.IsDir() {
					if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
						return err
					}
				} else {
					if err := copyFile(path, targetPath); err != nil {
						return err
					}
				}
				return nil
			})

			if err != nil {
				runtime.LogError(a.ctx, err.Error())
				continue
			}
		}

		if wordPath != "" {
			// Edit word document
			err = EditWordDocument(wordPath, wordFileNamePattern, headers, rows[i], targetFolderPath)

			if err != nil {
				runtime.LogError(a.ctx, err.Error())
				continue
			}
		}
	}

	a.SendNotification("Klasör oluşturma başarılı", "", strings.ReplaceAll(targetPath, "\\", "\\\\"), "success")

	return ""
}

func EditWordDocument(filePath string, wordFileNamePattern string, headers []string, row []string, targetPath string) error {
	r, err := docx.ReadDocxFile(filePath)

	if err != nil {
		runtime.LogError(appContext, err.Error())
		return err
	}

	docx1 := r.Editable()
	docx1Content := docx1.GetContent()

	for k, header := range headers {
		colCell := row[k]

		// remove whitespace
		colCell = strings.TrimSpace(colCell)

		// Replace placeholders in the pattern
		placeholder := "{" + header + "}"
		titlePlaceholder := "{{" + header + "}}"

		// Title case conversion for specific placeholders
		if strings.Contains(docx1Content, titlePlaceholder) {
			words := strings.Split(colCell, " ")
			for i, word := range words {
				words[i] = cases.Title(language.Turkish).String(word)
			}
			colCell = strings.Join(words, " ")
			err = docx1.Replace(titlePlaceholder, colCell, -1)

			if err != nil {
				runtime.LogError(appContext, err.Error())
				return err
			}
		}

		// Simple replacement
		err = docx1.Replace(placeholder, colCell, -1)

		if err != nil {
			runtime.LogError(appContext, err.Error())
			return err
		}
	}

	fileName := filepath.Base(wordFileNamePattern)
	for k, header := range headers {
		colCell := row[k]

		// remove whitespace
		colCell = strings.TrimSpace(colCell)

		// replace slashes
		colCell = strings.ReplaceAll(colCell, "/", "_")

		// Replace placeholders in the pattern
		placeholder := "{" + header + "}"
		titlePlaceholder := "{{" + header + "}}"

		// Title case conversion for specific placeholders
		if strings.Contains(fileName, titlePlaceholder) {
			words := strings.Split(colCell, " ")
			for i, word := range words {
				words[i] = cases.Title(language.Turkish).String(word)
			}
			colCell = strings.Join(words, "_")
			fileName = strings.ReplaceAll(fileName, titlePlaceholder, colCell)
		}

		// Simple replacement
		fileName = strings.ReplaceAll(fileName, placeholder, colCell)
	}

	docx1.WriteToFile(filepath.Join(targetPath, fileName) + ".docx")

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
