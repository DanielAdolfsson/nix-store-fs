package fs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"path/filepath"
	"syscall"
)

type DirStream struct {
	fs.DirStream
	next       *fuse.DirEntry
	loadErrno  syscall.Errno
	fileSystem *FileSystem
	path       string
}

func (d *DirStream) load() {
	if d.next != nil {
		return
	}
	for d.DirStream.HasNext() {
		ent, errno := d.DirStream.Next()
		if errno != 0 {
			d.loadErrno = errno
			return
		}
		if d.path == "" {
			if !d.fileSystem.CheckAccess(ent.Name) {
				continue
			}
		}
		d.next = &ent
		return
	}
	d.loadErrno = 0
}

func (d *DirStream) HasNext() bool {
	d.load()
	return d.next != nil
}

func (d *DirStream) Next() (fuse.DirEntry, syscall.Errno) {
	d.load()
	if d.next != nil {
		next := d.next
		d.next = nil
		return *next, 0
	} else {
		return fuse.DirEntry{}, d.loadErrno
	}
}

func (d *DirStream) Readdirent(context.Context) (*fuse.DirEntry, syscall.Errno) {
	if !d.HasNext() {
		return nil, 0
	}
	de, errno := d.Next()
	return &de, errno
}

func NewDirStream(fileSystem *FileSystem, name string) (fs.DirStream, syscall.Errno) {
	inner, errno := fs.NewLoopbackDirStream(filepath.Join(fileSystem.root.Path, name))
	if errno != 0 {
		return nil, errno
	}
	return &DirStream{DirStream: inner, fileSystem: fileSystem, path: name}, 0
}
