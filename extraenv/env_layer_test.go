package extraenv

import (
	"os"
	"testing"

	. "gopkg.in/fzerorubigd/onion.v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExtraEnvLoader(t *testing.T) {
	Convey("Extra env in config", t, func() {
		o := New()
		layer := NewExtraEnvLayer("test")
		o.AddLayer(layer)
		Convey("check data from env", func() {
			os.Setenv("TEST_DATA_NESTED", "TDN")
			So(o.GetString("data.nested"), ShouldEqual, "TDN")
		})

	})

}
