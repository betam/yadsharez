package yadsharez

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// OAuthResp represents correct response on OAuth request.
type OAuthResp struct {
	OperationID string `json:"operation_id,omit_empty"`
	Href        string `json:"href"`
	Method      string `json:"method"`
	Templated   bool   `json:"templated"`
}

// Upload request temporary URL for uploads
// and upload the file.
func Upload(filePath, token string) {
	uploadHref, err := getUploadHref(filePath, token)
	if err != nil {
		log.Fatalf("error creating uploading URL: %v\n", err)
	}

	log.Printf("Upload href: %v\n", uploadHref)

	uploadYandexDisk(filePath, uploadHref)
}

// getUploadHref return URL for direct file download.
func getUploadHref(filePath, token string) (string, error) {
	url := getUploadURL(filePath)
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

	uploadHref, err := getHref(res)
	if err != nil {
		log.Printf("error decoding http get request %v\n", err)
	}
	return uploadHref, nil
}

// uploadYandexDisk performs direct file uploading to the Yandex Disk storage.
func uploadYandexDisk(filePath, href string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("error opening file %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fieldName := "file"
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		log.Printf("error creating multipart form for file %s: %v", file.Name(), err)
	}

	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		log.Printf("error closing multipart form: %v", err)
	}

	log.Printf("Creting PUT request to upload to YaDisk\n")
	req, err := http.NewRequest(http.MethodPut, href, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("error uploading file: %v", err)
	}
	defer res.Body.Close()

	return nil
}

// getUploadURL forms correct URL for upload request.
func getUploadURL(file string) string {
	fp := filepath.Base(file)
	return fmt.Sprintf("https://cloud-api.yandex.net:443/v1/disk/resources/upload/?path=app:/%s&overwrite=true", fp)
}

// AuthHeader forms string for Authorization header.
func AuthHeader(token string) string {
	return fmt.Sprintf("OAuth %s", token)
}

// getHref decodes JSON response into OAuthResp struct and
// returns href field from it.
func getHref(res *http.Response) (string, error) {
	var oauth OAuthResp
	err := json.NewDecoder(res.Body).Decode(&oauth)
	if err != nil {
		return "", err
	}

	return oauth.Href, nil
}

// Help prints out the app usage.
func Help() {
	log.Fatalf("Usage: yadsharez [file]\n")
}
