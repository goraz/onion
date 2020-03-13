package onionwriter

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/goraz/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSerializeOnion(t *testing.T) {
	Convey("Serialize the onion into json", t, func() {
		m1 := map[string]interface{}{
			"key1": "HALLO",
			"key2": "JA",
			"key3": 100,
			"key4": []string{},
		}

		m2 := map[string]interface{}{
			"key1": "HELLO",
			"key4": []string{"a", "b", "c"},
		}

		o := onion.New(onion.NewMapLayer(m1), onion.NewMapLayer(m2))
		buf := &bytes.Buffer{}
		So(SerializeOnion(o, buf), ShouldBeNil)
		data := make(map[string]interface{})
		So(json.Unmarshal(buf.Bytes(), &data), ShouldBeNil)
		o2 := onion.New(onion.NewMapLayer(data))

		So(o2.GetString("key1"), ShouldEqual, o.GetString("key1"))
		So(o2.GetString("key2"), ShouldEqual, o.GetString("key2"))
		So(o2.GetInt("key3"), ShouldEqual, o.GetInt("key3"))
		So(o2.GetStringSlice("key4"), ShouldResemble, o.GetStringSlice("key4"))
	})

	Convey("Test layers merge", t, func() {
		lm1 := onion.NewMapLayer(onion.getMap("test", 1, true))
		lm2 := onion.NewMapLayer(onion.getMap("test", 2, false))
		os.Setenv("TEST0", "3")
		os.Setenv("TEST1", "True")
		os.Setenv("TEST2", "INVALIDBOOL")
		lm3 := onion.NewEnvLayer("_", "TEST0", "TEST1", "TEST2")

		o := onion.New(lm1)
		o.AddLayers(lm2)

		mergedLayers := MergeLayersOnion(o)
		So(mergedLayers["test0"], ShouldEqual, 2)
		So(mergedLayers["test1"], ShouldBeFalse)

		o.AddLayers(lm3) // Special case in ENV loader
		mergedLayers = MergeLayersOnion(o)
		So(mergedLayers["test0"], ShouldEqual, "3")
		So(mergedLayers["test1"], ShouldEqual, "True")
		So(mergedLayers["test2"], ShouldEqual, "INVALIDBOOL")
	})

	Convey("Test layers merge and decode to struct", t, func() {
		type Config struct {
			Test0 int64
			Test1 bool
			Test2 bool
			Test3 string
		}

		lm1 := onion.NewMapLayer(onion.getMap("test", 1, true))
		lm2 := onion.NewMapLayer(onion.getMap("test", 2, false))

		o := onion.New(lm1)
		o.AddLayers(lm2)

		var conf, conf2 Config

		DecodeOnion(o, &conf)
		So(conf.Test0, ShouldEqual, 2)
		So(conf.Test1, ShouldBeFalse)

		os.Setenv("TEST3", "ALongStringInSnakeCase")

		lm3 := onion.NewEnvLayer("_", "TEST3")
		o.AddLayers(lm3) // Special case in ENV loader

		o.MergeAndDecode(o, &conf2)
		So(conf2.Test0, ShouldEqual, 2)
		So(conf2.Test1, ShouldBeFalse)
		So(conf2.Test2, ShouldBeFalse)
		So(conf2.Test3, ShouldEqual, "ALongStringInSnakeCase")
	})

}
