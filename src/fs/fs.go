package fs

import (
	"github.com/hanwen/go-fuse/v2/fs"
	"path/filepath"
	"strings"
	"syscall"
)

type FileSystem struct {
	allow map[string]struct{}
	root  *fs.LoopbackRoot
}

func (m *FileSystem) CheckAccess(path string) bool {
	if path == "" {
		return true
	}
	var prefix string
	index := strings.IndexRune(path, '/')
	if index == -1 {
		prefix = path
	} else {
		prefix = path[:index]
	}
	_, ok := m.allow[prefix]
	return ok
}

func (m *FileSystem) RootNode() fs.InodeEmbedder {
	return m.root.RootNode
}

func (m *FileSystem) Allow(path string) {
	m.allow[path] = struct{}{}
}

type Node struct {
	fs.LoopbackNode
	fileSystem *FileSystem
}

func NewFileSystem(root string) (*FileSystem, error) {
	var err error
	if root, err = filepath.Abs(root); err != nil {
		return nil, err
	}

	var st syscall.Stat_t
	if err = syscall.Stat(root, &st); err != nil {
		return nil, err
	}

	fileSystem := &FileSystem{
		allow: make(map[string]struct{}),
	}

	var lb fs.LoopbackRoot
	lb.Path = root
	lb.Dev = st.Dev
	lb.NewNode = func(rootData *fs.LoopbackRoot, parent *fs.Inode, name string, st *syscall.Stat_t) fs.InodeEmbedder {
		return &Node{
			LoopbackNode: fs.LoopbackNode{
				RootData: &lb,
			},
			fileSystem: fileSystem,
		}
	}
	lb.RootNode = lb.NewNode(&lb, nil, "", &st)

	fileSystem.root = &lb
	return fileSystem, nil
}
