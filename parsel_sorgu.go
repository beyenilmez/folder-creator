package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// FeatureCollection represents the top-level structure of the JSON.
type FeatureCollection struct {
	Features []Feature `json:"features"`
	Type     string    `json:"type"`
	CRS      struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
}

// Feature represents each feature within the FeatureCollection.
type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

// Geometry represents the geometry of a feature.
type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// Properties represents the properties of a feature.
type Properties struct {
	ParselNo string `json:"ParselNo"`
	Alan     string `json:"Alan"`
	Mevkii   string `json:"Mevkii"`
	Nitelik  string `json:"Nitelik"`
	Ada      string `json:"Ada"`
	Il       string `json:"Il"`
	Ilce     string `json:"Ilce"`
	Pafta    string `json:"Pafta"`
	Mahalle  string `json:"Mahalle"`
}

type QueryParams struct {
	Province     string `json:"province"`
	District     string `json:"district"`
	Neighborhood string `json:"neighborhood"`
	Block        string `json:"block"`
	Parcel       string `json:"parcel"`
}

var parselSorguContext context.Context
var close1 context.CancelFunc
var close2 context.CancelFunc
var close3 context.CancelFunc
var downloadDir string

func (a *App) InitParselSorgu(headless bool) error {
	downloadDir = filepath.Join(appFolder, "downloads")

	err := create_folder(downloadDir)

	if err != nil {
		return err
	}

	// Create download directory if it doesn't exist
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	// Create a new context with logging allocator and headful mode
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless), // Run in headful mode
		chromedp.Flag("disable-gpu", headless),
		//chromedp.Flag("start-maximized", true),
	)
	var allocCtx context.Context
	allocCtx, close1 = chromedp.NewExecAllocator(context.Background(), opts...)

	// Create a new context with visible browser window
	var ctx context.Context
	ctx, close2 = chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	timeOut := 2400 * time.Second
	if headless {
		timeOut = 60 * time.Second
	}
	// Set up a timeout
	ctx, close3 = context.WithTimeout(ctx, timeOut)

	err = chromedp.Run(ctx,
		chromedp.Navigate(`https://parselsorgu.tkgm.gov.tr/`),

		// Accept terms and conditions
		chromedp.WaitReady(`//*[@id="terms-ok"]`, chromedp.BySearch),
		chromedp.Click(`//*[@id="terms-ok"]`, chromedp.BySearch),

		browser.SetDownloadBehavior(browser.SetDownloadBehaviorBehaviorAllowAndName).WithDownloadPath(downloadDir).WithEventsEnabled(true),
	)

	if err != nil {
		return err
	}

	parselSorguContext = ctx

	return nil
}

