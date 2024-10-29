package ucodesdk

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

var (
	baseUrl = "https://api.admin.u-code.io"
)

func TestEndToEnd(t *testing.T) {
	var (
		response      Response
		errorResponse ResponseError

		returnError = func(errorResponse ResponseError) string {
			response = Response{
				Status: "error",
				Data:   map[string]interface{}{"message": errorResponse.ClientErrorMessage, "error": errorResponse.ErrorMessage, "description": errorResponse.Description},
			}
			marshaledResponse, _ := json.Marshal(response)
			return string(marshaledResponse)
		}
		housesMongo    []map[string]interface{}
		housesPostgres []map[string]interface{}
		roomsPostgres  []map[string]interface{}
		roomsMongo     []map[string]interface{}
		mongoAppId     string
		postgresAppId  string
		roomsCount     = 4
		housesCount    = 2
	)

	// // check DoRequest method
	t.Run("TestDoRequest", func(t *testing.T) {
		ucodeApi := NewSDK(&Config{BaseURL: baseUrl})

		header := map[string]string{
			"authorization": "API-KEY",
			"X-API-KEY":     "test_app_id",
		}

		// Test successful request
		_, err := ucodeApi.DoRequest(baseUrl+"/test", "GET", nil, header)
		if err != nil {
			t.Errorf("Error on DoRequest: %v", err)
			return
		}

		// Test with invalid URL
		_, err = ucodeApi.DoRequest("invalid-url", "GET", nil, header)
		if err == nil {
			t.Error("Expected error for invalid URL, got nil")
			return
		}

		// Test with invalid request body
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, err = ucodeApi.DoRequest(baseUrl+"/test", "GET", MyStruct{}, header)
		if err == nil {
			t.Error("Expected error for invalid reuest body, got nil")
			return
		}

		// Test with invalid character in url
		_, err = ucodeApi.DoRequest("https://example.com/path\x7F", "", nil, header)
		if err == nil {
			t.Error("Expected error with invalid character in url, got nil")
			return
		}

		// Test with custom headers
		customHeaders := map[string]string{
			"Custom-Header": "TestValue",
		}
		_, err = ucodeApi.DoRequest(baseUrl+"/test", "GET", nil, customHeaders)
		if err != nil {
			t.Errorf("Error on DoRequest with custom headers: %v", err)
			return
		}
	})

	// getting app_id for mongodb and postgres
	t.Run("getAppId", func(t *testing.T) {
		err := godotenv.Load()
		if err != nil {
			t.Error("error loading .env file")
			return
		}

		mongoAppId = os.Getenv("MONGO_APP_ID")
		if mongoAppId == "" {
			t.Error("Error on setting MONGO_APP_ID from .env file")
			return
		}

		postgresAppId = os.Getenv("POSTGRES_APP_ID")
		if postgresAppId == "" {
			t.Error("Error on setting POSTGRES_APP_ID from .env file")
			return
		}
	})

	var (
		ucodeApi   = NewSDK(&Config{BaseURL: baseUrl, AppId: mongoAppId})
		ucodeApiPg = NewSDK(&Config{BaseURL: baseUrl, AppId: postgresAppId})
	)

	t.Run("createInMongo", func(t *testing.T) {
		// --------------------------CreateObject------------------------------
		// create houses
		createHousesRequest := map[string]interface{}{
			"name":       "house",
			"price":      15000,
			"room_count": 5,
		}

		for i := 0; i < housesCount; i++ {
			_, response, err := ucodeApi.Items("houses").Create(createHousesRequest).Exec()
			if err != nil {
				errorResponse.Description = response.Data["description"]
				errorResponse.ClientErrorMessage = "error on creating new hourse"
				errorResponse.ErrorMessage = err.Error()
				errorResponse.StatusCode = http.StatusInternalServerError
				t.Error(returnError(errorResponse))
				return
			}

			if response.Status != "done" {
				t.Error(response.Status, response.Data, response.Error)
			}
		}

		// check error case
		_, _, err := ucodeApi.Items("houses").Create(nil).Exec()
		if err == nil {
			t.Error("error: request not given but work")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}

		_, response, err = ucodeApi.Items("houses").Create(map[string]any{"guid": MyStruct{}}).Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("gethousesInMongo", func(t *testing.T) {
		// --------------------------GetList------------------------------
		// getting houses
		ExistObject, response, err := ucodeApi.Items("houses").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on useing GetList method"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
		housesMongo = ExistObject.Data.Data.Response

		if len(housesMongo) != housesCount {
			t.Errorf("error on created houses count = %d\nExpected = %d", len(housesMongo), housesCount)
		}

		// check data created fully
		for _, house := range housesMongo {
			assert.Equal(t, "house", cast.ToString(house["name"]))
			assert.Equal(t, 15000, cast.ToInt(house["price"]))
			assert.Equal(t, 5, cast.ToInt(house["room_count"]))
		}

		// Test with invalid parameters
		_, _, err = ucodeApi.Items("invalid_table").GetList().Page(-1).Limit(-1).Exec()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}

		_, _, err = ucodeApi.Items("houses").
			GetList().
			Filter(map[string]any{"guid": MyStruct{}}).
			Page(1).
			Limit(10).
			Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("createInPostgres", func(t *testing.T) {
		// --------------------------CreateObject------------------------------
		// create houses
		createHousesRequest := map[string]interface{}{
			"name":       "house",
			"price":      15000,
			"room_count": 5,
		}

		for i := 0; i < housesCount; i++ {
			_, response, err := ucodeApiPg.Items("houses").Create(createHousesRequest).Exec()
			if err != nil {
				errorResponse.Description = response.Data["description"]
				errorResponse.ClientErrorMessage = "error on creating new hourse"
				errorResponse.ErrorMessage = err.Error()
				errorResponse.StatusCode = http.StatusInternalServerError
				t.Error(returnError(errorResponse))
				return
			}

			if response.Status != "done" {
				t.Error(response.Status, response.Data, response.Error)
			}
		}

		// check error case
		_, _, err := ucodeApiPg.Items("houses").Create(nil).Exec()
		if err == nil {
			t.Error("error: request not given but work")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApiPg.Items("houses").Create(map[string]interface{}{"guid": MyStruct{}}).Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("gethousesInPostgres", func(t *testing.T) {
		// --------------------------GetList------------------------------
		// getting houses
		ExistObject, response, err := ucodeApiPg.Items("houses").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on useing GetList method"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
		housesPostgres = ExistObject.Data.Data.Response

		if len(housesPostgres) != housesCount {
			t.Errorf("error on created houses count = %d\nExpected = %d", len(housesPostgres), housesCount)
		}

		// check data created fully
		for _, house := range housesPostgres {
			assert.Equal(t, "house", cast.ToString(house["name"]))
			assert.Equal(t, 15000, cast.ToInt(house["price"]))
			assert.Equal(t, 5, cast.ToInt(house["room_count"]))
		}

		// Test with invalid parameters
		_, _, err = ucodeApiPg.Items("invalid_table").GetList().Page(-1).Limit(-1).Exec()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApiPg.Items("houses").
			GetList().
			Filter(map[string]any{"guid": MyStruct{}}).
			Page(1).
			Limit(10).
			Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("createRoomInMongo", func(t *testing.T) {
		// create rooms
		createRoomRequest := map[string]interface{}{
			"name": "room",
		}

		for i := 0; i < roomsCount; i++ {
			_, response, err := ucodeApi.Items("room").Create(createRoomRequest).Exec()
			if err != nil {
				errorResponse.Description = response.Data["description"]
				errorResponse.ClientErrorMessage = "error on creating new hourse"
				errorResponse.ErrorMessage = err.Error()
				errorResponse.StatusCode = http.StatusInternalServerError
				t.Error(returnError(errorResponse))
				return
			}

			if response.Status != "done" {
				t.Error(response.Status, response.Data, response.Error)
			}
		}
	})

	t.Run("getRoomsInMongo", func(t *testing.T) {
		// --------------------------GetListSlim------------------------------
		getListSlim, response, err := ucodeApi.Items("room").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on GetListSlim"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
		roomsMongo = getListSlim.Data.Data.Response

		if len(roomsMongo) != roomsCount {
			t.Errorf("error on created rooms count = %d\nExpected = %d", len(roomsMongo), roomsCount)
		}
		// check data is currect created
		for _, room := range roomsMongo {
			assert.Equal(t, "room", cast.ToString(room["name"]))
		}

		// Test with invalid Request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}

		_, _, err = ucodeApi.Items("houses").
			GetList().
			Page(1).Limit(10).
			Filter(map[string]any{"guid": MyStruct{}}).
			Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("createRoomInPostgres", func(t *testing.T) {
		// create rooms
		createRoomRequest := map[string]interface{}{
			"name": "room",
		}

		for i := 0; i < roomsCount; i++ {
			_, response, err := ucodeApiPg.Items("room").Create(createRoomRequest).Exec()
			if err != nil {
				errorResponse.Description = response.Data["description"]
				errorResponse.ClientErrorMessage = "error on creating new hourse"
				errorResponse.ErrorMessage = err.Error()
				errorResponse.StatusCode = http.StatusInternalServerError
				t.Error(returnError(errorResponse))
				return
			}

			if response.Status != "done" {
				t.Error(response.Status, response.Data, response.Error)
			}
		}
	})

	t.Run("getRoomsInPostgres", func(t *testing.T) {
		// --------------------------GetListSlim------------------------------
		getListSlim, response, err := ucodeApiPg.Items("room").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on GetListSlim"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
		roomsPostgres = getListSlim.Data.Data.Response

		if len(roomsPostgres) != roomsCount {
			t.Errorf("error on created rooms count = %d\nExpected = %d", len(roomsPostgres), roomsCount)
		}
		// check data is currect created
		for _, room := range roomsPostgres {
			assert.Equal(t, "room", cast.ToString(room["name"]))
		}

		// Test with invalid Request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApi.Items("houses").GetList().Filter(map[string]any{"guid": MyStruct{}}).Page(1).Limit(10).Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("updateHousesInMongo", func(t *testing.T) {
		// --------------------------UpdateObject------------------------------
		// update first house
		if len(housesMongo) < housesCount {
			t.Errorf("error houses count = %d\nExpected count = %d", len(housesMongo), housesCount)
			return
		}

		updateReq := map[string]interface{}{
			"guid":       housesMongo[0]["guid"],
			"room_count": 10,
		}
		resp, response, err := ucodeApi.Items("houses").Update(updateReq).ExecSingle()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "error on UpdateObject"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// check data updated successfully
		assert.Equal(t, 10, cast.ToInt(resp.Data.Data["room_count"]))
		assert.Equal(t, housesMongo[0]["guid"], cast.ToString(resp.Data.Data["guid"]))

		// Test with invalid parameters
		_, _, err = ucodeApi.Items("invalid_table").Update(map[string]interface{}{"guid": "invalid_guid"}).ExecSingle()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApi.Items("invalid_table").Update(map[string]interface{}{"guid": MyStruct{}}).ExecSingle()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("updateHousesInPostgres", func(t *testing.T) {
		// --------------------------UpdateObject------------------------------
		// update first house
		if len(housesPostgres) < housesCount {
			t.Errorf("error houses count = %d\nExpected count = %d", len(housesPostgres), housesCount)
			return
		}
		updateReq := map[string]interface{}{
			"guid":       housesPostgres[0]["guid"],
			"room_count": 10,
		}
		_, response, err := ucodeApiPg.Items("houses").Update(updateReq).ExecSingle()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "error on UpdateObject"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		t.Run("GetSingleInPostgres", func(t *testing.T) {
			// --------------------------GetSingle------------------------------
			// get the house info
			houseInfoo, response, err := ucodeApiPg.Items("houses").GetSingle(cast.ToString(housesPostgres[0]["guid"])).Exec()
			if err != nil {
				errorResponse.Description = response.Data["description"]
				errorResponse.ClientErrorMessage = "Error on getting single"
				errorResponse.ErrorMessage = err.Error()
				errorResponse.StatusCode = http.StatusInternalServerError
				t.Error(returnError(errorResponse))
				return
			}

			if response.Status != "done" {
				t.Error(response.Status, response.Data, response.Error)
			}

			if len(houseInfoo.Data.Data.Response) == 0 {
				t.Errorf("error GetSingle method not return data")
			}

			// check data updated successfully
			assert.Equal(t, 10, cast.ToInt(houseInfoo.Data.Data.Response["room_count"]))
			assert.Equal(t, housesPostgres[0]["guid"], cast.ToString(houseInfoo.Data.Data.Response["guid"]))

			// Test with invalid Request parameters
			_, _, err = ucodeApiPg.Items("invalid_table").GetSingle("").Exec()
			if err == nil {
				t.Error("error: invalid request given but work")
				return
			}
		})

		// Test with invalid parameters
		_, _, err = ucodeApiPg.Items("invalid_table").Update(map[string]interface{}{"guid": "invalid_guid"}).ExecSingle()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApiPg.Items("invalid_table").Update(map[string]interface{}{"guid": MyStruct{}}).ExecSingle()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("GetSingleInMongo", func(t *testing.T) {
		// Test with invalid parameters
		_, _, err := ucodeApi.Items("invalid_table").GetSingle("").Exec()
		if err == nil {
			t.Error("error: invalid parameters given but work")
			return
		}
	})

	t.Run("MultipleUpdate in mongo", func(t *testing.T) {
		// --------------------------MultipleUpdate------------------------------
		var (
			multipleUpdateRequest = []map[string]interface{}{}
			ids                   = make([]string, len(housesMongo))
		)

		for i, house := range housesMongo {
			ids[i] = cast.ToString(house["guid"])
			multipleUpdateRequest = append(multipleUpdateRequest, map[string]interface{}{
				"guid":       cast.ToString(house["guid"]),
				"room_count": 15,
			})
		}

		_, response, err := ucodeApi.Items("houses").Update(map[string]any{"objects": multipleUpdateRequest}).ExecMultiple()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on MultipleUpdate"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// --------------------------GetListSlim------------------------------
		getListSlim, response, err := ucodeApi.Items("houses").
			GetList().
			Page(1).
			Limit(100000).
			Filter(map[string]any{"ids": ids}).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on GetListSlim"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(getListSlim.Data.Data.Response) != housesCount {
			t.Errorf("error on houses count = %d\nExpected = %d", len(getListSlim.Data.Data.Response), housesCount)
		}

		// check data is currect created
		for _, house := range getListSlim.Data.Data.Response {
			assert.Equal(t, 15, cast.ToInt(house["room_count"]))
		}

		// Test with invalid parameters
		_, _, err = ucodeApi.Items("").Update(map[string]any{"objects": []map[string]interface{}{}}).ExecMultiple()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApi.Items("houses").Update(map[string]any{"objects": MyStruct{}}).ExecMultiple()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("MultipleUpdate in postgres", func(t *testing.T) {
		// --------------------------MultipleUpdate------------------------------
		var (
			multipleUpdateRequest = []map[string]interface{}{}
			ids                   = make([]string, len(housesPostgres))
		)

		for i, house := range housesPostgres {
			ids[i] = cast.ToString(house["guid"])
			multipleUpdateRequest = append(multipleUpdateRequest, map[string]interface{}{
				"guid":       cast.ToString(house["guid"]),
				"room_count": 15,
			})
		}

		_, response, err := ucodeApiPg.Items("houses").Update(map[string]any{"objects": multipleUpdateRequest}).ExecMultiple()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on MultipleUpdate"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// --------------------------GetListSlim------------------------------
		getListSlim, response, err := ucodeApiPg.Items("houses").
			GetList().
			Page(1).
			Filter(map[string]any{"ids": ids}).
			Limit(100000).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on GetListSlim"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(getListSlim.Data.Data.Response) != housesCount {
			t.Errorf("error on houses count = %d\nExpected = %d", len(getListSlim.Data.Data.Response), housesCount)
		}

		// check data is currect created
		for _, house := range getListSlim.Data.Data.Response {
			assert.Equal(t, 15, cast.ToInt(house["room_count"]))
		}

		// Test with invalid parameters
		_, _, err = ucodeApiPg.Items("").Update(map[string]any{"objects": []map[string]interface{}{}}).ExecMultiple()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, _, err = ucodeApiPg.Items("houses").Update(map[string]any{"objects": MyStruct{}}).ExecMultiple()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("GetListAggregation in mongo", func(t *testing.T) {
		// --------------------------GetListAggregation FOR MongoDB------------------------------
		getListAggregationPipeline := map[string]any{
			"pipelines": []map[string]any{{
				"$match": map[string]any{
					"price": map[string]any{
						"$exists": true,
						"$eq":     15000,
					},
				},
			},
			},
		}
		getListAggregationList, response, err := ucodeApi.Items("houses").
			GetList().
			Pipelines(getListAggregationPipeline).
			ExecAggregation()

		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "error on GetListAggregation"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(getListAggregationList.Data.Data.Data) != housesCount {
			t.Errorf("error on took houses count = %d\nExpected = %d", len(getListAggregationList.Data.Data.Data), housesCount)
		}

		// Test with invalid parameters
		_, _, err = ucodeApi.Items("houses").
			GetList().
			Pipelines(map[string]any{}).
			ExecAggregation()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}

		// Test with invalid request parameters
		type MyStruct struct {
			A int
			B func() // functions are not supported
		}
		_, response, err = ucodeApi.Items("houses").
			GetList().
			Pipelines(map[string]any{"pipelines": MyStruct{}}).
			ExecAggregation()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("GetSingleSlim in mongo", func(t *testing.T) {
		// --------------------------GetSingleSlim------------------------------
		var id = cast.ToString(roomsMongo[0]["guid"])
		courseResponse, response, err := ucodeApi.Items("room").GetSingle(id).ExecSlim()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on get-single course"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// check values are correct
		assert.Equal(t, "room", courseResponse.Data.Data.Response["name"])

		// Test with invalid parameters
		_, _, err = ucodeApi.Items("houses").GetSingle("invalid_guid").ExecSlim()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}
	})

	t.Run("GetSingleSlim in postgres", func(t *testing.T) {
		// --------------------------GetSingleSlim------------------------------
		var id = cast.ToString(roomsPostgres[0]["guid"])

		courseResponse, response, err := ucodeApiPg.Items("room").GetSingle(id).ExecSlim()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on get-single course"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// check values are correct
		assert.Equal(t, "room", courseResponse.Data.Data.Response["name"])

		// Test with invalid parameters
		_, _, err = ucodeApiPg.Items("houses").GetSingle("invalid_guid").ExecSlim()
		if err == nil {
			t.Error("Expected error for invalid parameters, got nil")
			return
		}
	})

	t.Run("Delete in mongo", func(t *testing.T) {
		// --------------------------Delete------------------------------
		response, err := ucodeApi.Items("houses").
			Delete().
			DisableFaas(true).
			Single(cast.ToString(housesMongo[0]["guid"])).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while Delete"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
	})

	t.Run("Delete in postgres", func(t *testing.T) {
		// --------------------------Delete------------------------------
		response, err := ucodeApiPg.Items("houses").
			Delete().DisableFaas(true).
			Single(cast.ToString(housesPostgres[0]["guid"])).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while Delete"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}
		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
	})

	t.Run("MultipleDelete in mongo", func(t *testing.T) {
		// --------------------------MultipleDelete------------------------------
		// deleting from houses
		var (
			idMultipleDeleteHouses = []string{}
		)
		for _, val := range housesMongo {
			idMultipleDeleteHouses = append(idMultipleDeleteHouses, cast.ToString(val["guid"]))
		}

		response, err := ucodeApi.Items("houses").
			Delete().DisableFaas(true).
			Multiple(idMultipleDeleteHouses).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while Delete"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// Test with invalid request parameters
		response, err = ucodeApi.Items("houses").
			Delete().
			DisableFaas(true).
			Multiple(nil).
			Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("MultipleDelete in postgres", func(t *testing.T) {
		// --------------------------MultipleDelete------------------------------
		// deleting from houses
		var (
			idMultipleDeleteHouses = []string{}
		)
		for _, val := range housesPostgres {
			idMultipleDeleteHouses = append(idMultipleDeleteHouses, cast.ToString(val["guid"]))
		}
		response, err := ucodeApiPg.Items("houses").
			Delete().
			DisableFaas(true).
			Multiple(idMultipleDeleteHouses).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while Delete"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		// Test with invalid request parameters
		response, err = ucodeApiPg.Items("houses").
			Delete().
			DisableFaas(true).
			Multiple(nil).
			Exec()
		if err == nil {
			t.Error("error: invalid request given but work")
			return
		}
	})

	t.Run("MultipleDelete in mongo", func(t *testing.T) {
		// --------------------------MultipleDelete------------------------------
		// deleting from rooms
		var (
			idMultipleDeleteRoom = []string{}
		)
		for _, val := range roomsMongo {
			idMultipleDeleteRoom = append(idMultipleDeleteRoom, cast.ToString(val["guid"]))
		}

		response, err := ucodeApi.Items("room").
			Delete().
			DisableFaas(true).
			Multiple(idMultipleDeleteRoom).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while Delete"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
	})

	t.Run("MultipleDelete in postgres", func(t *testing.T) {
		// --------------------------MultipleDelete------------------------------
		// deleting from rooms
		var (
			idMultipleDeleteRoom = []string{}
		)
		for _, val := range roomsPostgres {
			idMultipleDeleteRoom = append(idMultipleDeleteRoom, cast.ToString(val["guid"]))
		}

		response, err := ucodeApiPg.Items("room").
			Delete().
			DisableFaas(true).
			Multiple(idMultipleDeleteRoom).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error while Delete"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}
		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}
	})

	t.Run("Checking all houses were deleted", func(t *testing.T) {
		// --------------------------GetList------------------------------
		// getting houses in mongo
		ExistObject, response, err := ucodeApi.Items("houses").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on useing GetList method"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(ExistObject.Data.Data.Response) != 0 {
			t.Errorf("error on not all houses deleted\nHave count = %d", len(ExistObject.Data.Data.Response))
		}

		// --------------------------GetList------------------------------
		// getting houses in postgres
		ExistObject, response, err = ucodeApiPg.Items("houses").
			GetList().
			Page(1).
			Limit(100000).
			Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on useing GetList method"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(ExistObject.Data.Data.Response) != 0 {
			t.Errorf("error on not all houses deleted\nHave count = %d", len(ExistObject.Data.Data.Response))
		}
	})

	t.Run("Checking all rooms were deleted", func(t *testing.T) {
		// --------------------------GetList------------------------------
		// getting houses in mongo
		ExistObject, response, err := ucodeApi.Items("room").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on useing GetList method"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(ExistObject.Data.Data.Response) != 0 {
			t.Errorf("error on not all rooms deleted\nHave count = %d", len(ExistObject.Data.Data.Response))
		}

		// --------------------------GetList------------------------------
		// getting houses in postgres
		ExistObject, response, err = ucodeApiPg.Items("room").GetList().Page(1).Limit(100000).Exec()
		if err != nil {
			errorResponse.Description = response.Data["description"]
			errorResponse.ClientErrorMessage = "Error on useing GetList method"
			errorResponse.ErrorMessage = err.Error()
			errorResponse.StatusCode = http.StatusInternalServerError
			t.Error(returnError(errorResponse))
			return
		}

		if response.Status != "done" {
			t.Error(response.Status, response.Data, response.Error)
		}

		if len(ExistObject.Data.Data.Response) != 0 {
			t.Errorf("error on not all room deleted\nHave count = %d", len(ExistObject.Data.Data.Response))
		}
	})
}
