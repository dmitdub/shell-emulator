package vfs

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// Save writes the tree to a CSV file.
//
// The file is overwritten. Nodes are written in ascending order
// of their identifiers. Save returns an error if the file cannot
// be created or written.
func Save(path string, root *Node) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writeHeader(writer); err != nil {
		return err
	}

	nodes := collectNodes(root)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].node.ID < nodes[j].node.ID
	})

	for _, node := range nodes {
		if err := writeRow(writer, node); err != nil {
			return err
		}
	}
	return writer.Error()
}

// writeHeader writes the CSV header row.
func writeHeader(w *csv.Writer) error {
	return w.Write([]string{"id", "parent", "name", "type", "content"})
}

// collectNodes returns all nodes of the tree in depth-first order.
//
// Each node is paired with the identifier of its parent.
func collectNodes(root *Node) []entry {
	var entries []entry
	walk(root, rootParent, &entries)
	return entries
}

// walk appends the node and all its descendants to entries.
func walk(node *Node, parent int, entries *[]entry) {
	*entries = append(*entries, entry{node: node, parent: parent})
	for _, child := range node.Children {
		walk(child, node.ID, entries)
	}
}

// entry pairs a node with its parent identifier.
type entry struct {
	node   *Node
	parent int
}

// writeRow writes a single node as a CSV row.
func writeRow(w *csv.Writer, e entry) error {
	kind := "file"
	if e.node.Dir {
		kind = "dir"
	}

	record := []string{
		strconv.Itoa(e.node.ID),
		strconv.Itoa(e.parent),
		e.node.Name,
		kind,
		e.node.Content,
	}
	if err := w.Write(record); err != nil {
		return fmt.Errorf("write node %d: %w", e.node.ID, err)
	}
	return nil
}
