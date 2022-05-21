package core

import "github.com/oklookat/gocors"

func (i *Instance) setupCors() error {
	// cors.
	var corsed, err = gocors.New(i.Config.Security.CORS)
	if err != nil {
		return err
	}
	i.Cors = corsed
	return nil
}
