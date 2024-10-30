package ucodesdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func (u *object) Files() FilesI {
	return &APIFiles{
		config: u.config,
	}
}

type FilesI interface {
	/*
		Upload is a function that uploads a file to the server.

		Works for [Mongo, Postgres]

		sdk.Files().
			Upload("file_path").
			Exec()

		Use this method to store a file and obtain its metadata for retrieval or management.
	*/
	Upload(filePath string) *UploadFile
	/*
		Delete is a function that deletes a file from the server.

		Works for [Mongo, Postgres]

		sdk.Files().
			Delete("file_id").
			Exec()

		This method removes a file based on its unique identifier, allowing for clean file management.
	*/
	Delete(fileID string) *DeleteFile
}

func (f *APIFiles) Upload(filePath string) *UploadFile {
	return &UploadFile{
		config: f.config,
		path:   filePath,
	}
}

func (c *UploadFile) Exec() (CreateFileResponse, Response, error) {
	var (
		file          *os.File
		fileBuffer    bytes.Buffer
		writer        *multipart.Writer
		response      = Response{Status: "done"}
		createdObject CreateFileResponse
		url           = fmt.Sprintf("%s/v1/files/folder_upload?folder_name=Media", c.config.BaseURL)
	)

	file, err := os.Open(c.path)
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't open file by path", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}
	defer file.Close()

	writer = multipart.NewWriter(&fileBuffer)
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't create from file", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't copy file", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	err = writer.Close()
	if err != nil {
		response.Data = map[string]any{"description": string(c.path), "message": "can't close writer", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	var appId = c.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	createFileInByte, err := DoFileRequest(url, http.MethodPost, header, fileBuffer, writer)
	if err != nil {
		response.Data = map[string]any{"description": string(createFileInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	err = json.Unmarshal(createFileInByte, &createdObject)
	if err != nil {
		response.Data = map[string]any{"description": string(createFileInByte), "message": "Error while unmarshalling create file object", "error": err.Error()}
		response.Status = "error"
		return CreateFileResponse{}, response, err
	}

	return createdObject, response, nil
}

func (f *APIFiles) Delete(fileID string) *DeleteFile {
	return &DeleteFile{
		config: f.config,
		id:     fileID,
	}
}

func (a *DeleteFile) Exec() (Response, error) {
	var (
		response = Response{Status: "done"}
		url      = fmt.Sprintf("%s/v1/files/%s", a.config.BaseURL, a.id)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := DoRequest(url, http.MethodDelete, Request{Data: map[string]any{}}, header)
	if err != nil {
		response.Data = map[string]any{"message": "Error while deleting file", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func DoFileRequest(url, method string, headers map[string]string, body bytes.Buffer, writer *multipart.Writer) ([]byte, error) {
	request, err := http.NewRequest(method, url, &body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)

	return respByte, err
}
