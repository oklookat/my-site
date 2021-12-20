package core

import "io"

// controller - controls controllers.
type controller struct {
	controllers []Controller
}

// add - add controller.
func (c *controller) add(co Controller) {
	if co == nil {
		panic("[core/control]: controller nil pointer")
	}
	if c.controllers == nil {
		c.controllers = make([]Controller, 0)
	}
	c.controllers = append(c.controllers, co)
}

//SendMessage - send message to all controllers.
func (c *controller) SendMessage(message string) {
	if c.controllers == nil {
		return
	}
	go func() {
		for index := range c.controllers {
			c.controllers[index].SendMessage(message)
		}
	}()
}

// SendFile - send file to all controllers.
func (c *controller) SendFile(caption *string, filename string, reader io.Reader) {
	if c.controllers == nil {
		return
	}
	go func() {
		for index := range c.controllers {
			c.controllers[index].SendFile(caption, filename, reader)
		}
	}()
}
