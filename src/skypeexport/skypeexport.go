package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"flag"
	"os"
)

func main() {
	var dbFileName = flag.String("db", "main.db", "Path to the database file")
	var chatName = flag.String("chatname", "", "Path to the database file")
	flag.Parse()

	if *chatName == "" {
		fmt.Println("-chatname is required")
		os.Exit(11)
	}

	err := dumpLogs(*dbFileName, *chatName)

	if err != nil {
		// FIXME: respond with an error code beside 0 here
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func dumpLogs(dbFileName string, chatName string) (error) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `
	select timestamp, author, body_xml
	from Messages
	where chatname=?
	order by timestamp asc;
	`

	stm, err := db.Prepare(query)
	if err != nil {
		return err
	}

	rows, err := stm.Query(chatName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var timestamp int64
		var author string
		var msg string
		rows.Scan(&timestamp, &author, &msg)

		t := time.Unix(timestamp, 0)

		fmt.Println(
			t.Format("2006-01-02 15:04:05"),
			author, msg)
	}
	rows.Close()

	return nil
}