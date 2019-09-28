package storage 

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"

	// Tied to postgreSQL
	_ "github.com/lib/pq"
)

func nullPrintf(f string, a ...interface{}) (count int, err error) {
	return
}

// ExecuteSQLFile will run a .sql file against PostgreSQL
func ExecuteSQLFile(db *sql.DB, sqlFile string, verbose bool) error {
	log := nullPrintf
	if !verbose {
		log = fmt.Printf
	}
	file, err := ioutil.ReadFile(sqlFile)

	if err != nil {
		return err
	}

	requests := strings.Split(string(file), ";")

	for line, request := range requests {
		if line == len(requests)-1 {
			continue
		}
		log("\nExecuting SQL line %d: %s", line, request)
		result, err := db.Exec(request)
		if err != nil {
			log("\nSQL line %d Error: %s\n", line, err.Error())
		} else {
			log("\nSQL line %d: successful, ", line)
			if rowsCount, rcError := result.RowsAffected(); rcError == nil {
				log("%d Rows affected\n", rowsCount)
			} else {
				log("Unable to read rows affected: Error: %s\n", rcError.Error())
			}
		}
	}
	return nil
}

// SetupDatabase will setup a postgre backend
func SetupDatabase(backend, connectionString, setupScript string) (*sql.DB, error) {

	db, err := sql.Open(backend, connectionString)
	if err != nil {
		return nil, err
	}
	return db, ExecuteSQLFile(db, setupScript, true)
	// defer db.Close()
}
