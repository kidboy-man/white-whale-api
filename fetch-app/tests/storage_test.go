package test

import (
	_ "fetch-app/routers"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	logs.Info("file", file)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	logs.Info("apppath", apppath)
	beego.TestBeegoInit(apppath)
}

// TestGet is a sample to run an endpoint test
func TestGetStorages(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/storages", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Info("testing", "TestGetStorages", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Get Storages Endpoint\n", t, func() {
		Convey("Status Code Should Be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The Result Should Not Be Empty", func() {
			So(w.Body.Len(), ShouldBeGreaterThan, 0)
		})
	})
}
