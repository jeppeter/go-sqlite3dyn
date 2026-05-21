package sqlite3dyn

import (
    "github.com/jeppeter/go-sqlite3dyn/internal/dlfunc"
    "unsafe"
)

/*
#include <stdint.h>
extern int lx_ptr_callstk(int* ptr,int argc,char** argv,char** argcols);
*/
import "C"

//export lx_ptr_callstk
func lx_ptr_callstk(ptr *C.int, argc C.int, argv **C.char, argvcols **C.char) C.int {
    var stks []string = []string{}
    var cols []string = []string{}
    var curptr uintptr
    var i int
    var pval *execCallArgs
    var retval C.int = 0
    var err error

    args := (*[1 << 30]*byte)(unsafe.Pointer(argv))
    argcols := (*[1 << 30]*byte)(unsafe.Pointer(argvcols))
    pval = (*execCallArgs)(unsafe.Pointer(ptr))

    if pval.callback == nil {
        return
    }

    for i = 0; i < int(argc); i += 1 {
        curptr = uintptr(unsafe.Pointer(args[i]))
        stks = append(stks, dlfunc.MakeGoStringFromPointer(curptr))
        curptr = uintptr(unsafe.Pointer(argcols[i]))
        cols = append(cols, dlfunc.MakeGoStringFromPointer(curptr))
    }

    err = pval.callback(pval.innerarg, stks, cols)
    if err != nil {
        retval = 1
    }
    return retval
}

func new_callback_func() (retptr uintptr) {
    return uintptr(C.lx_ptr_callstk)
}
