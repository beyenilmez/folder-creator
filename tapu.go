package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	r "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Tapu struct {
	Cilt  int
	Sayfa int
	Mevki string
	Alan  float64
}

func (app *App) AddTapuToExcel(excelPath string, path string, tapuPathPattern string, ciltHeader string, sayfaHeader string, mevkiHeader string, alanHeader string) string {
	runtime.LogInfo(app.ctx, "Adding tapu to "+excelPath)

	headers, rows, excel, err := ReadExcel(excelPath)
	sheetName := excel.GetSheetList()[0]

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return err.Error()
	}

	ciltIndex := -1
	sayfaIndex := -1
	mevkiIndex := -1
	alanIndex := -1

	for i, header := range headers {
		header = strings.TrimSpace(header)

		if header == ciltHeader {
			ciltIndex = i
		}
		if header == sayfaHeader {
			sayfaIndex = i
		}
		if header == mevkiHeader {
			mevkiIndex = i
		}
		if header == alanHeader {
			alanIndex = i
		}
	}

	runtime.LogInfo(app.ctx, "Indexes: Cilt: "+fmt.Sprint(ciltIndex)+" Sayfa: "+fmt.Sprint(sayfaIndex)+" Mevki: "+fmt.Sprint(mevkiIndex)+" Alan: "+fmt.Sprint(alanIndex))

	for i, row := range rows {
		runtime.WindowExecJS(appContext, `window.setCiltMessage("`+fmt.Sprintf("%d/%d", i+1, len(rows))+`");`)
		runtime.LogDebug(app.ctx, "Generating pattern: "+tapuPathPattern)
		newPattern := generatePatternName(tapuPathPattern, headers, row)
		runtime.LogDebug(app.ctx, "Generated pattern: "+newPattern)

		wholePath := filepath.Join(path, newPattern)

		runtime.LogDebug(app.ctx, "Searching for: "+wholePath)

		matches, err := FilterDirs(wholePath)

		if err != nil {
			app.SendNotification("", strings.ReplaceAll(err.Error(), "\\", "\\\\"), "", "error")
			time.Sleep(time.Second * 1)
			continue
		}

		if len(matches) == 0 {
			runtime.LogInfo(app.ctx, "Tapu not found for row: "+fmt.Sprint(row))
			continue
		}

		runtime.LogDebug(app.ctx, "Found matches: "+fmt.Sprint(matches))

		for _, match := range matches {
			tapu, err := app.ParseTapu(match)

			if err != nil {
				continue
			}

			if ciltIndex != -1 {
				err = excel.SetCellInt(sheetName, alphabet[ciltIndex]+fmt.Sprint(i+2), tapu.Cilt)

				if err != nil {
					runtime.LogError(app.ctx, err.Error())
				}
			}
			if sayfaIndex != -1 {
				err = excel.SetCellInt(sheetName, alphabet[sayfaIndex]+fmt.Sprint(i+2), tapu.Sayfa)

				if err != nil {
					runtime.LogError(app.ctx, err.Error())
				}
			}
			if mevkiIndex != -1 {
				err = excel.SetCellStr(sheetName, alphabet[mevkiIndex]+fmt.Sprint(i+2), tapu.Mevki)

				if err != nil {
					runtime.LogError(app.ctx, err.Error())
				}
			}
			if alanIndex != -1 {
				err = excel.SetCellFloat(sheetName, alphabet[alanIndex]+fmt.Sprint(i+2), tapu.Alan, 2, 64)

				if err != nil {
					runtime.LogError(app.ctx, err.Error())
				}
			}
		}
	}

	// Save
	runtime.LogInfo(app.ctx, "Saving "+excelPath)
	err = excel.SaveAs(excelPath)

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return err.Error()
	}

	runtime.WindowExecJS(appContext, `window.setCiltMessage("Excel dosyası başarıyla güncellendi");`)

	return ""
}

func FilterDirs(path string) ([]string, error) {
	var matches []string

	base := filepath.Base(path)
	dir := filepath.Dir(path)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		match, err := filepath.Match(base, file.Name())
		if err != nil {
			return nil, err
		}
		if match {
			matches = append(matches, filepath.Join(dir, file.Name()))
		}
	}

	return matches, nil
}

