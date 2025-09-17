package report

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/xuri/excelize/v2"
	"irpl.com/kanban-commons/bash"
	"irpl.com/kanban-commons/utils"
)

type ByteReport struct{}

type Report struct {
	Template string
	Values   []string
	OutPath  string
}

type MultiReport struct {
	Pages []*Report
}

func (r *MultiReport) AddPage(page *Report) {
	r.Pages = append(r.Pages, page)
}

func (report *Report) CreateXLSXPage() (string, string) {
	var basename string

	if xlsx, templErr := excelize.OpenFile(report.Template); templErr != nil {
		log.Println(templErr.Error())
	} else {

		MASTER := "master"
		REPORT := "report"

		if val, _ := xlsx.GetSheetIndex(MASTER); val > 0 {
			for _, cell := range report.Values {
				if strings.HasPrefix(cell, "${") {
					re, _ := regexp.Compile(`\$\{(.*?)\}(.*)`)
					match := re.FindStringSubmatch(cell)
					if len(match) == 3 {
						cellCoordinate := match[1]
						value := match[2]
						if strings.HasPrefix(value, "image:") {
							imagePath := strings.TrimPrefix(value, "image:")
							err := addImageToCell(xlsx, REPORT, cellCoordinate, imagePath)
							if err != nil {
								log.Println("Error adding image:", err)
							}
						} else {
							var cellValue interface{} = value
							xlsx.SetCellValue(MASTER, cellCoordinate, cellValue)
						}
					}
				}
			}
			xlsx.UpdateLinkedValue()

			// basename = out + parse.UniqueFileName()
			basename = report.OutPath
			xlsx.SaveAs(basename + ".xlsx")
			os.Chmod(basename+".xlsx", 0777)
		}
	}
	return basename + ".xlsx", basename
}
func (parse *Report) UniqueFileName() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1e6))
}

// addImageToCell adds an image to a specific cell in the Excel file.
func addImageToCell(xlsx *excelize.File, sheet, cell, imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return err
	}

	tempFilePath := "temp_image.png"
	if err := os.WriteFile(tempFilePath, buffer.Bytes(), 0644); err != nil {
		return err
	}
	defer os.Remove(tempFilePath)

	err = xlsx.AddPicture(sheet, cell, tempFilePath, nil)
	return err
}

// ConvertExcelToPDF converts an Excel file to a PDF file using LibreOffice.
func ConvertExcelToPDF(sheetName, inputFilePath, tempFilePath, pdfOutDir string) (pdfFilePath string, err error) {
	// Open the Excel file
	f, err := excelize.OpenFile(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open Excel file: %s", err)
	}

	for _, sheet := range f.GetSheetMap() {
		if sheet != sheetName {
			if err := f.SetSheetVisible(sheet, false); err != nil {
				return "", fmt.Errorf("failed to hide sheet %s: %v", sheet, err)
			}
		}
	}

	if err := f.SaveAs(tempFilePath); err != nil {
		return "", fmt.Errorf("failed to save temporary Excel file: %s", err)
	}

	os.Chmod(tempFilePath, 0777)

	// Construct the LibreOffice command for conversion
	command := fmt.Sprintf("soffice --headless --norestore --quickstart --invisible --nodefault --nofirststartwizard --nolockcheck --nologo --convert-to pdf --outdir %s %s", pdfOutDir, tempFilePath)
	bash.NatsCli().SendCommand(command)

	pdfFileName := filepath.Base(tempFilePath[:len(tempFilePath)-len(filepath.Ext(tempFilePath))] + ".pdf")
	pdfFilePath = filepath.Join(pdfOutDir, pdfFileName)

	return pdfFilePath, nil
}

func MoveFile(sourcePath, destinationPath string) error {

	if err := os.Chmod(sourcePath, 0777); err != nil {
		log.Println(err)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destinationFile.Close()

	if err = os.Chmod(destinationPath, 0777); err != nil {
		log.Println(err)
	}

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	sourceFile.Close()
	destinationFile.Close()

	if err = os.Chmod(destinationPath, 0777); err != nil {
		log.Println(err)
	}

	if err := os.Remove(sourcePath); err != nil {
		return fmt.Errorf("failed to remove source file: %v", err)
	}

	return nil
}

type PdfDialog struct {
	Data       string
	RecordId   string
	PageSource string
}

func (p *PdfDialog) Build() string {

	downloadReportURL := ""

	downloadReportURL = utils.JoinStr(p.PageSource)

	//<!--html-->
	ret := utils.JoinStr(`

	<div id="report-dialog" class="modal fade" tabindex="-1" role="dialog" data-backdrop="static"  data-bs-backdrop="static"  data-bs-backdrop="static"  >
	<div class="modal-dialog modal-xl" role="document">
		<div class="modal-content">

			<div class="modal-body">
				<div class="container-fluid">
					<div class="row">
						<div class="col">
							<embed id="embed" width="100%" height="600" title="SamplePdf" type="application/pdf" src="`, p.PageSource, `" />
						</div>
					</div>
				</div>
			</div>

			<div class="modal-footer">
				<a class="btn btn-success" title="Download Report" href="`, downloadReportURL, `" download>
					<!-- <span id="report-spinner" class="spinner-border spinner-border-sm spinner-action" style="display:none" role="status" ></span> -->
					<i class="fa fa-download"></i> Download
				</a>
				<button type="button" class="btn btn-danger" data-bs-dismiss="modal"  data-dismiss="modal">Cancel</button>
			</div>
		</div>
	</div>
	<script>

		//var midTop = $("#report-dialog").offset().top
		//var wh = $(window).height() - (2*midTop);
		var left = $("#report-dialog").find("#embed").attr("height", $(window).height() - 200);

	</script>
	</div>
    `)
	//<!--!html-->
	return ret
}

// MergePdfs merges two or more pdf files into one
func (parse *ByteReport) MergePdfs(files []string, outPathName string) string {
	config := model.NewDefaultConfiguration()
	config.Reader15 = true // For compatibility with older PDFs

	err := api.MergeCreateFile(files, outPathName, false, config)
	if err != nil {
		log.Println("MergePdfs: Failed to merge pdf")
	}

	return outPathName
}

func HandleFileWithTimeout(pdfFilePath, tempFilePath string) bool {
	timeout := time.After(5 * time.Minute)
	for {
		select {
		case <-timeout:

			return true
		default:
			_, err := os.Stat(pdfFilePath)
			if err == nil {

				err = os.Chmod(pdfFilePath, 0777)
				if err != nil {
					return true
				}

				err = os.Remove(tempFilePath)
				if err != nil {
					return true
				}

				return true
			}

			if os.IsNotExist(err) {

				time.Sleep(5 * time.Second)
			} else {

				return true
			}
		}
	}
}

func Encode(file string) string {

	encoded := ""
	if f, err := os.Open(file); err == nil {

		reader := bufio.NewReader(f)
		content, _ := io.ReadAll(reader)

		encoded = base64.StdEncoding.EncodeToString(content)
	} else {
		log.Println(err)
	}

	return encoded

}
