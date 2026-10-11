package vfs

// Default returns the default virtual file system.
//
// The default tree contains a root directory, two subdirectories,
// and one regular file. It is used by the vfs-init command and
// when no VFS path is provided at startup.
func Default() *Node {
	root := &Node{ID: 0, Name: "/", Dir: true}

	home := &Node{ID: 1, Name: "home", Dir: true}
	tmp := &Node{ID: 2, Name: "tmp", Dir: true}

	readme := &Node{
		ID:      3,
		Name:    "readme.txt",
		Dir:     false,
		Content: "Default VFS.",
	}

	home.AddChild(readme)
	root.AddChild(home)
	root.AddChild(tmp)

	return root
}
