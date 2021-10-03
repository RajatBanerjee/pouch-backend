package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"

	_ "github.com/go-sql-driver/mysql"
)

type UserData struct {
	UserId          int    `json:"userId" db:"USER_ID"`
	FileName        string `json:"fileName" db:"FILE_NAME"`
	InsertTs        string `json:"insTs" db:"INS_TS"`
	UpdateTs        string `json:"updTs" db:"UPD_TS"`
	FileDescription string `json:"fileDesc" db:"FILE_DESC"`
}

func HandleLambdaEvent() ([]UserData, error) {
	rds_host := os.Getenv("rds_host")
	name := os.Getenv("rds_user_name")
	password := os.Getenv("rds_password")
	db_name := os.Getenv("rds_db_name")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", name, password, rds_host, db_name))
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)

	rows, err := db.Query("SELECT USER_ID,FILE_NAME, INS_TS, UPD_TS,FILE_DESC FROM pouch_db.USER_UPLOADS;")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var resp []UserData

	for rows.Next() {
		var a int
		var b string
		var c string
		var d string
		var e string
		err := rows.Scan(&a, &b, &c, &d, &e)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		resp = append(resp, UserData{a, b, c, d, e})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return resp, nil
}
func main() {
	lambda.Start(HandleLambdaEvent)
}
