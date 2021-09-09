package cmd

import (
"fmt"
"os"
	"servus/core"
)

var servus *core.Servus

func scanner(writeTo *string){
	_, err := fmt.Fscan(os.Stdin, writeTo)
	if err != nil {
		panic(err)
	}
}

func BootCmd(_servus *core.Servus) {
	if len(os.Args) > 1 {
		for index, element := range os.Args {
			if index == 1 && element == "-menu" {
				servus = _servus
				menuMaster()
				break
			}
			if index > 2 {
				return
			}
		}
	} else {
		return
	}
}

func menuDrawer(title string, items []string){
	var upper = fmt.Sprintf( "--- %v ---", title)
	fmt.Println(upper)
	for index, element := range items {
		var format = fmt.Sprintf("%v. %v", index, element)
		fmt.Println(format)
	}
	var titleLen = len(upper)
	var bottom string
	for i := 0; i < titleLen; i++ {
		bottom += "-"
	}
	fmt.Println(bottom)
	fmt.Println("Type: ")
}

func menuMaster() {
	menuDrawer("Servus", []string{"Exit", "Continue booting", "elven: create superuser"})
	var selected string
	scanner(&selected)
	switch selected {
	case "0": os.Exit(1)
	case "1": return
	case "2": createSuperuser()
	default:
		fmt.Println("Wrong selection. Try again.")
		menuMaster()
	}
	return
}


func createSuperuser(){
	fmt.Println("Username: ")
	var selected string
	scanner(&selected)

}

