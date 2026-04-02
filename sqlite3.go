// Copyright (C) 2016 Samuel Melrose <sam@infitialis.com>.
//
// Based on work by Yasuhiro Matsumoto <mattn.jp@gmail.com>
// https://github.com/mattn/go-sqlite3
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sqlite3dyn

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func init() {
	registerLibrary()

	sql.Register("sqlite3", &SQLiteDriver{})
}

func Version() (libVersion string, libVersionNumber int, sourceID string) {
	return sqlite3_libversion(), sqlite3_libversion_number(), sqlite3_sourceid()
}

func errorString(err Error) string {
	return sqlite3_errstr(int(err.Code))
}

func (d *SQLiteDriver) Open(dsn string) (driver.Conn, error) {
	logTrace(" ")
	if !libraryRegistered {
		if libraryRegisterError != nil {
			return nil, libraryRegisterError
		} else {
			return nil, fmt.Errorf("Failed to register SQLite3 library, perhaps not found?")
		}
	}

	logTrace(" ")
	if sqlite3_threadsafe() == 0 {
		return nil, fmt.Errorf("sqlite library was not compiled for thread-safe operation")
	}

	var loc *time.Location
	logTrace(" ")
	txlock := "BEGIN"
	busyTimeout := 5000
	pos := strings.IndexRune(dsn, '?')
	if pos >= 1 {
		logTrace(" ")
		params, err := url.ParseQuery(dsn[pos+1:])
		if err != nil {
			return nil, err
		}

		// _loc
		if val := params.Get("_loc"); val != "" {
			if val == "auto" {
				loc = time.Local
			} else {
				loc, err = time.LoadLocation(val)
				if err != nil {
					return nil, fmt.Errorf("Invalid _loc: %v: %v", val, err)
				}
			}
		}

		logTrace(" ")
		// _busy_timeout
		if val := params.Get("_busy_timeout"); val != "" {
			iv, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid _busy_timeout: %v: %v", val, err)
			}
			busyTimeout = int(iv)
		}

		// _txlock
		if val := params.Get("_txlock"); val != "" {
			switch val {
			case "immediate":
				txlock = "BEGIN IMMEDIATE"
			case "exclusive":
				txlock = "BEGIN EXCLUSIVE"
			case "deferred":
				txlock = "BEGIN"
			default:
				return nil, fmt.Errorf("Invalid _txlock: %v", val)
			}
		}

		logTrace(" ")
		if !strings.HasPrefix(dsn, "file:") {
			dsn = dsn[:pos]
		}
		logTrace(" ")
	}

	var db sqlite3
	logTrace("dsn %s", dsn)
	rv := sqlite3_open_v2(dsn, &db, SQLITE_OPEN_FULLMUTEX|SQLITE_OPEN_READWRITE|SQLITE_OPEN_CREATE, "")
	if rv != 0 {
		logTrace("rv %d", rv)
		err := Error{Code: ErrNo(rv)}
		return nil, fmt.Errorf("Error Opening DB: %s", err)
	}
	logTrace(" ")
	if db == sqlite3(uintptr(0)) {
		return nil, fmt.Errorf("sqlite succeeded without returning a database")
	}

	logTrace(" ")
	rv = sqlite3_busy_timeout(db, busyTimeout)
	if rv != SQLITE_OK {
		err := Error{Code: ErrNo(rv)}
		return nil, fmt.Errorf("Error setting busy timeout: %s", err)
	}

	logTrace(" ")
	conn := &SQLiteConn{db: db, loc: loc, txlock: txlock}

	logTrace(" ")
	if d.ConnectHook != nil {
		if err := d.ConnectHook(conn); err != nil {
			return nil, fmt.Errorf("Error in Connect Hook: %s", err)
		}
	}
	logTrace(" ")
	runtime.SetFinalizer(conn, (*SQLiteConn).Close)
	return conn, nil
}
