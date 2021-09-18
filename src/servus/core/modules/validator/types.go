package validator

type validatorI interface {
	IsEmpty(text string) bool
	IsAlphanumeric(text string) bool
	IsAlphanumericWithSymbols(text string) bool
}
