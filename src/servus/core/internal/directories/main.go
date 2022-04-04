package directories

import (
	"os"
	"path/filepath"
)

type Instance struct {
}

func (i *Instance) Boot() error {
	// check data dir and create if not exists.
	var dataDir, err = i.GetData()
	if err != nil {
		return err
	}
	err = os.MkdirAll(dataDir, 0644)
	return err
}

func (i *Instance) GetExecution() (path string, err error) {
	// get executable dir (binary file or main.go if debug).
	path, err = os.Executable()
	if err != nil {
		err = wrapError(err, "failed to get execution directory. Error")
		return
	}

	// remove filename from path.
	path = filepath.Dir(path)

	// go to symlinks.
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		err = wrapError(err, "failed symlink follow. Error")
		return
	}
	path = filepath.ToSlash(path)
	return
}

func (i *Instance) GetData() (string, error) {
	// get exec dir.
	var exec, err = i.GetExecution()
	if err != nil {
		return exec, err
	}

	// get data dir.
	var dataDir = exec + "/data"
	return dataDir, err
}
