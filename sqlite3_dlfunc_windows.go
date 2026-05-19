package sqlite3dyn

import (
	"github.com/jeppeter/go-sqlite3dyn/internal/dlfunc"
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
	ptr = nil
	err = nil
	if sqlite3_open_v2 == nil {
		err = fmt.Errorf("not call InitDll succ")
		return
	}

	
}
