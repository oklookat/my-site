package elven

import (
	"fmt"
	"io/ioutil"
	"os"
	"servus/core"
	"servus/core/modules/ancientUI"
)

// bootCmd - when correct args specified, it calls functions.
func bootCmd() {
	// create superuser
	var isSuperuser = ancientUI.ReadArg("elven:superuser")
	if isSuperuser {
		cmdSuperuser()
	}
	// delete all from db
	var isRunMigration = ancientUI.ReadArg("elven:migrate")
	if isRunMigration {
		cmdMigrate()
	}
	// create tables in db
	var isRollbackMigration = ancientUI.ReadArg("elven:rollback")
	if isRollbackMigration {
		cmdRollback()
	}
}

// cmdSuperuser - create superuser (user with admin rights).
func cmdSuperuser() {
	core.Logger.Info("--- create superuser (CTRL + D to exit)")
	var username = ancientUI.AddInput(ancientUI.InputItem{Title: "Username"})
	var err = validatorUsername(username)
	if err != nil {
		core.Logger.Error(err.Error())
		cmdSuperuser()
	}
	user, err := dbUserFindBy(username)
	// if user exists
	if user != nil && err == nil {
		var deleteHim = ancientUI.AddQuestion(ancientUI.QuestionItem{Question: "Username exists. Delete?"})
		if !deleteHim {
			os.Exit(1)
		}
		err = dbUserDeleteBy(user.ID)
		if err != nil {
			println("Error while deleting user.")
			cmdSuperuser()
		}
		var createNew = ancientUI.AddQuestion(ancientUI.QuestionItem{Question: "Create new user?"})
		if !createNew {
			os.Exit(1)
		}
	}
	// continue
	var password = ancientUI.AddInput(ancientUI.InputItem{Title: "Password"})
	err = validatorPassword(password)
	if err != nil {
		core.Logger.Error(err.Error())
		cmdSuperuser()
	}
	// create
	_, err = dbUserCreate(ModelUser{Role: "admin", Username: username, Password: password})
	if err != nil {
		println("Error while creating user.")
		cmdSuperuser()
	}
	core.Logger.Info("User created.")
	os.Exit(1)
}

// cmdMigrate - create tables from SQL file.
func cmdMigrate() {
	var sqlPath = fmt.Sprintf("%v/settings/sql/elven.sql", core.Utils.GetExecuteDir())
	sqlPath = core.Utils.FormatPath(sqlPath)
	sql, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		core.Logger.Panic(err)
	}
	_, err = core.Database.Exec(string(sql))
	if err != nil {
		core.Logger.Panic(err)
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
		core.Logger.Panic(err)
	}
	core.Logger.Info("elven: database rollback successful")
	os.Exit(1)
}
