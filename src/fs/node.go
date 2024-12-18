// SPDX-License-Identifier: MIT

package fs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"path/filepath"
	"syscall"
)

func (n *Node) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	if !n.fileSystem.CheckAccess(filepath.Join(n.Path(nil), name)) {
		return nil, syscall.ENOENT
	}
	return n.LoopbackNode.Lookup(ctx, name, out)
}

func (n *Node) Mknod(context.Context, string, uint32, uint32, *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOTSUP
}

func (n *Node) Mkdir(context.Context, string, uint32, *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOTSUP
}

func (n *Node) Rmdir(context.Context, string) syscall.Errno {
	return syscall.ENOTSUP
}

func (n *Node) Unlink(context.Context, string) syscall.Errno {
	return syscall.ENOTSUP
}

func (n *Node) Rename(context.Context, string, fs.InodeEmbedder, string, uint32) syscall.Errno {
	return syscall.ENOTSUP
}

func (n *Node) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {
	if !n.fileSystem.CheckAccess(filepath.Join(n.Path(nil), name)) {
		return nil, nil, 0, syscall.ENOENT
	}
	if flags&syscall.O_ACCMODE != syscall.O_RDONLY {
		return nil, nil, 0, syscall.ENOTSUP
	}
	return n.LoopbackNode.Create(ctx, name, flags, mode, out)
}

func (n *Node) Symlink(context.Context, string, string, *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOTSUP
}

func (n *Node) Link(context.Context, fs.InodeEmbedder, string, *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOTSUP
}

func (n *Node) Readlink(ctx context.Context) ([]byte, syscall.Errno) {
	if !n.fileSystem.CheckAccess(n.Path(nil)) {
		return nil, syscall.ENOENT
	}
	return n.LoopbackNode.Readlink(ctx)
}

func (n *Node) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	if !n.fileSystem.CheckAccess(n.Path(nil)) {
		return nil, 0, syscall.ENOENT
	}
	if flags&syscall.O_ACCMODE != syscall.O_RDONLY {
		return nil, 0, syscall.ENOTSUP
	}
	return n.LoopbackNode.Open(ctx, flags)
}

func (n *Node) OpendirHandle(context.Context, uint32) (fs.FileHandle, uint32, syscall.Errno) {
	if !n.fileSystem.CheckAccess(n.Path(nil)) {
		return nil, 0, syscall.ENOENT
	}
	ds, errno := NewDirStream(n.fileSystem, n.Path(nil))
	if errno != 0 {
		return nil, 0, errno
	}
	return ds, 0, errno
}

func (n *Node) Readdir(context.Context) (fs.DirStream, syscall.Errno) {
	if !n.fileSystem.CheckAccess(n.Path(nil)) {
		return nil, syscall.ENOENT
	}
	return NewDirStream(n.fileSystem, n.Path(nil))
}

func (n *Node) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if !n.fileSystem.CheckAccess(n.Path(nil)) {
		return syscall.ENOENT
	}
	return n.LoopbackNode.Getattr(ctx, f, out)
}

func (n *Node) Setattr(context.Context, fs.FileHandle, *fuse.SetAttrIn, *fuse.AttrOut) syscall.Errno {
	return syscall.ENOTSUP
}

func (n *Node) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	if !n.fileSystem.CheckAccess(n.Path(nil)) {
		return 0, syscall.ENOENT
	}
	return n.LoopbackNode.Getxattr(ctx, attr, dest)
}

func (n *Node) Setxattr(context.Context, string, []byte, uint32) syscall.Errno {
	return syscall.ENOTSUP
}

func (n *Node) Removexattr(context.Context, string) syscall.Errno {
	return syscall.ENOTSUP
}

func (n *Node) CopyFileRange(
	context.Context, fs.FileHandle, uint64, *fs.Inode,
	fs.FileHandle, uint64, uint64, uint64,
) (uint32, syscall.Errno) {
	return 0, syscall.ENOTSUP
}
