package filewatchlayer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"testing"

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
		fl := cfg.Name() + ".json" // In go1.9 the TempFile function behave differently
		defer func() { _ = os.Remove(fl) }()

		So(cfg.Close(), ShouldBeNil)
		So(writeJson(fl, map[string]interface{}{"hi": 100}), ShouldBeNil)
		l, err := NewFileWatchLayer(fl, nil)
		So(err, ShouldBeNil)
		o := onion.New(l)
		So(o.GetInt("hi"), ShouldEqual, 100)
		w := sync.WaitGroup{}
		w.Add(1)
		go func() {
			defer w.Done()
			<-o.ReloadWatch()
		}()
		So(writeJson(fl, map[string]interface{}{"hi": 200}), ShouldBeNil)
		w.Wait()
		So(o.GetInt("hi"), ShouldEqual, 200)
	})
}
