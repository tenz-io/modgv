// This is a modified version of Modgraphviz created by the Go authors.
// Original Modgraphviz resides in the experimental repository.
// https://github.com/golang/exp/tree/master/cmd/modgraphviz

// modgv converts “go mod graph” output into Graphviz's DOT language,
// for use with Graphviz visualization and analysis tools like dot, dotty, and sccmap.
//
// Usage:
//
//	go mod graph | modgv > graph.dot
//	go mod graph | modgv | dot -Tpng -o graph.png
//
// modgv takes no options or arguments; it reads a graph in the format
// generated by “go mod graph” on standard input and writes DOT language
// on standard output.
//
// For each module, the node representing the greatest version (i.e., the
// version chosen by Go's minimal version selection algorithm) is colored green.
// Other nodes, which aren't in the final build list, are colored grey.
//
// See http://www.graphviz.org/doc/info/lang.html for details of the DOT language
// and http://www.graphviz.org/about/ for Graphviz itself.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tenz-io/modgv/internal/render"
)

const (
	banner = ` _ __ ___   ___   __| | __ ___   __
| '_ ' _ \ / _ \ / _' |/ _' \ \ / /
| | | | | | (_) | (_| | (_| |\ V / 
|_| |_| |_|\___/ \__,_|\__, | \_/  v{{VERSION}}
                       |___/       
forked from- https://github.com/lucasepe/modgv - By Luca Sepe`
)

var (
	version = "1.1.0"
	dstNode = os.Getenv("MODGV_DST_NODE")
)

func main() {

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 0 {
		usage()
	}

	if err := render.Render(os.Stdin, os.Stdout, dstNode); err != nil {
		exitOnErr(err)
	}
}

func usage() {
	printBanner()

	fmt.Fprintf(os.Stderr, "Converts 'go mod graph' output into Graphviz's DOT language.\n\n")

	fmt.Fprintf(os.Stderr, "  * takes no options or arguments\n")
	fmt.Fprintf(os.Stderr, "  * it reads the output generated by “go mod graph” on stdin\n")
	fmt.Fprintf(os.Stderr, "  * generates a DOT language and writes to stdout\n")
	fmt.Fprintf(os.Stderr, "\n")

	fmt.Fprintf(os.Stderr, "USAGE:\n\n")
	fmt.Fprintf(os.Stderr, "  go mod graph | %s | dot -Tpng -o graph.png\n\n", appName())

	fmt.Fprintf(os.Stderr, "For each module:\n")
	fmt.Fprintf(os.Stderr, "  * the node representing the greatest version (i.e., the version ")
	fmt.Fprintf(os.Stderr, "chosen by Go's MVS algorithm) is colored green\n")
	fmt.Fprintf(os.Stderr, "  * other nodes, which aren't in the final build list, are colored grey\n")

	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(2)
}

func printBanner() {
	str := strings.Replace(banner, "{{VERSION}}", version, 1)
	fmt.Fprint(os.Stderr, str)
	fmt.Fprint(os.Stderr, "\n\n")
}

func appName() string {
	return filepath.Base(os.Args[0])
}

// exitOnErr check for an error and eventually exit
func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
