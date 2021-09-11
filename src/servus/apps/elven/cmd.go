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
	var username = ancientUI.AddInput(ancientUI.InputItem{Title: "username"})
	var err = validatorUsername(username)
	if err != nil {
		servus.Logger.Error(err.Error())
		cmdSuperuserRunForm()
	}
	user, err := dbSearchUserBy(username)
	// if user exists
	if len(user.id) > 1 {
		var deleteHim = ancientUI.AddQuestion(ancientUI.QuestionItem{Question: "This user exists. Delete him?"})
		if !deleteHim {
			os.Exit(1)
		}
		err = dbDeleteUserBy(user.id)
		if err != nil {
			println("Error while deleting user.")
			cmdSuperuserRunForm()
		}
		var createNew = ancientUI.AddQuestion(ancientUI.QuestionItem{Question: "Continue creating user?"})
		if !createNew {
			os.Exit(1)
		}
	}
	// continue
	var password = ancientUI.AddInput(ancientUI.InputItem{Title: "password"})
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
}

