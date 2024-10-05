package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Function to upload the file to Nextcloud via the WebDAV path
func uploadToNextcloud(filePath, nextcloudURL, username, password string) error {
	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create an HTTP PUT request for the upload
	req, err := http.NewRequest("PUT", nextcloudURL, file)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add HTTP Basic Auth
	req.SetBasicAuth(username, password)

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error with HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error uploading file: status code: %d, response: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("file uploaded successfully: %s\n", filePath)

	return nil
}

func main() {
	// Retrieve environment variables
	filePath := os.Getenv("INPUT_FILE_PATH")
	nextcloudURL := os.Getenv("INPUT_NEXTCLOUD_URL")
	username := os.Getenv("INPUT_USERNAME")
	password := os.Getenv("INPUT_PASSWORD")

	// Check if all variables are set
	if filePath == "" || nextcloudURL == "" || username == "" || password == "" {
		fmt.Println("missing inputs! please ensure all necessary parameters are provided.")
		os.Exit(1)
	}

	// Display input values
	fmt.Printf("file path: %s\n", filePath)
	fmt.Printf("Nextcloud URL (WebDAV path): %s\n", nextcloudURL)
	fmt.Printf("username: %s\n", username)

	// Perform the upload
	err := uploadToNextcloud(filePath, nextcloudURL, username, password)
	if err != nil {
		fmt.Printf("an error occurred: %v\n", err)
		os.Exit(1)
	}
}
