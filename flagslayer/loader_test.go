package flagslayer

import (
	"flag"
	"os"
	"testing"
	"time"

	. "github.com/fzerorubigd/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestYamlLoader(t *testing.T) {
	Convey("Load flag data in config", t, func() {
		o := New()
		flagset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		layer := NewFlagLayer(flagset)
		o.AddLayer(layer)
	})

	Convey("Load flag data in config default", t, func() {
		o := New()
		layer := NewFlagLayer(nil)
		layer.SetBool("bool", "bool", false, "usage")
		layer.SetDuration("duration", "duration", time.Minute, "usage")
		layer.SetInt64("int", "int", 1, "usage")
		layer.SetString("str", "str", "test", "usage")
		o.AddLayer(layer)

		So(o.GetBool("bool"), ShouldBeFalse)
		So(o.GetDuration("duration"), ShouldEqual, time.Minute)
		So(o.GetInt64("int"), ShouldEqual, 1)
		So(o.GetString("str"), ShouldEqual, "test")
	})

	Convey("Load flag data in config with mock flagset", t, func() {
		o := New()
		flagset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		layer := NewFlagLayer(flagset)
		layer.SetBool("bool", "bool", false, "usage")
		layer.SetDuration("duration", "duration", time.Minute, "usage")
		layer.SetInt64("int", "int", 1, "usage")
		layer.SetString("str", "str", "test", "usage")
		// Mock the layer
		args := []string{"-bool=true", "-duration=1h2m3s", "-int=22", "-str=stringtest"}
		flagset.Parse(args)

		o.AddLayer(layer)
		So(o.GetBool("bool"), ShouldBeTrue)
		d, _ := time.ParseDuration("1h2m3s")
		So(o.GetDuration("duration"), ShouldEqual, d)
		So(o.GetInt64("int"), ShouldEqual, 22)
		So(o.GetString("str"), ShouldEqual, "stringtest")
	})

	Convey("Load flag data in config with mock flagset", t, func() {
		o := New()
		flagset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

		layer := NewFlagLayer(flagset)
		layer.SetDelimiter("")
		So(layer.GetDelimiter(), ShouldEqual, ".")

		layer.SetBool("bool-nested", "bool", false, "usage")
		layer.SetDelimiter("-")
		// Mock the layer
		args := []string{"-bool=true"}
		flagset.Parse(args)

		o.AddLayer(layer)
		So(o.GetBool("bool.nested"), ShouldBeTrue)
	})
}
