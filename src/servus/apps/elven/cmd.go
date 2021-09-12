package elven

import (
	"os"
	"servus/core/ancientUI"
)

func bootCmd() {
	cmdSuperuser()
}

// create superuser
func cmdSuperuser() {
	var isCmd = ancientUI.ReadArg("elven:superuser")
	if !isCmd {
		return
	}
	cmdSuperuserRunForm()
}

func cmdSuperuserRunForm() {
	servus.Logger.Info("--- create superuser (CTRL + D to exit)")
	var username = ancientUI.AddInput(ancientUI.InputItem{Title: "Username"})
	var err = validatorUsername(username)
	if err != nil {
		servus.Logger.Error(err.Error())
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
		servus.Logger.Error(err.Error())
		cmdSuperuserRunForm()
	}
	// create
	err = dbCreateUser(modelUser{role: "admin", username: username, password: password})
	if err != nil {
		println("Error while creating user.")
		cmdSuperuserRunForm()
	}
	servus.Logger.Info("User created.")
	os.Exit(1)
}

