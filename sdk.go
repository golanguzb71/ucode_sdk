package ucodesdk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type UcodeApis interface {
	/*
		Items returns an interface to interact with items within a specified collection.

		Items are objects within a Collection which contain values for one or more fields.
		Each item represents a record in your database, allowing CRUD operations.

		Usage:
		sdk.Items("collection_name").
			Create(data).
			Exec()

		This enables you to manage items in collections across databases, such as MongoDB and PostgreSQL.
	*/
	Items(collection string) ItemsI
	/*
		Auth returns an interface for handling user authentication and authorization operations.

		Use this interface to manage user registration, login, password resets, and other
		authentication-related tasks.

		Usage:
		sdk.Auth().
			Register(data).
			Exec()

		Supports various authentication workflows compatible with both MongoDB and PostgreSQL.
	*/
	Auth() AuthI
	/*
		Files returns an interface for file management operations.

		Use this interface to upload, delete, or manage files stored on the server, allowing
		for easy integration of file-based data alongside other operations.

		Usage:
		sdk.Files().
			Upload("file_path").
			Exec()

		Designed for compatibility with both MongoDB and PostgreSQL for consistent file management.
	*/
	Files() FilesI
	/*
		Function returns an interface for invoking server-side functions.

		This interface enables the execution of predefined or custom server functions,
		facilitating complex data processing and automation workflows.

		Usage:
		sdk.Function("function_path").
			Invoke(data).
			Exec()

		Supported across MongoDB and PostgreSQL, providing flexibility for backend processing.
	*/
	Function(path string) FunctionI
	Config() *Config
	DoRequest(url string, method string, body any, headers map[string]string) ([]byte, error)
}

func New(cfg *Config) UcodeApis {
	return &object{
		config: cfg,
	}
}

// UcodeAPI struct implements UcodeAPIInterface
type object struct {
	config *Config
}

func (u *object) Config() *Config {
	return u.config
}

func DoRequest(url string, method string, body any, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Add headers from the map
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)

	return respByte, err
}

func (a *object) DoRequest(url string, method string, body any, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(&body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Add headers from the map
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	return respByte, err
}
