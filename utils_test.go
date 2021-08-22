package onion

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeLayer(t *testing.T) {
	Convey("Test merge data", t, func() {
		data1 := map[string]interface{}{
			"1": 10,
			"2": 20,
			"3": 30,
		}

		data2 := map[string]interface{}{
			"1": 100,
		}

		data3 := map[string]interface{}{
			"1": 1000,
			"3": 3000,
		}

		merged := mergeLayersData(data1, data2, data3)
		So(merged["1"], ShouldEqual, 1000)
		So(merged["2"], ShouldEqual, 20)
		So(merged["3"], ShouldEqual, 3000)

		merged = mergeLayersData(data3, data2, data1)
		So(merged["1"], ShouldEqual, 10)
		So(merged["2"], ShouldEqual, 20)
		So(merged["3"], ShouldEqual, 30)

		merged = mergeLayersData(data3, data1, data2)
		So(merged["1"], ShouldEqual, 100)
		So(merged["2"], ShouldEqual, 20)
		So(merged["3"], ShouldEqual, 30)

		merged = mergeLayersData(data1)
		So(merged["1"], ShouldEqual, 10)
		So(merged["2"], ShouldEqual, 20)
		So(merged["3"], ShouldEqual, 30)

		merged = mergeLayersData()
		So(len(merged), ShouldEqual, 0)
	})
}
