package text

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func Test_readAsAdjacent(t *testing.T) {
	type args struct {
		in io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    adjacent
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				in: func() io.Reader {
					return strings.NewReader(`
											A B
											A C
											B D
											`)
				}(),
			},
			want: adjacent{
				"A": {"B", "C"},
				"B": {"D"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readAsAdjacent(tt.args.in)
			t.Logf("err: %v", err)
			t.Logf("got: %v", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("readAsAdjacent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAsAdjacent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdjacent_findAllPaths(t *testing.T) {
	type args struct {
		node string
		dst  string
		path []string
	}
	tests := []struct {
		name string
		adj  adjacent
		args args
		want [][]string
	}{
		{
			name: "find all paths from A to D",
			adj: adjacent{
				"A": {"B", "C", "D"},
				"B": {"D"},
				"C": {"D", "E"},
			},
			args: args{
				node: "A",
				dst:  "D",
				path: []string{},
			},
			want: [][]string{
				{"A", "B", "D"},
				{"A", "C", "D"},
				{"A", "D"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.adj.findAllPaths(tt.args.node, tt.args.dst, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findAllPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitPath(t *testing.T) {
	type args struct {
		path []string
	}
	tests := []struct {
		name string
		args args
		want []edge
	}{
		{
			name: "split path",
			args: args{
				path: []string{"A", "B", "C"},
			},
			want: []edge{
				{"A", "B"},
				{"B", "C"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitPathAsEdges(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deduplicateEdges(t *testing.T) {
	type args struct {
		edges []edge
	}
	tests := []struct {
		name string
		args args
		want []edge
	}{
		{
			name: "deduplicate edges",
			args: args{
				edges: []edge{
					{"A", "B"},
					{"B", "C"},
					{"A", "B"},
				},
			},
			want: []edge{
				{"A", "B"},
				{"B", "C"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deduplicateEdges(tt.args.edges); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deduplicateEdges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args struct {
		in  io.Reader
		end string
	}
	type setup func()
	tests := []struct {
		name    string
		args    args
		setup   setup
		wantOut string
		wantErr bool
	}{
		{
			name: "filter input",
			args: args{
				in: func() io.Reader {
					return strings.NewReader(`
						A B
						A C
						B D
						E F
					`)
				}(),
				end: "D",
			},
			setup: func() {
				isRootFunc = func(s string) bool {
					return s == "A"
				}
			},
			wantOut: "A B\nB D\n",
			wantErr: false,
		},
		{
			name: "filter input 2",
			args: args{
				in: func() io.Reader {
					return strings.NewReader(`
						A B
						B C
						B C
						B D
						D E
						A E
						E F
						A F
						G H
					`)
				}(),
				end: "E",
			},
			setup: func() {
				isRootFunc = func(s string) bool {
					return s == "A"
				}
			},
			wantOut: "A B\nB D\nD E\nA E\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			out := &bytes.Buffer{}
			err := Filter(tt.args.in, out, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Filter() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
