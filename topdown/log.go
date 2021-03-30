package topdown

import (
	"time"

	"github.com/open-policy-agent/opa/ast"
)

const (
	timeLayout = "2006-01-02 15:04:05 MST"
)

// Logger defines the interface for logging in the top-down evaluation engine.
type Logger interface {
	Log(time time.Time, loc *ast.Location, statement interface{})
}

// BufferLogger implements the Logger interface by
// simply buffering all log statements received.
type BufferLogger []*LogStatement

// NewBufferLogger returns a new BufferLogger.
func NewBufferLogger() *BufferLogger {
	return &BufferLogger{}
}

// Log adds the given log statement to the buffer.
func (bl *BufferLogger) Log(time time.Time, loc *ast.Location, statement interface{}) {
	*bl = append(*bl, &LogStatement{
		Statement: statement,
		Timestamp: time,
		Location:  loc,
	})
}

// All returns all logs in the buffer.
func (bl *BufferLogger) All() []interface{} {
	all := make([]interface{}, len(*bl))

	for i, l := range *bl {
		all[i] = l
	}

	return all
}

// LogStatement is a simple log statement.
type LogStatement struct {
	Timestamp time.Time     `json:"timestamp"`
	Location  *ast.Location `json:"location"`
	Statement interface{}   `json:"statement"`
}

func builtinLog(bctx BuiltinContext, args []*ast.Term, iter func(*ast.Term) error) error {
	var epoch int64
	if err := ast.As(bctx.Time.Value, &epoch); err != nil {
		return handleBuiltinErr(ast.Log.Name, bctx.Location, err)
	}

	st, err := ast.JSON(args[0].Value)
	if err != nil {
		return handleBuiltinErr(ast.Log.Name, bctx.Location, err)
	}

	for i := range bctx.Loggers {
		bctx.Loggers[i].Log(time.Unix(0, epoch), bctx.Location, st)
	}

	return iter(ast.BooleanTerm(true))
}

func init() {
	RegisterBuiltinFunc(ast.Log.Name, builtinLog)
}
