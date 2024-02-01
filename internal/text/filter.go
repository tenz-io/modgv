package text

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type adjacent map[string][]string

type edge [2]string

type edgeWriter struct {
	edges []edge
}

func Filter(in io.Reader, out io.Writer, end string) error {
	adj, err := readAsAdjacent(in)
	if err != nil {
		return err
	}
	root := adj.getRoot()
	if root == "" {
		return fmt.Errorf("no root found")
	}

	paths := adj.findAllPaths(root, end, []string{})
	if len(paths) == 0 {
		return fmt.Errorf("no path found")
	}

	var edges []edge
	for _, path := range paths {
		edges = append(edges, splitPathAsEdges(path)...)
	}

	return writeEdges(out, deduplicateEdges(edges))
}

func writeEdges(out io.Writer, edges []edge) error {
	var ew = edgeWriter{edges: edges}
	for _, e := range ew.edges {
		_, err := fmt.Fprintf(out, "%s %s\n", e[0], e[1])
		if err != nil {
			return fmt.Errorf("failed to write edge %v: %w", e, err)
		}
	}
	return nil
}

func readAsAdjacent(in io.Reader) (adjacent, error) {
	var (
		scanner = bufio.NewScanner(in)
		adj     = make(adjacent)
	)

	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		if l == "" {
			continue
		}

		parts := strings.Fields(l)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 words in line, but got %d: %s", len(parts), l)
		}

		src := parts[0]
		dst := parts[1]

		if _, ok := adj[src]; !ok {
			adj[src] = []string{dst}
		} else {
			adj[src] = append(adj[src], dst)
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(adj) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	return adj, nil
}

// findAllPaths returns all paths from src to dst
func (adj adjacent) findAllPaths(src, dst string, path []string) [][]string {
	path = append(path, src)

	if strings.Contains(src, dst) {
		return [][]string{path}
	}

	if _, ok := adj[src]; !ok {
		return [][]string{}
	}

	var paths [][]string
	for _, neighbor := range adj[src] {
		if !contains(path, neighbor) {
			newPaths := adj.findAllPaths(neighbor, dst, path)
			for _, np := range newPaths {
				paths = append(paths, np)
			}
		}
	}
	return paths
}

func (adj adjacent) getRoot() string {
	for src, _ := range adj {
		if len(strings.Split(src, "@")) == 1 {
			return src
		}
	}
	return ""
}

// splitPathAsEdges splits a path into a list of edges.
// eg: ['A', 'B', 'C'] -> [('A', 'B'), ('B', 'C')]
func splitPathAsEdges(path []string) []edge {
	var paths []edge
	if len(path) < 2 {
		return paths
	}

	for i := 0; i < len(path)-1; i++ {
		paths = append(paths, [2]string{path[i], path[i+1]})
	}
	return paths
}

// deduplicateEdges reduplicates edges in a list of edges.
// eg: [('A', 'B'), ('B', 'C'), ('A', 'B')] -> [('A', 'B'), ('B', 'C')]
func deduplicateEdges(edges []edge) []edge {
	var result []edge
	for _, path := range edges {
		if !containsEdge(result, path) {
			result = append(result, path)
		}
	}
	return result
}

func containsEdge(edges []edge, edge edge) bool {
	for _, p := range edges {
		if reflect.DeepEqual(p, edge) {
			return true
		}
	}
	return false
}

func contains(path []string, node string) bool {
	for _, p := range path {
		if p == node {
			return true
		}
	}
	return false
}
