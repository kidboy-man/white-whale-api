package controllers

import (
	"fetch-app/datatransfers"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/pagination"
)

type BaseController struct {
	beego.Controller
	JSONResponse *JSONResponse
}

type JSONResponse struct {
	Success     bool        `json:"success"`
	Status      int         `json:"status"` // http status code
	Data        interface{} `json:"data"`   //
	Error       error       `json:"error"`  //
	CurrentPage int         `json:"currentPage"`
	TotalPages  int         `json:"totalPages"`
	DataPerPage int         `json:"dataPerPage"`
	HasNextPage bool        `json:"hasNextPage"`
	HasPrevPage bool        `json:"hasPrevPage"`
}

func doReturnOK(response *JSONResponse, obj interface{}) {

	response.Success = true
	response.Status = http.StatusOK
	response.Data = obj
}

func doReturnNotOK(response *JSONResponse, err error) {
	response.Success = false
	response.Status = http.StatusInternalServerError
	if v, ok := err.(*datatransfers.CustomError); ok {
		response.Error = v
		response.Status = v.Status

		return
	}

	response.Error = err
}

func (c *BaseController) ReturnJSONResponse(obj interface{}, err error) *JSONResponse {
	c.JSONResponse = &JSONResponse{}
	if err != nil {
		doReturnNotOK(c.JSONResponse, err)
	} else {
		doReturnOK(c.JSONResponse, obj)
	}

	c.Ctx.Output.SetStatus(c.JSONResponse.Status)

	return c.JSONResponse
}

func (c *BaseController) setPagination(ctx *context.Context, totalData int64, limit, page int) {
	paginator := pagination.SetPaginator(ctx, limit, totalData)
	c.JSONResponse.CurrentPage = page
	c.JSONResponse.TotalPages = paginator.PageNums()
	c.JSONResponse.DataPerPage = limit
	c.JSONResponse.HasNextPage = paginator.HasNext()
	c.JSONResponse.HasPrevPage = paginator.HasPrev()
}

func (c *BaseController) GetUserIDFromToken() (uid string) {
	uid = c.Ctx.Input.GetData("uid").(string)
	return
}

func (c *BaseController) ReturnJSONListResponse(obj interface{}, cnt int64, limit, page int, err error) *JSONResponse {
	c.JSONResponse = c.ReturnJSONResponse(obj, err)
	c.setPagination(c.Ctx, cnt, limit, page)
	return c.JSONResponse
}
