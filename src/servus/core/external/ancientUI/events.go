package ancientUI

import (
	"fmt"
	"strconv"
	"strings"
)

func onSelectItem(title string, items []string) (selected string, err error) {
	for {
		fmt.Printf("--- %v\n", title)
		var counter = 1
		for index := range items {
			var format = fmt.Sprintf("%v. %v", counter, items[index])
			fmt.Println(format)
			counter++
		}
		fmt.Println("- Choose: ")
		err = scanner(&selected)
		if err != nil {
			return "", err
		}
		selectedNumber, err := strconv.Atoi(selected)
		if err != nil {
			return "", err
		}
		var selectedNumberInSlice = selectedNumber - 1
		if selectedNumberInSlice >= len(items) || selectedNumberInSlice < 0 {
			fmt.Println("- Wrong selection. Try again.")
			continue
		}
		return items[selectedNumberInSlice], err
	}
}

func onInputItem(title string) (input string, err error) {
	// ex: - username:.
	fmt.Println(fmt.Sprintf("- %v:", title))
	err = scanner(&input)
	return input, err
}

func onQuestionItem(question string) (result bool, err error) {
	for {
		// ex: - you drink water? (Y/N).
		fmt.Println(fmt.Sprintf("- %v (Y/N)", question))
		var input string
		err = scanner(&input)
		if err != nil {
			return false, err
		}
		input = strings.ToUpper(input)
		switch input {
		case "Y", "y":
			return true, err
		case "N", "n":
			return false, err
		default:
			fmt.Println(`- Wrong answer. Type «Y» or «N».`)
			break
		}
	}
}
