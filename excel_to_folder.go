package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

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

func (a *App) CreateFolders(excelPath string, copyFolderPath string, targetPath string, folderNamePattern string) {
	excelFile, err := excelize.OpenFile(excelPath)

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}

	sheetName := excelFile.GetSheetList()[0]

	// Get all the rows in the Sheet1.
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}

	var folderNames []string

	var headers []string

	for i, row := range rows {
		if i == 0 {
			for _, colCell := range row {
				headers = append(headers, colCell)
				continue
			}

			runtime.LogInfo(a.ctx, "Headers: "+strings.Join(headers, ", "))
		}

		// folderNamePattern: {Dosya No}_{{Mahalle}}_{Ada/Parsel}({Kurum})

		folderName := folderNamePattern
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

	for _, folderName := range folderNames {
		targetFolderPath := filepath.Join(targetPath, folderName)

		if _, err := os.Stat(targetFolderPath); os.IsNotExist(err) {
			err = os.MkdirAll(targetFolderPath, 0o755)
			if err != nil {
				runtime.LogError(a.ctx, err.Error())
				return
			}
		}

		if copyFolderPath == "" {
			continue
		}

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
			return
		}
	}

	a.SendNotification("Klasör oluşturma başarılı", "", strings.ReplaceAll(targetPath, "\\", "\\\\"), "success")
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
