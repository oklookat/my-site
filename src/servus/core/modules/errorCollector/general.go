package errorCollector

// New - create new errorCollector instance.
func New() *ErrorCollector {
	var ec = ErrorCollector{Errors: struct {
		ErrorsSlice `json:"errors"`
	}(struct{ ErrorsSlice }{ErrorsSlice: make([]interface{}, 0)})}
	return &ec
}
