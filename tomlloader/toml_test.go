package tomlloader

import (
	"bytes"
	"testing"

	. "github.com/fzerorubigd/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestYamlLoader(t *testing.T) {
	Convey("Load a yaml structure into a json", t, func() {
		buf := bytes.NewBufferString(`
  str = "string_data"
  bool = true
  integer = 10
  [nested]
    key1 = "string"
    key2 = 100
`)
		bufInvalid := bytes.NewBufferString(`invalid toml file`)
		Convey("Check if the file is loaded correctly ", func() {
			fl, err := NewStreamLayer(buf, "toml")
			So(err, ShouldBeNil)
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldBeNil)
			So(o.GetStringDefault("str", ""), ShouldEqual, "string_data")
			So(o.GetStringDefault("nested.key1", ""), ShouldEqual, "string")
			So(o.GetIntDefault("nested.key2", 0), ShouldEqual, 100)
			So(o.GetBoolDefault("bool", false), ShouldBeTrue)

			a := New() // Just for test load again
			So(a.AddLayer(fl), ShouldBeNil)
		})

		Convey("Check for the invalid file content", func() {
			_, err := NewStreamLayer(bufInvalid, "toml")
			So(err, ShouldNotBeNil)
		})
	})
}
