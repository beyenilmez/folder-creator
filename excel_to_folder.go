package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

func (a *App) GetFileDialog() string {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Dosya seçin",
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "UDF Dosyası",
				Pattern:     "*.udf",
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

func ReadExcel(excelPath string) ([]string, [][]string, *excelize.File, error) {
	excelFile, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, nil, nil, err
	}
	sheetName := excelFile.GetSheetList()[0]

	// Get all the rows in the Sheet1.
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return nil, nil, nil, err
	}

	return rows[0], rows[1:], excelFile, nil
}

func generatePatternName(pattern string, headers []string, row []string) string {
	for i, header := range headers {
		if i >= len(row) {
			break
		}
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
		err1 := docx.Replace(titlePlaceholder, replacement, -1)
		err2 := docx.ReplaceFooter("{"+placeholder+"}", replacement)
		err3 := docx.ReplaceHeader("{"+placeholder+"}", replacement)

		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		if err3 != nil {
			return err3
		}
	}
	err1 := docx.Replace("{"+placeholder+"}", replacement, -1)
	err2 := docx.ReplaceFooter("{"+placeholder+"}", replacement)
	err3 := docx.ReplaceHeader("{"+placeholder+"}", replacement)

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if err3 != nil {
		return err3
	}
	return nil
}

func replacePlaceholdersUdf(content, placeholder, replacement string) string {
	titlePlaceholder := "{{" + placeholder + "}}"
	if strings.Contains(content, titlePlaceholder) {
		replacement = toTitleCaseWord(replacement)
		content = strings.ReplaceAll(content, titlePlaceholder, replacement)
	}
	return strings.ReplaceAll(content, "{"+placeholder+"}", replacement)
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

func (a *App) CreateFolders(excelPath string, wordPath string, copyFolderPath string, targetPath string, folderNamePattern string, createFolderConfig bool, wordFileNamePattern string, fileNamePattern string, filePath string, wordReplaceRules string) string {
	headers, rows, err := ReadExcelRows(excelPath)

	runtime.LogDebug(a.ctx, "Headers: "+strings.Join(headers, ","))

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err.Error()
	}

	folderNames := generateFolderNames(folderNamePattern, headers, rows)

	for i, folderName := range folderNames {
		var targetFolderPath string

		if createFolderConfig {
			targetFolderPath = filepath.Join(targetPath, folderName)
		} else {
			targetFolderPath = targetPath
		}

		if err := createFolder(targetFolderPath); err != nil {
			runtime.LogError(a.ctx, "Failed to create folder: "+err.Error())
			continue
		}

		if copyFolderPath != "" {
			if err := copyFolderContents(copyFolderPath, targetFolderPath); err != nil {
				runtime.LogError(a.ctx, "Failed to copy folder contents: "+err.Error())
			}
		}

		if wordPath != "" {
			if err := createWordDocument(wordPath, wordFileNamePattern, headers, rows[i], targetFolderPath, wordReplaceRules); err != nil {
				runtime.LogError(a.ctx, "Failed to create word document: "+err.Error())
			}
		}

		if filePath != "" {
			if err := createUdfDocument(filePath, fileNamePattern, headers, rows[i], targetFolderPath); err != nil {
				runtime.LogError(a.ctx, "Failed to create udf document: "+err.Error())
			}
		}

		runtime.WindowExecJS(appContext, `window.setExcelMessage("`+fmt.Sprintf("%d/%d", i+1, len(folderNames))+`");`)
	}

	a.SendNotification("Klasör oluşturma başarılı", "", strings.ReplaceAll(targetPath, "\\", "\\\\"), "success")

	return ""
}

func (a *App) CreateFoldersV2(excelPath string, copyFolderPath string, targetPath string) string {
	headers, rows, err := ReadExcelRows(excelPath)

	runtime.LogDebug(a.ctx, "Headers: "+strings.Join(headers, ","))

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err.Error()
	}

	folderNamePattern := filepath.Base(copyFolderPath)
	folderNames := generateFolderNames(folderNamePattern, headers, rows)

	for i, folderName := range folderNames {
		targetFolderPath := filepath.Join(targetPath, folderName)

		if err := createFolder(targetFolderPath); err != nil {
			runtime.LogError(a.ctx, "Failed to create folder: "+err.Error())
			continue
		}

		if copyFolderPath != "" {
			if err := copyFolderContentsV2(copyFolderPath, targetFolderPath, headers, rows[i]); err != nil {
				runtime.LogError(a.ctx, err.Error())
				continue
			}
		}
		runtime.WindowExecJS(appContext, `window.setExcelMessage("`+fmt.Sprintf("%d/%d", i+1, len(folderNames))+`");`)
	}

	a.SendNotification("Klasör oluşturma başarılı", "", strings.ReplaceAll(targetPath, "\\", "\\\\"), "success")

	return ""
}

