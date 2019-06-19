package filewatchlayer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/fzerorubigd/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func writeJson(fl string, data map[string]interface{}) error {
	f, err := os.Create(fl)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return json.NewEncoder(f).Encode(data)
}

func TestNewFileWatchLayerContext(t *testing.T) {
	Convey("Test read file system", t, func() {
		cfg, err := ioutil.TempFile(os.TempDir(), "*.json")
		So(err, ShouldBeNil)
		fl := cfg.Name()
		defer func() { _ = os.Remove(fl) }()

		So(cfg.Close(), ShouldBeNil)
		So(writeJson(fl, map[string]interface{}{"hi": 100}), ShouldBeNil)
		l, err := NewFileWatchLayer(fl)
		o := onion.New(l)
		So(o.GetInt("hi"), ShouldEqual, 100)
		So(writeJson(fl, map[string]interface{}{"hi": 200}), ShouldBeNil)
		// TODO : Event channel
		time.Sleep(time.Second)
		So(o.GetInt("hi"), ShouldEqual, 200)
	})
}
