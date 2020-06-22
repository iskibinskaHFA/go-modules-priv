// Package connect enables connection to Postgres Database and Returns Connection handle.
package connect

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jinzhu/gorm"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type dbInfo struct {
	host     string
	port     int
	user     string
	password string
	database string
}

//GetDBConnection returns Postgres connection handle
func GetDBConnection(test bool) *sql.DB {
	psqlInfo := getDBInfo(test)
	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)

	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func getSession() (session.Session, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-east-1")},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		panic(err)
	}
	return *sess, nil
}

func getValue(keyname string) (string, error) {
	sess, _ := getSession()
	ssmsvc := ssm.New(&sess, aws.NewConfig().WithRegion("us-east-1"))
	withDecryption := true
	param, err := ssmsvc.GetParameter(&ssm.GetParameterInput{
		Name:           &keyname,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return strings.TrimSpace(aws.StringValue(param.Parameter.Value)), nil
}

func getDBInfo(test bool) string {
	dbValues := getDBValues(test)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbValues.host, dbValues.port, dbValues.user, dbValues.password, dbValues.database)
	return psqlInfo
}

func getDBValues(test bool) *dbInfo {
	var info dbInfo
	var port string

	if info.host = os.Getenv("HOST"); !test {
		info.host, _ = getValue("/royalties_sys/MLC_HOST")
	}

	if info.user = os.Getenv("USER"); !test {
		info.user, _ = getValue("/royalties_sys/MLC_USER")
	}

	if info.password = os.Getenv("PASSWORD"); !test {
		info.password, _ = getValue("/royalties_sys/MLC_PASSWORD")
	}

	if port = os.Getenv("PORT"); !test {
		port, _ = getValue("/royalties_sys/MLC_PORT")
	}
	info.port, _ = strconv.Atoi(port)

	if info.database = os.Getenv("DATABASE"); !test {
		info.database, _ = getValue("/royalties_sys/MLC_DATABASE")
	}
	return &info
}

//GormConnection returns Gorm connection handle for either unitary test/regular run
func GormConnection(test bool) (*gorm.DB, *sql.DB) {
	dbValues := getDBValues(test)
	dsn := url.URL{
		User:     url.UserPassword(dbValues.user, dbValues.password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", dbValues.host, dbValues.port),
		Path:     dbValues.database,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err := gorm.Open("postgres", dsn.String())

	if err != nil {
		log.Println(err.Error())
	}

	dbase := db.DB()
	err = dbase.Ping()
	if err != nil {
		log.Println(err.Error())
	}

	defineHandler(db)

	return db, dbase
}

func defineHandler(*gorm.DB) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		switch defaultTableName {
		case "headers":
			fallthrough
		case "usage_summaries":
			fallthrough
		case "releases":
			fallthrough
		case "resources":
			fallthrough
		case "works":
			fallthrough
		case "works_writer":
			fallthrough
		case "usages":
			fallthrough
		case "writers":
			return "usage." + defaultTableName
		}
		return "royalty." + defaultTableName
	}
}
