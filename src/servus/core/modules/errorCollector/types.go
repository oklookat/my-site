package errorCollector

type errorCollector struct {
	errorCollectorI
	errorsArray []interface{}
}

type errorCollectorI interface {
	HasErrors() bool
	GetErrors() string
	AddEUnknown(issuers []string, message string)
	AddECustom(issuers []string, message string, data interface{})
	AddEAuthIncorrect(issuers []string)
	AddEAuthForbidden(issuers []string)
	AddENotFound(issuers []string)
	AddEValidationMinMax(issuers []string, min int, max int)
	AddEValidationEmpty(issuers []string)
	AddEValidationAllowed(issuers []string, allowed []string)
	AddEValidationInvalid(issuers []string, message string)
}