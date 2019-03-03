package mssqllog

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorizedOutput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "select statement",
			input:  "SELECT 1",
			output: "\x1b[96m\x1b[40mSELECT\x1b[0m",
		},
		{
			name:   "update statement",
			input:  "UPDATE foo SET x = 1 WHERE x = 2",
			output: "\x1b[96m\x1b[40mUPDATE\x1b[0m",
		},
		{
			name:   "insert statement",
			input:  "INSERT INTO foo (bar, baz) VALUES (1)",
			output: "\x1b[96m\x1b[40mINSERT\x1b[0m",
		},
		{
			name:   "delete statement",
			input:  "DELETE FROM foo WHERE bar=-838381",
			output: "\x1b[96m\x1b[40mDELETE\x1b[0m",
		},
		{
			name:   "execute statement",
			input:  "EXECUTE dbo.uspGetEmployeeManagers @BusinessEntityID = 50;",
			output: "\x1b[96m\x1b[40mEXECUTE\x1b[0m",
		},
		{
			name:   "execute ('exec') statement",
			input:  "EXEC dbo.uspGetEmployeeManagers @BusinessEntityID = 50;",
			output: "\x1b[96m\x1b[40mEXEC\x1b[0m",
		},
		{
			name:   "declare statement",
			input:  "DECLARE\n @foo nvarchar(128) = 345",
			output: "\x1b[96m\x1b[40mDECLARE\x1b[0m",
		},
		{
			name:   "param values",
			input:  "@p1 ignore_this",
			output: "\x1b[97m\x1b[40m@p1\x1b[0m",
		},
		{
			name:   "whitespace ignored",
			input:  "  \n\n SELECT 1",
			output: "\x1b[96m\x1b[40mSELECT\x1b[0m",
		},
		{
			name:   "non-SQL output",
			input:  "some log statement",
			output: "some log statement",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := strings.Builder{}
			log := &SQLLogger{
				Logger:    log.New(&buf, "", 0),
				ignoreTTY: true,
			}

			log.Printf(tc.input)
			if !assert.Contains(t, buf.String(), tc.output) {
				assert.Equal(t, "", buf.String())
			}
		})
	}
}

func TestTransactionColoring(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "begin transaction",
			input:  "BEGIN TRANSACTION 001",
			output: "\x1b[94mBEGIN TRANSACTION",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := strings.Builder{}
			log := &SQLLogger{
				Logger:    log.New(&buf, "", 0),
				ignoreTTY: true,
			}

			log.Printf(tc.input)
			if !assert.Contains(t, buf.String(), tc.output) {
				assert.Equal(t, "", buf.String())
			}
		})
	}
}

func TestLogNotColorizedToFile(t *testing.T) {
	input := "SELECT 1"

	buf, _ := ioutil.TempFile("", "TestLogNotColorized_*.out")
	defer os.Remove(buf.Name())

	log := &SQLLogger{Logger: log.New(buf, "", 0)}

	log.Println(input)

	output, _ := ioutil.ReadFile(buf.Name())
	assert.Contains(t, string(output), "SELECT 1")
}
