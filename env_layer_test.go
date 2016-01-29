package onion

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEnvLayer(t *testing.T) {
	Convey("EnvLayer test ", t, func() {
		// Just for passing the errcheck errors :)
		// Unsetenv is not available in < 1.3
		//So(os.Unsetenv("TEST1"), ShouldBeNil)
		//So(os.Unsetenv("TEST2"), ShouldBeNil)
		//So(os.Unsetenv("TEST3"), ShouldBeNil)
		So(os.Setenv("BLACK", "blacklisted"), ShouldBeNil)
		//o := New()

		Convey("Check if there is anything loaded", func() {
			el := NewEnvLayer("TEST1", "test2", "Test3")

			data, err := el.Load()
			So(err, ShouldBeNil)
			So(len(data), ShouldEqual, 0)
		})

		Convey("Check if the variable is loaded correctly", func() {
			el := NewEnvLayer("BLACK")

			data, err := el.Load()
			So(err, ShouldBeNil)
			So(len(data), ShouldEqual, 1)
			So(data["BLACK"], ShouldEqual, "blacklisted")

			Convey("Check if the onion handle it correctly", func() {
				o := New()
				err := o.AddLayer(el)
				So(err, ShouldBeNil)

				So(o.GetStringDefault("black", "no!"), ShouldEqual, "blacklisted")
			})
		})
	})

}
