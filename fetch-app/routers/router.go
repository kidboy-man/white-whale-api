// @APIVersion 1.0.0
// @Title Fetch-app API
// @Description API Documentation
// @Contact ingunawandra@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"fetch-app/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/private/storages",
			// beego.NSBefore(middlewares.VerifyTokenAdmin),
			beego.NSInclude(
				&controllers.StoragePrivateController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
