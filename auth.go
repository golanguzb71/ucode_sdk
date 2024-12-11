package ucodesdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (u *object) Auth() AuthI {
	return &APIAuth{
		config: u.config,
	}
}

type AuthI interface {
	/*
		Register is a function that registers a new user with the provided data.

		Works for [Mongo, Postgres]

		sdk.Auth().
			Register(data).
			Exec()

		Use this method to create new users with basic or custom fields for authentication.
	*/
	Register(data map[string]any) *Register
	/*
		ResetPassword is a function that resets a user's password with the provided data.

		Works for [Mongo, Postgres]

		sdk.Auth().
			ResetPassword(data).
			Exec()

		This method initiates a password reset process, often requiring additional validation
		such as email or phone verification before allowing the reset.
	*/
	ResetPassword(data map[string]any) *ResetPassword
	Login(body map[string]any) *Login
	SendCode(data map[string]any) *SendCode
}

func (a *APIAuth) Register(data map[string]any) *Register {
	return &Register{
		config: a.config,
		data:   AuthRequest{Body: data},
	}
}

func (a *Register) Headers(headers map[string]string) *Register {
	a.data.Headers = headers
	return a
}

func (a *Register) Exec() (RegisterResponse, Response, error) {
	var (
		response = Response{
			Status: "done",
		}
		registerObject RegisterResponse
		url            = fmt.Sprintf("%s/v2/register?project-id=%s", a.config.BaseAuthUrl, a.config.ProjectId)
	)

	registerResponseInByte, err := DoRequest(url, http.MethodPost, a.data.Body, a.data.Headers)
	if err != nil {
		response.Data = map[string]any{"description": string(registerResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return RegisterResponse{}, response, err
	}

	err = json.Unmarshal(registerResponseInByte, &registerObject)
	if err != nil {
		response.Data = map[string]any{"description": string(registerResponseInByte), "message": "Error while unmarshalling register object", "error": err.Error()}
		response.Status = "error"
		return RegisterResponse{}, response, err
	}

	return registerObject, response, nil
}

func (a *APIAuth) ResetPassword(data map[string]any) *ResetPassword {
	return &ResetPassword{
		config: a.config,
		data:   AuthRequest{Body: data},
	}
}

func (a *ResetPassword) Headers(headers map[string]string) *ResetPassword {
	a.data.Headers = headers
	return a
}

func (a *ResetPassword) Exec() (Response, error) {
	var (
		response = Response{Status: "done"}
		url      = fmt.Sprintf("%s/v2/reset-password", a.config.BaseAuthUrl)
	)

	var appId = a.config.AppId

	header := map[string]string{
		"authorization": "API-KEY",
		"X-API-KEY":     appId,
	}

	_, err := DoRequest(url, http.MethodPut, a.data.Body, header)
	if err != nil {
		response.Data = map[string]any{"message": "Error while reset password", "error": err.Error()}
		response.Status = "error"
		return response, err
	}

	return response, nil
}

func (a *APIAuth) Login(data map[string]any) *Login {
	return &Login{
		config: a.config,
		data:   AuthRequest{Body: data},
	}
}

func (a *Login) Headers(headers map[string]string) *Login {
	a.data.Headers = headers
	return a
}

func (a *Login) Exec() (LoginResponse, Response, error) {
	var (
		response    = Response{Status: "done"}
		loginObject LoginResponse
		url         = fmt.Sprintf("%s/v2/login", a.config.BaseAuthUrl)
	)

	if a.data.Body["project_id"] == nil {
		a.data.Body["project_id"] = a.config.ProjectId
	}

	loginResponseInByte, err := DoRequest(url, http.MethodPost, a.data.Body, a.data.Headers)
	if err != nil {
		response.Data = map[string]any{"description": string(loginResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return LoginResponse{}, response, err
	}

	err = json.Unmarshal(loginResponseInByte, &loginObject)
	if err != nil {
		response.Data = map[string]any{"description": string(loginResponseInByte), "message": "Error while unmarshalling login object", "error": err.Error()}
		response.Status = "error"
		return LoginResponse{}, response, err
	}

	return loginObject, response, nil
}

func (a *Login) ExecWithOption() (LoginWithOptionResponse, Response, error) {
	var (
		response    = Response{Status: "done"}
		loginObject LoginWithOptionResponse
		url         = fmt.Sprintf("%s/v2/login/with-option?project-id=%s", a.config.BaseAuthUrl, a.config.ProjectId)
	)

	loginResponseInByte, err := DoRequest(url, http.MethodPost, a.data.Body, a.data.Headers)
	if err != nil {
		response.Data = map[string]any{"description": string(loginResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return LoginWithOptionResponse{}, response, err
	}

	err = json.Unmarshal(loginResponseInByte, &loginObject)
	if err != nil {
		response.Data = map[string]any{"description": string(loginResponseInByte), "message": "Error while unmarshalling login with option object", "error": err.Error()}
		response.Status = "error"
		return LoginWithOptionResponse{}, response, err
	}

	return loginObject, response, nil
}

func (a *APIAuth) SendCode(data map[string]any) *SendCode {
	return &SendCode{
		config: a.config,
		data:   AuthRequest{Body: data},
	}
}

func (a *SendCode) Headers(headers map[string]string) *SendCode {
	a.data.Headers = headers
	return a
}

func (a *SendCode) Exec() (SendCodeResponse, Response, error) {
	var (
		response   = Response{Status: "done"}
		codeObject SendCodeResponse
		url        = fmt.Sprintf("%s/v2/send-code", a.config.BaseAuthUrl)
	)

	codeResponseInByte, err := DoRequest(url, http.MethodPost, a.data.Body, a.data.Headers)
	if err != nil {
		response.Data = map[string]any{"description": string(codeResponseInByte), "message": "Can't send request", "error": err.Error()}
		response.Status = "error"
		return SendCodeResponse{}, response, err
	}

	err = json.Unmarshal(codeResponseInByte, &codeObject)
	if err != nil {
		response.Data = map[string]any{"description": string(codeResponseInByte), "message": "Error while unmarshalling send code object", "error": err.Error()}
		response.Status = "error"
		return SendCodeResponse{}, response, err
	}

	return codeObject, response, nil
}
