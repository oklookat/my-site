package elven

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"servus/core/external/ancientUI"
)

const (
	cmdFlagSuperuser = "elven:superuser"
	cmdFlagUser      = "elven:user"
	cmdFlagRollback  = "elven:rollback"
	cmdFlagMigrate   = "elven:migrate"
)

// objectCmd - commandline methods.
type objectCmd struct {
}

// boot - call methods depending on startup arguments.
func (c *objectCmd) boot() {
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

// createUser - create user or superuser (ModelUser).
func (c *objectCmd) createUser(flag string) {
	var err error
	if flag != cmdFlagSuperuser && flag != cmdFlagUser {
		instance.Logger.Error("elven: wrong flag at cmd createUser")
		os.Exit(1)
	}
	if flag == cmdFlagSuperuser {
		instance.Logger.Info("elven: create superuser (CTRL + D to exit)")
	} else if flag == cmdFlagUser {
		instance.Logger.Info("elven: create user (CTRL + D to exit)")
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
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	err = eUser.validatorUsername(username)
	if err != nil {
		instance.Logger.Error(fmt.Sprintf("elven: validation failed. Error: %v", err.Error()))
		goto chooseUsername
	}
	var user = ModelUser{Username: username}
	found, err := user.findByUsername()
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: failed to find user. Error")
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	// if user exists
	if found {
		deleteHim, err := ancientUI.AddQuestion("Username exists. Delete?")
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: failed to scan answer. Error")
			instance.Logger.Panic(errPretty)
			os.Exit(1)
		}
		if !deleteHim {
			os.Exit(1)
		}
		err = user.deleteByID()
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: delete user failed. Error")
			instance.Logger.Panic(errPretty)
			os.Exit(1)
		}
		createNew, err := ancientUI.AddQuestion(fmt.Sprintf("%v successfully deleted. Create new?", sign))
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: scan answer failed. Error")
			instance.Logger.Panic(errPretty)
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
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	err = eUser.validatorPassword(password)
	if err != nil {
		instance.Logger.Error(fmt.Sprintf("elven: validation failed. Error: %v", err.Error()))
		goto choosePassword
	}
	// create
	user = ModelUser{Role: role, Username: username, Password: password}
	err = user.create()
	if err != nil {
		var errPretty = errors.Wrap(err, fmt.Sprintf("elven: error while creating %v. Error", sign))
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	instance.Logger.Info(fmt.Sprintf("elven: %v successfully created", sign))
	os.Exit(1)
}

// migrate - create tables in database from SQL file.
func (c *objectCmd) migrate() {
	var sqlPath = fmt.Sprintf("%v/settings/sql/elven.sql", instance.Utils.GetExecutionDir())
	sqlPath = instance.Utils.FormatPath(sqlPath)
	script, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Read SQL file error")
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	_, err = instance.DB.Conn.Exec(string(script))
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Failed to execute SQL file)")
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	instance.Logger.Info("elven: database migrate successful")
	os.Exit(1)
}

// rollback - delete tables from database.
func (c *objectCmd) rollback() {
	_, err := instance.DB.Conn.Exec(`
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;
	`)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: rollback failed. Failed to execute drop script")
		instance.Logger.Panic(errPretty)
		os.Exit(1)
	}
	instance.Logger.Info("elven: database rollback successful")
	os.Exit(1)
}
