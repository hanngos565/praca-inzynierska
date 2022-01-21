package algorithm

/*
import (
	"mime/multipart"
	"reflect"
	"testing"
)

func TestAlgorithm_GetID(t *testing.T) {
	type fields struct {
		ID  string
		URL string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Algorithm{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			if got := a.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_RunSimulation(t *testing.T) {
	type fields struct {
		ID  string
		URL string
	}
	type args struct {
		opType string
		data   []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Algorithm{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			got, err := a.RunSimulation(tt.args.opType, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunSimulation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunSimulation() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlgorithm_UploadModel(t *testing.T) {
	type fields struct {
		ID  string
		URL string
	}
	type args struct {
		modelFile   multipart.File
		modelHeader *multipart.FileHeader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Algorithm{
				ID:  tt.fields.ID,
				URL: tt.fields.URL,
			}
			if err := a.UploadModel(tt.args.modelFile, tt.args.modelHeader); (err != nil) != tt.wantErr {
				t.Errorf("UploadModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAlgorithm(t *testing.T) {
	type args struct {
		id  string
		url string
	}
	tests := []struct {
		name string
		args args
		want Algorithm
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAlgorithm(tt.args.id, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAlgorithm() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
