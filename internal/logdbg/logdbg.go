package logdbg

import (
	"fmt"
	"os"
	"reflect"
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

func format_out_string_singal(level int, fmtstr string) string {
	outstr := format_out_stack((level + 1))
	outstr += fmt.Sprintf(fmtstr)
	return outstr
}

func format_out_stack_data(data []byte, fmtstr string, a ...interface{}) string {
	var i, lasti int
	//var r rune
	//var p []byte
	var j int
	outstr := ""
	if fmtstr != "" {
		outstr += fmt.Sprintf(fmtstr, a...)
	}
	lasti = 0
	//p = make([]byte, 1)
	for i = 0; i < len(data); i++ {
		if (i % 16) == 0 {
			if i > 0 {
				outstr += "    "
				for i != lasti {
					//p[0] = data[lasti]
					//r, j = utf8.DecodeRune(p)
					if data[lasti] < ' ' || data[lasti] > '~' {
						outstr += "."
					} else {
						outstr += fmt.Sprintf("%c", data[lasti])
					}
					lasti++
				}
			}
			outstr += fmt.Sprintf("\n[0x%08x]", i)
		}
		outstr += fmt.Sprintf(" 0x%02x", data[i])
	}

	if lasti != i {
		j = i
		for (j % 16) != 0 {
			outstr += "     "
			j++
		}

		outstr += "    "
		for lasti != i {
			//p[0] = data[lasti]
			//r, j = utf8.DecodeRune(p)
			//if j != 1 || !strconv.IsPrint(r) {
			if data[lasti] < ' ' || data[lasti] > '~' {
				outstr += "."
			} else {
				outstr += fmt.Sprintf("%c", data[lasti])
			}
			lasti++
		}
		outstr += "\n"
	}

	return outstr
}

const (
	def_stacklevel_added = 3
)

func format_out_string_cap(a ...interface{}) string {
	var stacklevel int = def_stacklevel_added
	var vaargs []interface{}
	var fmtstr string = ""
	var ct string
	if len(a) > 0 {
		switch v := a[0].(type) {
		case int:
			stacklevel = a[0].(int)
			stacklevel += def_stacklevel_added
			if len(a) > 2 {
				vaargs = a[2:]
			}
			if len(a) > 1 {
				ct = reflect.TypeOf(a[1]).Name()
				if ct == "string" {
					fmtstr = a[1].(string)
				} else {
					fmtstr = "unknown type string"
				}
			}
		case string:
			if len(a) > 1 {
				vaargs = a[1:]
			}
			fmtstr = a[0].(string)
		default:
			fmtstr = fmt.Sprintf("unknown type [%s]", v)
		}
	}

	outstr := format_out_stack(stacklevel)
	if len(vaargs) == 0 {
		outstr += fmt.Sprintf(fmtstr)
	} else {
		outstr += fmt.Sprintf(fmtstr, vaargs...)
	}
	return outstr
}

func format_out_data_cap(a ...interface{}) string {
	var stacklevel int = def_stacklevel_added
	var vaargs []interface{}
	var data []byte = []byte{}
	var fmtstr string = ""
	var ct string
	if len(a) > 0 {
		switch v := a[0].(type) {
		case int:
			stacklevel = a[0].(int)
			stacklevel += def_stacklevel_added
			if len(a) > 3 {
				vaargs = a[3:]
			}
			if len(a) > 1 {
				switch a[1].(type) {
				case []byte:
					data = a[1].([]byte)
				}
			}

			if len(a) > 2 {
				ct = reflect.TypeOf(a[2]).Name()
				if ct == "string" {
					fmtstr = a[2].(string)
				} else {
					fmtstr = fmt.Sprintf("unknown type string [%s]", ct)
				}
			}
		case []byte:
			if len(a) > 2 {
				vaargs = a[2:]
			}
			data = a[0].([]byte)
			if len(a) > 1 {
				ct = reflect.TypeOf(a[1]).Name()
				if ct == "string" {
					fmtstr = a[1].(string)
				} else {
					fmtstr = fmt.Sprintf("unknown type [%s]", ct)
				}
			}

		default:
			fmtstr = fmt.Sprintf("unknown type [%s]", v)
			fmt.Printf("%s\n", fmtstr)
		}
	}

	outstr := format_out_stack(stacklevel)
	if len(vaargs) == 0 {
		outstr += format_out_stack_data(data, fmtstr)
	} else {
		outstr += format_out_stack_data(data, fmtstr, vaargs...)
	}
	return outstr
}

func output_stderr(level int, outs string) {
	fmt.Fprintf(os.Stderr, "%s\n", outs)
	os.Stderr.Sync()
}

func Error(a ...interface{}) int {
	var retval int = 0
	outstr := "<ERROR>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	output_stderr(0, outstr)
	return retval
}

func ErrorBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<ERROR>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	output_stderr(0, outstr)
	return retval
}

func Warn(a ...interface{}) int {
	var retval int = 0
	outstr := "<WARN>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	output_stderr(1, outstr)
	return retval
}

func WarnBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<WARN>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	output_stderr(1, outstr)
	return retval
}

func Info(a ...interface{}) int {
	var retval int = 0
	outstr := "<INFO>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	output_stderr(2, outstr)
	return retval
}

func InfoBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<INFO>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	output_stderr(2, outstr)
	return retval
}

func Debug(a ...interface{}) int {
	var retval int = 0
	outstr := "<DEBUG>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	output_stderr(3, outstr)
	return retval
}

func DebugBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<DEBUG>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	output_stderr(3, outstr)
	return retval
}

func Trace(a ...interface{}) int {
	var retval int = 0
	outstr := "<TRACE>"
	outstr += format_out_string_cap(a...)
	retval = len(outstr)
	output_stderr(4, outstr)
	return retval
}

func TraceBuffer(a ...interface{}) int {
	var retval int = 0
	outstr := "<TRACE>"
	outstr += format_out_data_cap(a...)
	retval = len(outstr)
	output_stderr(4, outstr)
	return retval
}
