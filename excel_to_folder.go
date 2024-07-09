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

func (a *App) CreateFolders(excelPath string, copyFolderPath string, targetPath string) {
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

	for i, row := range rows {
		if i == 0 {
			continue
		}
		var folderNameElements []string

		for j, colCell := range row {
			if j == 0 {
				continue
			}
			colCell = strings.TrimSpace(colCell)

			words := strings.Split(colCell, " ")
			for k, word := range words {
				word = cases.Title(language.Turkish).String(word)
				words[k] = word
			}
			colCell = strings.Join(words, "_")
			colCell = strings.ReplaceAll(colCell, "/", "_")

			folderNameElements = append(folderNameElements, colCell)
		}

		// join except last element
		folderName := strings.Join(folderNameElements[:len(folderNameElements)-1], "_")
		// append last element
		folderName = folderName + "(" + folderNameElements[len(folderNameElements)-1] + ")"

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
