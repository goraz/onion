package properties

import (
	"bytes"
	"testing"

	. "github.com/fzerorubigd/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPropertiesLoader(t *testing.T) {
	Convey("Load a properties structure into a config", t, func() {

		buf := bytes.NewBufferString(`---
str=string_data
bool=true
integer=10
`)
		bufInvalid := bytes.NewBufferString(`
=ss=jw=ishwi======
w
www
10*188`)

		Convey("Check if the file is loaded correctly ", func() {
			fl, err := NewStreamLayer(buf, "props", nil)
			So(err, ShouldBeNil)
			o := New(fl)
			So(o.GetStringDefault("str", ""), ShouldEqual, "string_data")
			So(o.GetBoolDefault("bool", false), ShouldBeTrue)
		})

		Convey("Check for the invalid file content", func() {
			_, err := NewStreamLayer(bufInvalid, "properties", nil)
			So(err, ShouldNotBeNil)
		})
	})
}