func (app *App) ParselSorgu(params QueryParams) (Properties, error) {
	params.Province = cases.Title(language.Turkish).String(params.Province)
	params.District = cases.Title(language.Turkish).String(params.District)
	params.Neighborhood = cases.Title(language.Turkish).String(params.Neighborhood)

	runtime.LogInfo(app.ctx, fmt.Sprintf("Parsel Sorgu: %s,%s,%s,%s,%s", params.Province, params.District, params.Neighborhood, params.Block, params.Parcel))

	// Clear the download directory
	err := os.RemoveAll(downloadDir)
	if err != nil {
		return Properties{}, err
	}
	err = create_folder(downloadDir)
	if err != nil {
		return Properties{}, err
	}

	if params.Province == "" || params.District == "" || params.Neighborhood == "" || params.Block == "" || params.Parcel == "" {
		return Properties{}, fmt.Errorf("province, district, neighborhood, block and parcel are required")
	}

	err = chromedp.Run(parselSorguContext,
		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		// Select province
		chromedp.ActionFunc(func(ctx context.Context) error {
			runtime.LogInfo(app.ctx, "Waiting for province selection")
			chromedp.WaitVisible(`//*[@id="province-select"]`, chromedp.BySearch)

			// Find the option based on its text content to get its value attribute
			var optionValue string
			err := chromedp.Run(ctx, chromedp.Evaluate(`document.evaluate('//*[@id="province-select"]/option[contains(text(), "`+params.Province+`")]', document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.value;`, &optionValue))
			if err != nil {
				return fmt.Errorf("failed to find option value for province %s: %w", params.Province, err)
			}
			runtime.LogInfo(app.ctx, "Option value: "+optionValue)

			// Set the value of the <select> element to the retrieved optionValue
			err = chromedp.Run(ctx, chromedp.SetValue(`#province-select`, optionValue, chromedp.ByID))
			if err != nil {
				return fmt.Errorf("failed to set province value for %s: %w", params.Province, err)
			}

			return nil
		}),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		// Select district
		chromedp.ActionFunc(func(ctx context.Context) error {
			runtime.LogInfo(app.ctx, "Waiting for district selection")
			chromedp.WaitVisible(`//*[@id="district-select"]`, chromedp.BySearch)

			// Find the option based on its text content to get its value attribute
			var optionValue string
			err := chromedp.Run(ctx, chromedp.Evaluate(`document.evaluate('//*[@id="district-select"]/option[contains(text(), "`+params.District+`")]', document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.value;`, &optionValue))
			if err != nil {
				return fmt.Errorf("failed to find option value for district %s: %w", params.District, err)
			}
			runtime.LogInfo(app.ctx, "Option value: "+optionValue)

			// Set the value of the <select> element to the retrieved optionValue
			err = chromedp.Run(ctx, chromedp.SetValue(`#district-select`, optionValue, chromedp.ByID))
			if err != nil {
				return fmt.Errorf("failed to set district value for %s: %w", params.District, err)
			}

			return nil
		}),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		// Select neighborhood
		chromedp.ActionFunc(func(ctx context.Context) error {
			runtime.LogInfo(app.ctx, "Waiting for neighborhood selection")
			chromedp.WaitVisible(`//*[@id="neighborhood-select"]`, chromedp.BySearch)

			// Find the option based on its text content to get its value attribute
			var optionValue string
			err := chromedp.Run(ctx, chromedp.Evaluate(`document.evaluate('//*[@id="neighborhood-select"]/option[contains(text(), "`+params.Neighborhood+`")]', document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue.value;`, &optionValue))
			if err != nil {
				return fmt.Errorf("failed to find option value for neighborhood %s: %w", params.Neighborhood, err)
			}
			runtime.LogInfo(app.ctx, "Option value: "+optionValue)

			// Set the value of the <select> element to the retrieved optionValue
			err = chromedp.Run(ctx, chromedp.SetValue(`#neighborhood-select`, optionValue, chromedp.ByID))
			if err != nil {
				return fmt.Errorf("failed to set neighborhood value for %s: %w", params.Neighborhood, err)
			}

			return nil
		}),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		// Wait for other elements to be ready and interact with them as needed
		chromedp.WaitReady(`#block-input`, chromedp.ByID),
		chromedp.SendKeys(`#block-input`, params.Block),

		chromedp.WaitReady(`#parcel-input`, chromedp.ByID),
		chromedp.SendKeys(`#parcel-input`, params.Parcel),

		chromedp.Sleep(100*time.Millisecond),

		chromedp.WaitReady(`#administrative-query-btn`, chromedp.ByID),
		chromedp.Click(`#administrative-query-btn`, chromedp.ByID),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		chromedp.Sleep(5*time.Second),

		// stop runs if .close-btn exists
		chromedp.ActionFunc(func(ctx context.Context) error {
			runtime.LogInfo(app.ctx, "Checking for close button")
			var exists bool
			err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelector("#close-btn") !== null`, &exists))
			if err != nil {
				return fmt.Errorf("failed to check for close button: %w", err)
			}
			if exists {
				runtime.LogInfo(app.ctx, "Close button found, clicking it and stopping execution")
				err = chromedp.Run(ctx, chromedp.Click("#close-btn", chromedp.ByID))
				if err != nil {
					return fmt.Errorf("failed to click close button: %w", err)
				}
				return fmt.Errorf("execution stopped due to close button")
			} else {
				runtime.LogInfo(app.ctx, "Close button not found, continuing execution")
			}
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// Evaluate the XPath to find the last path element
			var pathCount int
			err := chromedp.Run(ctx, chromedp.Evaluate(`document.querySelectorAll("#map-canvas > div.leaflet-pane.leaflet-map-pane > div.leaflet-pane.leaflet-overlay-pane > svg > g > path").length`, &pathCount))
			if err != nil {
				return fmt.Errorf("failed to count paths: %w", err)
			}

			// Select the last path element
			pathSelector := fmt.Sprintf(`#map-canvas > div.leaflet-pane.leaflet-map-pane > div.leaflet-pane.leaflet-overlay-pane > svg > g > path:nth-child(%d)`, pathCount)
			err = chromedp.Run(ctx, chromedp.Click(pathSelector))
			if err != nil {
				return fmt.Errorf("failed to click on last path element: %w", err)
			}

			return nil
		}),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		chromedp.WaitReady(".dropdown-toggle", chromedp.BySearch),
		chromedp.Click(".dropdown-toggle", chromedp.BySearch),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		chromedp.WaitReady(`/html/body/div[3]/div[3]/div/div/div[1]/div/ul/li[11]/a`, chromedp.BySearch),
		chromedp.Click(`/html/body/div[3]/div[3]/div/div/div[1]/div/ul/li[11]/a`, chromedp.BySearch),

		// wait while loading
		chromedp.WaitNotPresent(`.nprogress-busy`, chromedp.BySearch),

		chromedp.WaitReady(`/html/body/div[3]/div[3]/div/div/div[2]/div[1]/div/div[3]/table/tbody/tr/td[3]/input`, chromedp.BySearch),
		chromedp.Click(`/html/body/div[3]/div[3]/div/div/div[2]/div[1]/div/div[3]/table/tbody/tr/td[3]/input`, chromedp.BySearch),

		// Download button
		chromedp.WaitReady(`#export-data`, chromedp.ByID),
		chromedp.Click(`#export-data`, chromedp.ByID),

		// wait while loading
		chromedp.Sleep(3*time.Second),

		// Close button
		chromedp.WaitReady(`#close-btn`, chromedp.ByID),
		chromedp.Click(`#close-btn`, chromedp.ByID),

		// wait while loading
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return Properties{}, err
	}

	// Select the only file in downloads folder
	files, err := os.ReadDir(downloadDir)
	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return Properties{}, err
	}

	for {
		if len(files) == 0 {
			runtime.LogWarning(app.ctx, "Downloaded file not found")
			time.Sleep(200 * time.Millisecond)
			// Select the only file in downloads folder
			files, err = os.ReadDir(downloadDir)
			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return Properties{}, err
			}
		} else {
			break
		}
	}

	filePath := filepath.Join(downloadDir, files[0].Name())

	if len(files) != 1 {
		runtime.LogError(app.ctx, "More than one downloaded file found")

		// Use the last added file
		// Sort files by added date
		sort.Slice(files, func(i, j int) bool {
			infoI, err := os.Stat(path.Join(downloadDir, files[i].Name()))
			if err != nil {
				return false
			}

			infoJ, err := os.Stat(path.Join(downloadDir, files[j].Name()))
			if err != nil {
				return false
			}

			return infoI.ModTime().After(infoJ.ModTime())
		})

		runtime.LogInfo(app.ctx, "Last added file: "+files[0].Name())
		filePath = filepath.Join(downloadDir, files[0].Name())
	}

	runtime.LogInfo(app.ctx, "Processing file: "+filePath)

	// Read file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return Properties{}, fmt.Errorf("failed to read downloaded file: %w", err)
	}

	time.Sleep(2 * time.Second)

	// Delete the downloaded file
	err = os.Remove(filePath)
	if err != nil {
		return Properties{}, fmt.Errorf("failed to delete downloaded file: %w", err)
	}

	// Unmarshal JSON content into FeatureCollection struct
	var featureCollection FeatureCollection
	if err := json.Unmarshal(fileContent, &featureCollection); err != nil {
		return Properties{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	feature := featureCollection.Features[0]
	properties := feature.Properties

	runtime.LogInfo(app.ctx, "Feature properties: "+fmt.Sprint(properties))

	time.Sleep(2 * time.Second)

	return properties, nil
}

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL", "AM", "AN", "AO", "AP", "AQ", "AR", "AS", "AT", "AU", "AV", "AW", "AX", "AY", "AZ"}

func (app *App) AddParselSorguFields(excelPath string, ilHeader, ilceHeader, mahalleHeader, adaHeader, parselHeader, alanHeader, paftaHeader, cinsHeader, mevkiHeader string, headless bool) error {
	if !headless {
		err := app.InitParselSorgu(headless)

		if err != nil {
			runtime.LogError(app.ctx, err.Error())
			return err
		}
	}

	headers, rows, excel, err := ReadExcel(excelPath)
	sheetName := excel.GetSheetList()[0]

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return err
	}

	var il, ilce, mahalle, ada, parsel string

	ilIndex, ilceIndex, mahalleIndex, adaIndex, parselIndex, alanIndex, paftaIndex, cinsIndex, mevkiIndex := -1, -1, -1, -1, -1, -1, -1, -1, -1

	runtime.LogInfo(app.ctx, "Headers: "+fmt.Sprint(headers))

	for i, header := range headers {
		header = strings.TrimSpace(header)

		if header == ilHeader {
			ilIndex = i
			runtime.LogInfo(app.ctx, "Matched ilHeader: "+header)
		} else if header == ilceHeader {
			ilceIndex = i
			runtime.LogInfo(app.ctx, "Matched ilceHeader: "+header)
		} else if header == mahalleHeader {
			mahalleIndex = i
			runtime.LogInfo(app.ctx, "Matched mahalleHeader: "+header)
		} else if header == adaHeader {
			adaIndex = i
			runtime.LogInfo(app.ctx, "Matched adaHeader: "+header)
		} else if header == parselHeader {
			parselIndex = i
			runtime.LogInfo(app.ctx, "Matched parselHeader: "+header)
		} else if header == alanHeader {
			alanIndex = i
			runtime.LogInfo(app.ctx, "Matched alanHeader: "+header)
		} else if header == paftaHeader {
			paftaIndex = i
			runtime.LogInfo(app.ctx, "Matched paftaHeader: "+header)
		} else if header == cinsHeader {
			cinsIndex = i
			runtime.LogInfo(app.ctx, "Matched cinsHeader: "+header)
		} else if header == mevkiHeader {
			mevkiIndex = i
			runtime.LogInfo(app.ctx, "Matched mevkiHeader: "+header)
		}
	}

	if ilIndex == -1 && ilHeader != "" {
		runtime.LogInfo(app.ctx, "Setting all il headers to: "+ilHeader)
		il = ilHeader
	}

	if ilceIndex == -1 && ilceHeader != "" {
		runtime.LogInfo(app.ctx, "Setting all ilce headers to: "+ilceHeader)
		ilce = ilceHeader
	}

	if mahalleIndex == -1 && mahalleHeader != "" {
		runtime.LogInfo(app.ctx, "Setting all mahalle headers to: "+mahalleHeader)
		mahalle = mahalleHeader
	}

	if adaIndex == -1 && adaHeader != "" {
		runtime.LogInfo(app.ctx, "Setting all ada headers to: "+adaHeader)
		ada = adaHeader
	}

	if parselIndex == -1 && parselHeader != "" {
		runtime.LogInfo(app.ctx, "Setting all parsel headers to: "+parselHeader)
		parsel = parselHeader
	}

	for i := 0; i < len(rows); i++ {
		if headless {
			err := app.InitParselSorgu(headless)

			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return err
			}
		}

		runtime.WindowExecJS(appContext, `window.setParselMessage("`+fmt.Sprintf("%d/%d", i+1, len(rows))+`");`)

		row := rows[i]

		if ilIndex != -1 {
			il = row[ilIndex]
		}
		if ilceIndex != -1 {
			ilce = row[ilceIndex]
		}
		if mahalleIndex != -1 {
			mahalle = row[mahalleIndex]
		}
		if adaIndex != -1 {
			ada = row[adaIndex]
		}
		if parselIndex != -1 {
			parsel = row[parselIndex]
		}

		properties, err := app.ParselSorgu(QueryParams{Province: il, District: ilce, Neighborhood: mahalle, Block: ada, Parcel: parsel})

		if err != nil {
			runtime.LogError(app.ctx, err.Error())
			continue
		}

		if alanHeader != "" {
			alan := properties.Alan
			alan = strings.ReplaceAll(alan, ".", "")
			alan = strings.ReplaceAll(alan, ",", ".")

			floa64Alan, err := strconv.ParseFloat(alan, 64)
			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				continue
			}
			err = excel.SetCellFloat(sheetName, alphabet[alanIndex]+fmt.Sprint(i+2), floa64Alan, 2, 32)

			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return err
			}
		}

		if paftaHeader != "" {
			pafta := properties.Pafta

			err = excel.SetCellStr(sheetName, alphabet[paftaIndex]+fmt.Sprint(i+2), pafta)

			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return err
			}
		}

		if cinsHeader != "" {
			cins := properties.Nitelik
			err = excel.SetCellStr(sheetName, alphabet[cinsIndex]+fmt.Sprint(i+2), cins)

			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return err
			}
		}

		if mevkiHeader != "" {
			mevki := properties.Mevkii
			err = excel.SetCellStr(sheetName, alphabet[mevkiIndex]+fmt.Sprint(i+2), mevki)

			if err != nil {
				runtime.LogError(app.ctx, err.Error())
				return err
			}
		}

		if headless {
			close1()
			close2()
			close3()
		}
	}

	err = excel.SaveAs(excelPath)

	if err != nil {
		runtime.LogError(app.ctx, err.Error())
		return err
	}

	app.SendNotification("Excel dosyası başarıyla güncellendi", "", "", "success")

	runtime.WindowExecJS(appContext, `window.setParselMessage("Excel dosyası başarıyla güncellendi");`)

	if !headless {
		close1()
		close2()
		close3()
	}

	return nil
}
