package main

import (
	"database/sql"
	"fmt"
	"github.com/jeppeter/go-sqlite3dyn"
	"os"
)

func main() {
	path := "./test.db"

	//f, err := os.Create(path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	f.Close()

	fmt.Printf("file [%s]\n", path)
	var libver, libnum, sourceid = sqlite3dyn.Version()
	fmt.Printf("sqlite3 version [%v %v %v]\n", libver, libnum, sourceid)

	db, err := sql.Open(`sqlite3`, path)
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
	//os.Remove(`./test.db`)
}
