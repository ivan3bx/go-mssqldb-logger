# go-mssqldb-logger (colorized logs for SQLServer driver for Go)

This is a simple wrapper around the logger interface for [go-mssqldb](https://github.com/denisenkom/go-mssqldb/), adding syntax highlighting to SQL output.

## Installation

`go get -u github.com/ivan3bx/go-mssqldb-logger`

## Features

* Colorized SQL output for debugging!
* Suppresses colorization for non-TTY logs (i.e. writing to a file or dumb terminal)
* Special handling for logrus (will obey colorization option in logrus.TextFormatter)

## Sample Code

See the docs for `go-mssqldb` for how to set up a DB connection, particularly how to set the logging level to get generated output in your logs.

This (below) is a simple use of this package (within the context of sqlx, just to illustrate how to access the driver):

```go
    db = Connection{sqlx.MustOpen("mssql", connectionString)}

    if mssql, ok := db.Driver().(*mssql.Driver); ok {
        mssql.SetLogger(&mssqllog.SQLLogger{})
    }
```

### Variation #1 (standard log package)

```go
    logger := log.New(os.Stdout, "", log.LstdFlags)
    log := &mssqllog.SQLLogger{Logger: logger}
    mssql.SetLogger(log)
```

### Variation #2 (logrus)

```go
    logger := logrus.New()
    log := &mssqllog.SQLLogger{Logger: logger}
    mssql.SetLogger(log)
```
