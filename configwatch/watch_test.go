package configwatch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/goraz/onion"
)

type streamReload interface {
	onion.Layer
	Reload(context.Context, io.Reader, string) error
}

func mapToJson(m map[string]interface{}) io.Reader {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func TestRefWatch(t *testing.T) {
	Convey("Test Watch", t, func() {
		i32 := RegisterInt("i32", 32)
		i64 := RegisterInt64("i64", 64)
		d := RegisterDuration("du", time.Hour)
		b := RegisterBool("b", true)
		s := RegisterString("s", "string")
		f32 := RegisterFloat32("f32", 10.0)
		f64 := RegisterFloat64("f64", 100.0)

		data := make(map[string]interface{})
		l, err := onion.NewStreamLayer(mapToJson(data), "json", nil)
		So(err, ShouldBeNil)
		ly := l.(streamReload)
		o := onion.New(l)
		Watch(o)

		So(i32.Int(), ShouldEqual, 32)
		So(i64.Int64(), ShouldEqual, 64)
		So(d.Duration(), ShouldEqual, time.Hour)
		So(b.Bool(), ShouldEqual, true)
		So(s.String(), ShouldEqual, "string")
		So(f32.Float32(), ShouldEqual, 10.0)
		So(f64.Float64(), ShouldEqual, 100.0)

		data["i32"] = 132
		data["i64"] = 164
		data["du"] = "1m"
		data["b"] = false
		data["s"] = "diff"
		data["f32"] = 99.0
		data["f64"] = 9999.0

		c := o.ReloadWatch()
		So(ly.Reload(context.Background(), mapToJson(data), "json"), ShouldBeNil)
		// This is just a hack for the load to finish
		<-c
		time.Sleep(3 * time.Second)

		So(i32.Int(), ShouldEqual, 132)
		So(i64.Int64(), ShouldEqual, 164)
		So(d.Duration(), ShouldEqual, time.Minute)
		So(b.Bool(), ShouldEqual, false)
		So(s.String(), ShouldEqual, "diff")
		So(f32.Float32(), ShouldEqual, 99.0)
		So(f64.Float64(), ShouldEqual, 9999.0)

	})
}
