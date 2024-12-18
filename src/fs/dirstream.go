// SPDX-License-Identifier: MIT

package fs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"path/filepath"
	"sync"
	"syscall"
)

type DirStream struct {
	lo         fs.DirStream
	mu         sync.Mutex
	todo       *fuse.DirEntry
	todoErrno  syscall.Errno
	fileSystem *FileSystem
	path       string
}

func (d *DirStream) load() {
	if d.todo != nil {
		return
	}
	for d.lo.HasNext() {
		ent, errno := d.lo.Next()
		if errno != 0 {
			d.todoErrno = errno
			return
		}
		if d.path == "" {
			if !d.fileSystem.CheckAccess(ent.Name) {
				continue
			}
		}
		d.todo = &ent
		return
	}
	d.todoErrno = 0
}

func (d *DirStream) HasNext() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.load()
	return d.todo != nil
}

func (d *DirStream) Next() (fuse.DirEntry, syscall.Errno) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.load()
	if d.todo != nil {
		next := d.todo
		d.todo = nil
		return *next, 0
	} else {
		return fuse.DirEntry{}, d.todoErrno
	}
}

func (d *DirStream) Close() {
	d.lo.Close()
}

func (d *DirStream) Readdirent(context.Context) (*fuse.DirEntry, syscall.Errno) {
	if !d.HasNext() {
		return nil, 0
	}
	de, errno := d.Next()
	return &de, errno
}

func (d *DirStream) Seekdir(ctx context.Context, off uint64) syscall.Errno {
	d.mu.Lock()
	defer d.mu.Unlock()
	errno := d.lo.(fs.FileSeekdirer).Seekdir(ctx, off)
	if errno != 0 {
		return errno
	}
	d.todo = nil
	d.todoErrno = 0
	d.load()
	return 0
}

func (d *DirStream) Fsyncdir(ctx context.Context, flags uint32) syscall.Errno {
	return d.lo.(fs.FileFsyncdirer).Fsyncdir(ctx, flags)
}

func (d *DirStream) Releasedir(context.Context, uint32) {
	d.lo.Close()
}

func NewDirStream(fileSystem *FileSystem, name string) (fs.DirStream, syscall.Errno) {
	inner, errno := fs.NewLoopbackDirStream(filepath.Join(fileSystem.root.Path, name))
	if errno != 0 {
		return nil, errno
	}
	return &DirStream{lo: inner, fileSystem: fileSystem, path: name}, 0
}
