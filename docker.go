package main

import (
	"fmt"
	"io/ioutil"

	"github.com/docker/go-plugins-helpers/volume"
)

func memoryCgroupPath(req *volume.Request) string {
	id, err := ioutil.ReadFile(req.Options[CidFile])
	if err != nil {
		return "/docker/"
	}
	return fmt.Sprintf("/docker/%s", string(id))
}
