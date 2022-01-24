package api

import (
	"backend/internal/structure"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

//go:generate mockery --name=IDatabase
type IDatabase interface {
	Set(key string, value string) error
	Get(key string) (interface{}, error)
	Keys(pattern string) ([]string, error)
}

//go:generate mockery --name=IAlgorithm
type IAlgorithm interface {
	GetID() string
	UploadModel(modelFile multipart.File, modelHeader *multipart.FileHeader) error
	RunSimulation(opType string, data []byte) (int, error)
}

type Handler struct {
	iDatabase         IDatabase
	iAlgorithm        []IAlgorithm
	getImagesEndpoint string
	putImageEndpoint  string

	getModelsEndpoint string
	postModelEndpoint string

	postSimulationResultsEndpoint string
	putSimulationResultsEndpoint  string
	getSimulationResultsEndpoint  string
}

func (h Handler) InitializeEndpoints(mux *mux.Router) {
	mux.HandleFunc(h.getImagesEndpoint, h.GetImages).Methods("GET")
	mux.HandleFunc(h.putImageEndpoint, h.AddImage).Methods("PUT")
	mux.HandleFunc(h.getModelsEndpoint, h.GetModels).Methods("GET")
	mux.HandleFunc(h.postModelEndpoint, h.UploadModel).Methods("PUT")
	mux.HandleFunc(h.postSimulationResultsEndpoint, h.RunSimulation).Methods("POST")
	mux.HandleFunc(h.putSimulationResultsEndpoint, h.UpdateResults).Methods("PUT")
	mux.HandleFunc(h.getSimulationResultsEndpoint, h.GetResults).Methods("GET")
}

func NewHandler(iDatabase IDatabase, iAlgorithm []IAlgorithm) Handler {
	return Handler{
		iDatabase:                     iDatabase,
		iAlgorithm:                    iAlgorithm,
		getImagesEndpoint:             "/v1/images",
		putImageEndpoint:              "/v1/images",
		getModelsEndpoint:             "/v1/models/{alg}",
		postSimulationResultsEndpoint: "/v1/simulation-results/{type}",
		putSimulationResultsEndpoint:  "/v1/simulation-results",
		getSimulationResultsEndpoint:  "/v1/simulation-results/{type}/{alg}",
		postModelEndpoint:             "/v1/models",
	}
}

func (h Handler) Config() error {
	model := map[string][]string{}
	for _, item := range h.iAlgorithm {
		model[item.GetID()] = []string{"default"}
	}

	if _, err := h.iDatabase.Get("models"); err != nil {
		if err.Error() != "key does not exist" {
			return err
		}
	}
	models := structure.Algorithm{Models: model}
	jsonModels, err := json.Marshal(models)
	if err != nil {
		return err
	}
	if err = h.iDatabase.Set("models", string(jsonModels)); err != nil {
		return err
	}
	if _, err = h.iDatabase.Get("images"); err != nil {
		if err.Error() != "key does not exist" {
			return err
		}
	}
	images := structure.Images{}
	jsonImages, err := json.Marshal(images)
	if err = h.iDatabase.Set("images", string(jsonImages)); err != nil {
		return err
	}

	return nil
}

//PUT /v1/images/
func (h Handler) AddImage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != h.putImageEndpoint {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data structure.Data
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, "failed to unmarshal body", http.StatusBadRequest)
		return
	}

	fromDB, err := h.iDatabase.Get("images")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var allImages structure.Images
	if err = json.Unmarshal([]byte(fromDB.(string)), &allImages); err != nil {
		http.Error(w, "failed to unmarshal images", http.StatusInternalServerError)
		return
	}
	allImages.Images = append(allImages.Images, data.Content)

	jsonAllImages, err := json.Marshal(allImages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.iDatabase.Set("images", string(jsonAllImages)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//GET /v1/images
func (h Handler) GetImages(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != h.getImagesEndpoint {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	fromDB, err := h.iDatabase.Get("images")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var images structure.Images
	if err = json.Unmarshal([]byte(fromDB.(string)), &images); err != nil {
		http.Error(w, "failed to unmarshal", http.StatusInternalServerError)
		return
	}
	jsonImages, err := json.Marshal(images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = fmt.Fprint(w, string(jsonImages)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//PUT /v1/models
func (h Handler) UploadModel(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != h.postModelEndpoint {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "ParseMultipartForm"+err.Error(), http.StatusBadRequest)
		return
	}

	modelFile, modelHeader, err := r.FormFile("model")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer modelFile.Close()

	name := r.PostFormValue("name")
	modelHeader.Filename = name

	id := r.PostFormValue("id")
	index := -1
	for i, item := range h.iAlgorithm {
		if algID := item.GetID(); algID == id {
			index = i
			break
		}
	}
	if index < 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.iAlgorithm[index].UploadModel(modelFile, modelHeader); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fromDB, err := h.iDatabase.Get("models")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var allModels structure.Algorithm
	if err = json.Unmarshal([]byte(fromDB.(string)), &allModels); err != nil {
		http.Error(w, "failed to unmarshal "+err.Error(), http.StatusInternalServerError)
		return
	}
	models := allModels.Models[id]
	models = append(models, name)
	allModels.Models[id] = models

	jsonAllModels, err := json.Marshal(allModels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = h.iDatabase.Set("models", string(jsonAllModels)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

//GET /v1/models/{alg}
func (h Handler) GetModels(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["alg"]
	url := strings.Replace(h.getModelsEndpoint, "{alg}", id, 1)
	if r.URL.Path != url {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	fromDB, err := h.iDatabase.Get("models")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var allModels structure.Algorithm
	err = json.Unmarshal([]byte(fromDB.(string)), &allModels)
	if err != nil {
		http.Error(w, "failed to unmarshal", http.StatusInternalServerError)
		return
	}

	models := allModels.Models[id]
	jsonModels, err := json.Marshal(models)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = fmt.Fprint(w, string(jsonModels)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//POST /v1/simulation-results/{type}
func (h Handler) RunSimulation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	opType := params["type"]
	url := strings.Replace(h.postSimulationResultsEndpoint, "{type}", opType, 1)
	if r.URL.Path != url {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	bd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var body structure.Body
	if err = json.Unmarshal(bd, &body); err != nil {
		http.Error(w, "failed to unmarshal body", http.StatusBadRequest)
		return
	}

	timeStamp := time.Now()
	dbID := timeStamp.Format(time.RFC3339) + body.ID + opType

	index := -1
	for i, item := range h.iAlgorithm {
		if algID := item.GetID(); algID == body.ID {
			index = i
			break
		}
	}
	if index < 0 {
		http.Error(w, "algorithm with this id does not exists", http.StatusInternalServerError)
		return
	}

	sendData := structure.Body{ID: dbID, Model: body.Model, Image: body.Image}
	jsonSendData, err := json.Marshal(sendData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respCode, err := h.iAlgorithm[index].RunSimulation(opType, jsonSendData)
	if err != nil || respCode != 200 {
		http.Error(w, "failed to run simulation", http.StatusInternalServerError)
		return
	}

	results := structure.Results{
		Algorithm: body.ID,
		Model:     body.Model,
		Image:     body.Image,
		TimeStamp: timeStamp.Format(time.RFC3339),
		Status:    "in-progress",
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.iDatabase.Set(dbID, string(jsonResults)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

//PUT /v1/simulation-results
func (h Handler) UpdateResults(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != h.putSimulationResultsEndpoint {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusBadRequest)
		return
	}

	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data structure.Data
	if err = json.Unmarshal(value, &data); err != nil {
		http.Error(w, "failed to unmarshal "+err.Error(), http.StatusBadRequest)
		return
	}

	fromDB, err := h.iDatabase.Get(data.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var results structure.Results
	if err = json.Unmarshal([]byte(fromDB.(string)), &results); err != nil {
		http.Error(w, "failed to unmarshal "+err.Error(), http.StatusInternalServerError)
		return
	}

	if data.Content == "\"error\"" {
		results.Status = "error"
	} else {
		results.Result = data.Content
		results.Status = "finished"
	}

	jsonResults, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.iDatabase.Set(data.ID, string(jsonResults)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//GET /v1/simulation-results/{type}/{alg}
func (h Handler) GetResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	opType := params["type"]
	id := params["alg"]
	url := strings.Replace(strings.Replace(h.getSimulationResultsEndpoint, "{type}", opType, 1), "{alg}", id, 1)
	if r.URL.Path != url {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	pattern := fmt.Sprintf("20[0-9][0-9]-[0-9][0-9]-[0-9][0-9]T[0-9][0-9]:[0-9][0-9]:[0-9][0-9]Z%s%s", id, opType)

	keys, err := h.iDatabase.Keys(pattern)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var result structure.Results
	var results []structure.Results
	for _, key := range keys {
		fromDB, err := h.iDatabase.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal([]byte(fromDB.(string)), &result); err != nil {
			http.Error(w, "failed to unmarshal", http.StatusInternalServerError)
			return
		}

		results = append(results, result)
	}
	jsonResults, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(jsonResults))
}
