package onion

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultLayer(t *testing.T) {
	Convey("Default layer basic test, some kind of coverage bitch :) ", t, func() {
		l := NewDefaultLayer()
		data, err := l.Load()
		So(err, ShouldBeNil)
		So(len(data), ShouldEqual, 0)

		err = l.SetDefault("layer1", 1)
		So(err, ShouldBeNil)
		err = l.SetDefault("layer1.layer2", 42)
		So(err, ShouldNotBeNil)

		err = l.SetDefault("p1.p2", true)
		So(err, ShouldBeNil)

		err = l.SetDefault("map", make(map[interface{}]interface{}))
		So(err, ShouldBeNil)

		err = l.SetDefault("map.in", "inside")
		So(err, ShouldBeNil)

		tmp := make(map[interface{}]interface{})
		tmp["data"] = "data"

		tmp["map3"] = make(map[interface{}]interface{})
		err = l.SetDefault("map.map2", tmp)
		So(err, ShouldBeNil)

		err = l.SetDefault("map.map2.another.int", 101)
		So(err, ShouldBeNil)

		err = l.SetDefault("map.map2.map3.int", 100)
		So(err, ShouldBeNil)

		err = l.SetDefault("map.map2.map3.int.other", 100)
		So(err, ShouldNotBeNil)

		tmp2 := make(map[string]interface{})
		tmp2["data"] = "data"

		tmp2["map3"] = make(map[string]interface{})
		err = l.SetDefault("map.map5", tmp2)
		So(err, ShouldBeNil)

		err = l.SetDefault("map.map5.map3.int", 100)
		So(err, ShouldBeNil)

		err = l.SetDefault("map.map5.map3.int.other", 100)
		So(err, ShouldNotBeNil)

		So(l.GetDelimiter(), ShouldEqual, ".")
		l.SetDelimiter("-")
		So(l.GetDelimiter(), ShouldEqual, "-")
		So(l.SetDefault("test-delim", 1), ShouldBeNil)
		l.SetDelimiter("")
		So(l.GetDelimiter(), ShouldEqual, ".")

		o := New()
		err = o.AddLayer(l)
		So(err, ShouldBeNil)

		So(o.GetInt("layer1"), ShouldEqual, 1)
		So(o.GetBool("p1.p2"), ShouldBeTrue)
		So(o.GetString("map.in"), ShouldEqual, "inside")
		So(o.GetString("map.map2.data"), ShouldEqual, "data")
		So(o.GetInt("map.map2.map3.int"), ShouldEqual, 100)
		So(o.GetInt("map.map2.another.int"), ShouldEqual, 101)
		So(o.GetInt("test.delim"), ShouldEqual, 1)
	})
}
