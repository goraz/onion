package consulloader

import (
	"testing"

	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
	. "github.com/smartystreets/goconvey/convey"
	. "gopkg.in/fzerorubigd/onion.v3"
)

func TestConsulLoader(t *testing.T) {
	Convey("Extra env in config, when there is no server ", t, func() {
		client, err := api.NewClient(
			api.DefaultConfig(),
		)
		So(err, ShouldBeNil)

		o := New()
		layer := NewConsulLayer(client, "prefix")
		o.AddLazyLayer(layer)
		Convey("check when the server is gone", func() {
			So(o.GetString("a.data"), ShouldEqual, "")
			So(o.GetString(""), ShouldEqual, "")
		})

	})

	// TODO : I get a panic on consul NewTestServer and I have no time to investigate
	SkipConvey("consul in config", t, func() {
		// Create a test Consul server
		srv1, err := testutil.NewTestServer()
		if err != nil {
			t.Fatal(err)
		}
		defer srv1.Stop()
		client, err := api.NewClient(
			&api.Config{
				HttpClient: srv1.HTTPClient,
				Address:    srv1.HTTPAddr,
			},
		)
		So(err, ShouldBeNil)

		o := New()
		layer := NewConsulLayer(client, "prefix")
		o.AddLazyLayer(layer)
		Convey("check data from consul", func() {
			srv1.SetKV(t, "prefix/data/nested", []byte("TDN"))
			fmt.Println(string(srv1.GetKV(t, "prefix/data/nested")))
			So(o.GetString("data.nested"), ShouldEqual, "TDN")
		})

		Convey("check data not in consul", func() {
			So(o.GetString("not.valid.data"), ShouldEqual, "")
			So(o.GetString(""), ShouldEqual, "")
		})

	})

}
