package ucodesdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (u *object) Function(path string) FunctionI {
	return &APIFunction{
		config: u.config,
		path:   path,
	}
}

// Function interface defines methods for invoking functions
type FunctionI interface {
	Invoke(data map[string]any) *APIFunction
}

// APIFunction struct implements FunctionInterface

func (f *APIFunction) Invoke(data map[string]any) *APIFunction {
	return &APIFunction{
		config:  f.config,
		request: Request{Data: data},
	}
}

func (f *APIFunction) Exec() (FunctionResponse, Response, error) {
	var (
		response     = Response{Status: "done"}
		invokeObject FunctionResponse
		url          = fmt.Sprintf("%s/v1/invoke_function/%s", f.config.BaseURL, f.path)
	)

	var appId = f.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	invokeFunctionResponseInByte, err := DoRequest(url, http.MethodPost, f.request, header)
	if err != nil {
		response.Data = map[string]any{"description": string(invokeFunctionResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return FunctionResponse{}, response, err
	}

	err = json.Unmarshal(invokeFunctionResponseInByte, &invokeObject)
	if err != nil {
		response.Data = map[string]any{"description": string(invokeFunctionResponseInByte), "message": "Error while unmarshalling invoke function", "error": err.Error()}
		response.Status = "error"
		return FunctionResponse{}, response, err
	}

	return invokeObject, response, nil
}
