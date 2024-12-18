// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	fusefs "github.com/hanwen/go-fuse/v2/fs"
	"log"
	"nix-store-fs/fs"
	"nix-store-fs/nix"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	daemonSocketPath := flag.String(
		"daemon-socket-path",
		"/nix/var/nix/daemon-socket/socket",
		"Path to the daemon socket",
	)

	storePath := flag.String(
		"store-path",
		"/nix/store",
		"Path to the NIX store",
	)

	flag.Parse()

	if flag.NArg() != 2 {
		_, _ = os.Stderr.WriteString("Usage:\nnix-store-fs [options] <item> <mountpoint>")
		os.Exit(1)
	}

	derivation := flag.Arg(0)
	mountPoint := flag.Arg(1)

	conn, err := nix.Connect(*daemonSocketPath)
	if err != nil {
		log.Fatalln(err)
	}

	refs, err := conn.GetAllReferences(filepath.Join(*storePath, derivation))
	if err != nil {
		log.Fatalln(err)
	}

	var opts fusefs.Options
	opts.FsName = "nix: " + derivation
	opts.Name = "nix-store-fs"

	fileSystem, err := fs.NewFileSystem(*storePath)
	if err != nil {
		log.Fatalln(err)
	}

	for _, ref := range refs {
		fileSystem.Allow(ref[11:])
	}

	opts.AllowOther = true

	println("mounting: ", mountPoint)

	server, err := fusefs.Mount(mountPoint, fileSystem.RootNode(), &opts)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmount if the program is terminated.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		println("SIGTERM: unmounting: ", mountPoint)
		_ = server.Unmount()
	}()

	server.Wait()

	println("fini")
}
