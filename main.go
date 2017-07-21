package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"os/user"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/zieckey/goini"

	_ "github.com/SAP/go-hdb/driver"
)

const (
	driverName = "hdb"
)

var iniLocation string
var backupLocation string
var backupPrefix string

func main() {
	currentUser, _ := user.Current()
	iniLocation = fmt.Sprintf("%s/.hdbsql", currentUser.HomeDir)
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
	backupLocation, ok := ini.SectionGet("backup", "location")
	if !ok {
		fmt.Printf("Failed to read password from %s\n", iniLocation)
		return
	}
	backupPrefix, ok := ini.SectionGet("backup", "prefix")
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

	currentTime := time.Now()
	backupname := fmt.Sprintf("%s-%d-%d-%d_%d.%d", backupPrefix, currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute())
	backupquery := fmt.Sprintf("BACKUP DATA USING FILE ('%s','%s')", backupLocation, backupname)

	fmt.Println(backupquery)
	rows, err := db.Query(backupquery)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	uploadBackups(backupLocation, backupname)
}

func uploadBackups(backuplocation string, backupname string) {
	ini := goini.New()
	err := ini.ParseFile(iniLocation)
	if err != nil {
		fmt.Printf("parse INI file %v failed : %v\n", iniLocation, err.Error())
		return
	}
	region, ok := ini.SectionGet("s3", "region")
	if !ok {
		fmt.Printf("Failed to read s3 region from %s\n", iniLocation)
		return
	}
	bucket, ok := ini.SectionGet("s3", "bucket")
	if !ok {
		fmt.Printf("Failed to read s3 bucket from %s\n", iniLocation)
		return
	}
	fileFilter := fmt.Sprintf("%s/%s*", backuplocation, backupname)
	files, _ := filepath.Glob(fileFilter)
	for _, file := range files {
		binFile, err := os.Open(file)
		if err != nil {
			fmt.Printf("err opening file: %s", err)
		}
		defer binFile.Close()

		fileInfo, _ := binFile.Stat()
		var size = fileInfo.Size()
		buffer := make([]byte, size)
		binFile.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)
		fmt.Printf("Uploading %s to %s:%s\n", file, region, bucket)
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String(region)},
		}))
		uploader := s3manager.NewUploader(sess)

		upload, err := uploader.Upload(&s3manager.UploadInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(file),
			Body:        fileBytes,
			ContentType: aws.String(fileType),
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Upload URL: %s\n", upload.Location)
	}

}
