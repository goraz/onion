package directorywatchlayer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/goraz/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDirectoryWatchLayerContext(t *testing.T) {
	Convey("Test watch directory change", t, func() {
		dir, errMkdir := ioutil.TempDir(os.TempDir(), "hexagon-test-onion-*")
		So(errMkdir, ShouldBeNil)
		defer func() {
			_ = os.RemoveAll(dir)
		}()

		filenames := make([]string, 3)
		for k := range filenames {
			fh, errTouch := ioutil.TempFile(dir, "*.json")
			So(errTouch, ShouldBeNil)

			filenames[k] = fh.Name()
			So(fh.Close(), ShouldBeNil)

			So(writeJson(filenames[k], getCfgMapFixture(1, k)), ShouldBeNil)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ll, errInit := NewDirectoryWatchLayerContext(ctx, dir, nil, "json")
		So(errInit, ShouldBeNil)

		cfg := onion.New(ll...)
		So(cfg.GetInt(cfgKey(0)), ShouldEqual, 100)
		So(cfg.GetInt(cfgKey(1)), ShouldEqual, 101)
		So(cfg.GetInt(cfgKey(2)), ShouldEqual, 102)

		wg := sync.WaitGroup{}
		wg.Add(1)

		watch := cfg.ReloadWatch()
		go func() {
			defer wg.Done()

			<-watch
		}()

		<-time.After(time.Second)
		So(writeJson(filenames[0], getCfgMapFixture(2, 0)), ShouldBeNil)

		wg.Wait()
		So(cfg.GetInt(cfgKey(0)), ShouldEqual, 200)
	})
}

func writeJson(fn string, data map[string]interface{}) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	return json.NewEncoder(f).Encode(data)
}

func getCfgMapFixture(group, index int) map[string]interface{} {
	return map[string]interface{}{
		cfgKey(index): (group * 100) + index,
	}
}

func cfgKey(index int) string {
	return fmt.Sprintf("cfgtkey%d", index)
}
