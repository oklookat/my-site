package cmd

import (
	"errors"
	"fmt"

	"github.com/oklookat/argument"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Config struct {
	Logger Logger
}

var cfg *Config

func Boot(c *Config) error {

	if c == nil {
		return errors.New("cmd: empty config")
	}

	cfg = c

	var arguments = argument.New()

	// help.
	arguments.Add("elven-help", "eh", func(values []string) {
		fmt.Println(helpPage)
		exit(false)
	})

	// create superuser.
	arguments.Add("create-superuser", "csu", func(values []string) {
		cfg.Logger.Info("cmd: create superuser")
		var err = createUser(true, values)
		afterCommand(err)
	})

	// create user.
	arguments.Add("create-user", "cu", func(values []string) {
		cfg.Logger.Info("cmd: create user")
		var err = createUser(false, values)
		afterCommand(err)
	})

	// create tables in db.
	arguments.Add("migrate", "mg", func(values []string) {
		cfg.Logger.Info("cmd: migrate")
		var err = migrate(values)
		afterCommand(err)
	})

	// delete all from db.
	arguments.Add("rollback", "rb", func(values []string) {
		cfg.Logger.Info("cmd: rollback")
		var err = rollback()
		afterCommand(err)
	})

	arguments.Start()

	return nil
}
