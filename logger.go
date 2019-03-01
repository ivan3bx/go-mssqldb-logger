package syntaxlogger

import (
	"fmt"
	"os"
	"regexp"

	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
)

const (
	theme = "monokai"
)

var (
	rxpTxPass = regexp.MustCompile("(?i)^(begin|commit) transaction")
	rxpTxFail = regexp.MustCompile("(?i)^rollback transaction")
	rxpSql    = regexp.MustCompile("(?i)(SELECT|UPDATE|INSERT|EXECUTE|@p)")
)

type SQLLogger struct{}

func (s *SQLLogger) Printf(format string, v ...interface{}) {
	output := fmt.Sprintf(format, v...)

	switch {
	case rxpTxPass.MatchString(output):
		color.New(color.FgHiBlue).Fprintln(os.Stdout, output)
	case rxpTxFail.MatchString(output):
		color.New(color.FgRed).Fprintln(os.Stdout, output)
	case rxpSql.MatchString(output):
		quick.Highlight(os.Stdout, output, "Transact-SQL", "terminal", theme)
	default:
		os.Stdout.WriteString(output)
	}
}

func (s *SQLLogger) Println(v ...interface{}) {
	output := fmt.Sprintln(v...)
	s.Printf(output)
}
