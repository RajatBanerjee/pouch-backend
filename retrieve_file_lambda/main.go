package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strconv"

	"time"
)

type FileInfo struct {
	Id              int    `json:"id" db:"id"`
	UserId          string `json:"userId" db:"user_id"`
	Username        string `json:"userName" db:"user_name"`
	FileName        string `json:"fileName" db:"file_name"`
	InsertTs        string `json:"insTs" db:"ins_ts"`
	UpdateTs        string `json:"updTs" db:"upd_ts"`
	FileDescription string `json:"fileDesc" db:"file_desc"`
	FilePath string `json:"filePath" db:"file_path"`
}

func HandleLambdaEvent(ctx context.Context, event map[string]interface{}) (*FileInfo, error) {
	rds_host := os.Getenv("rds_host")
	name := os.Getenv("rds_user_name")
	password := os.Getenv("rds_password")
	db_name := os.Getenv("rds_db_name")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", name, password, rds_host, db_name))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)



	v := event["params"].(map[string]interface{})

	y := v["path"].(map[string]interface{})

	id, err := strconv.Atoi(y["id"].(string))
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println("id====", id)
	// perform a db.Query insert
	result := db.QueryRow("select id ,user_id,user_name,file_name, ins_ts, upd_ts,file_desc, file_path FROM pouch_db.USER_UPLOADS WHERE id= ?", id)

	// if there is an error inserting, handle it
	if err != nil {
		log.Fatal(err)
	}

	var a int
	var b string
	var c string
	var d string
	var e string
	var f string
	var g string
	var h string
	err = result.Scan(&a, &b, &c, &d, &e, &f, &g, &h)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	resp := FileInfo{
		Id:              a,
		UserId:          b,
		Username:        c,
		FileName:        d,
		InsertTs:        e,
		UpdateTs:        f,
		FileDescription: g,
		FilePath: h,
	}
	return &resp, nil
}

func main (){
	lambda.Start(HandleLambdaEvent)
}