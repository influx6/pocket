package guviews_test

import (
	"testing"

	"github.com/influx6/gu/gutrees"
	"github.com/influx6/gu/gutrees/attrs"
	"github.com/influx6/gu/gutrees/elems"
	"github.com/influx6/gu/guviews"
)

var success = "\u2713"
var failed = "\u2717"

var treeRenderlen = 275

type videoList []map[string]string

func (v videoList) Render(m ...string) gutrees.Markup {
	dom := elems.Div()
	for _, data := range v {
		gutrees.Augment(dom, elems.Video(
			attrs.Src(data["src"]),
			elems.Text(data["name"]),
		))
	}
	return dom
}

func TestView(t *testing.T) {
	videos := guviews.View("video-vabbs", videoList([]map[string]string{
		map[string]string{
			"src":  "https://youtube.com/xF5R32YF4",
			"name": "Joyride Lewis!",
		},
		map[string]string{
			"src":  "https://youtube.com/dox32YF4",
			"name": "Wonderlust Bombs!",
		},
	}))

	bo := videos.RenderHTML()

	if len(bo) != treeRenderlen {
		t.Fatalf("\t%s\t Rendered result with invalid length, expected %d but got %d -> \n %s", failed, treeRenderlen, len(bo), bo)
	}

	t.Logf("\t%s\t Rendered result accurated with length %d", success, treeRenderlen)
}
