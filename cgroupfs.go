package main

import (
	"fmt"
	"os"
	"time"

	"bazil.org/fuse"
	fusefs "bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"

	"github.com/chenchun/cgroupfs/fs"
)

func start(mountPoint, cgroupDir string) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- serve(mountPoint, cgroupDir)
	}()
	select {
	case <- time.After(100 * time.Millisecond):
		return nil
	case err := <- errChan:
		return err
	}
}

func serve(mountPoint, cgroupDir string) error {
	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("cgroupfs"),
		fuse.Subtype("cgroupfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("cgroup volume"),
		fuse.AllowOther(),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	var srv *fusefs.Server
	if os.Getenv("FUSE_DEBUG") != "" {
		srv = fusefs.New(c, &fusefs.Config{
			Debug: func(msg interface{}) {
				fmt.Printf("%s\n", msg)
			},
		})
	} else {
		srv = fusefs.New(c, nil)
	}

	err = srv.Serve(fs.FS{cgroupDir})
	if err != nil {
		return err
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}

	return nil
}

func stop(mountPoint string) error {
	if err := fuse.Unmount(mountPoint); err != nil {
		return fmt.Errorf("Error umounting %s: %s", mountPoint, err)
	}
	return nil
}

