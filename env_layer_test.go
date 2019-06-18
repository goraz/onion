package onion

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewEnvLayer(t *testing.T) {
	Convey("ENV Tests", t, func() {

		os.Setenv("KEY_TEST_SEP", "1")
		l := NewEnvLayer("_", "KEY_TEST_SEP")
		o := New(l)
		So(o.GetInt("key.test.sep"), ShouldEqual, 1)

		l2 := NewEnvLayerPrefix("_", "key")
		o2 := New(l2)
		So(o2.GetInt("test.sep"), ShouldEqual, 1)
	})
}
