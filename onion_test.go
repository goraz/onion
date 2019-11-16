package onion

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func getMap(prefix string, s ...interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	for i := range s {
		tmp[fmt.Sprintf("%s%d", prefix, i)] = s[i]
	}
	return tmp
}

type dummyWatch struct {
	data map[string]interface{}
	c    chan map[string]interface{}
}

func (d *dummyWatch) Load() map[string]interface{} {
	return d.data
}

func (d *dummyWatch) Watch() <-chan map[string]interface{} {
	return d.c
}

func newDummy(data map[string]interface{}) *dummyWatch {
	return &dummyWatch{
		data: data,
		c:    make(chan map[string]interface{}),
	}
}

func TestOnion(t *testing.T) {
	Convey("Onion basic functionality", t, func() {
		data := getMap("key", 42, "universe", "answer", true, float32(20.88), float64(200.123), int64(100))
		data["nested"] = getMap("n", "a", 99, true)

		t1 := make(map[interface{}]interface{})
		t1["str1"] = 1
		t1["str2"] = "hi"
		t1["other"] = struct{}{}
		t1["what"] = getMap("n", "a")
		t2 := make(map[interface{}]interface{})
		for k, v := range t1 {
			t2[k] = v
		}
		t1["nested"] = t2
		t2[true] = false

		data["yes"] = t1
		data["slice1"] = []string{"a", "b", "c"}
		data["slice2"] = []interface{}{"a", "b", "c"}
		data["slice3"] = []interface{}{"a", "b", true}
		data["slice4"] = []int{1, 2, 3}
		data["ignored"] = "ignore me"

		data["dur"] = time.Minute
		data["durstring"] = "1h2m3s"
		data["durstringinvalid"] = "ertyuiop"
		data["durint"] = 100000000
		data["durint64"] = int64(100000000)
		data["booldur"] = true

		lm := NewMapLayer(data)
		AddLayers(lm)
		Convey("Get direct variable", func() {
			So(GetInt("key0"), ShouldEqual, 42)
			So(GetString("key1"), ShouldEqual, "universe")
			So(GetString("key2"), ShouldEqual, "answer")
			So(GetBool("key3"), ShouldBeTrue)
			So(GetInt("key4"), ShouldEqual, 20)
			So(GetFloat32("key4"), ShouldEqual, 20.88)
			So(GetInt("key5"), ShouldEqual, 200)
			So(GetFloat64("key5"), ShouldEqual, 200.123)
			So(GetInt("key6"), ShouldEqual, 100)

			So(GetInt64("key0"), ShouldEqual, 42)
			So(GetInt64("key4"), ShouldEqual, 20)
			So(GetInt64("key5"), ShouldEqual, 200)
			So(GetInt64("key6"), ShouldEqual, 100)
			d, _ := time.ParseDuration("1h2m3s")
			So(GetDuration("durstring"), ShouldEqual, d)
			So(GetDuration("durstringinvalid"), ShouldEqual, 0)
			So(GetDuration("not-set-value"), ShouldEqual, 0)
			So(GetDuration("durint"), ShouldEqual, time.Duration(100000000))
			So(GetDuration("durint64"), ShouldEqual, time.Duration(100000000))
			So(GetDuration("booldur"), ShouldEqual, 0)
			So(GetDuration("dur"), ShouldEqual, time.Minute)
		})

		Convey("Get default value", func() {
			So(GetIntDefault("key1", 0), ShouldEqual, 0)
			So(GetIntDefault("nokey1", 0), ShouldEqual, 0)
			So(GetStringDefault("key0", ""), ShouldEqual, "")
			So(GetStringDefault("nokey0", ""), ShouldEqual, "")
			So(GetBoolDefault("key0", false), ShouldBeFalse)
			So(GetBoolDefault("nokey0", false), ShouldBeFalse)
			So(GetInt64Default("key1", 0), ShouldEqual, 0)
			So(GetInt64Default("nokey1", 0), ShouldEqual, 0)

			So(GetInt64Default("", 0), ShouldEqual, 0) // Empty key
			So(GetInt64Default("key3", 10000), ShouldEqual, 10000)
			So(GetFloat32Default("", 0), ShouldEqual, 0) // Empty key
			So(GetFloat32Default("key3", 10000), ShouldEqual, 10000)
			So(GetFloat64Default("", 0.123), ShouldEqual, 0.123) // Empty key
			So(GetFloat64Default("key3", 10000.123), ShouldEqual, 10000.123)

			So(GetDurationDefault("not-set-value", time.Minute), ShouldEqual, time.Minute)
		})

		Convey("Get nested variable", func() {
			So(GetStringDefault("nested.n0", ""), ShouldEqual, "a")
			So(GetInt64Default("nested.n1", 0), ShouldEqual, 99)
			So(GetIntDefault("nested.n1", 0), ShouldEqual, 99)
			So(GetBoolDefault("nested.n2", false), ShouldEqual, true)
			So(GetIntDefault("yes.str1", 0), ShouldEqual, 1)
			So(GetStringDefault("yes.str2", ""), ShouldEqual, "hi")
			So(GetStringDefault("yes.nested.str2", ""), ShouldEqual, "hi")
			So(GetStringDefault("yes.what.n0", ""), ShouldEqual, "a")
		})

		Convey("Get nested default variable", func() {
			So(GetStringDefault("nested.n01", ""), ShouldEqual, "")
			So(GetStringDefault("key0.n01", ""), ShouldEqual, "")
			So(GetInt64Default("nested.n11", 0), ShouldEqual, 0)
			So(GetIntDefault("nested.n11", 0), ShouldEqual, 0)
			So(GetBoolDefault("nested.n21", false), ShouldEqual, false)

			So(GetStringDefault("yes.nested.no", "def"), ShouldEqual, "def")
			So(GetStringDefault("yes.nested.other.key", "def"), ShouldEqual, "def")
			So(GetStringDefault("yes.what.no", "def"), ShouldEqual, "def")
		})

		Convey("change delimiter", func() {
			So(GetDelimiter(), ShouldEqual, ".")
			SetDelimiter("/")
			So(GetDelimiter(), ShouldEqual, "/")
			Convey("get with modified delimiter", func() {
				So(GetStringDefault("nested/n0", ""), ShouldEqual, "a")
				So(GetInt64Default("nested/n1", 0), ShouldEqual, 99)
				So(GetIntDefault("nested/n1", 0), ShouldEqual, 99)
				So(GetBoolDefault("nested/n2", false), ShouldEqual, true)
				So(GetStringDefault("nested.n0", ""), ShouldEqual, "")
				So(GetInt64Default("nested.n1", 0), ShouldEqual, 0)
				So(GetIntDefault("nested.n1", 0), ShouldEqual, 0)
				So(GetBoolDefault("nested.n2", false), ShouldEqual, false)
				So(GetStringDefault("key0/n01", ""), ShouldEqual, "")
			})

			SetDelimiter("")
			So(GetDelimiter(), ShouldEqual, ".")
		})

		Convey("slice test", func() {
			So(reflect.DeepEqual(GetStringSlice("slice1"), []string{"a", "b", "c"}), ShouldBeTrue)
			So(reflect.DeepEqual(GetStringSlice("slice2"), []string{"a", "b", "c"}), ShouldBeTrue)
			So(GetStringSlice("slice3"), ShouldBeNil)
			So(GetStringSlice("notslice3"), ShouldBeNil)
			So(GetStringSlice("yes.str1"), ShouldBeNil)
			So(GetStringSlice("slice4"), ShouldBeNil)
		})
	})

	Convey("Test layer overwrite", t, func() {
		lm1 := NewMapLayer(getMap("test", 1, true))
		lm2 := NewMapLayer(getMap("test", 2, false))
		os.Setenv("TEST0", "3")
		os.Setenv("TEST1", "True")
		os.Setenv("TEST2", "INVALIDBOOL")
		lm3 := NewEnvLayer("_", "TEST0", "TEST1", "TEST2")

		o := New(lm1)
		So(o.GetInt64Default("test0", 0), ShouldEqual, 1)
		So(o.GetBoolDefault("test1", false), ShouldBeTrue)
		o.AddLayers(lm2)
		So(o.GetInt64Default("test0", 0), ShouldEqual, 2)
		So(o.GetBoolDefault("test1", true), ShouldBeFalse)
		o.AddLayers(lm3) // Special case in ENV loader
		So(o.GetInt64Default("test0", 0), ShouldEqual, 3)
		So(o.GetBoolDefault("test1", false), ShouldBeTrue)
		So(o.GetBoolDefault("test2", false), ShouldBeFalse)
	})

	Convey("test direct creation", t, func() {
		o := &Onion{}
		So(o.GetIntDefault("empty", 1000), ShouldEqual, 1000)
		lm := NewMapLayer(getMap("test", 1, true))
		o1 := &Onion{}
		o1.AddLayers(lm)
		So(o1.GetIntDefault("test0", 0), ShouldEqual, 1)
	})
}

