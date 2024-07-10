package main

import (
	"fmt"
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

func ReadExcelRows(excelPath string) ([]string, [][]string, error) {
	excelFile, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, nil, err
	}
	sheetName := excelFile.GetSheetList()[0]

	// Get all the rows in the Sheet1.
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, nil, err
	}

	return rows[0], rows[1:], nil
}

func generatePatternName(pattern string, headers []string, row []string) string {
	for i, header := range headers {
		colCell := sanitizeCellFolder(row[i])
		pattern = replacePlaceholdersFolder(pattern, header, colCell)
	}
	return pattern
}

func generateFolderNames(folderNamePattern string, headers []string, rows [][]string) []string {
	var folderNames []string
	for _, row := range rows {
		folderName := generatePatternName(folderNamePattern, headers, row)
		folderNames = append(folderNames, strings.TrimSpace(folderName))
	}
	return folderNames
}

func sanitizeCellFolder(cell string) string {
	cell = strings.TrimSpace(cell)
	cell = strings.ReplaceAll(cell, "\r\n", "_")
	cell = strings.ReplaceAll(cell, "\n", "_")
	cell = strings.ReplaceAll(cell, "\t", "_")
	cell = strings.ReplaceAll(cell, "/", "_")
	return cell
}

func sanitizeCellWord(cell string) string {
	cell = strings.TrimSpace(cell)
	cell = strings.ReplaceAll(cell, "\r\n", " ")
	cell = strings.ReplaceAll(cell, "\n", " ")
	cell = strings.ReplaceAll(cell, "\t", " ")
	return cell
}

func replacePlaceholdersFolder(text, placeholder, replacement string) string {
	titlePlaceholder := "{{" + placeholder + "}}"
	if strings.Contains(text, titlePlaceholder) {
		replacement = toTitleCaseFolder(replacement)
		text = strings.ReplaceAll(text, titlePlaceholder, replacement)
	}
	return strings.ReplaceAll(text, "{"+placeholder+"}", replacement)
}

func replacePlaceholdersWord(docx *docx.Docx, placeholder, replacement string) error {
	titlePlaceholder := "{{" + placeholder + "}}"
	if strings.Contains(docx.GetContent(), titlePlaceholder) {
		replacement = toTitleCaseWord(replacement)
		return docx.Replace(titlePlaceholder, replacement, -1)
	}
	return docx.Replace("{"+placeholder+"}", replacement, -1)
}

func toTitleCaseFolder(text string) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		words[i] = cases.Title(language.Turkish).String(word)
	}
	return strings.Join(words, "_")
}

func toTitleCaseWord(text string) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		words[i] = cases.Title(language.Turkish).String(word)
	}
	return strings.Join(words, " ")
}

func (a *App) CreateFolders(excelPath string, wordPath string, copyFolderPath string, targetPath string, folderNamePattern string, wordFileNamePattern string) string {
	headers, rows, err := ReadExcelRows(excelPath)

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err.Error()
	}

	folderNames := generateFolderNames(folderNamePattern, headers, rows)

	for i, folderName := range folderNames {
		targetFolderPath := filepath.Join(targetPath, folderName)

		if err := createFolder(targetFolderPath); err != nil {
			runtime.LogError(a.ctx, "Failed to create folder: "+err.Error())
			continue
		}

		if copyFolderPath != "" {
			if err := copyFolderContents(copyFolderPath, targetFolderPath); err != nil {
				runtime.LogError(a.ctx, err.Error())
				continue
			}
		}

		if wordPath != "" {
			if err := createWordDocument(wordPath, wordFileNamePattern, headers, rows[i], targetFolderPath); err != nil {
				runtime.LogError(a.ctx, err.Error())
				continue
			}
		}

		runtime.WindowExecJS(appContext, `window.setExcelMessage("Klasörler oluşturuluyor: `+fmt.Sprintf("%d/%d", i+1, len(folderNames))+`");`)
	}

	a.SendNotification("Klasör oluşturma başarılı", "", strings.ReplaceAll(targetPath, "\\", "\\\\"), "success")

	return ""
}

func createWordDocument(filePath string, wordFileNamePattern string, headers []string, row []string, targetPath string) error {
	r, err := docx.ReadDocxFile(filePath)

	if err != nil {
		runtime.LogError(appContext, err.Error())
		return err
	}

	docx1 := r.Editable()

	for i, header := range headers {
		colCell := sanitizeCellWord(row[i])
		replacePlaceholdersWord(docx1, header, colCell)
	}

	fileName := generatePatternName(wordFileNamePattern, headers, row)

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

func createFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func copyFolderContents(src, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dest, relativePath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath)
	})
}
