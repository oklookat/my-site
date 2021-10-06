package elven

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"servus/core"
	"servus/core/modules/ancientUI"
)

const (
	cmdFlagSuperuser = "elven:superuser"
	cmdFlagRollback = "elven:rollback"
	cmdFlagMigrate = "elven:migrate"
)

// bootCmd - when correct args specified, it calls functions.
func bootCmd() {
	// create superuser.
	if ancientUI.ArgumentExists(cmdFlagSuperuser) {
		cmdSuperuser()
	}
	// delete all from db.
	if ancientUI.ArgumentExists(cmdFlagRollback) {
		cmdRollback()
	}
	// create tables in db.
	if ancientUI.ArgumentExists(cmdFlagMigrate) {
		cmdMigrate()
	}
}

// cmdSuperuser - create superuser (user with admin rights).
func cmdSuperuser() {
	var err error
	core.Logger.Info("elven: create superuser (CTRL + D to exit)")
	chooseUsername:
		username, err := ancientUI.AddInput("Username")
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: failed to read username. Error")
			core.Logger.Panic(errPretty)
			os.Exit(1)
		}
		err = validatorUsername(username)
		if err != nil {
			core.Logger.Error("elven: validation failed. Error: %v", err.Error())
			goto chooseUsername
		}
		user, err := dbUserFindBy(username)
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: failed to find user. Error")
			core.Logger.Panic(errPretty)
			os.Exit(1)
		}
		// if user exists
		if user != nil {
			deleteHim, err := ancientUI.AddQuestion("Username exists. Delete?")
			if err != nil {
				var errPretty = errors.Wrap(err, "elven: failed to scan answer. Error")
				core.Logger.Panic(errPretty)
				os.Exit(1)
			}
			if !deleteHim {
				os.Exit(1)
			}
			err = dbUserDeleteBy(user.ID)
			if err != nil {
				var errPretty = errors.Wrap(err, "elven: delete user failed. Error")
				core.Logger.Panic(errPretty)
				os.Exit(1)
			}
			createNew, err := ancientUI.AddQuestion("User successfully deleted. Create new?")
			if err != nil {
				var errPretty = errors.Wrap(err, "elven: scan answer failed. Error")
				core.Logger.Panic(errPretty)
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
			core.Logger.Panic(errPretty)
			os.Exit(1)
		}
		err = validatorPassword(password)
		if err != nil {
			core.Logger.Error("elven: validation failed. Error: %v", err.Error())
			goto choosePassword
		}
		// create
		_, err = dbUserCreate(ModelUser{Role: "admin", Username: username, Password: password})
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: error while creating user. Error")
			core.Logger.Panic(errPretty)
			os.Exit(1)
		}
		core.Logger.Info("elven: user successfully created")
		os.Exit(1)
}

// cmdMigrate - create tables from SQL file.
func cmdMigrate() {
	var sqlPath = fmt.Sprintf("%v/settings/sql/elven.sql", core.Utils.GetExecuteDir())
	sqlPath = core.Utils.FormatPath(sqlPath)
	script, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Read SQL file error")
		core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	_, err = core.Database.Exec(string(script))
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: migration failed. Failed to execute SQL file)")
		core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	core.Logger.Info("elven: database migrate successful")
	os.Exit(1)
}

// cmdRollback - delete tables.
func cmdRollback() {
	_, err := core.Database.Exec(`
	DROP SCHEMA public CASCADE;
	CREATE SCHEMA public;
	GRANT ALL ON SCHEMA public TO postgres;
	GRANT ALL ON SCHEMA public TO public;
	`)
	if err != nil {
		var errPretty = errors.Wrap(err, "elven: rollback failed. Failed to execute drop script")
		core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	core.Logger.Info("elven: database rollback successful")
	os.Exit(1)
}
