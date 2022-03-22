package core

import "servus/core/internal/directories"

func (i *Instance) setupDirectories() {
	var dirs = directories.Instance{}
	var err = dirs.Boot()
	if err != nil {
		panic(err)
	}
	i.Dirs = &dirs
}
