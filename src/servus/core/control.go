package core

import "io"

type control struct {
	controllers []Controller
}

func (c *control) GetEnabled() bool {
	return true
}

// addController - add controller.
func (c *control) addController(co Controller) {
	if c.controllers == nil {
		c.controllers = make([]Controller, 0)
	}
	c.controllers = append(c.controllers, co)
}

func (c *control) SendMessage(message string) {
	for index := range c.controllers {
		if !c.controllers[index].GetEnabled() {
			continue
		}
		c.controllers[index].SendMessage(message)
	}
}

func (c *control) SendFile(caption *string, filename string, reader io.Reader) {
	for index := range c.controllers {
		if !c.controllers[index].GetEnabled() {
			continue
		}
		c.controllers[index].SendFile(caption, filename, reader)
	}
}
