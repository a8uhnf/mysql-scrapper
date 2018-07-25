package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"database/sql"

	"github.com/ghodss/yaml"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Mysql struct {
		Username string `json:"username"`
		Dbname   string `json:"dbname"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Query    string `json:"query"`
	} `json:"mysql"`
}

func main() {
	log.Println("---- Start ----")

	b, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}
	cfg := &Config{}
	err = yaml.Unmarshal(b, cfg)
	fmt.Println("----- Check -----", fmt.Sprintf("%s:%s@tcp(%s)/%s",
	cfg.Mysql.Username,
	cfg.Mysql.Password,
	cfg.Mysql.Address,
	cfg.Mysql.Dbname))
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Address,
		cfg.Mysql.Dbname))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(1000)
	db.SetMaxOpenConns(-1)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return

	rows, err := db.Query(cfg.Mysql.Query)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
}
