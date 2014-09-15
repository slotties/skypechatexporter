package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"flag"
	"os"
	"os/user"
	"errors"
	"io/ioutil"
)

func main() {
	var dbFileNameParam = flag.String("db", "", "Path to the database file")
	var chatName = flag.String("chatname", "", "Path to the database file")
	flag.Parse()

	if *chatName == "" {
		fmt.Println("-chatname is required")
		os.Exit(11)
	}

	dbFileName, err := locateDatabase(*dbFileNameParam)
	if err != nil {
		fmt.Println(err)
		os.Exit(12)
	}

	err = dumpLogs(dbFileName, *chatName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func locateDatabase(providedFileName string) (string, error) {
	if providedFileName != "" {
		_, err := os.Stat(providedFileName)
		return providedFileName, err
	}

	// Search for main.db in various locations:
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir
	possibleLocations := []string {
		// Skype on Windows for Desktop
		homeDir + "/AppData/Roaming/Skype",
		// Skype on Windows 8 (App)
		// TODO
		// Skype on Linux
		// TODO
	}

	for _, dirName := range possibleLocations {
		// Search for a main.db in the account directories.
		files, err := ioutil.ReadDir(dirName)
		// Skip any non-existent directory.
		if err == nil {
			for _, file := range files {
				mainDbFile, exists := containsMainDb(dirName, file)
				if exists {
					return mainDbFile, nil
				}
			}
		}
	}

	return "", errors.New("Could not locate main.db")
}

func containsMainDb(rootDir string, dir os.FileInfo) (string, bool) {
	if dir.IsDir() {
		mainDbFile := rootDir + "/" + dir.Name() + "/main.db"
		if _, err := os.Stat(mainDbFile); os.IsNotExist(err) {
			return "", false
		} else {
			return mainDbFile, true
		}
	} else {
		return "", false
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