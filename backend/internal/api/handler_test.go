package api

import (
	"backend/internal/api/mocks"
	"backend/internal/structure"
	"bou.ke/monkey"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func fixedTime() string {
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2009, 11, 10, 20, 34, 58, 651387237, time.UTC)
	})
	return time.Now().Format(time.RFC3339)
}

func TestHandler_AddImage(t *testing.T) {
	data := structure.Data{ID: "algID", Content: "image"}
	jsonData, err := json.Marshal(data)
	assert.NoError(t, err)
	img := structure.Images{Images: []string{"image"}}
	jsonImg, err := json.Marshal(img)
	images := structure.Images{Images: []string{"image", "image"}}
	jsonImages, err := json.Marshal(images)
	assert.NoError(t, err)
	tests := []struct {
		testName         string
		requestURL       string
		body             io.Reader
		contentType      string
		getReturned      string
		getError         error
		insertData       string
		insertError      error
		assertNoOfInsert int
		assertNoOfGet    int
		bodyContains     string
		statusCode       int
	}{
		{
			testName:     "should return 404 when url is wrong",
			requestURL:   "/v1/images/wrong",
			bodyContains: "404 not found",
			statusCode:   http.StatusNotFound,
		},
		{
			testName:     "should return 400 when Content-Type is incorrect",
			requestURL:   "/v1/images",
			contentType:  "application/wrong",
			bodyContains: "invalid content type",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:     "should return 400 error when request body is not json",
			requestURL:   "/v1/images",
			body:         bytes.NewBuffer([]byte("string")),
			contentType:  "application/json",
			bodyContains: "failed to unmarshal",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:      "should return 500 when database does not respond",
			requestURL:    "/v1/images",
			body:          bytes.NewBuffer(jsonData),
			contentType:   "application/json",
			getError:      errors.New("database not respond error"),
			assertNoOfGet: 1,
			bodyContains:  "database not respond error",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when there is no `images` in database",
			requestURL:    "/v1/images",
			body:          bytes.NewBuffer(jsonData),
			contentType:   "application/json",
			getError:      errors.New("key images does not exist"),
			assertNoOfGet: 1,
			bodyContains:  "key images does not exist",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when `images` have no value assigned",
			requestURL:    "/v1/images",
			body:          bytes.NewBuffer(jsonData),
			contentType:   "application/json",
			getReturned:   "",
			getError:      errors.New("for key images value is empty"),
			assertNoOfGet: 1,
			bodyContains:  "for key images value is empty",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when `images` are not in json format",
			requestURL:    "/v1/images",
			body:          bytes.NewBuffer(jsonData),
			contentType:   "application/json",
			getReturned:   "not json",
			assertNoOfGet: 1,
			bodyContains:  "failed to unmarshal",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:         "should return 500 when failed to insert `images` to database",
			requestURL:       "/v1/images",
			body:             bytes.NewBuffer(jsonData),
			contentType:      "application/json",
			getReturned:      string(jsonImg),
			assertNoOfGet:    1,
			insertData:       string(jsonImages),
			insertError:      errors.New("failed to insert json to database"),
			assertNoOfInsert: 1,
			bodyContains:     "failed to insert json to database",
			statusCode:       http.StatusInternalServerError,
		},
		{
			testName:         "should return 200 when images where take from database correctly",
			requestURL:       "/v1/images",
			body:             bytes.NewBuffer(jsonData),
			contentType:      "application/json",
			getReturned:      string(jsonImg),
			insertData:       string(jsonImages),
			assertNoOfGet:    1,
			assertNoOfInsert: 1,
			statusCode:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			r, err := http.NewRequest("PUT", tt.requestURL, tt.body)
			assert.NoError(t, err)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := []IAlgorithm{&mocks.IAlgorithm{}}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMock)
			iDatabaseMock.On("Get", "images").Return(tt.getReturned, tt.getError)
			iDatabaseMock.On("Set", "images", tt.insertData).Return(tt.insertError)

			//when
			testSubject.AddImage(w, r)

			//then
			iDatabaseMock.AssertNumberOfCalls(t, "Get", tt.assertNoOfGet)
			iDatabaseMock.AssertNumberOfCalls(t, "Set", tt.assertNoOfInsert)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandler_GetImages(t *testing.T) {
	img := structure.Images{Images: []string{"image", "image"}}
	jsonImg, err := json.Marshal(img)
	assert.NoError(t, err)
	tests := []struct {
		testName      string
		requestURL    string
		getReturned   string
		getError      error
		assertNoOfGet int
		bodyContains  string
		statusCode    int
	}{
		{
			testName:      "should return 404 when url is wrong",
			requestURL:    "/v1/images/wrong",
			assertNoOfGet: 0,
			bodyContains:  "404 not found",
			statusCode:    http.StatusNotFound,
		},
		{
			testName:      "should return 500 when database does not respond",
			requestURL:    "/v1/images",
			getError:      errors.New("database not respond error"),
			assertNoOfGet: 1,
			bodyContains:  "database not respond error",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when there is no `images` in database",
			requestURL:    "/v1/images",
			getError:      errors.New("key does not exist"),
			assertNoOfGet: 1,
			bodyContains:  "key does not exist",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when `images` are not in json format",
			requestURL:    "/v1/images",
			getReturned:   "not json",
			assertNoOfGet: 1,
			bodyContains:  "failed to unmarshal",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 200 when images where take from database correctly",
			requestURL:    "/v1/images",
			getReturned:   string(jsonImg),
			assertNoOfGet: 1,
			bodyContains:  string(jsonImg),
			statusCode:    http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			r, err := http.NewRequest("GET", tt.requestURL, nil)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := []IAlgorithm{&mocks.IAlgorithm{}}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMock)
			iDatabaseMock.On("Get", "images").Return(tt.getReturned, tt.getError)

			//when
			testSubject.GetImages(w, r)

			//then
			iDatabaseMock.AssertNumberOfCalls(t, "Get", tt.assertNoOfGet)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandler_GetModels(t *testing.T) {
	id := "algID"
	allModels := structure.Algorithm{
		Models: map[string][]string{
			id:       {"model", "model"},
			id + "2": {"model"}}}
	jsonAllModels, err := json.Marshal(allModels)
	assert.NoError(t, err)
	models := allModels.Models[id]
	jsonModels, err := json.Marshal(models)
	assert.NoError(t, err)
	tests := []struct {
		testName      string
		requestURL    string
		getReturned   string
		getError      error
		assertNoOfGet int
		bodyContains  string
		statusCode    int
	}{
		{
			testName:     "should return 404 when url is wrong",
			requestURL:   "/v1/models/{id}/wrong",
			bodyContains: "404 not found",
			statusCode:   http.StatusNotFound,
		},
		{
			testName:      "should return 500 when database does not respond",
			requestURL:    "/v1/models/" + id,
			getError:      errors.New("database not respond error"),
			assertNoOfGet: 1,
			bodyContains:  "database not respond error",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when there is no `models` in database",
			requestURL:    "/v1/models/" + id,
			getError:      errors.New("key models does not exist"),
			assertNoOfGet: 1,
			bodyContains:  "key models does not exist",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when `models` are not in json format",
			requestURL:    "/v1/models/" + id,
			getReturned:   "not json",
			assertNoOfGet: 1,
			bodyContains:  "failed to unmarshal",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 200 when images where take from database correctly",
			requestURL:    "/v1/models/" + id,
			getReturned:   string(jsonAllModels),
			assertNoOfGet: 1,
			bodyContains:  string(jsonModels),
			statusCode:    http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			r, err := http.NewRequest("GET", tt.requestURL, nil)
			vars := map[string]string{
				"alg": id,
			}
			r = mux.SetURLVars(r, vars)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := []IAlgorithm{&mocks.IAlgorithm{}}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMock)
			iDatabaseMock.On("Get", "models").Return(tt.getReturned, tt.getError)

			//when
			testSubject.GetModels(w, r)

			//then
			iDatabaseMock.AssertNumberOfCalls(t, "Get", tt.assertNoOfGet)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandler_UploadModel(t *testing.T) {
	id := "algID"
	modelFile := "model.h5"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	allModels := structure.Algorithm{
		Models: map[string][]string{
			id: {modelFile, modelFile},
		}}
	jsonAllModels, err := json.Marshal(allModels)
	assert.NoError(t, err)
	newAllModels := structure.Algorithm{
		Models: map[string][]string{
			id: {modelFile, modelFile, modelFile},
		}}
	jsonNewAllModels, err := json.Marshal(newAllModels)
	assert.NoError(t, err)
	tests := []struct {
		testName              string
		requestURL            string
		body                  io.Reader
		contentType           string
		formModel             string
		getIDReturned         string
		uploadError           error
		getReturned           string
		getError              error
		insertData            string
		insertError           error
		assertNoOfInsert      int
		assertNoOfGet         int
		assertNoOfGetID       int
		assertNoOfUploadModel int
		bodyContains          string
		statusCode            int
	}{
		{
			testName:     "should return 404 when url is wrong",
			requestURL:   "/v1/models/wrong",
			bodyContains: "404 not found",
			statusCode:   http.StatusNotFound,
		},
		{
			testName:     "should return 400 error when Content-type is wrong",
			requestURL:   "/v1/models",
			contentType:  "multipart/wrong",
			bodyContains: "Content-Type isn't multipart/form-data",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:     "should return 400 error when there is no boundary param",
			requestURL:   "/v1/models",
			contentType:  "multipart/form-data",
			bodyContains: "no multipart boundary param in Content-Type",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:     "should return 400 error when there is no `model` file",
			requestURL:   "/v1/models",
			contentType:  writer.FormDataContentType(),
			formModel:    "-",
			bodyContains: "no such file",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:              "should return 500 when failed to upload model",
			requestURL:            "/v1/models",
			contentType:           writer.FormDataContentType(),
			formModel:             "model",
			getIDReturned:         id,
			uploadError:           errors.New("failed to upload model"),
			assertNoOfGetID:       1,
			assertNoOfUploadModel: 1,
			assertNoOfGet:         1,
			bodyContains:          "failed to upload model",
			statusCode:            http.StatusInternalServerError,
		},
		{
			testName:              "should return 500 when database does not respond",
			requestURL:            "/v1/models",
			contentType:           writer.FormDataContentType(),
			formModel:             "model",
			getIDReturned:         id,
			getError:              errors.New("database not respond error"),
			assertNoOfGetID:       1,
			assertNoOfUploadModel: 1,
			assertNoOfGet:         1,
			bodyContains:          "database not respond error",
			statusCode:            http.StatusInternalServerError,
		},
		{
			testName:              "should return 500 error when failed to unmarshal",
			requestURL:            "/v1/models",
			contentType:           writer.FormDataContentType(),
			formModel:             "model",
			getIDReturned:         id,
			assertNoOfGetID:       1,
			assertNoOfUploadModel: 1,
			assertNoOfGet:         1,
			bodyContains:          "failed to unmarshal",
			statusCode:            http.StatusInternalServerError,
		},
		{
			testName:              "should return 500 when failed to insert to db",
			requestURL:            "/v1/models",
			contentType:           writer.FormDataContentType(),
			formModel:             "model",
			getIDReturned:         id,
			getReturned:           string(jsonAllModels),
			insertData:            string(jsonNewAllModels),
			insertError:           errors.New("failed to insert to db"),
			assertNoOfGetID:       1,
			assertNoOfUploadModel: 1,
			assertNoOfGet:         1,
			assertNoOfInsert:      1,
			bodyContains:          "failed to insert to db",
			statusCode:            http.StatusInternalServerError,
		},
		{
			testName:              "should return 200 when model was correctly uploaded",
			requestURL:            "/v1/models",
			contentType:           writer.FormDataContentType(),
			formModel:             "model",
			getIDReturned:         id,
			getReturned:           string(jsonAllModels),
			insertData:            string(jsonNewAllModels),
			assertNoOfGetID:       1,
			assertNoOfUploadModel: 1,
			assertNoOfGet:         1,
			assertNoOfInsert:      1,
			statusCode:            http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			writer.WriteField("id", id)
			writer.WriteField("name", modelFile)
			part, _ := writer.CreateFormFile(tt.formModel, modelFile)
			part.Write([]byte(`sample`))
			writer.Close()
			r, err := http.NewRequest("PUT", tt.requestURL, body)
			assert.NoError(t, err)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := mocks.IAlgorithm{}
			iAlgorithmMocks := []IAlgorithm{&iAlgorithmMock, &iAlgorithmMock}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMocks)
			modelFile, modelHeader, err := r.FormFile("model")
			iAlgorithmMock.On("GetID").Return(tt.getIDReturned)
			iAlgorithmMock.On("UploadModel", modelFile, modelHeader).Return(tt.uploadError)
			iDatabaseMock.On("Get", "models").Return(tt.getReturned, tt.getError)
			iDatabaseMock.On("Set", "models", tt.insertData).Return(tt.insertError)

			//when
			testSubject.UploadModel(w, r)

			//then
			iAlgorithmMock.AssertNumberOfCalls(t, "GetID", tt.assertNoOfGetID)
			iAlgorithmMock.AssertNumberOfCalls(t, "UploadModel", tt.assertNoOfUploadModel)
			iDatabaseMock.AssertNumberOfCalls(t, "Get", tt.assertNoOfGet)
			iDatabaseMock.AssertNumberOfCalls(t, "Set", tt.assertNoOfInsert)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandler_RunSimulation(t *testing.T) {
	id := "algID"
	model := "model.h5"
	image := "image"
	opType := "demo"
	body := structure.Body{ID: id, Model: model, Image: image}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)
	timeStamp := fixedTime()
	dbID := timeStamp + body.ID + opType
	sendData := structure.Body{ID: dbID, Model: body.Model, Image: body.Image}
	jsonSendData, err := json.Marshal(sendData)
	assert.NoError(t, err)
	results := structure.Results{Algorithm: id, Model: body.Model, Image: body.Image, TimeStamp: timeStamp, Status: "in-progress"}
	jsonResults, err := json.Marshal(results)
	assert.NoError(t, err)
	tests := []struct {
		testName                string
		requestURL              string
		body                    io.Reader
		getIDReturned           string
		runSimulationData       []byte
		runSimulationReturned   int
		runSimulationError      error
		insertData              string
		insertError             error
		assertNoOfGetID         int
		assertNoOfRunSimulation int
		assertNoOfInsert        int
		bodyContains            string
		statusCode              int
	}{
		{
			testName:     "should return 404 when url is wrong",
			requestURL:   "/v1/simulation-results/{type}/wrong",
			bodyContains: "404 not found",
			statusCode:   http.StatusNotFound,
		},
		{
			testName:     "should return 400 when failed to unmarshal body",
			requestURL:   "/v1/simulation-results/",
			body:         bytes.NewBuffer([]byte{}),
			bodyContains: "failed to unmarshal",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:        "should return 500 when failed to find algorithm",
			requestURL:      "/v1/simulation-results/",
			body:            bytes.NewBuffer(jsonBody),
			assertNoOfGetID: 2,
			bodyContains:    "algorithm with this id does not exists",
			statusCode:      http.StatusInternalServerError,
		},
		{
			testName:                "should return 500 when run simulation response is not 200",
			requestURL:              "/v1/simulation-results/",
			body:                    bytes.NewBuffer(jsonBody),
			getIDReturned:           id,
			runSimulationData:       jsonSendData,
			runSimulationReturned:   500,
			assertNoOfGetID:         1,
			assertNoOfRunSimulation: 1,
			bodyContains:            "failed to run simulation",
			statusCode:              http.StatusInternalServerError,
		},
		{
			testName:                "should return 500 when failed to run simulation",
			requestURL:              "/v1/simulation-results/",
			body:                    bytes.NewBuffer(jsonBody),
			getIDReturned:           id,
			runSimulationData:       jsonSendData,
			runSimulationError:      errors.New("failed to run simulation"),
			runSimulationReturned:   200,
			assertNoOfGetID:         1,
			assertNoOfRunSimulation: 1,
			bodyContains:            "failed to run simulation",
			statusCode:              http.StatusInternalServerError,
		},
		{
			testName:                "should return 500 when failed to insert data to database",
			requestURL:              "/v1/simulation-results/",
			body:                    bytes.NewBuffer(jsonBody),
			getIDReturned:           id,
			runSimulationData:       jsonSendData,
			runSimulationReturned:   200,
			insertData:              string(jsonResults),
			insertError:             errors.New("failed to insert"),
			assertNoOfGetID:         1,
			assertNoOfRunSimulation: 1,
			assertNoOfInsert:        1,
			bodyContains:            "failed to insert",
			statusCode:              http.StatusInternalServerError,
		},
		{
			testName:                "should return 202 when simulation started and data inserted to database",
			requestURL:              "/v1/simulation-results/",
			body:                    bytes.NewBuffer(jsonBody),
			getIDReturned:           id,
			runSimulationData:       jsonSendData,
			runSimulationReturned:   200,
			insertData:              string(jsonResults),
			assertNoOfGetID:         1,
			assertNoOfRunSimulation: 1,
			assertNoOfInsert:        1,
			statusCode:              http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			r, err := http.NewRequest("POST", tt.requestURL+opType, tt.body)
			vars := map[string]string{
				"type": opType,
			}
			r = mux.SetURLVars(r, vars)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := mocks.IAlgorithm{}
			iAlgorithmMocks := []IAlgorithm{&iAlgorithmMock, &iAlgorithmMock}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMocks)
			iAlgorithmMock.On("GetID").Return(tt.getIDReturned)
			iAlgorithmMock.On("RunSimulation", opType, tt.runSimulationData).Return(tt.runSimulationReturned, tt.runSimulationError)
			iDatabaseMock.On("Set", dbID, tt.insertData).Return(tt.insertError)

			//when
			testSubject.RunSimulation(w, r)

			//then
			iAlgorithmMock.AssertNumberOfCalls(t, "GetID", tt.assertNoOfGetID)
			iAlgorithmMock.AssertNumberOfCalls(t, "RunSimulation", tt.assertNoOfRunSimulation)
			iDatabaseMock.AssertNumberOfCalls(t, "Set", tt.assertNoOfInsert)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandler_UpdateResults(t *testing.T) {
	id := "algID"
	opType := "demo"
	model := "model.h5"
	image := "image"
	bb := structure.Body{ID: id, Model: model, Image: image}
	timeStamp := fixedTime()
	dbID := timeStamp + bb.ID + opType
	results := "results"
	body := structure.Data{ID: dbID, Content: results}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)
	beforeResults := structure.Results{Algorithm: id, Model: bb.Model, Image: bb.Image, TimeStamp: timeStamp, Status: "in-progress"}
	jsonBeforeResults, err := json.Marshal(beforeResults)
	assert.NoError(t, err)
	afterResults := structure.Results{Algorithm: id, Model: bb.Model, Image: bb.Image, TimeStamp: timeStamp, Status: "finished", Result: results}
	jsonAfterResults, err := json.Marshal(afterResults)
	assert.NoError(t, err)
	errResults := "\"error\""
	errBody := structure.Data{ID: dbID, Content: errResults}
	jsonErrBody, err := json.Marshal(errBody)
	assert.NoError(t, err)
	errorResults := structure.Results{Algorithm: id, Model: bb.Model, Image: bb.Image, TimeStamp: timeStamp, Status: "error"}
	jsonErrorResults, err := json.Marshal(errorResults)
	assert.NoError(t, err)
	tests := []struct {
		testName         string
		requestURL       string
		body             io.Reader
		contentType      string
		getReturned      string
		getError         error
		insertData       string
		insertError      error
		assertNoOfGet    int
		assertNoOfInsert int
		bodyContains     string
		statusCode       int
	}{
		{
			testName:     "should return 404 when url is wrong",
			requestURL:   "/v1/simulation-results/wrong",
			bodyContains: "404 not found",
			statusCode:   http.StatusNotFound,
		},
		{
			testName:     "should return 400 when content type is invalid",
			requestURL:   "/v1/simulation-results",
			body:         bytes.NewBuffer([]byte{}),
			contentType:  "application/wrong",
			bodyContains: "invalid content type",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:     "should return 400 when failed to unmarshal body",
			requestURL:   "/v1/simulation-results",
			body:         bytes.NewBuffer([]byte{}),
			contentType:  "application/json",
			bodyContains: "failed to unmarshal",
			statusCode:   http.StatusBadRequest,
		},
		{
			testName:      "should return 500 when failed to get data from database",
			requestURL:    "/v1/simulation-results",
			body:          bytes.NewBuffer(jsonBody),
			contentType:   "application/json",
			getError:      errors.New("failed to get data from database"),
			assertNoOfGet: 1,
			bodyContains:  "failed to get data from database",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:      "should return 500 when failed to unmarshal data from database",
			requestURL:    "/v1/simulation-results",
			body:          bytes.NewBuffer(jsonBody),
			contentType:   "application/json",
			getReturned:   "string",
			assertNoOfGet: 1,
			bodyContains:  "failed to unmarshal",
			statusCode:    http.StatusInternalServerError,
		},
		{
			testName:         "should return 200 when results were updated - with error ",
			requestURL:       "/v1/simulation-results",
			body:             bytes.NewBuffer(jsonErrBody),
			contentType:      "application/json",
			getReturned:      string(jsonBeforeResults),
			insertData:       string(jsonErrorResults),
			assertNoOfGet:    1,
			assertNoOfInsert: 1,
			statusCode:       http.StatusOK,
		},
		{
			testName:         "should return 500 when failed to insert results to database",
			requestURL:       "/v1/simulation-results",
			body:             bytes.NewBuffer(jsonBody),
			contentType:      "application/json",
			getReturned:      string(jsonBeforeResults),
			insertData:       string(jsonAfterResults),
			insertError:      errors.New("failed to insert"),
			assertNoOfGet:    1,
			assertNoOfInsert: 1,
			bodyContains:     "failed to insert",
			statusCode:       http.StatusInternalServerError,
		},
		{
			testName:         "should return 200 when results where updated",
			requestURL:       "/v1/simulation-results",
			body:             bytes.NewBuffer(jsonBody),
			contentType:      "application/json",
			getReturned:      string(jsonBeforeResults),
			insertData:       string(jsonAfterResults),
			assertNoOfGet:    1,
			assertNoOfInsert: 1,
			statusCode:       http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			r, err := http.NewRequest("PUT", tt.requestURL, tt.body)
			vars := map[string]string{
				"type": opType,
			}
			r = mux.SetURLVars(r, vars)
			assert.NoError(t, err)
			r.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := mocks.IAlgorithm{}
			iAlgorithmMocks := []IAlgorithm{&iAlgorithmMock, &iAlgorithmMock}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMocks)
			iDatabaseMock.On("Get", dbID).Return(tt.getReturned, tt.getError)
			iDatabaseMock.On("Set", dbID, tt.insertData).Return(tt.insertError)

			//when
			testSubject.UpdateResults(w, r)

			//then
			iDatabaseMock.AssertNumberOfCalls(t, "Get", tt.assertNoOfGet)
			iDatabaseMock.AssertNumberOfCalls(t, "Set", tt.assertNoOfInsert)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

func TestHandler_GetResults(t *testing.T) {
	alg := "algID"
	opType := "demo"
	model := "model.h5"
	image := "image"
	results := "results"
	timeStamp := fixedTime()
	keys := []string{alg, alg}
	dbResult := structure.Results{Algorithm: alg, Model: model, Image: image, TimeStamp: timeStamp, Status: "finished", Result: results}
	jsonDBResult, err := json.Marshal(dbResult)
	assert.NoError(t, err)
	dbResults := []structure.Results{dbResult, dbResult}
	jsonDBResults, err := json.Marshal(dbResults)
	assert.NoError(t, err)
	tests := []struct {
		testName         string
		requestURL       string
		keysReturned     []string
		keysError        error
		getReturned      string
		getError         error
		assertNoOfKeys   int
		assertNoOfGet    int
		assertNoOfInsert int
		bodyContains     string
		statusCode       int
	}{
		{
			testName:     "should return 404 when url is wrong",
			requestURL:   "/v1/simulation-results/wrong",
			keysError:    errors.New("failed to get keys from database"),
			bodyContains: "404 not found",
			statusCode:   http.StatusNotFound,
		},
		{
			testName:       "should return 500 when failed to get keys from database",
			requestURL:     "/v1/simulation-results/" + opType + "/" + alg,
			keysError:      errors.New("failed to get keys from database"),
			assertNoOfKeys: 1,
			bodyContains:   "failed to get keys from database",
			statusCode:     http.StatusInternalServerError,
		},
		{
			testName:       "should return 500 when failed to get results from database",
			requestURL:     "/v1/simulation-results/" + opType + "/" + alg,
			keysReturned:   keys,
			getError:       errors.New("failed to get results from database"),
			assertNoOfKeys: 1,
			assertNoOfGet:  1,
			bodyContains:   "failed to get results from database",
			statusCode:     http.StatusInternalServerError,
		},
		{
			testName:       "should return 500 when failed to unmarshal",
			requestURL:     "/v1/simulation-results/" + opType + "/" + alg,
			keysReturned:   keys,
			getReturned:    "",
			assertNoOfKeys: 1,
			assertNoOfGet:  1,
			bodyContains:   "failed to unmarshal",
			statusCode:     http.StatusInternalServerError,
		},
		{
			testName:       "should return 200 when results were properly taken from database",
			requestURL:     "/v1/simulation-results/" + opType + "/" + alg,
			keysReturned:   keys,
			getReturned:    string(jsonDBResult),
			assertNoOfKeys: 1,
			assertNoOfGet:  2,
			bodyContains:   string(jsonDBResults),
			statusCode:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			//given
			r, err := http.NewRequest("GET", tt.requestURL, nil)
			assert.NoError(t, err)
			vars := map[string]string{
				"type": opType,
				"alg":  alg,
			}
			r = mux.SetURLVars(r, vars)
			w := httptest.NewRecorder()
			iDatabaseMock := mocks.IDatabase{}
			iAlgorithmMock := mocks.IAlgorithm{}
			iAlgorithmMocks := []IAlgorithm{&iAlgorithmMock, &iAlgorithmMock}
			testSubject := NewHandler(&iDatabaseMock, iAlgorithmMocks)
			iDatabaseMock.On("Keys", "20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z"+alg+opType).Return(tt.keysReturned, tt.keysError)
			iDatabaseMock.On("Get", alg).Return(tt.getReturned, tt.getError)
			iDatabaseMock.On("Get", alg).Return(tt.getReturned, tt.getError)

			//when
			testSubject.GetResults(w, r)

			//then
			iDatabaseMock.AssertNumberOfCalls(t, "Keys", tt.assertNoOfKeys)
			iDatabaseMock.AssertNumberOfCalls(t, "Get", tt.assertNoOfGet)
			assert.Contains(t, w.Body.String(), tt.bodyContains)
			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}
