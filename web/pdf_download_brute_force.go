// Written for Intelligence on HTB
package main

import (
	"fmt"
	"github.com/unidoc/unipdf/v3/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
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
	if pdfInfo.Creator != nil {
		fmt.Println(pdfInfo.Creator.Decoded())
	}
	return pdfInfo.Creator.Decoded()
}

// DumpUsersFromPdfs extracts the author names from all PDF files in a directory and writes them to a file
func DumpUsersFromPdfs(dir string) {
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
		fileName := "users.txt"
		fp, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating user output:", err)
			return
		}
		defer fp.Close()
		for _, user := range users {
			fmt.Fprintln(fp, user)
		}
	}
	return
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
	DumpUsersFromPdfs(dir)
}
