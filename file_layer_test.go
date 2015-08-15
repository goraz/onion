package onion

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileLayer(t *testing.T) {
	Convey("FileLayer test ", t, func() {
		// Make sure there is two ile available in the tmp file system
		tmp := os.TempDir()
		dir := tmp + "/onion_f"
		So(os.MkdirAll(dir, 0744), ShouldBeNil)
		path := dir + "/test.json"
		path2 := dir + "/invalid.json"

		buf := bytes.NewBufferString(`{"slicestr":["a","b","c"],"str" : "string_data","bool" : true,"integer" : 10 ,"nested" : {"key1" : "string","key2" : 100}}`)
		bufInvalid := bytes.NewBufferString(`invalid{json}[]`)
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
			So(reflect.DeepEqual(o.GetStringSlice("slicestr"), []string{"a", "b", "c"}), ShouldBeTrue)

			a := New() // Just for test load again
			So(a.AddLayer(fl), ShouldBeNil)
		})

		Convey("Check for the invalid ext", func() {
			f, err := os.Create(path + ".invalid_ext")
			So(err, ShouldBeNil)
			defer func() {
				_ = os.Remove(path + ".invalid_ext")
			}()
			So(f.Close(), ShouldBeNil)
			fl := NewFileLayer(path + ".invalid_ext")
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldNotBeNil)
		})

		Convey("Check for the invalid file content", func() {
			fl := NewFileLayer(path2)
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldNotBeNil)
		})

		Convey("Check for the invalid file", func() {
			fl := NewFileLayer(path + "imnothere.json")
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldNotBeNil)
		})
	})
}
