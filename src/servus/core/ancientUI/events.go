package ancientUI

import (
	"fmt"
	"strconv"
	"strings"
)

func onSelectItem(item SelectItem) string {
	fmt.Printf("--- %v\n", item.Title)
	var selected string
	var counter = 1
	for _, element := range item.Items {
		var format = fmt.Sprintf("%v. %v", counter, element)
		fmt.Println(format)
		counter++
	}
	fmt.Println("Choose: ")
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

func onInputItem(items InputItem) string{
	fmt.Println(fmt.Sprintf("%v:", items.Title))
	var input string
	scanner(&input)
	return input
}

func onQuestionItem(item QuestionItem) bool {
	fmt.Println(fmt.Sprintf("%v (Y/N)", item.Question))
	var input string
	scanner(&input)
	input = strings.ToUpper(input)
	switch input {
	case "Y":
		return true
	case "N":
		return false
	default:
		fmt.Println("Wrong answer. Type 'Y' or 'N'")
		onQuestionItem(item)
	}
	return false
}
