package core

import "servus/core/internal/banhammer"

func (i *Instance) setupBanhammer() {
	// get dir.
	var dataDir, err = i.Dirs.GetData()
	if err != nil {
		i.Logger.Panic(err)
	}

	// boot banhammer.
	var hammer = banhammer.Instance{}
	err = hammer.Boot(dataDir, 3)
	if err != nil {
		i.Logger.Panic(err)
	}
	i.Banhammer = hammer
}
