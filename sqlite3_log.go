package sqlite3dyn

import (
	"fmt"
	"os"
	"runtime"
)

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_out_string_total(level int, fmtstr string, a ...interface{}) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr, a...)
	return outstr
}

func log_inner_function(dbglvl int, level int, fmtstr string, a ...interface{}) {
	outstr := format_out_string_total(level+1, fmtstr, a...)
	if dbglvl < 30 {
		fmt.Fprintf(os.Stderr, "%s\n", outstr)
	}
	return
}
func logTrace(fmt string, a ...interface{}) {
	log_inner_function(10, 2, fmt, a...)
}
