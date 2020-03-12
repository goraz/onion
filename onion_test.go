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
		o := New(lm)
		Convey("Get direct variable", func() {
			So(o.GetInt("key0"), ShouldEqual, 42)
			So(o.GetString("key1"), ShouldEqual, "universe")
			So(o.GetString("key2"), ShouldEqual, "answer")
			So(o.GetBool("key3"), ShouldBeTrue)
			So(o.GetInt("key4"), ShouldEqual, 20)
			So(o.GetFloat32("key4"), ShouldEqual, 20.88)
			So(o.GetInt("key5"), ShouldEqual, 200)
			So(o.GetFloat64("key5"), ShouldEqual, 200.123)
			So(o.GetInt("key6"), ShouldEqual, 100)

			So(o.GetInt64("key0"), ShouldEqual, 42)
			So(o.GetInt64("key4"), ShouldEqual, 20)
			So(o.GetInt64("key5"), ShouldEqual, 200)
			So(o.GetInt64("key6"), ShouldEqual, 100)
			d, _ := time.ParseDuration("1h2m3s")
			So(o.GetDuration("durstring"), ShouldEqual, d)
			So(o.GetDuration("durstringinvalid"), ShouldEqual, 0)
			So(o.GetDuration("not-set-value"), ShouldEqual, 0)
			So(o.GetDuration("durint"), ShouldEqual, time.Duration(100000000))
			So(o.GetDuration("durint64"), ShouldEqual, time.Duration(100000000))
			So(o.GetDuration("booldur"), ShouldEqual, 0)
			So(o.GetDuration("dur"), ShouldEqual, time.Minute)
		})

		Convey("Get default value", func() {
			So(o.GetIntDefault("key1", 0), ShouldEqual, 0)
			So(o.GetIntDefault("nokey1", 0), ShouldEqual, 0)

			So(o.GetStringDefault("key0", ""), ShouldEqual, "")
			So(o.GetStringDefault("nokey0", ""), ShouldEqual, "")

			So(o.GetBoolDefault("key0", false), ShouldBeFalse)
			So(o.GetBoolDefault("nokey0", false), ShouldBeFalse)

			So(o.GetInt64Default("key1", 0), ShouldEqual, 0)
			So(o.GetInt64Default("nokey1", 0), ShouldEqual, 0)

			So(o.GetInt64Default("", 0), ShouldEqual, 0) // Empty key
			So(o.GetInt64Default("key3", 10000), ShouldEqual, 10000)

			So(o.GetFloat32Default("", 0), ShouldEqual, 0) // Empty key
			So(o.GetFloat32Default("key3", 10000), ShouldEqual, 10000)

			So(o.GetFloat64Default("", 0.123), ShouldEqual, 0.123) // Empty key
			So(o.GetFloat64Default("key3", 10000.123), ShouldEqual, 10000.123)
		})

		Convey("Get nested variable", func() {
			So(o.GetStringDefault("nested.n0", ""), ShouldEqual, "a")
			So(o.GetInt64Default("nested.n1", 0), ShouldEqual, 99)
			So(o.GetIntDefault("nested.n1", 0), ShouldEqual, 99)
			So(o.GetBoolDefault("nested.n2", false), ShouldEqual, true)

			So(o.GetIntDefault("yes.str1", 0), ShouldEqual, 1)
			So(o.GetStringDefault("yes.str2", ""), ShouldEqual, "hi")

			So(o.GetStringDefault("yes.nested.str2", ""), ShouldEqual, "hi")
			So(o.GetStringDefault("yes.what.n0", ""), ShouldEqual, "a")
		})

		Convey("Get nested default variable", func() {
			So(o.GetStringDefault("nested.n01", ""), ShouldEqual, "")
			So(o.GetStringDefault("key0.n01", ""), ShouldEqual, "")
			So(o.GetInt64Default("nested.n11", 0), ShouldEqual, 0)
			So(o.GetIntDefault("nested.n11", 0), ShouldEqual, 0)
			So(o.GetBoolDefault("nested.n21", false), ShouldEqual, false)

			So(o.GetStringDefault("yes.nested.no", "def"), ShouldEqual, "def")
			So(o.GetStringDefault("yes.nested.other.key", "def"), ShouldEqual, "def")
			So(o.GetStringDefault("yes.what.no", "def"), ShouldEqual, "def")
		})

		Convey("change delimiter", func() {
			So(o.GetDelimiter(), ShouldEqual, ".")
			o.SetDelimiter("/")
			So(o.GetDelimiter(), ShouldEqual, "/")
			Convey("get with modified delimiter", func() {
				So(o.GetStringDefault("nested/n0", ""), ShouldEqual, "a")
				So(o.GetInt64Default("nested/n1", 0), ShouldEqual, 99)
				So(o.GetIntDefault("nested/n1", 0), ShouldEqual, 99)
				So(o.GetBoolDefault("nested/n2", false), ShouldEqual, true)
				So(o.GetStringDefault("nested.n0", ""), ShouldEqual, "")
				So(o.GetInt64Default("nested.n1", 0), ShouldEqual, 0)
				So(o.GetIntDefault("nested.n1", 0), ShouldEqual, 0)
				So(o.GetBoolDefault("nested.n2", false), ShouldEqual, false)
				So(o.GetStringDefault("key0/n01", ""), ShouldEqual, "")
			})

			o.SetDelimiter("")
			So(o.GetDelimiter(), ShouldEqual, ".")
		})

		Convey("slice test", func() {
			So(reflect.DeepEqual(o.GetStringSlice("slice1"), []string{"a", "b", "c"}), ShouldBeTrue)
			So(reflect.DeepEqual(o.GetStringSlice("slice2"), []string{"a", "b", "c"}), ShouldBeTrue)
			So(o.GetStringSlice("slice3"), ShouldBeNil)
			So(o.GetStringSlice("notslice3"), ShouldBeNil)
			So(o.GetStringSlice("yes.str1"), ShouldBeNil)
			So(o.GetStringSlice("slice4"), ShouldBeNil)
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

	Convey("Test layers merge", t, func() {
		lm1 := NewMapLayer(getMap("test", 1, true))
		lm2 := NewMapLayer(getMap("test", 2, false))
		os.Setenv("TEST0", "3")
		os.Setenv("TEST1", "True")
		os.Setenv("TEST2", "INVALIDBOOL")
		lm3 := NewEnvLayer("_", "TEST0", "TEST1", "TEST2")

		o := New(lm1)
		o.AddLayers(lm2)

		mergedLayers := o.MergedLayersData()
		So(mergedLayers["test0"], ShouldEqual, 2)
		So(mergedLayers["test1"], ShouldBeFalse)

		o.AddLayers(lm3) // Special case in ENV loader
		mergedLayers = o.MergedLayersData()
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

		lm1 := NewMapLayer(getMap("test", 1, true))
		lm2 := NewMapLayer(getMap("test", 2, false))

		o := New(lm1)
		o.AddLayers(lm2)

		var conf, conf2 Config

		o.MergeAndDecode(&conf)
		So(conf.Test0, ShouldEqual, 2)
		So(conf.Test1, ShouldBeFalse)

		os.Setenv("TEST3", "ALongStringInSnakeCase")

		lm3 := NewEnvLayer("_", "TEST3")
		o.AddLayers(lm3) // Special case in ENV loader

		o.MergeAndDecode(&conf2)
		So(conf2.Test0, ShouldEqual, 2)
		So(conf2.Test1, ShouldBeFalse)
		So(conf2.Test2, ShouldBeFalse)
		So(conf2.Test3, ShouldEqual, "ALongStringInSnakeCase")
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
