package main

import (
	"database/sql"
	"fmt"
	"github.com/jeppeter/go-sqlite3dyn"
	"time"
)

func main() {
	var filename string
	resetTime := time.Now()

	fmt.Println(sqlite3dyn.Version())
	filename = "file:" + resetTime.Format("2006-01-02") + "?mode=memory&cache=shared"
	fmt.Printf("fname [%s]\n", filename)

	db, err := sql.Open(`sqlite3`, filename)
	if err != nil {
		panic(err)
	}

	r, err := db.Exec(`CREATE TABLE test (
		id integer PRIMARY KEY NOT NULL,
		name varchar(30)
	)`)
	if err != nil {
		panic(err)
	}

	_ = r

	r, err = db.Exec(`INSERT INTO test(name) VALUES ('first') `)
	if err != nil {
		panic(err)
	}
	_, err = r.LastInsertId()
	if err != nil {
		panic(err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		panic(err)
	}

	r, err = db.Exec(`INSERT INTO test(name) VALUES ('second') `)
	if err != nil {
		panic(err)
	}
	_, err = r.LastInsertId()
	if err != nil {
		panic(err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		panic(err)
	}

	db.Close()
}
