package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context/param"
)

func init() {

    beego.GlobalControllerRouter["fetch-app/controllers:StoragePrivateController"] = append(beego.GlobalControllerRouter["fetch-app/controllers:StoragePrivateController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("limit"),
				param.New("page"),
			),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["fetch-app/controllers:StoragePrivateController"] = append(beego.GlobalControllerRouter["fetch-app/controllers:StoragePrivateController"],
        beego.ControllerComments{
            Method: "GetAllAggregated",
            Router: "/aggregated",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(
				param.New("limit"),
				param.New("page"),
			),
            Filters: nil,
            Params: nil})

}
