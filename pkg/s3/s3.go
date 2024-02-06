package s3

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mylittleboxy/backup/pkg/configType"
)

func SendFiles(config configType.Config, fileNames []string, archiveName string) error {

	// Create output file
	out, err := os.Create(archiveName)
	if err != nil {
		log.Fatalln("Error writing archive:", err)
	}
	defer out.Close()

	// Write the uncompressed data to the gzip writer
	err = createArchive(fileNames, out)
	if err != nil {
		log.Fatalln("Error creating archive:", err)
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
	archiveFile, err := os.Open(archiveName)
	// Upload the file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(archiveName),
		Body:   archiveFile,
	})
	if err != nil {
		fmt.Println("Failed to upload file", err)
		return err
	}
	fmt.Println("File uploaded successfully")
	return nil
}
func createArchive(files []string, buf io.Writer) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := addToArchive(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToArchive(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
