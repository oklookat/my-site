package core

import (
	"errors"
	"io"
	"servus/core/external/utils"
	"servus/core/internal/controlTelegram"
	"time"
)

// setup parent controller.
func (i *Instance) setupControl() error {
	var err error

	// create parent controller.
	var parent = &parentController{}
	parent.new()

	// add Telegram bot.
	var tgEnabled = i.Config.Control.Telegram.Enabled
	if tgEnabled {
		var controlTG = controlTelegram.Controller{}
		if err = controlTG.New(i.Config.Control.Telegram, i.Logger); err != nil {
			return err
		}
		if err = parent.add(&controlTG); err != nil {
			return err
		}
	}

	// set.
	i.Control = parent
	return err
}

// control controllers / ops on all controllers.
type parentController struct {
	debouncer   func(callback func())
	controllers []Controller
}

func (c *parentController) new() {
	// create 5-second debouncer to avoid controllers spam.
	c.debouncer = utils.Debounce(5*time.Second, true)
}

// add controller.
func (c *parentController) add(co Controller) error {
	var err error
	if co == nil {
		return errors.New("[core/control]: controller nil pointer")
	}
	if c.controllers == nil {
		c.controllers = make([]Controller, 0)
	}
	c.controllers = append(c.controllers, co)
	return err
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
	if c.controllers == nil {
		return
	}
	c.fetchControllers(func(c Controller) {
		c.AddCommand(command, callback)
	}, false)
}
