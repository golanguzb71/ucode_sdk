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
	baseUrl     = "https://api.admin.u-code.io"
	authBaseURL = "https://auth-api.ucode.run"
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

		gg := sdk.NewSDK(&sdk.Config{
			BaseURL:     baseUrl,
			AppId:       "P-bgh4cmZxaWTXWscpH6sUa9gGlsuvKyZO",
			AuthBaseURL: authBaseURL,
			ProjectId:   "462baeca-37b0-4355-addc-b8ae5d26995d",
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

		// createResp, response, err := gg.Items("order_abdurahmon").Create(body).Exec()
		// if err != nil {
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }
		// marssss, _ := json.Marshal(createResp)
		// fmt.Println("CREATE RESP: ",string(marssss))

		// getListResp, _, err := gg.Items("order_abdurahmon").
		// 	GetList().
		// 	Page(1).
		// 	Limit(20).
		// 	Sort(map[string]any{"title": 1}).
		// 	Exec()
		// if err != nil {
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }

		// getListSlimResp, _, err := gg.Items("order_abdurahmon").
		// 	GetListSlim().Page(1).Limit(20).WithRelations(true).Exec()
		// if err != nil {
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }
		// ressss, _ := json.Marshal(getListSlimResp.Data.Data.Response)
		// fmt.Println("LENGTH: ", string(ressss))

		// getSingleSlimResp, _, err := gg.Items("order_abdurahmon").
		// 	GetSingle("20200c76-deb4-4646-a754-af4695857243").
		// 	ExecSlim()
		// if err != nil {
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }
		// ressss, _ := json.Marshal(getSingleSlimResp.Data.Data.Response)
		// fmt.Println("LENGTH: ", string(ressss))

		// set timeout for request

		body = map[string]any{
			"data": map[string]any{
				"type":           "phone",
				"Name":           fmt.Sprintf("%s %d", "otashjkee", time.Now().Unix()),
				"phone":          "+967000000001",
				"client_type_id": "1d75cd99-577d-4d84-8d08-c4f87507a452",
				"role_id":        "eba0211b-bb79-4c92-ba49-4ffcb1c9caac",
			},
		}
		heders := map[string]string{
			"Resource-Id":    "05df5e41-1066-474e-8435-3781e0841603",
			"Environment-Id": "ad41c493-8697-4f23-979a-341722465748",
		}
		registerResp, _, err := gg.Auth().Register(body).Headers(heders).Exec()
		if err != nil {
			errorResponse.ClientErrorMessage = "Error on getting request body"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
			return
		}

		marssss, _ = json.Marshal(registerResp)
		fmt.Println("REGISTER RESP: ", string(marssss))

		// fileResp, _, err := gg.Files().Upload("models.go").Exec()
		// if err != nil {
		// 	fmt.Println("ERROR: ", err)
		// 	errorResponse.ClientErrorMessage = "Error on getting request body"
		// 	errorResponse.ErrorMessage = err.Error()
		// 	errorResponse.StatusCode = http.StatusInternalServerError
		// 	handleResponse(w, returnError(errorResponse), http.StatusBadRequest)
		// 	return
		// }
		// marssss, _ := json.Marshal(fileResp)
		// fmt.Println("FILE RESP: ", string(marssss))

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
