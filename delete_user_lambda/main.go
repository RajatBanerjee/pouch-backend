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

func HandleLambdaEvent(ctx context.Context, event map[string]interface{}) (int64, error) {
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

	// perform a db.Query insert
	result, err := db.Exec("DELETE FROM  pouch_db.USER_UPLOADS WHERE id= ?", id)

	// if there is an error inserting, handle it
	if err != nil {
		log.Fatal(err)
	}

	no, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return no, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
