package controllers

import (
	usecase "fetch-app/usecases"
)

type StorageAdminController struct {
	BaseController
	storageUcase usecase.StorageUsecase
}

func (c *StorageAdminController) Prepare() {
	c.storageUcase = usecase.NewStorageUsecase()
}

// @Title Get All Storages
// @Description Get All Storages
// @Summary Get All Storages
// @Param limit query int false "limit of this request"
// @Param page query int false "page of this request"
// @Param authorization header string true "bearer token in jwt"
// @Success 200
// @Failure 403
// @router /aggregated [get]
func (c *StorageAdminController) GetAllAggregated(limit, page int) *JSONResponse {
	storages, err := c.storageUcase.GetAggregatedStorages()
	return c.ReturnJSONResponse(storages, err)
}
