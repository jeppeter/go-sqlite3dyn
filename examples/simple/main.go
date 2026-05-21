package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"github.com/jeppeter/go-sqlite3dyn"
	"github.com/tebeka/atexit"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	Sqlexec_handler(nil, nil, nil)
}

func LoadParser(parser *extargsparse.ExtArgsParse) (err error) {
	var commandline_fmt string
	var commandline string
	var dllfile string
	commandline_fmt = `{
		"dllfile" : "%s",
		"sqlexec<Sqlexec_handler>##dllname funcname numfunc params ... to call functions##" : {
			"$" : "+"
		}
	}`

	if runtime.GOOS == "windows" {
		dllfile = ".\\sqlite3.dll"
		dllfile = strings.Replace(dllfile, "\\", "\\\\", -1)
	} else {
		dllfile = filepath.Join(".", "libsqlite3.so")
	}

	commandline = fmt.Sprintf(commandline_fmt, dllfile)
	err = parser.LoadCommandLineString(commandline)
	return
}

func callback_func(c uintptr, cols []string, names []string) (err error) {
	var i int
	for i = 0; i < len(names); i += 1 {
		if i > 0 {
			fmt.Fprintf(os.Stdout, " ")
		}
		fmt.Fprintf(os.Stdout, "%s", names[i])
	}
	fmt.Fprintf(os.Stdout, "\n")
	for i = 0; i < len(cols); i += 1 {
		if i > 0 {
			fmt.Fprintf(os.Stdout, " ")
		}
		fmt.Fprintf(os.Stdout, "%s", cols[i])
	}
	fmt.Fprintf(os.Stdout, "\n")
	err = nil
	return
}

func Sqlexec_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) (err error) {
	var dllfile string
	var sarr []string
	var conn *sqlite3dyn.Sqlite3BaseConn = nil
	err = nil
	if ns == nil {
		return
	}

	dllfile = ns.GetString("dllfile")

	err = sqlite3dyn.InitDll(dllfile)
	if err != nil {
		return
	}

	sarr = ns.GetArray("subnargs")
	if len(sarr) < 2 {
		err = fmt.Errorf("need dbfile sql")
		return
	}

	conn, err = sqlite3dyn.ConnSqlite3(sarr[0])
	if err != nil {
		return
	}

	err = conn.Exec(sarr[1], uintptr(0), callback_func)
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "exec [%s].[%s] succ\n", sarr[0], sarr[1])
	err = nil
	return
}

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		atexit.Exit(5)
	}

	err = LoadParser(parser)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		atexit.Exit(5)
	}

	_, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		atexit.Exit(4)
	}
	atexit.Exit(0)
}
