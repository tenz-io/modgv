package text

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// adjacent is a map of nodes to their adjacent nodes.
// eg: {"A": ["B", "C", "E"], "B": ["D"], "D": ["E"], "E": ["F"], "G": ["H"]}
type adjacent map[string][]string

type edge [2]string

var (
	// isRootFunc is a function that returns true if a string is a root node.
	// used for testing purposes.
	isRootFunc = func(s string) bool {
		return !strings.Contains(s, "@")
	}
)

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

	paths := adj.findAllPaths(root, end)
	if len(paths) == 0 {
		return fmt.Errorf("no path found")
	}

	var edges []edge
	for _, path := range paths {
		edges = append(edges, splitPathAsEdges(path)...)
	}

	return writeEdges(out, deduplicate(edges))
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
func (adj adjacent) findAllPaths(src, dst string) [][]string {
	var result [][]string
	visited := make(map[string]bool)
	var path []string
	adj.dfs(src, dst, visited, &path, &result)
	return result
}

// dfs performs a depth-first search.
func (adj adjacent) dfs(src, dst string, visited map[string]bool, path *[]string, result *[][]string) {
	if visited[src] {
		return
	}

	if strings.Contains(src, dst) {
		*path = append(*path, src)
		newPath := make([]string, len(*path))
		copy(newPath, *path)
		*result = append(*result, newPath)
		*path = (*path)[:len(*path)-1]
		return
	}

	visited[src] = true
	*path = append(*path, src)

	for _, next := range adj[src] {
		adj.dfs(next, dst, visited, path, result)
	}

	visited[src] = false
	*path = (*path)[:len(*path)-1]
}

func (adj adjacent) getRoot() string {
	for src, _ := range adj {
		if isRootFunc(src) {
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

// deduplicate reduplicates edges in a list of edges.
// eg: [('A', 'B'), ('B', 'C'), ('A', 'B')] -> [('A', 'B'), ('B', 'C')]
func deduplicate[T comparable](list []T) []T {
	var result []T
	var seen = make(map[T]bool)
	for _, item := range list {
		if _, ok := seen[item]; !ok {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func contains(path []string, node string) bool {
	for _, p := range path {
		if p == node {
			return true
		}
	}
	return false
}
