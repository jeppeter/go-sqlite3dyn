package sqlite3dyn

import (
	"syscall"
)

func win_ptr_callstk(ptr uintptr, argc uintptr, argv uintptr, colname uintptr) uintptr {
	var stks []string = []string{}
	var cols []string = []string{}
	var curptr uintptr
	var i int
	var pval *PtrValues
	var retval uintptr = 0
	var err error

	args := (*[1 << 30]*byte)(unsafe.Pointer(argv))
	argcols := (*[1 << 30]*byte)(unsafe.Pointer(colname))
	pval = (*PtrValues)(unsafe.Pointer(ptr))
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
	return uintptr(syscall.NewCallback(win_ptr_callstk))
}
