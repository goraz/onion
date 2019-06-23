package structlayer

import (
	"testing"

	"github.com/fzerorubigd/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewStructLayer(t *testing.T) {
	Convey("Test structure layer", t, func() {
		data := struct {
			Hi     string                 `mapstructure:"hi"`
			Int    int                    `mapstructure:"int"`
			Bool   bool                   `mapstructure:"bool"`
			Nested map[string]interface{} `mapstructure:"nested"`
		}{
			Hi:   "hello",
			Int:  100,
			Bool: true,
			Nested: map[string]interface{}{
				"in":  1000,
				"out": 88,
			},
		}

		l, err := NewStructLayer(data)
		So(err, ShouldBeNil)

		o := onion.New(l)
		So(o.GetString("hi"), ShouldEqual, "hello")

		_, err = NewStructLayer("string")
		So(err, ShouldNotBeNil)
	})
}
