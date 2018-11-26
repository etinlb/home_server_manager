package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func initDatabase(database string) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
       CREATE TABLE IF NOT EXISTS plex_servers (url text NOT NULL PRIMARY KEY, token text);
    `

	_, err = db.Exec(sqlStmt)

	if err != nil {
		return err
	}

	return nil
}

func insertPlexData(database, url, token string) error {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return errors.New("Can't open database")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return errors.New("Can't begin transaction")
	}

	stmt, err := tx.Prepare("insert into plex_servers(url, token) values(?, ?)")
	if err != nil {
		return errors.New("Can't prepare statement")
	}

	defer stmt.Close()

	_, err = stmt.Exec(url, token)
	if err != nil {
		return errors.New("Can't execute statement ")
	}
	tx.Commit()

	return nil
}

func getPlexData(database string) (string, string, error) {

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// sqlStmt := `
	//    create table foo (id integer not null primary key, name text);
	//    delete from foo;
	//    `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Printf("%q: %s\n", err, sqlStmt)
	// 	return
	// }

	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// for i := 0; i < 100; i++ {
	// 	_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// tx.Commit()

	rows, err := db.Query("select token, url from plexData")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var token string
		var url string
		err = rows.Scan(&token, &url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(token, url)
		return token, url, nil
	}
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// stmt, err = db.Prepare("select name from foo where id = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// var name string
	// err = stmt.QueryRow("3").Scan(&name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(name)

	// _, err = db.Exec("delete from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// rows, err = db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return "", "", errors.New("No plex data")
}
