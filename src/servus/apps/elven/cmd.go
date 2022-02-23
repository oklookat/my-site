package elven

import (
	"fmt"
	"io/ioutil"
	"os"
	"servus/apps/elven/model"
	"servus/core/external/argument"

	"github.com/pkg/errors"
)

const (
	cmdFlagUsername       = "-username"
	cmdFlagPassword       = "-password"
	cmdFlagDeleteIfExists = "-die"
	cmdFlagSQLPath        = "-sql"
	cmdFlagSuperuser      = "el:su"
	cmdFlagUser           = "el:tu"
	cmdFlagRollback       = "el:rb"
	cmdFlagMigrate        = "el:mg"
)

// commandline methods.
type cmd struct {
}

// call methods depending on arguments.
func (c *cmd) boot() {
	// create superuser.
	var createSU = argument.Get(cmdFlagSuperuser)
	if createSU != nil {
		call.Logger.Info("cmd: create superuser")
		c.createUser(cmdFlagSuperuser)
		os.Exit(0)
	}
	// create user.
	var createTU = argument.Get(cmdFlagUser)
	if createTU != nil {
		call.Logger.Info("cmd: create user")
		c.createUser(cmdFlagUser)
		os.Exit(0)
	}
	// delete all from db.
	if argument.Get(cmdFlagRollback) != nil {
		call.Logger.Info("cmd: rollback")
		c.rollback()
		os.Exit(0)
	}
	// create tables in db.
	if argument.Get(cmdFlagMigrate) != nil {
		call.Logger.Info("cmd: migrate")
		c.migrate()
		os.Exit(0)
	}
}

// create user or superuser (UserModel).
func (c *cmd) createUser(flag string) {
	// validate args.
	var usernameArg = argument.Get(cmdFlagUsername)
	var passwordArg = argument.Get(cmdFlagPassword)
	if usernameArg == nil || passwordArg == nil {
		call.Logger.Error("if you need to create superuser or user, set username and password")
		return
	}
	if usernameArg.Value == nil || passwordArg.Value == nil {
		call.Logger.Error("username or password cannot be empty")
		return
	}
	// main.
	var err error
	defer func() {
		if err != nil {
			call.Logger.Error(err.Error())
		}
	}()
	var username = *usernameArg.Value
	var password = *passwordArg.Value
	var deleteIfExists = argument.Get(cmdFlagDeleteIfExists) != nil
	var role string
	switch flag {
	default:
		call.Logger.Error("wrong flag provided to createUser")
	case cmdFlagSuperuser:
		role = "admin"
	case cmdFlagUser:
		role = "user"
	}
	var user = model.User{}
	user.Username = username
	err = user.ValidateUsername()
	if err != nil {
		return
	}
	found, err := user.FindByUsername()
	if err != nil {
		return
	}
	// if user exists
	if found {
		if !deleteIfExists {
			call.Logger.Info("user already exists")
			return
		}
		err = user.DeleteByID()
		if err != nil {
			return
		}
	}
	user.Password = password
	err = user.ValidatePassword()
	if err != nil {
		return
	}
	// create
	user.Role = role
	err = user.Create()
	if err != nil {
		return
	}
	call.Logger.Info("done")
}

// create tables in database from SQL file.
func (c *cmd) migrate() {
	var sqlPath string
	var sqlPathVal = argument.Get(cmdFlagSQLPath)
	if sqlPathVal.Value == nil {
		var executionDir, _ = call.Utils.GetExecutionDir()
		sqlPath = fmt.Sprintf("%v/settings/sql/elven.sql", executionDir)
	} else {
		sqlPath = *sqlPathVal.Value
	}
	sqlPath = call.Utils.FormatPath(sqlPath)
	script, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Read SQL file error")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	_, err = model.StringAdapter.Exec(string(script))
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Failed to execute SQL file")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	call.Logger.Info("elven: database migrate successful")
	os.Exit(1)
}

// delete tables from database.
func (c *cmd) rollback() {
	_, err := model.StringAdapter.Exec(`
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;
	`)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: rollback failed. Failed to execute drop script")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	call.Logger.Info("elven: database rollback successful")
	os.Exit(1)
}
