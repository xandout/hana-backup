package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"os/user"

	"github.com/zieckey/goini"

	_ "github.com/SAP/go-hdb/driver"
)

const (
	driverName = "hdb"
)

func main() {
	currentUser, _ := user.Current()
	iniLocation := fmt.Sprintf("%s/.hdbsql", currentUser.HomeDir)
	ini := goini.New()
	err := ini.ParseFile(iniLocation)
	if err != nil {
		fmt.Printf("parse INI file %v failed : %v\n", iniLocation, err.Error())
		return
	}
	host, ok := ini.SectionGet("hana", "host")
	if !ok {
		fmt.Printf("Failed to read host from %s\n", iniLocation)
		return
	}
	port, ok := ini.SectionGet("hana", "port")
	if !ok {
		fmt.Printf("Failed to read port from %s\n", iniLocation)
		return
	}
	user, ok := ini.SectionGet("hana", "user")
	if !ok {
		fmt.Printf("Failed to read user from %s\n", iniLocation)
		return
	}
	password, ok := ini.SectionGet("hana", "password")
	if !ok {
		fmt.Printf("Failed to read password from %s\n", iniLocation)
		return
	}

	var hdbDsn = fmt.Sprintf("hdb://%s:%s@%s:%s", user, password, host, port)

	db, err := sql.Open(driverName, hdbDsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	currentTtime := time.Now()
	backupname := fmt.Sprintf("%d-%d-%d_%d.%d_BACKUP", currentTtime.Year(), currentTtime.Month(), currentTtime.Day(), currentTtime.Hour(), currentTtime.Minute())
	backupquery := fmt.Sprintf("BACKUP DATA USING FILE ('%s')", backupname)

	fmt.Println(backupquery)
	rows, err := db.Query(backupquery)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	// defer rows.Close()
	// colNames, err := rows.Columns()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// cols := make([]interface{}, len(colNames))
	// colPtrs := make([]interface{}, len(colNames))
	// for i := 0; i < len(colNames); i++ {
	// 	colPtrs[i] = &cols[i]
	// }
	// for rows.Next() {
	// 	var myMap = make(map[string]interface{})

	// 	err = rows.Scan(colPtrs...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	for i, col := range cols {
	// 		myMap[colNames[i]] = col
	// 	}
	// 	// Do something with the map
	// 	for key, val := range myMap {
	// 		fmt.Println("Key:", key, "Value Type:", reflect.TypeOf(val))
	// 	}
	// }
}
