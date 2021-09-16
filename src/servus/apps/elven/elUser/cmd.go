package elUser

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"servus/core"
	ancientUI "servus/core/modules/ancientUI"
)

// create superuser
func bootCmd() {
	var isSuperuser = ancientUI.ReadArg("elven:superuser")
	if isSuperuser {
		cmdSuperuserRunForm()
	}
	var isRunMigration = ancientUI.ReadArg("elven:migrate")
	if isRunMigration {
		cmdMigrate()
	}
	var isRollbackMigration = ancientUI.ReadArg("elven:rollback")
	if isRollbackMigration {
		cmdRollback()
	}
}

func cmdSuperuserRunForm() {
	core.Logger.Info("--- create superuser (CTRL + D to exit)")
	var username = ancientUI.AddInput(ancientUI.InputItem{Title: "Username"})
	var err = validatorUsername(username)
	if err != nil {
		core.Logger.Error(err.Error())
		cmdSuperuserRunForm()
	}
	user, err := dbFindUserBy(username)
	// if user exists
	if len(user.id) > 1 {
		var deleteHim = ancientUI.AddQuestion(ancientUI.QuestionItem{Question: "Username exists. Delete?"})
		if !deleteHim {
			os.Exit(1)
		}
		err = dbDeleteUserBy(user.id)
		if err != nil {
			println("Error while deleting user.")
			cmdSuperuserRunForm()
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
		cmdSuperuserRunForm()
	}
	// create
	err = dbCreateUser(modelUser{role: "admin", username: username, password: password})
	if err != nil {
		println("Error while creating user.")
		cmdSuperuserRunForm()
	}
	core.Logger.Info("User created.")
	os.Exit(1)
}

func cmdMigrate(){
	var sqlPath = fmt.Sprintf("%v/apps/elven/elSQL/all.sql", core.Utils.GetExecuteDir())
	sqlPath = core.Utils.FormatPath(sqlPath)
	sql, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		core.Logger.Panic(err)
	}
	_, err = core.Database.Connection.Exec(context.Background(), string(sql))
	if err != nil {
		core.Logger.Panic(err)
	}
	core.Logger.Info("elven: database migrate successful")
	os.Exit(1)
}

func cmdRollback(){
	_, err := core.Database.Connection.Exec(context.Background(), `
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