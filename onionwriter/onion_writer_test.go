package onionwriter

import (
	"bytes"
	"encoding/json"
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
}
