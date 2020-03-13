package onion

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const testFile1 = `
{
	"string-not-to-override": "pippo",
	"string-to-override" : "this will be overridden",
	"number" : 100,
	"nested" : {
		"bool" : true
	}
}
`
const testFile2 = `
{
	"string-to-override" : "This string will override",
	"number" : 101,
	"nested" : {
		"bool" : false
	}
}
`

func TestNewFolderLayer(t *testing.T) {
	Convey("Test folder layer", t, func() {
		folderName, err := ioutil.TempDir("", "onion-test-*")
		if err != nil {
			t.Error("Something went wrong creating temp directory")
		}

		for i, testFile := range []string{testFile1, testFile2} {
			ioutil.WriteFile(folderName+"/test"+strconv.Itoa(i)+".json", []byte(testFile), 0644)
		}

		folderLayer, err := NewFolderLayer(folderName, "json")
		o := New(folderLayer)
		So(o.GetString("string-not-to-override"), ShouldEqual, "pippo")
		So(o.GetString("string-to-override"), ShouldEqual, "This string will override")
		So(o.GetInt("number"), ShouldEqual, 101)
		So(o.GetBool("nested.bool"), ShouldEqual, false)

		os.RemoveAll(folderName)
	})
}
