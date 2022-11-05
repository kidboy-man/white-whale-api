package controllers

import (
	usecase "fetch-app/usecases"
)

type StoragePrivateController struct {
	BaseController
	storageUcase usecase.StorageUsecase
}

func (c *StoragePrivateController) Prepare() {
	c.storageUcase = usecase.NewStorageUsecase()
}

// @Title Get All Storages
// @Description Get All Storages
// @Summary Get All Storages
// @Param limit query int false "limit of this request"
// @Param page query int false "page of this request"
// @Success 200
// @Failure 403
// @router / [get]
func (c *StoragePrivateController) GetAll(limit, page int) *JSONResponse {
	storages, err := c.storageUcase.GetStorages()
	return c.ReturnJSONResponse(storages, err)
}

// @Title Get All Storages
// @Description Get All Storages
// @Summary Get All Storages
// @Param limit query int false "limit of this request"
// @Param page query int false "page of this request"
// @Success 200
// @Failure 403
// @router /agregated [get]
func (c *StoragePrivateController) GetAllAgregates(limit, page int) *JSONResponse {
	storages, err := c.storageUcase.GetAgregatedStorages()
	return c.ReturnJSONResponse(storages, err)
}
