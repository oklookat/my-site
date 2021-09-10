package ancientUI

import (
	"fmt"
	"os"
	"strconv"
)

const selectItemID = 0

type SelectItem struct {
	Title string
	Items []string
}

type items struct {
	itemID     int
	selectItem SelectItem
}

type IAncientUI interface {
	Start()
}

type AncientUI struct {
	items
	IAncientUI
}

func AddSelect(item SelectItem) AncientUI {
	var ancient = AncientUI{}
	ancient.items.selectItem = item
	ancient.itemID = selectItemID
	return ancient
}

func (a *AncientUI) Start() string {
	switch a.itemID {
	case selectItemID:
		return onSelectItem(a.selectItem)
	}
	return ""
}


func onSelectItem(item SelectItem) string {
	// TODO: add title
	var selected string
	var counter = 1
	for _, element := range item.Items {
		var format = fmt.Sprintf("%v. %v", counter, element)
		fmt.Println(format)
		counter++
	}
	scanner(&selected)
	selectedNumber, err := strconv.Atoi(selected)
	var selectedNumberInSlice = selectedNumber - 1
	if err != nil || selectedNumberInSlice >= len(item.Items) || selectedNumberInSlice < 0 {
		fmt.Println("Wrong selection. Try again.")
		onSelectItem(item)
	}
	var selectedInSlice = selectedNumber - 1
	return item.Items[selectedInSlice]
}

func scanner(writeTo *string) {
	_, err := fmt.Fscan(os.Stdin, writeTo)
	if err != nil {
		panic(err)
	}
}

//func BootCmd(_servus *core.Servus) {
//	if len(os.Args) > 1 {
//		for index, element := range os.Args {
//			if index == 1 && element == "-menu" {
//				servus = _servus
//				menuMaster()
//				break
//			}
//			if index > 2 {
//				return
//			}
//		}
//	} else {
//		return
//	}
//}
