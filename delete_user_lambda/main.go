package main

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

type UserData struct {
	UserId          string `json:"userId" db:"user_id"`
	FileName        string `json:"fileName" db:"file_name"`
	InsertTs        string `json:"insTs" db:"ins_ts"`
	UpdateTs        string `json:"updTs" db:"upd_ts"`
	FileDescription string `json:"fileDesc" db:"file_desc"`
}

func HandleLambdaEvent(userData UserData) (int64, error) {
	rds_host := os.Getenv("rds_host")
	name := os.Getenv("rds_user_name")
	password := os.Getenv("rds_password")
	db_name := os.Getenv("rds_db_name")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", name, password, rds_host, db_name))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)

	// perform a db.Query insert
	result, err := db.Exec("DELETE FROM  pouch_db.USER_UPLOADS WHERE user_id=? and file_name= ?", userData.UserId, userData.FileName)

	// if there is an error inserting, handle it
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return id, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
