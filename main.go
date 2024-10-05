package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Function to upload the file to Nextcloud via the WebDAV path
func uploadToNextcloud(filePath, nextcloudURL, username, password string, override bool) error {
	// Check if the file exists on Nextcloud
	if !override {
		req, err := http.NewRequest("HEAD", nextcloudURL, nil)
		if err != nil {
			return fmt.Errorf("error creating HEAD request: %v", err)
		}

		// Add HTTP Basic Auth
		req.SetBasicAuth(username, password)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error checking if file exists: %v", err)
		}
		defer resp.Body.Close()

		// If the file exists (status 200), we don't want to overwrite it
		if resp.StatusCode == http.StatusOK {
			return fmt.Errorf("file already exists at %s, set override to true to overwrite", nextcloudURL)
		}
	}

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
	override := os.Getenv("INPUT_OVERRIDE")

	// Parse override flag
	overrideFlag := false
	if override == "true" {
		overrideFlag = true
	}

	// Check if all variables are set
	if filePath == "" || nextcloudURL == "" || username == "" || password == "" {
		fmt.Println("missing inputs! please ensure all necessary parameters are provided.")
		os.Exit(1)
	}

	// Display input values
	fmt.Printf("file path: %s\n", filePath)
	fmt.Printf("Nextcloud URL (WebDAV path): %s\n", nextcloudURL)
	fmt.Printf("username: %s\n", username)
	fmt.Printf("override: %v\n", overrideFlag)

	// Perform the upload
	err := uploadToNextcloud(filePath, nextcloudURL, username, password, overrideFlag)
	if err != nil {
		fmt.Printf("an error occurred: %v\n", err)
		os.Exit(1)
	}
}
