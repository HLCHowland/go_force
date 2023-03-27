// Written for Intelligence on HTB
package main

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/model"
	"code.sajari.com/docconv" //extracts pdf text
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"log"
)

// GenerateFileNames generates a list of PDFs with names corresponding to a certain time period
func GenerateFileNames() []string {
	var fileNames []string
	startDate := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2020, time.December, 31, 0, 0, 0, 0, time.UTC)
	for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		fileNames = append(fileNames, (d.Format("2006-01-02") + "-upload.pdf"))
	}
	return fileNames
}

// ExtractPdfCreator extracts the creator name from a PDF file
func ExtractPdfCreator(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return ""
	}
	pdfInfo, err := pdfReader.GetPdfInfo()
	if err != nil {
		return ""
	}
	if pdfInfo.Creator == nil {
		return ""
	}
	return pdfInfo.Creator.Decoded()
}

// DumpUsersFromPdfsXmp extracts the author names from all PDF files in a directory and writes them to a file
func DumpUsersFromPdfsXmp(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Unable to read PDF directory:", dir)
		return
	}
	pattern := ".*\\.pdf$"
	var users []string
	for _, file := range files {
		matched, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			fmt.Println(err)
		}
		if matched {
			path := filepath.Join(dir, file.Name())
			// In the future we could possibly look for more values like author
			users = append(users, ExtractPdfCreator(path))
		}
	}
	// wait to collect before dumping to file to avoid IO slowdown for little data write
	outFileName := "users.txt"
	fp, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Error creating user output:", err)
		return
	}
	defer fp.Close()
	for _, user := range users {
		fmt.Fprintln(fp, user)
	}
	fmt.Println("Wrote all potential users to:", outFileName)
}

// ConvertPdfToText does what it says on the tin
func ConvertPdfToText(path string) string {
    // This should be replaced at some point in time, this library is reliant on you having poppler-tools installed
	// and actually calls another executable to do the parsing, very lazy, but all I could find that worked
	res, err := docconv.ConvertPath(path)
    if err != nil {
        log.Fatal(err)
    }
	return res.Body
}

// ConvertPdfsToTextFile converts a whole directory of PDFs into a single text file 
func ConvertPdfsToTextFile(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Unable to read PDF directory:", dir)
	}
	outFileName := "collated_pdfs.txt"
	fp, err := os.Create(outFileName)
	if err != nil {
		fmt.Println("Unable to create collated pdfs file")
	}
	pattern := ".*\\.pdf"
	for _, file := range files {
		match, err := regexp.MatchString(pattern, file.Name())
		if err != nil {
			fmt.Println("Issue performing regex")
			return
		}
		if match {
			path := filepath.Join(dir, file.Name())
			pdfText := ConvertPdfToText(path)
			// Write to disk immidiately as to not hog memory with multiple possibly large PDFs of text
			fmt.Fprintln(fp, pdfText) 
		}
	}
	fmt.Println("Wrote all PDF text to:", outFileName)
}

// RequestAndSavePdf requests a PDF file from a URL and saves it to a directory
func RequestAndSavePdf(url string, fileName string, dir string) {
	// The third argument of NewRequest is normally the request body
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	// need to find out what exactly this is doing
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return
	}
	fmt.Println("Found file:", fileName)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	// The reason we use filepath.Join here is so that the filepath will be OS
	// independent.
	filePath := filepath.Join(dir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}
}

func main() {
	dir := "files_found"
	fileNames := GenerateFileNames()
	for _, fileName := range fileNames {
		RequestAndSavePdf(("http://intelligence.htb/Documents/" + fileName), fileName, dir)
	}
	DumpUsersFromPdfsXmp(dir)
	ConvertPdfsToTextFile(dir)
}
