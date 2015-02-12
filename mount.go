// Copyright 2015 Google Inc. All Rights Reserved.
// Author: jacobsa@google.com (Aaron Jacobs)

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jacobsa/gcsfuse/fs"
	"github.com/jacobsa/gcsfuse/fuseutil"
	"golang.org/x/net/context"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s [flags] <mount-point>\n", os.Args[0])
	flag.PrintDefaults()
}

var fBucketName = flag.String("bucket", "", "Name of GCS bucket to mount.")

func getBucketName() string {
	s := *fBucketName
	if s == "" {
		fmt.Println("You must set -bucket.")
		os.Exit(1)
	}

	return s
}

func main() {
	// Set up flags.
	flag.Usage = usage
	flag.Parse()

	// Enable debugging, if requested.
	initDebugging()

	// Grab the mount point.
	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	mountPoint := flag.Arg(0)

	// Set up a GCS connection.
	log.Println("Initializing GCS connection.")
	conn, err := getConn()
	if err != nil {
		log.Fatal("Couldn't get GCS connection: ", err)
	}

	// Create a file system.
	fileSystem, err := fs.NewFuseFS(conn.GetBucket(getBucketName()))
	if err != nil {
		log.Fatal("fs.NewFuseFS:", err)
	}

	// Mount the file system.
	mountedFS := fuseutil.MountFileSystem(mountPoint, fileSystem)

	if err := mountedFS.WaitForReady(context.Background()); err != nil {
		log.Fatal("MountedFileSystem.WaitForReady:", err)
	}

	log.Println("File system has been successfully mounted.")

	// Wait for it to be unmounted.
	if err := mountedFS.Join(context.Background()); err != nil {
		log.Fatal("MountedFileSystem.Join:", err)
	}

	log.Println("Successfully unmounted.")
}