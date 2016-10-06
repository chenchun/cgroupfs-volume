package main

import (
	"flag"
	"path/filepath"
	"os"

	"github.com/docker/go-plugins-helpers/volume"
)

var (
	flagRootPath = flag.String("root", filepath.Join(volume.DefaultDockerRootDirectory, "cgroupfs"), "root path of cgroupfs volume plugin")
)

const CidFile = "cidfile"

func main() {
	d := &volumeServer{volumes: make(map[string]*volume.Request)}
	h := volume.NewHandler(d)
	h.ServeTCP("cgroupfs_volume", ":8083")
}

type volumeServer struct {
	volumes map[string]*volume.Request
}

func (s *volumeServer) Create(req volume.Request) volume.Response {
	s.volumes[req.Name] = &req
	return volume.Response{}
}

func (s *volumeServer) List(req volume.Request) volume.Response {
	var volumes []*volume.Volume
	for v, _ := range s.volumes {
		volumes = append(volumes, &volume.Volume{Name: v, Mountpoint: mountPoint(v)})
	}
	return volume.Response{Volumes: volumes}
}

func (s *volumeServer) Get(req volume.Request) volume.Response {
	if _, ok := s.volumes[req.Name]; ok {
		return volume.Response{Volume: &volume.Volume{Name: req.Name, Mountpoint: mountPoint(req.Name)}}
	}
	return volume.Response{}
}

func (s *volumeServer) Remove(req volume.Request) volume.Response {
	delete(s.volumes, req.Name)
	return volume.Response{}
}

func (s *volumeServer) Path(req volume.Request) volume.Response {
	if _, ok := s.volumes[req.Name]; ok {
		return volume.Response{Mountpoint: mountPoint(req.Name)}
	}
	return volume.Response{}
}

func (s *volumeServer) Mount(req volume.MountRequest) volume.Response {
	resp := volume.Response{}
	if r, ok := s.volumes[req.Name]; ok {
		if err := mount(req.Name, memoryCgroupPath(r)); err != nil {
			resp.Err = err.Error()
		} else {
			resp.Mountpoint = mountPoint(req.Name)
		}
	} else {
		resp.Err = "volume does not exist"
	}
	return resp
}

func (s *volumeServer) Unmount(req volume.UnmountRequest) volume.Response {
	resp :=  volume.Response{}
	if err := umount(req.Name); err != nil {
		resp.Err = err.Error()
	}
	return resp
}

func (s *volumeServer) Capabilities(req volume.Request) volume.Response {
	return volume.Response{Capabilities: volume.Capability{Scope: "local"}}
}

func mount(name, memoryCgroupPath string) error {
	os.MkdirAll(mountPath(name), 0755)
	return start(mountPath(name), memoryCgroupPath)
}

func umount(name string) error {
	defer os.Remove(mountPath(name))
	return stop(mountPath(name))
}

func mountPath(name string) string {
	return filepath.Join(*flagRootPath, name)
}

func mountPoint(name string) string {
	return filepath.Join(*flagRootPath, name, "meminfo")
}
