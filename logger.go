package syntaxlogger

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/fatih/color"
)

const (
	THEME = "monokai"
)

type SQLLogger struct {
}

func (s *SQLLogger) Printf(format string, v ...interface{}) {
	check := strings.ToLower(format)
	output := fmt.Sprintf(format, v...)

	if strings.Index(check, "begin transaction") == 0 ||
		strings.Index(check, "commit transaction") == 0 {

		// transaction start/end
		color.New(color.FgHiBlue).Fprintln(os.Stdout, output)
	} else if strings.Index(check, "rollback transaction") == 0 {

		// transaction rollback
		color.New(color.FgRed).Fprintln(os.Stdout, output)
	} else if strings.Index(check, "select") > -1 ||
		strings.Index(check, "update") > -1 ||
		strings.Index(check, "insert") > -1 ||
		strings.Index(check, "execute") > -1 ||
		strings.Index(strings.TrimSpace(check), "@p") == 0 {

		// T-SQL
		quick.Highlight(os.Stdout, output, "Transact-SQL", "terminal", THEME)
	} else {

		// Anything else
		os.Stdout.WriteString(output)
	}
}

func (s *SQLLogger) Println(v ...interface{}) {
	output := fmt.Sprintln(v...)
	s.Printf(output)
}