func createWordDocument(filePath string, wordFileNamePattern string, headers []string, row []string, targetPath string, wordReplaceRules string) error {
	r, err := docx.ReadDocxFile(filePath)

	if err != nil {
		runtime.LogError(appContext, "Failed to read docx file: "+err.Error())
		return err
	}

	docx1 := r.Editable()

	for i, header := range headers {
		if i >= len(row) {
			runtime.LogWarning(appContext, "Row is shorter than headers")
			break
		}
		colCell := sanitizeCellWord(row[i])
		replacePlaceholdersWord(docx1, header, colCell)
	}

	// Replace rules
	splittedRules := strings.Split(wordReplaceRules, ",")
	for _, rule := range splittedRules {
		splittedRule := strings.Split(rule, "->")

		if len(splittedRule) != 2 {
			return errors.New("wordReplaceRules is not valid")
		}

		if splittedRule[1] == `""` {
			splittedRule[1] = ""
		}

		err1 := docx1.Replace(splittedRule[0], splittedRule[1], -1)
		err2 := docx1.ReplaceFooter(splittedRule[0], splittedRule[1])
		err3 := docx1.ReplaceHeader(splittedRule[0], splittedRule[1])

		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		if err3 != nil {
			return err3
		}
	}

	// Strip the file extension
	wordFileNamePattern = strings.TrimSuffix(wordFileNamePattern, filepath.Ext(wordFileNamePattern))

	fileName := generatePatternName(wordFileNamePattern, headers, row)

	docx1.WriteToFile(filepath.Join(targetPath, fileName) + ".docx")

	err = r.Close()

	if err != nil {
		runtime.LogError(appContext, "Failed to close docx file: "+err.Error())
		return err
	}

	return nil
}

func createUdfDocument(filePath string, fileNamePattern string, headers []string, row []string, targetPath string) error {
	// extract filePath to temp folder (zip)
	tempFolder := filepath.Join(os.TempDir(), "folder-creator-"+uuid.NewString())

	if err := unzip(filePath, tempFolder); err != nil {
		runtime.LogError(appContext, "Failed to unzip file: "+err.Error())
		return err
	}

	contentPath := filepath.Join(tempFolder, "content.xml")

	if _, err := os.Stat(contentPath); err != nil {
		runtime.LogError(appContext, "Failed to find content.xml: "+err.Error())
		return err
	}

	// find and change content of content.xml
	content, err := os.ReadFile(contentPath)
	if err != nil {
		runtime.LogError(appContext, "Failed to read content.xml: "+err.Error())
		return err
	}

	strContent := string(content)

	for i, header := range headers {
		if i >= len(row) {
			runtime.LogWarning(appContext, "Row is shorter than headers")
			break
		}
		colCell := sanitizeCellWord(row[i])
		strContent = replacePlaceholdersUdf(strContent, header, colCell)
	}

	// write content.xml
	if err := os.WriteFile(contentPath, []byte(strContent), 0644); err != nil {
		runtime.LogError(appContext, "Failed to write content.xml: "+err.Error())
		return err
	}

	// Generate file name
	fileNamePattern = strings.TrimSuffix(fileNamePattern, filepath.Ext(fileNamePattern))
	fileName := generatePatternName(fileNamePattern, headers, row) + ".udf"

	// Create a zip at targetPath
	if err := createZip(tempFolder, filepath.Join(targetPath, fileName)); err != nil {
		runtime.LogError(appContext, "Failed to zip folder: "+err.Error())
		return err
	}

	// Delete temp folder
	if err := os.RemoveAll(tempFolder); err != nil {
		runtime.LogError(appContext, "Failed to delete temp folder: "+err.Error())
		return err
	}

	return nil
}

func createZip(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source folder itself
		if path == source {
			return nil
		}

		// Create a zip file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Ensure the file path in the zip file is relative to the source folder
		header.Name = filepath.Base(path)

		// If it's a directory, create the header with a trailing slash
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// Create a writer for the zip file header
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// If it's a directory, return here
		if info.IsDir() {
			return nil
		}

		// Open the file to be zipped
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copy the file contents to the zip writer
		_, err = io.Copy(writer, file)
		return err
	})

	return err
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

func copyFolderContentsV2(src, dest string, headers []string, row []string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if filepath.Ext(relativePath) == ".docx" {
			return createWordDocument(path, filepath.Base(path), headers, row, dest, "")
		} else if filepath.Ext(relativePath) == ".udf" {
			return createUdfDocument(path, filepath.Base(path), headers, row, dest)
		}

		relativePath = generatePatternName(relativePath, headers, row)

		targetPath := filepath.Join(dest, relativePath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath)
	})
}
