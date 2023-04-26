package s3

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mylittleboxy/backup/pkg/configType"
)

func SendFile(config configType.Config, fileName string) error {

	// Open the file to be uploaded
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Failed to open file", err)
		return err
	}
	defer file.Close()

	// Read the uncompressed file data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to read file", err)
		return err
	}

	// Create a new buffer to store the compressed data
	var compressedData bytes.Buffer
	gz := gzip.NewWriter(&compressedData)
	defer gz.Close()

	// Write the uncompressed data to the gzip writer
	if _, err := gz.Write(data); err != nil {
		fmt.Println("Failed to write compressed data to gzip writer", err)
		return err
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("fr-par"),                                                           // replace with your Scaleway region
		Credentials: credentials.NewStaticCredentials(config.S3.AccessKey, config.S3.SecretKey, ""), // replace with your Scaleway access and secret keys
		Endpoint:    aws.String("https://s3.fr-par.scw.cloud"),                                      // replace with the endpoint for Scaleway's S3-compatible service
	})
	// Create an S3 client
	svc := s3.New(sess)

	// Specify the S3 bucket and file path
	bucketName := "mylittleboxy-backup"

	// Upload the file to S3
	out, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("%s.gz", fileName)),
		Body:   bytes.NewReader(compressedData.Bytes()),
	})
	if err != nil {
		fmt.Println("Failed to upload file", err)
		return err
	}
	out.GoString()
	fmt.Println("File uploaded successfully")
	return nil
}
