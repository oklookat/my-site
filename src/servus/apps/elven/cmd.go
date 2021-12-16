package elven

import (
	"fmt"
	"io/ioutil"
	"os"
	"servus/core/external/ancientUI"

	"github.com/pkg/errors"
)

const (
	cmdFlagSuperuser = "elven:superuser"
	cmdFlagUser      = "elven:user"
	cmdFlagRollback  = "elven:rollback"
	cmdFlagMigrate   = "elven:migrate"
)

// cmd - commandline methods.
type cmd struct {
}

// boot - call methods depending on startup arguments.
func (c *cmd) boot() {
	// create superuser.
	if ancientUI.ArgumentExists(cmdFlagSuperuser) {
		c.createUser(cmdFlagSuperuser)
	}
	if ancientUI.ArgumentExists(cmdFlagUser) {
		c.createUser(cmdFlagUser)
	}
	// delete all from db.
	if ancientUI.ArgumentExists(cmdFlagRollback) {
		c.rollback()
	}
	// create tables in db.
	if ancientUI.ArgumentExists(cmdFlagMigrate) {
		c.migrate()
	}
}

// createUser - create user or superuser (UserModel).
func (c *cmd) createUser(flag string) {
	var err error
	if flag != cmdFlagSuperuser && flag != cmdFlagUser {
		call.Logger.Error("elven: wrong flag at cmd createUser")
		os.Exit(1)
	}
	if flag == cmdFlagSuperuser {
		call.Logger.Info("elven: create superuser (CTRL + D to exit)")
	} else if flag == cmdFlagUser {
		call.Logger.Info("elven: create user (CTRL + D to exit)")
	}
	var isSuperuser = flag == cmdFlagSuperuser
	var role string
	var sign string
	if isSuperuser {
		role = "admin"
		sign = "superuser"
	} else {
		role = "user"
		sign = "user"
	}
chooseUsername:
	username, err := ancientUI.AddInput("Username")
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: failed to read username. Error")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	var user = UserModel{Username: username}
	err = user.validateUsername()
	if err != nil {
		call.Logger.Error(fmt.Sprintf("elven: validation failed. Error: %v", err.Error()))
		goto chooseUsername
	}
	found, err := user.findByUsername()
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: failed to find user. Error")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	// if user exists
	if found {
		deleteHim, err := ancientUI.AddQuestion("Username exists. Delete?")
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: failed to scan answer. Error")
			call.Logger.Panic(errPretty)
			os.Exit(1)
		}
		if !deleteHim {
			os.Exit(1)
		}
		err = user.deleteByID()
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: delete user failed. Error")
			call.Logger.Panic(errPretty)
			os.Exit(1)
		}
		createNew, err := ancientUI.AddQuestion(fmt.Sprintf("%v successfully deleted. Create new?", sign))
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: scan answer failed. Error")
			call.Logger.Panic(errPretty)
			os.Exit(1)
		}
		if !createNew {
			os.Exit(1)
		}
	}
choosePassword:
	password, err := ancientUI.AddInput("Password")
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: failed to scan password. Error")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	user = UserModel{Role: role, Username: username, Password: password}
	err = user.validatePassword()
	if err != nil {
		call.Logger.Error(fmt.Sprintf("elven: validation failed. Error: %v", err.Error()))
		goto choosePassword
	}
	// create
	err = user.create()
	if err != nil {
		var errPretty = errors.Wrap(err, fmt.Sprintf("elven: error while creating %v. Error", sign))
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	call.Logger.Info(fmt.Sprintf("elven: %v successfully created", sign))
	os.Exit(1)
}

// migrate - create tables in database from SQL file.
func (c *cmd) migrate() {
	var executionDir, _ = call.Utils.GetExecutionDir()
	var sqlPath = fmt.Sprintf("%v/settings/sql/elven.sql", executionDir)
	sqlPath = call.Utils.FormatPath(sqlPath)
	script, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Read SQL file error")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	_, err = call.DB.Conn.Exec(string(script))
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Failed to execute SQL file)")
		call.Logger.Panic(errPretty)
		os.Exit(1)
	}
	call.Logger.Info("elven: database migrate successful")
	os.Exit(1)
}

// rollback - delete tables from database.
func (c *cmd) rollback() {
	_, err := call.DB.Conn.Exec(`
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
