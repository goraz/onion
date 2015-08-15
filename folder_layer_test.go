package onion

import (
	"bytes"
	"io"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFolderLayer(t *testing.T) {
	Convey("FolderLayer test ", t, func() {
		// Make sure there is two ile available in the tmp file system
		tmp := os.TempDir()
		dir := tmp + "/onion_folder"
		So(os.MkdirAll(dir, 0744), ShouldBeNil)
		path := dir + "/test.json"

		data := `{"str" : "string_data","bool" : true,"integer" : 10 ,"nested" : {"key1" : "string","key2" : 100}}`
		buf := bytes.NewBufferString(data)

		f, err := os.Create(path)
		So(err, ShouldBeNil)
		_, err = io.Copy(f, buf)
		So(err, ShouldBeNil)
		defer func() {
			_ = os.Remove(path)
		}()
		So(f.Close(), ShouldBeNil)

		Convey("Check if the folder is loaded correctly ", func() {
			fl := NewFolderLayer(dir, "test")
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

		Convey("Check for the invalid ext", func() {
			fl := NewFolderLayer(dir, "nofile")
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldNotBeNil)
		})

		Convey("Check for the invalid config file name", func() {
			fl := NewFolderLayer(dir, "[^")
			o := New()
			err = o.AddLayer(fl)
			So(err, ShouldNotBeNil)
		})
	})

}
