package sqlite3dyn

import (
	"fmt"
	"./internal/dlfunc"
)

var (
	sqlite3_dll     *dlfunc.DllLib  = nil
	sqlite3_open_v2 *dlfunc.DllFunc = nil
	sqlite3_exec    *dlfunc.DllFunc = nil
	sqlite3_close   *dlfunc.DllFunc = nil
	sqlite3_free    *dlfunc.DllFunc = nil
)

func InitDll(dllname string) (err error) {
	defer func() {
		if err != nil {
			sqlite3_open_v2 = nil
			sqlite3_exec = nil
			sqlite3_close = nil
			sqlite3_free = nil
			sqlite3_dll = nil
		}
	}

	sqlite3_dll, err = dlfunc.LoadDll(dllname)
	if err != nil {
		return
	}

	sqlite3_open_v2, err = sqlite3_dll.GetFunc("sqlite3_open_v2")
	if err != nil {
		return
	}

	sqlite3_exec, err = sqlite3_dll.GetFunc("sqlite3_exec")
	if err != nil {
		return
	}

	sqlite3_close, err = sqlite3_dll.GetFunc("sqlite3_close")
	if err != nil {
		return
	}

	sqlite3_free, err = sqlite3_dll.GetFunc("sqlite3_free")
	if err != nil {
		return
	}
	return
}

type Sqlite3BaseConn struct {
	dbconn uintptr
}

func (ptr *Sqlite3BaseConn) Close() {
	if sqlite3_close == nil {
		return
	}

	if ptr.dbconn != uintptr(0) {
		sqlite3_close.CallN(1,ptr.dbconn)
		ptr.dbconn = uintptr(0)
	}
	return
}

func ConnSqlite3(dsn string) (ptr *Sqlite3BaseConn, err error) {
	var flags uintptr = uintptr(SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE)
	var pdb uintptr = uintptr(0)
	var ppdb uintptr = uintptr(unsafe.Pointer(&pdb))
	var dbname uintptr = dlfunc.MakeCString(dsn)
	ptr = nil
	err = nil
	if sqlite3_open_v2 == nil {
		err = fmt.Errorf("not call InitDll succ")
		return
	}

	retval, err = sqlite3_open_v2.CallN(4,dbname,ppdb,flags,uintptr(0))
	if err != nil {
		return
	}

	if retval != uintptr(SQLITE_OK) {
		err = fmt.Errorf("open %s error %d", dsn,retval)
		return
	}

	ptr = &Sqlite3BaseConn{}
	ptr.dbconn = pdb
	err = nil
	runtime.SetFinalizer(ptr,(*Sqlite3BaseConn).Close)
	return	
}

type execCallArgs struct {
	innerarg uintptr
	callback func(uintptr,[]string,[]string) error
}

func new_exec_args(arg uintptr, callback func(uintptr,[]string,[]string) error) (retp *execCallArgs,err error) {
	retp = &execCallArgs{}
	retp.innerarg = arg
	retp.callback = callback
	err = nil
	return
}

func (ptr *Sqlite3BaseConn) Exec(sqlstr string,callarg uintptr, callback func(uintptr,[]string,[]string) error) ( err error) {
	var narg uintptr
	var execarg *execCallArgs = nil
	var retval uintptr
	var errmsg uintptr = uintptr(0)
	var perrmsg uintptr = uintptr(unsafe.Pointer(&errmsg))
	var sqlchar uintptr 

	if sqlite3_exec == nil {
		err = fmt.Errorf("not call InitDll succ")
		return
	}

	execarg, err = new_exec_args(callarg,callback)
	if err != nil {
		return
	}
	narg = uintptr(unsafe.Pointer(execarg))
	sqlchar = dlfunc.MakeCString(sqlstr)

	retval,err = sqlite3_exec.CallN(5,ptr.dbconn,sqlchar,new_callback_func(),narg,perrmsg)
	if err != nil {
		if errmsg != uintptr(0) {
			sqlite3_free.CallN(1,errmsg)
			errmsg = uintptr(0)
		}
		return
	}

	if retval != uintptr(SQLITE_OK) {
		err = fmt.Errorf("exec [%s] error %d", sqlstr, retval)
		if errmsg != uintptr(0) {
			sqlite3_free.CallN(1,errmsg)
			errmsg = uintptr(0)
		}
		return
	}

	/*ok for the return value*/
	return
}

func (*ptr Sqlite3BaseConn) Close() {
	if ptr != nil && sqlite3_close != nil && ptr.dbconn != uintptr(0) {
		sqlite3_close.CallN(1,ptr.dbconn)
		ptr.dbconn = uintptr(0)
	}
	return
}
