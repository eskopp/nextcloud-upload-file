package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

// Function to rename a file physically
func renameFile(originalPath, newPath string) error {
	err := os.Rename(originalPath, newPath)
	if err != nil {
		return fmt.Errorf("error renaming file: %v", err)
	}
	fmt.Printf("File renamed successfully: %s -> %s\n", originalPath, newPath)
	return nil
}

// Function to zip a file
func zipFile(sourceFile, destinationZip string) error {
	zipfile, err := os.Create(destinationZip)
	if err != nil {
		return fmt.Errorf("error creating zip file: %v", err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	fileToZip, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("error opening file to zip: %v", err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("error creating zip header: %v", err)
	}

	header.Name = filepath.Base(sourceFile)
	writer, err := archive.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("error creating zip writer: %v", err)
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return fmt.Errorf("error writing to zip: %v", err)
	}

	fmt.Printf("File zipped successfully: %s\n", destinationZip)
	return nil
}

// Function to append date and/or time to the filename and physically rename the file
func appendDateTimeAndRename(filePath string, dateFlag, timeFlag bool) (string, error) {
	ext := filepath.Ext(filePath)
	baseName := filePath[0 : len(filePath)-len(ext)]

	// Get current date and time
	currentTime := time.Now()
	if dateFlag {
		baseName += currentTime.Format("_2006_01_02") // Append date
	}
	if timeFlag {
		baseName += currentTime.Format("_15_04_05") // Append time
	}

	// New file name with date/time
	newPath := baseName + ext

	// Rename the file physically
	err := renameFile(filePath, newPath)
	if err != nil {
		return "", err
	}

	return newPath, nil
}

// Function to upload the file to Nextcloud via WebDAV
func uploadToNextcloud(filePath, nextcloudURL, username, password string, override, zipFlag bool) error {
	// Get the base name of the file (e.g., "search.css")
	fileName := path.Base(filePath)

	// Ensure nextcloudURL ends with a "/"
	if nextcloudURL[len(nextcloudURL)-1] != '/' {
		nextcloudURL += "/"
	}

	// If zipFlag is true, zip the file
	if zipFlag {
		zipFilePath := filePath + ".zip"
		err := zipFile(filePath, zipFilePath)
		if err != nil {
			return fmt.Errorf("error zipping file: %v", err)
		}
		filePath = zipFilePath       // Use the zipped file for upload
		fileName = fileName + ".zip" // Use the zipped file name for upload
		fmt.Printf("Uploading zipped file: %s\n", fileName)
	}

	// Construct the full URL for the file in Nextcloud
	uploadURL := nextcloudURL + fileName

	// Check if the file exists on Nextcloud
	if !override {
		req, err := http.NewRequest("HEAD", uploadURL, nil)
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
			return fmt.Errorf("file already exists at %s, set override to true to overwrite", uploadURL)
		}
	}

	// Open the local file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create an HTTP PUT request for the upload
	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Add HTTP Basic Auth
	req.SetBasicAuth(username, password)

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/octet-stream")

	// Execute the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error with HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		fmt.Printf("File uploaded successfully: %s\n", filePath)
		return nil
	}

	// Print the response body for further debugging
	fmt.Printf("Response body: %s\n", string(body))
	return fmt.Errorf("error uploading file: status code: %d, response: %s", resp.StatusCode, string(body))
}

func main() {
	// Retrieve environment variables
	filePath := os.Getenv("INPUT_FILE_PATH")
	nextcloudURL := os.Getenv("INPUT_NEXTCLOUD_URL")
	username := os.Getenv("INPUT_USERNAME")
	password := os.Getenv("INPUT_PASSWORD")
	overrideStr := os.Getenv("INPUT_OVERRIDE")
	rename := os.Getenv("INPUT_RENAME")
	zipStr := os.Getenv("INPUT_ZIP")
	dateStr := os.Getenv("INPUT_DATE")
	timeStr := os.Getenv("INPUT_TIME")

	// Parse override flag
	overrideFlag, err := strconv.ParseBool(overrideStr)
	if err != nil {
		fmt.Printf("Invalid value for override flag. Must be true or false, received: %s\n", overrideStr)
		os.Exit(1)
	}

	// Parse zip flag
	zipFlag, err := strconv.ParseBool(zipStr)
	if err != nil {
		fmt.Printf("Invalid value for zip flag. Must be true or false, received: %s\n", zipStr)
		os.Exit(1)
	}

	// Parse date and time flags
	dateFlag, err := strconv.ParseBool(dateStr)
	if err != nil {
		fmt.Printf("Invalid value for date flag. Must be true or false, received: %s\n", dateStr)
		os.Exit(1)
	}

	timeFlag, err := strconv.ParseBool(timeStr)
	if err != nil {
		fmt.Printf("Invalid value for time flag. Must be true or false, received: %s\n", timeStr)
		os.Exit(1)
	}

	// Check if all variables are set
	if filePath == "" || nextcloudURL == "" || username == "" || password == "" {
		fmt.Println("missing inputs! please ensure all necessary parameters are provided.")
		os.Exit(1)
	}

	// If the rename flag is set and not "false", rename the file
	if rename != "" && rename != "false" {
		newPath := filepath.Join(filepath.Dir(filePath), rename)
		err := renameFile(filePath, newPath)
		if err != nil {
			fmt.Printf("An error occurred while renaming the file: %v\n", err)
			os.Exit(1)
		}
		filePath = newPath
	}

	// Append date/time to the file name and physically rename the file if flags are true
	if dateFlag || timeFlag {
		newFilePath, err := appendDateTimeAndRename(filePath, dateFlag, timeFlag)
		if err != nil {
			fmt.Printf("An error occurred while renaming the file with date/time: %v\n", err)
			os.Exit(1)
		}
		filePath = newFilePath
	}

	// Perform the upload (with zip if necessary)
	err = uploadToNextcloud(filePath, nextcloudURL, username, password, overrideFlag, zipFlag)
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}
}
