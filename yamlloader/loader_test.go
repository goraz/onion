package yamlloader

import (
	"bytes"
	"io"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "gopkg.in/fzerorubigd/onion.v3"
)

func TestYamlLoader(t *testing.T) {
	Convey("Load a yaml structure into a json", t, func() {
		// Make sure there is two ile available in the tmp file system
		tmp := os.TempDir()
		dir := tmp + "/onion_f"
		So(os.MkdirAll(dir, 0744), ShouldBeNil)
		path := dir + "/test.yaml"
		path2 := dir + "/invalid.yml"

		buf := bytes.NewBufferString(`---
  str: "string_data"
  bool: true
  integer: 10
  nested:
    key1: "string"
    key2: 100
`)
		bufInvalid := bytes.NewBufferString(`---
str: - inv
  lid
 s
ALALA`)
		f, err := os.Create(path)
		So(err, ShouldBeNil)
		_, err = io.Copy(f, buf)
		So(err, ShouldBeNil)

		f2, err := os.Create(path2)
		So(err, ShouldBeNil)
		_, err = io.Copy(f2, bufInvalid)
		So(err, ShouldBeNil)

		defer func() {
			_ = os.Remove(path)
			_ = os.Remove(path2)
		}()
		So(f.Close(), ShouldBeNil)
		So(f2.Close(), ShouldBeNil)

		Convey("Check if the file is loaded correctly ", func() {
			fl := NewFileLayer(path)
			o := New()
			err := o.AddLayer(fl)
			So(err, ShouldBeNil)
			So(o.GetStringDefault("str", ""), ShouldEqual, "string_data")
			So(o.GetStringDefault("nested.key1", ""), ShouldEqual, "string")
			So(o.GetIntDefault("nested.key2", 0), ShouldEqual, 100)
			So(o.GetBoolDefault("bool", false), ShouldBeTrue)

			a := New() // Just for test load again
			So(a.AddLayer(fl), ShouldBeNil)
		})

		Convey("Check for the invalid file content", func() {
			fl := NewFileLayer(path2)
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldNotBeNil)
		})
	})
}
