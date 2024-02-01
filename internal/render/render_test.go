package render

import (
	"bytes"
	"testing"
)

func TestRender(t *testing.T) {
	t.Run("without end node", func(t *testing.T) {
		out := &bytes.Buffer{}
		in := bytes.NewBuffer([]byte(`
test.com/A@v1.0.0 test.com/B@v1.2.3
test.com/B@v1.0.0 test.com/C@v4.5.6
`))
		if err := Render(in, out, ""); err != nil {
			t.Fatal(err)
		}

		gotGraph := out.String()
		wantGraph := `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/A@v1.0.0" -> "test.com/B@v1.2.3"
	"test.com/B@v1.0.0" -> "test.com/C@v4.5.6"
	"test.com/A@v1.0.0" [style = filled, fillcolor = green]
	"test.com/B@v1.2.3" [style = filled, fillcolor = green]
	"test.com/C@v4.5.6" [style = filled, fillcolor = green]
	"test.com/B@v1.0.0" [style = filled, fillcolor = gray]
}
`
		if gotGraph != wantGraph {
			t.Fatalf("\ngot: %s\nwant: %s", gotGraph, wantGraph)
		}
	})

	t.Run("with end node", func(t *testing.T) {
		out := &bytes.Buffer{}
		in := bytes.NewBuffer([]byte(`
test.com/A test.com/B@v1.2.3
test.com/A test.com/C@v4.5.7
test.com/B@v1.2.3 test.com/C@v4.5.6
test.com/C@v4.5.6 test.com/D@v1.2.3
test.com/C@v4.5.7 test.com/D@v1.2.4
test.com/C@v4.5.7 test.com/E@v1.2.4
`))
		if err := Render(in, out, "test.com/D"); err != nil {
			t.Fatal(err)
		}

		gotGraph := out.String()
		wantGraph := `digraph gomodgraph {
	node [ shape=rectangle fontsize=12 ]
	"test.com/A@v1.0.0" -> "test.com/B@v1.2.3"
	"test.com/B@v1.0.0" -> "test.com/C@v4.5.6"
	"test.com/A@v1.0.0" [style = filled, fillcolor = green]
	"test.com/B@v1.2.3" [style = filled, fillcolor = green]
	"test.com/C@v4.5.6" [style = filled, fillcolor = green]
	"test.com/B@v1.0.0" [style = filled, fillcolor = gray]
}
`
		if gotGraph != wantGraph {
			t.Fatalf("\ngot: %s\nwant: %s", gotGraph, wantGraph)
		}
	})
}
