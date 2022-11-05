package repository

import (
	"encoding/json"
	"fetch-app/models"
	"net/http"
)

type StorageRepository interface {
	FetchStorages() (storages []*models.Storage, err error)
}

type storageRepository struct {
	client *http.Client
}

func NewStorageRepository() StorageRepository {
	return &storageRepository{
		client: http.DefaultClient,
	}
}

func (r *storageRepository) FetchStorages() (storages []*models.Storage, err error) {

	request, err := http.NewRequest("GET", "https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list", nil)
	if err != nil {
		return
	}

	response, err := r.client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&storages)
	if err != nil {
		return
	}

	return
}
