package vfs

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// csvColumns is the expected number of columns in every row.
const csvColumns = 5

// rootParent marks the parent identifier of the root node.
const rootParent = -1

// Load reads a CSV file and builds the virtual file system.
//
// The file is opened read-only and is not modified. Load returns
// an error if the file cannot be read, the format is invalid, or
// the tree cannot be assembled.
func Load(path string) (*Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := readRecords(file)
	if err != nil {
		return nil, err
	}
	return buildTree(records)
}

// readRecords reads all data rows from the CSV file, skipping the header.
func readRecords(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = csvColumns

	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read rows: %w", err)
	}
	return records, nil
}

// row holds a parsed CSV row ready for tree assembly.
type row struct {
	id      int
	parent  int
	name    string
	dir     bool
	content string
}

// buildTree converts CSV rows into a tree of nodes.
func buildTree(records [][]string) (*Node, error) {
	rows, err := parseRows(records)
	if err != nil {
		return nil, err
	}

	nodes := make(map[int]*Node, len(rows))
	for _, r := range rows {
		nodes[r.id] = &Node{
			ID:      r.id,
			Name:    r.name,
			Dir:     r.dir,
			Content: r.content,
		}
	}

	var root *Node
	for _, r := range rows {
		node := nodes[r.id]
		if r.parent == rootParent {
			if root != nil {
				return nil, fmt.Errorf("multiple roots found")
			}
			root = node
			continue
		}
		parent, ok := nodes[r.parent]
		if !ok {
			return nil, fmt.Errorf("node %d: parent %d not found", r.id, r.parent)
		}
		parent.AddChild(node)
	}

	if root == nil {
		return nil, fmt.Errorf("root node not found")
	}
	return root, nil
}

// parseRows converts raw CSV records into typed rows.
func parseRows(records [][]string) ([]row, error) {
	rows := make([]row, 0, len(records))
	for i, record := range records {
		r, err := parseRow(record)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+2, err)
		}
		rows = append(rows, r)
	}
	return rows, nil
}

// parseRow converts a single CSV record into a typed row.
func parseRow(record []string) (row, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return row{}, fmt.Errorf("invalid id %q", record[0])
	}

	parent, err := strconv.Atoi(record[1])
	if err != nil {
		return row{}, fmt.Errorf("invalid parent %q", record[1])
	}

	dir, err := parseType(record[3])
	if err != nil {
		return row{}, err
	}

	return row{
		id:      id,
		parent:  parent,
		name:    record[2],
		dir:     dir,
		content: record[4],
	}, nil
}

// parseType reports whether the given type string means a directory.
func parseType(value string) (bool, error) {
	switch value {
	case "dir":
		return true, nil
	case "file":
		return false, nil
	default:
		return false, fmt.Errorf("invalid type %q", value)
	}
}
