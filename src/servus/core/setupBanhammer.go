package core

import "servus/core/internal/banhammer"

func (i *Instance) setupBanhammer() error {
	// get dir.
	var dataDir, err = i.Dirs.GetData()
	if err != nil {
		return err
	}

	// boot banhammer.
	var hammer = banhammer.Instance{}
	if err = hammer.Boot(dataDir, 3); err != nil {
		return err
	}

	// set.
	i.Banhammer = hammer
	return err
}
