package core

import "io"

type control struct {
	controllers []Controller
}

// add - add controller.
func (c *control) add(co Controller) {
	if co == nil {
		panic("[core/control]: controller nil pointer")
	}
	if c.controllers == nil {
		c.controllers = make([]Controller, 0)
	}
	c.controllers = append(c.controllers, co)
}

func (c *control) SendMessage(message string) {
	if c.controllers == nil {
		return
	}
	go func() {
		for index := range c.controllers {
			c.controllers[index].SendMessage(message)
		}
	}()
}

func (c *control) SendFile(caption *string, filename string, reader io.Reader) {
	if c.controllers == nil {
		return
	}
	go func() {
		for index := range c.controllers {
			c.controllers[index].SendFile(caption, filename, reader)
		}
	}()
}
