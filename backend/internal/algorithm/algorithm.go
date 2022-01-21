package algorithm

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Algorithm struct {
	ID  string
	URL string
}

func NewAlgorithm(id string, url string) Algorithm {
	return Algorithm{
		ID:  id,
		URL: url,
	}
}

func (a Algorithm) UploadModel(modelFile multipart.File, modelHeader *multipart.FileHeader) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("model", modelHeader.Filename)
	if err != nil {
		return errors.New("CreateFormFile" + err.Error())
	}
	if _, err = io.Copy(fw, modelFile); err != nil {
		return errors.New("ioCopy" + err.Error())
	}
	writer.Close()

	req, err := http.NewRequest("POST", a.URL+"/upload_model", bytes.NewReader(body.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}

func (a Algorithm) RunSimulation(opType string, data []byte) (int, error) {
	req, err := http.NewRequest("POST", a.URL+"/"+opType, bytes.NewReader(data))
	if err != nil {
		return 0, errors.New("NewRequest" + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("DoRequest" + err.Error())
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func (a Algorithm) GetID() string {
	return a.ID
}
