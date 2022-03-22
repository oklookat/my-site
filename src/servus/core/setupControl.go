package core

import (
	"io"
	"servus/core/internal/controlTelegram"
	"time"
)

// setup parent controller.
func (i *Instance) setupControl() {
	var parent = &parentController{}
	parent.new(i.Utils)

	// add Telegram bot.
	var tgEnabled = i.Config.Control.Telegram.Enabled
	if tgEnabled {
		var controlTG = controlTelegram.Controller{}
		controlTG.New(i.Config.Control.Telegram, i.Logger)
		parent.add(&controlTG)
	}

	// set.
	i.Control = parent
}

// control controllers / ops on all controllers.
type parentController struct {
	utils       Utils
	debouncer   func(callback func())
	controllers []Controller
}

func (c *parentController) new(u Utils) {
	c.utils = u

	// create 5-second debouncer to avoid controllers spam.
	c.debouncer = u.Debounce(5 * time.Second)
}

// add controller.
func (c *parentController) add(co Controller) {
	if co == nil {
		panic("[core/control]: controller nil pointer")
	}
	if c.controllers == nil {
		c.controllers = make([]Controller, 0)
	}
	c.controllers = append(c.controllers, co)
}

// run callback on every controller.
func (c *parentController) fetchControllers(callback func(c Controller), withDebounce bool) {
	if c.controllers == nil {
		return
	}
	var fetch = func() {
		for index := range c.controllers {
			callback(c.controllers[index])
		}
	}
	// debounce - avoid spam.
	if withDebounce {
		c.debouncer(func() {
			fetch()
		})
	} else {
		fetch()
	}
}

// send message to all controllers.
func (c *parentController) SendMessage(message string) {
	go c.fetchControllers(func(c Controller) {
		c.SendMessage(message)
	}, true)
}

// send file to all controllers.
func (c *parentController) SendFile(caption *string, filename string, data io.Reader) {
	if caption == nil || data == nil {
		return
	}
	go c.fetchControllers(func(c Controller) {
		c.SendFile(caption, filename, data)
	}, true)
}

// add command to all controllers.
func (c *parentController) AddCommand(command string, callback func(args []string)) {
	c.fetchControllers(func(c Controller) {
		c.AddCommand(command, callback)
	}, false)
}
