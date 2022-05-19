package logger

// logging level.
const (
	// no messages.
	LevelSilent = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
)

type leveler interface {
	// get error level number. Used for internal comparison.
	getLevel() int
	// example: debug level must return "debug". Or debug level can return whatever. But better, level must return self name. Used for write to file and to console.
	getLevelWord() string
	// get console color for level.
	getColor() string
	// get user message for level.
	getMessage() string
}

// log level.
type level struct {
	number  int
	word    string
	color   string
	message string
}

func (l *level) getLevel() int {
	return l.number
}

func (l *level) getLevelWord() string {
	return l.word
}

func (l *level) getColor() string {
	return l.color
}

func (l *level) getMessage() string {
	return l.message
}
