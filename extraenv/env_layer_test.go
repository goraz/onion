package extraenv

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "gopkg.in/fzerorubigd/onion.v3"
)

func TestExtraEnvLoader(t *testing.T) {
	Convey("Extra env in config", t, func() {
		o := New()
		layer := NewExtraEnvLayer("test")
		o.AddLazyLayer(layer)
		Convey("check data from env", func() {
			os.Setenv("TEST_DATA_NESTED", "TDN")
			So(o.GetString("data.nested"), ShouldEqual, "TDN")
		})

		Convey("check data not in env", func() {
			So(o.GetString("not.valid.data"), ShouldEqual, "")
			So(o.GetString(""), ShouldEqual, "")
		})

	})

}
