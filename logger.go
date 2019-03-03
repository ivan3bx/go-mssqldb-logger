package mssqllog

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/quick"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
)

const (
	// see https://github.com/alecthomas/chroma/tree/master/styles
	theme = "monokai"
)

var (
	rxpTxPass = regexp.MustCompile("(?i)^(begin|commit) transaction")
	rxpTxFail = regexp.MustCompile("(?i)^rollback transaction")
	rxpSql    = regexp.MustCompile("(?i)^\\s*(DECLARE|SELECT|UPDATE|INSERT|DELETE|EXEC(UTE)?|@p)")
)

type SQLLogger struct {
	Logger     mssql.Logger
	forceColor bool
	ignoreTTY  bool
}

func (s *SQLLogger) Printf(format string, v ...interface{}) {
	var (
		colorize = true
		out      = &strings.Builder{}
	)

	switch ll := s.Logger.(type) {
	case *log.Logger:
		colorize = isLoggerColorEnabled(ll, s.ignoreTTY)
	case *logrus.Logger:
		colorize = isLogrusColorEnabled(ll, s.ignoreTTY)
	}

	if colorize {
		s.colorizeOutput(out, format, v...)
	} else {
		fmt.Fprintf(out, format, v...)
	}

	if s.Logger != nil {
		s.Logger.Printf(out.String())
	} else {
		log.Printf(out.String())
	}
}

func (s *SQLLogger) colorizeOutput(out io.Writer, format string, v ...interface{}) {
	var c *color.Color

	switch {
	case rxpSql.MatchString(format):
		//
		// sql logging
		//
		source := fmt.Sprintf(format, v...)
		quick.Highlight(out, source, "Transact-SQL", "terminal", theme)
		return
	case rxpTxPass.MatchString(format):
		//
		// transaction begin/end
		//
		c = color.New(color.FgHiBlue)
	case rxpTxFail.MatchString(format):
		//
		// transaction rollback
		//
		c = color.New(color.FgRed)
	default:
		c = color.New(color.FgHiYellow, color.Italic)
	}

	if c != nil {
		if s.ignoreTTY {
			c.EnableColor()
		}
		c.Fprintf(out, format, v...)
	}
}

func (s *SQLLogger) Println(v ...interface{}) {
	s.Printf(fmt.Sprint(v...))
}

func isLoggerColorEnabled(ll *log.Logger, ignoreTTY bool) bool {
	if w, ok := ll.Writer().(*os.File); !ignoreTTY && (!ok || !isatty.IsTerminal(w.Fd())) {
		return false
	}

	return true
}

func isLogrusColorEnabled(ll *logrus.Logger, ignoreTTY bool) bool {
	if w, ok := ll.Out.(*os.File); !ignoreTTY && (!ok || !isatty.IsTerminal(w.Fd())) {
		return false
	}

	switch fmtt := ll.Formatter.(type) {
	case *logrus.TextFormatter:
		return !fmtt.DisableColors
	case *logrus.JSONFormatter:
		return false
	}
	return true
}