func (app *App) ParseTapu(path string) (Tapu, error) {
	var tapu Tapu

	runtime.LogInfo(app.ctx, "Parsing "+path)

	err := installXpdf()
	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return Tapu{}, err
	}

	content, err := ReadPlainTextFromPDF(path)

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return Tapu{}, err
	}

	found := false

	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "Cilt") {
			split := strings.Split(line, ":")
			splitRight := strings.Split(split[1], "/")

			cilt := strings.TrimSpace(splitRight[0])
			sayfa := strings.TrimSpace(splitRight[1])

			runtime.LogInfo(app.ctx, "Cilt: "+cilt+" Sayfa: "+sayfa)

			numberSayfa, err := strconv.Atoi(sayfa)
			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return Tapu{}, err
			}

			numberCilt, err := strconv.Atoi(cilt)
			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return Tapu{}, err
			}

			tapu.Cilt = numberCilt
			tapu.Sayfa = numberSayfa

			found = true
		} else if strings.Contains(line, "Mevki") {
			split := strings.Split(line, ":")

			tapu.Mevki = strings.TrimSpace(split[1])
			tapu.Mevki = toTitleCaseWord(tapu.Mevki)

			found = true
		} else if strings.Contains(line, "Yüzölçüm") {
			splitSpace := strings.Split(line, " ")

			// find Yüzölçüm in split
			for i, word := range splitSpace {
				if word == "Yüzölçüm" && i+2 < len(splitSpace) {
					alanString := strings.TrimSpace(splitSpace[i+2])
					alanString = strings.ReplaceAll(alanString, "m2", "")
					alanString = strings.TrimSpace(alanString)
					alanString = strings.ReplaceAll(alanString, ".", "")
					alanString = strings.ReplaceAll(alanString, ",", ".")

					alan, err := strconv.ParseFloat(alanString, 64)

					if err != nil {
						runtime.LogError(app.ctx, err.Error())
						return Tapu{}, err
					}

					tapu.Alan = alan

					found = true
				}
			}

		}
	}

	if !found {
		return tapu, fmt.Errorf("Tapu bilgisi bulunamadı")
	}

	return tapu, nil
}

func ReadPlainTextFromPDF(pdfpath string) (text string, err error) {
	runtime.LogInfo(appContext, "Reading text from PDF")
	cmd := exec.Command(pdfToTextPath, "-simple2", "-enc", "UTF-8", pdfpath, filepath.Join(appFolder, "temp.txt"))
	err = cmd.Run()

	if err != nil {
		return "", err
	}

	bytes, err := os.ReadFile(filepath.Join(appFolder, "temp.txt"))

	if err != nil {
		return "", err
	}

	defer os.Remove(filepath.Join(appFolder, "temp.txt"))

	return string(bytes), nil
}

// installXpdf checks if Xpdf (pdftotext) is installed and installs it if not.
func installXpdf() error {
	// Check if pdftotext command is available
	_, err := exec.LookPath(pdfToTextPath)
	if err == nil {
		fmt.Println("Xpdf (pdftotext) is already installed.")
		return nil
	}

	app.SendNotification("Xpdf kuruluyor", "Bu işlem birkaç dakika sürebilir", "", "info")

	// Determine download URL based on architecture
	var downloadURL string
	switch r.GOARCH {
	case "amd64":
		downloadURL = "https://dl.xpdfreader.com/xpdf-tools-win-4.05.zip"
	case "386":
		downloadURL = "https://dl.xpdfreader.com/xpdf-tools-mac-4.05.tar.gz"
	default:
		return fmt.Errorf("unsupported architecture: %s", r.GOARCH)
	}

	// Download and extract Xpdf tools
	zipFilePath := filepath.Join(appFolder, "xpdf-tools.zip")
	defer os.Remove(zipFilePath)

	err = downloadFile(downloadURL, zipFilePath)
	if err != nil {
		return fmt.Errorf("error downloading Xpdf tools: %v", err)
	}

	// Unzip the downloaded file
	unzipDir := filepath.Join(appFolder, "xpdf-tools")
	err = unzip(zipFilePath, unzipDir)
	if err != nil {
		return fmt.Errorf("error extracting Xpdf tools: %v", err)
	}

	fmt.Println("Xpdf (pdftotext) has been installed successfully.")
	return nil
}

// downloadFile downloads a file from URL and saves it to filePath.
func downloadFile(url, filePath string) error {
	fmt.Printf("Downloading %s...\n", url)

	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write content to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// unzip extracts a ZIP file to the specified directory.
func unzip(zipFile, destDir string) error {
	fmt.Printf("Extracting %s to %s...\n", zipFile, destDir)

	// Open the zip archive for reading
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create destination directory if it does not exist
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Extract files
	for _, f := range r.File {
		// Open file from zip archive
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create file in destination directory
		path := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			os.MkdirAll(filepath.Dir(path), os.ModePerm)
			w, err := os.Create(path)
			if err != nil {
				return err
			}
			defer w.Close()

			// Copy file contents
			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
