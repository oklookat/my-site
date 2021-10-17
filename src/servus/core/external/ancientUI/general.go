package ancientUI

func AddSelect(title string, items []string) (selected string, err error) {
	return onSelectItem(title, items)
}

func AddInput(title string) (input string, err error) {
	return onInputItem(title)
}

func AddQuestion(question string) (result bool, err error) {
	return onQuestionItem(question)
}
