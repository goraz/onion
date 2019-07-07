package etcdlayer

import (
	"context"
	"sync"
	"testing"

	"github.com/etcd-io/etcd/client"
	"github.com/goraz/onion"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewEtcdLayerContext(t *testing.T) {
	Convey("Test etcd layer", t, func() {
		endPoints := []string{"http://127.0.0.1:2379"}
		cl, err := client.New(client.Config{
			Endpoints: endPoints,
		})
		So(err, ShouldBeNil)
		api := client.NewKeysAPI(cl)
		_, err = api.Set(context.Background(), "/app/config", `{"hi": 100}`, nil)
		So(err, ShouldBeNil)

		l, err := NewEtcdLayer("/app/config", "json", endPoints, nil)
		So(err, ShouldBeNil)
		o := onion.New(l)
		w := sync.WaitGroup{}
		w.Add(1)
		watch := o.ReloadWatch()
		go func() {
			defer w.Done()
			<-watch
		}()
		So(o.GetInt("hi"), ShouldEqual, 100)
		_, err = api.Set(context.Background(), "/app/config", `{"hi": 200}`, nil)
		So(err, ShouldBeNil)

		// Wait for reload channel
		w.Wait()
		So(o.GetInt("hi"), ShouldEqual, 200)
	})

}
