package yadsh

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/fatih/color"
)

// Download prints out download link for the particular file stored on the YandexDisk.
func Download(filePath, token string) {
	downloadHref, err := getDownloadHref(filePath, token)
	if err != nil {
		log.Fatalf("error creating uploading URL: %v\n", err)
	}

	log.Printf("Copy the download link:")
	d := color.New(color.FgGreen, color.Bold)
	d.Printf("%s\n", downloadHref)
}

// getDownloadHref return URL for direct file download.
func getDownloadHref(filePath, token string) (string, error) {
	url := getDownloadURL(filePath)
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error creating http get request %v\n", err)
		return "", err
	}
	req.Header.Set("Authorization", AuthHeader(token))

	res, err := client.Do(req)
	if err != nil {
		log.Printf("error doing http get request %v\n", err)
		return "", err
	}
	defer res.Body.Close()

	downloadHref, err := getHref(res)
	if err != nil {
		log.Printf("error decoding http get request %v\n", err)
	}
	return downloadHref, nil
}

// getDownloadURL returns API URL endpoint for download request.
func getDownloadURL(file string) string {
	fp := filepath.Base(file)
	return fmt.Sprintf("https://cloud-api.yandex.net:443/v1/disk/resources/download/?path=app:/%s&overwrite=true", fp)
}
