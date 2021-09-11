package ancientUI


func AddSelect(item SelectItem) string {
	return onSelectItem(item)
}

func AddInput(item InputItem) string {
	return onInputItem(item)
}

func AddQuestion(item QuestionItem) bool {
	return onQuestionItem(item)
}