func BenchmarkOion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benconion.GetInt("key0")
		benconion.GetString("key1")
		benconion.GetString("key2")
		benconion.GetBool("key3")
		benconion.GetInt("key4")
		benconion.GetInt("key5")
		benconion.GetInt("key6")

		benconion.GetInt64("key0")
		benconion.GetInt64("key4")
		benconion.GetInt64("key5")
		benconion.GetInt64("key6")
		benconion.GetDuration("durstring")
		benconion.GetDuration("durstringinvalid")
		benconion.GetDuration("not-set-value")
		benconion.GetDuration("durint")
		benconion.GetDuration("durint64")
		benconion.GetDuration("booldur")
		benconion.GetDuration("dur")
	}
}

var benconion = New()

func init() {

	data := getMap("key", 42, "universe", "answer", true, float32(20.88), float64(200), int64(100))
	data["nested"] = getMap("n", "a", 99, true)
	t1 := make(map[interface{}]interface{})
	t1["str1"] = 1
	t1["str2"] = "hi"
	t1["other"] = struct{}{}
	t1["what"] = getMap("n", "a")
	t2 := make(map[interface{}]interface{})
	for k, v := range t1 {
		t2[k] = v
	}
	t1["nested"] = t2
	t2[true] = false

	data["yes"] = t1
	data["slice1"] = []string{"a", "b", "c"}
	data["slice2"] = []interface{}{"a", "b", "c"}
	data["slice3"] = []interface{}{"a", "b", true}
	data["slice4"] = []int{1, 2, 3}
	data["ignored"] = "ignore me"

	data["dur"] = time.Minute
	data["durstring"] = "1h2m3s"
	data["durstringinvalid"] = "ertyuiop"
	data["durint"] = 100000000
	data["durint64"] = int64(100000000)
	data["booldur"] = true

	benconion.AddLayers(NewMapLayer(data))
}

func TestLayersData(t *testing.T) {
	Convey("Test merge", t, func() {
		o := New()
		l1 := NewMapLayer(map[string]interface{}{
			"key1": 1,
			"key2": 2,
		})
		l2 := NewMapLayer(map[string]interface{}{
			"key1": 10,
			"key3": 3,
		})

		o.AddLayers(l1, l2)

		ret := []map[string]interface{}{
			map[string]interface{}{
				"key1": 1,
				"key2": 2,
			},
			map[string]interface{}{
				"key1": 10,
				"key3": 3,
			},
		}

		So(o.LayersData(), ShouldResemble, ret)
	})
}

func TestWatch(t *testing.T) {
	Convey("Test watch", t, func() {
		data := map[string]interface{}{
			"k1": "10.0",
		}

		l := newDummy(data)
		o := New(l)
		ch := o.ReloadWatch()
		So(o.GetFloat32("k1"), ShouldEqual, 10.0)
		data["k1"] = "100.0"
		l.c <- data
		<-ch
		So(o.GetFloat32("k1"), ShouldEqual, 100.0)
	})
}
