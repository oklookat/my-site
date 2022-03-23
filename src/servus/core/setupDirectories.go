package core

import "servus/core/internal/directories"

func (i *Instance) setupDirectories() error {
	var err error

	// boot.
	var dirs = directories.Instance{}
	if err = dirs.Boot(); err != nil {
		return err
	}

	// set.
	i.Dirs = &dirs
	return err
}
