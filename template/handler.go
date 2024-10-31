package function

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cast"
	sdk "github.com/ucode-io/ucode_sdk"
)

var (
	baseUrl = "https://api.client.u-code.io"
)

/*
Answer below questions before starting the function.

When the function invoked?
  - table_slug -> AFTER | BEFORE | HTTP -> CREATE | UPDATE | MULTIPLE_UPDATE | DELETE | APPEND_MANY2MANY | DELETE_MANY2MANY

What does it do?
- Explain the purpose of the function.(O'zbekcha yozilsa ham bo'ladi.)
*/
// func main() {
// 	data := `{"data":{"app_id":"P-CgtoLQxIfoXuz081FuZCenSJbUSMCjOf","object_data":{"test_id":"41574168-4d2f-481a-8c6f-bc60be37e674"}}}`
// 	resp := Handle([]byte(data))
// 	fmt.Println(resp)
// }

// Handle a serverless request
func Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			request       sdk.Request
			response      sdk.Response
			errorResponse sdk.ResponseError
			returnError   = func(errorResponse sdk.ResponseError) string {
				response = sdk.Response{
					Status: "error",
					Data:   map[string]any{"message": errorResponse.ClientErrorMessage, "error": errorResponse.ErrorMessage, "description": errorResponse.Description},
				}
				marshaledResponse, _ := json.Marshal(response)
				return string(marshaledResponse)
			}
		)

		gg := sdk.New(&sdk.Config{
			BaseURL: baseUrl,
			AppId:   "P-bgh4cmZxaWTXWscpH6sUa9gGlsuvKyZO",
			// AuthBaseURL: authBaseURL,
			ProjectId: "f05fdd8d-f949-4999-9593-5686ac272993",
		})
		body := map[string]any{
			"title": fmt.Sprintf("%d", time.Now().Unix()),
		}

		createResp, _, err := gg.Items("order_abdurahmon").Create(body).DisableFaas(true).Exec()
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		marssss, _ := json.Marshal(createResp)
		fmt.Println("CREATE RESP: ", string(marssss))

		updateBody := map[string]any{
			"title": fmt.Sprintf("%d %s", time.Now().Unix(), "updated"),
			"guid":  createResp.Data.Data.Data["guid"],
		}

		updateResp, _, err := gg.Items("order_abdurahmon").Update(updateBody).DisableFaas(true).ExecSingle()
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}
		marssss, _ = json.Marshal(updateResp)
		fmt.Println("UPDATE RESP: ", string(marssss))

		_, err = gg.Items("order_abdurahmon").Delete().Single(cast.ToString(createResp.Data.Data.Data["guid"])).DisableFaas(true).Exec()
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		getListResp, _, err := gg.Items("order_abdurahmon").
			GetList().
			Page(1).
			Limit(20).
			Sort(map[string]any{"title": 1}).
			Exec()
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		marssss, _ = json.Marshal(getListResp)
		fmt.Println("GETLIST RESP: ", string(marssss))

		fileResp, _, err := gg.Files().Upload("models.go").Exec()
		if err != nil {
			fmt.Println("ERROR: ", err)
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}
		marssss, _ = json.Marshal(fileResp)
		fmt.Println("FILE RESP: ", string(marssss))

		// ucodeApi := sdk.NewSDK(&sdk.Config{
		// 	BaseURL: baseUrl,
		// 	AppId:   "P-kL7M9h0NarpDfsSTzBPhDGOE4H9rUPl5",
		// })

		// datalens1-new-template-nats-publisher
		// faasResp, _, err := ucodeApi.Function("datalens1-new-template-nats-publisher").Invoke(map[string]any{}).Exec()
		// if err != nil {
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }
		// marssss, _ = json.Marshal(faasResp)
		// fmt.Println("FAAS RESP: ", string(marssss))

		heders := map[string]string{
			"Resource-Id":    "b74a3b18-6531-45fc-8e05-0b9709af8faa",
			"Environment-Id": "e8b82a93-b87f-4103-abc4-b5a017f540a4",
		}

		// body = map[string]any{
		// 	"data": map[string]any{
		// 		"type":           "phone",
		// 		"name":           fmt.Sprintf("%s %d", "otashjkee", time.Now().Unix()),
		// 		"phone":          "+998490000010",
		// 		"client_type_id": "1ade0441-4798-4183-839f-40a71e3dcad8",
		// 		"role_id":        "b4112b2b-82db-4942-9122-f3f8c58db34a",
		// 	},
		// }

		// registerResp, _, err := gg.Auth().Register(body).Headers(heders).Exec()
		// if err != nil {
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }

		// marssss, _ = json.Marshal(registerResp)
		// fmt.Println("REGISTER RESP: ", string(marssss))

		body = map[string]any{
			"data": map[string]any{
				"sms_id":         "",
				"phone":          "+998490000010",
				"otp":            "111111",
				"role_id":        "b4112b2b-82db-4942-9122-f3f8c58db34a",
				"client_type_id": "1ade0441-4798-4183-839f-40a71e3dcad8",
			},
			"login_strategy": "PHONE_OTP",
		}

		loginWithOptionResp, _, err := gg.Auth().Login(body).Headers(heders).ExecWithOption()
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		marssss, _ = json.Marshal(loginWithOptionResp)
		fmt.Println("LOGIN WITH OPTION RESP: ", string(marssss))

		body = map[string]any{
			"username":    "integrationtestgo",
			"password":    "integrationtestgo",
			"project_id":  "f05fdd8d-f949-4999-9593-5686ac272993",
			"client_type": "10debeef-b5b9-415d-bfe8-dbd8646e2fd4",
		}

		loginResp, _, err := gg.Auth().Login(body).Headers(heders).Exec()
		if err != nil {
			fmt.Println("login error: ", err)
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		marssss, _ = json.Marshal(loginResp)
		fmt.Println("LOGIN RESP: ", string(marssss))

		send := sdk.New(&sdk.Config{
			BaseURL:   baseUrl,
			ProjectId: "0f111e78-3a93-4bec-945a-2a77e0e0a82d",
		})

		body = map[string]any{
			"recipient": "+998998136254",
			"text":      "code",
			"type":      "PHONE",
		}
		heders = map[string]string{
			"Resource-Id":    "491a431c-b6fe-4882-a7e4-9894f564835a",
			"Environment-Id": "2f7e62ee-3fba-4092-8a16-3d8e587e993d",
		}
		sendResp, _, err := send.Auth().SendCode(body).Headers(heders).Exec()
		if err != nil {
			fmt.Println("send ocde error: ", err)
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		marssss, _ = json.Marshal(sendResp)
		fmt.Println("CODE RESP: ", string(marssss))

		requestByte, err := io.ReadAll(r.Body)
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestByte, &request)
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on unmarshal request"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusInternalServerError)
			return
		}

		response.Status = "done"
		handleResponse(w, response, http.StatusOK)
	}
}

func handleResponse(w http.ResponseWriter, body any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	bodyByte, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`
			{
				"error": "Error marshalling response"
			}
		`))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(bodyByte)
}
