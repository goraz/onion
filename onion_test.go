package onion

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type layerMock struct {
	data map[string]interface{}
}

type anonNested struct {
	Key4      int64
	Key5      int64
	Key5Again int `onion:"key5"`
}

type nested struct {
	N0 string
	N1 int
	N2 bool
}

type structExample struct {
	Key0     int
	Universe string `onion:"key1"`
	Key2     string
	Key3     bool

	anonNested
	nested `onion:"nested"`

	Another nested `onion:"nested"`
}

func (lm layerMock) Load() (map[string]interface{}, error) {
	return lm.data, nil
}

func getMap(prefix string, s ...interface{}) map[string]interface{} {
	tmp := make(map[string]interface{})
	for i := range s {
		tmp[fmt.Sprintf("%s%d", prefix, i)] = s[i]
	}
	return tmp
}

func TestOnion(t *testing.T) {
	Convey("Onion basic functionality", t, func() {
		lm := &layerMock{}
		lm.data = getMap("key", 42, "universe", "answer", true, float32(20.88), float64(200), int64(100))
		lm.data["nested"] = getMap("n", "a", 99, true)
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

		lm.data["yes"] = t1
		lm.data["slice1"] = []string{"a", "b", "c"}
		lm.data["slice2"] = []interface{}{"a", "b", "c"}
		lm.data["slice3"] = []interface{}{"a", "b", true}
		lm.data["slice4"] = []int{1, 2, 3}

		o := New()
		So(o.AddLayer(lm), ShouldBeNil)
		Convey("Get direct variable", func() {
			So(o.GetInt("key0"), ShouldEqual, 42)
			So(o.GetString("key1"), ShouldEqual, "universe")
			So(o.GetString("key2"), ShouldEqual, "answer")
			So(o.GetBool("key3"), ShouldBeTrue)
			So(o.GetInt("key4"), ShouldEqual, 20)
			So(o.GetInt("key5"), ShouldEqual, 200)
			So(o.GetInt("key6"), ShouldEqual, 100)

			So(o.GetInt64("key0"), ShouldEqual, 42)
			So(o.GetInt64("key4"), ShouldEqual, 20)
			So(o.GetInt64("key5"), ShouldEqual, 200)
			So(o.GetInt64("key6"), ShouldEqual, 100)
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

		Convey("delegate to structure", func() {
			s := structExample{}
			o.GetStruct("", &s)
			ex := structExample{
				Key0:     42,
				Universe: "universe",
				Key2:     "answer",
				Key3:     true,
				anonNested: anonNested{
					Key4:      20,
					Key5:      200,
					Key5Again: 200,
				},
				nested: nested{
					N0: "a",
					N1: 99,
					N2: true,
				},
				Another: nested{
					N0: "a",
					N1: 99,
					N2: true,
				},
			}
			So(reflect.DeepEqual(s, ex), ShouldBeTrue)
			var tmp []string
			o.GetStruct("", tmp)
			So(tmp, ShouldBeNil)
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
		lm1 := &layerMock{getMap("test", 1, true)}
		lm2 := &layerMock{getMap("test", 2, false)}
		os.Setenv("TEST0", "3")
		os.Setenv("TEST1", "True")
		os.Setenv("TEST2", "INVALIDBOOL")
		lm3 := NewEnvLayer("TEST0", "TEST1", "TEST2")

		o := New()
		o.AddLayer(lm1)
		So(o.GetInt64Default("test0", 0), ShouldEqual, 1)
		So(o.GetBoolDefault("test1", false), ShouldBeTrue)
		o.AddLayer(lm2)
		So(o.GetInt64Default("test0", 0), ShouldEqual, 2)
		So(o.GetBoolDefault("test1", true), ShouldBeFalse)
		o.AddLayer(lm3) // Special case in ENV loader
		So(o.GetInt64Default("test0", 0), ShouldEqual, 3)
		So(o.GetBoolDefault("test1", false), ShouldBeTrue)
		So(o.GetBoolDefault("test2", false), ShouldBeFalse)
	})

	Convey("test direct creation", t, func() {
		o := &Onion{}
		So(o.GetIntDefault("empty", 1000), ShouldEqual, 1000)
		lm := &layerMock{getMap("test", 1, true)}
		o1 := &Onion{}
		o1.AddLayer(lm)
		So(o1.GetIntDefault("test0", 0), ShouldEqual, 1)
	})
}
