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
		core.Logger.Error("elven: wrong flag at cmd createUser")
		os.Exit(1)
	}
	if flag == cmdFlagSuperuser {
		core.Logger.Info("elven: create superuser (CTRL + D to exit)")
	} else if flag == cmdFlagUser {
		core.Logger.Info("elven: create user (CTRL + D to exit)")
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
		core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	err = validatorUsername(username)
	if err != nil {
		core.Logger.Error("elven: validation failed. Error: %v", err.Error())
		goto chooseUsername
	}
	user, err := eUser.databaseFind(username)
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
		err = eUser.databaseDelete(user.ID)
		if err != nil {
			var errPretty = errors.Wrap(err, "elven: delete user failed. Error")
			core.Logger.Panic(errPretty)
			os.Exit(1)
		}
		createNew, err := ancientUI.AddQuestion(fmt.Sprintf("%v successfully deleted. Create new?", sign))
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
	user = &ModelUser{Role: role, Username: username, Password: password}
	err = eUser.databaseCreate(user)
	if err != nil {
		var errPretty = errors.Wrap(err, fmt.Sprintf("elven: error while creating %v. Error", sign))
		core.Logger.Panic(errPretty)
		os.Exit(1)
	}
	core.Logger.Info("elven: %v successfully created", sign)
	os.Exit(1)
}

// migrate - create tables in database from SQL file.
func (c *objectCmd) migrate() {
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

// rollback - delete tables from database.
func (c *objectCmd) rollback() {
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
