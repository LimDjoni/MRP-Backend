package awshelper

import (
	"ajebackend/helper"
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadDocument(file *multipart.FileHeader, fileName string) (*s3manager.UploadOutput, error) {
	fileBody, openFileErr := file.Open()

	if openFileErr != nil {
		return nil, openFileErr
	}

	// Save file to root directory:
	sess, sessErr := helper.ConnectAws()

	if sessErr != nil {
		return nil, sessErr
	}

	uploader := s3manager.NewUploader(sess)
	MyBucket := helper.GetEnvWithKey("AWS_BUCKET_NAME")

	contentType := "application/pdf"
	splitFile := strings.Split(fileName, "/")
	contentDisposition := fmt.Sprintf("inline; filename=\"%s\"", splitFile[len(splitFile)-1])

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:             aws.String(MyBucket),
		Key:                aws.String(fileName),
		Body:               fileBody,
		ContentType:        &contentType,
		ContentDisposition: &contentDisposition,
	})

	return up, err
}

func DeleteDocument(key string) (bool, error) {
	AccessKeyID := helper.GetEnvWithKey("AWS_ID")
	SecretAccessKey := helper.GetEnvWithKey("AWS_SECRET_KEY")

	newSession, errNewSession := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			AccessKeyID,
			SecretAccessKey,
			"", // a token will be created when the session it's used.
		),
	})

	if errNewSession != nil {
		return false, errNewSession
	}

	svc := s3.New(newSession)

	bucket := helper.GetEnvWithKey("AWS_BUCKET_NAME")
	keyFile := key
	request := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &keyFile,
	}

	// Save file to root directory:
	_, err := svc.DeleteObject(request)

	if err != nil {
		return false, err
	}
	return true, err
}

func DeleteDocumentBatch(key string) (bool, error) {
	ctx := context.Background()
	AccessKeyID := helper.GetEnvWithKey("AWS_ID")
	SecretAccessKey := helper.GetEnvWithKey("AWS_SECRET_KEY")

	newSession, errNewSession := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			AccessKeyID,
			SecretAccessKey,
			"", // a token will be created when the session it's used.
		),
	})

	if errNewSession != nil {
		return false, errNewSession
	}

	svc := s3.New(newSession)

	bucket := helper.GetEnvWithKey("AWS_BUCKET_NAME")
	keyFile := key
	request := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(keyFile),
	}

	iter := s3manager.NewDeleteListIterator(svc, request)
	//_, err := svc.DeleteObject(request)

	// use the iterator to delete the files.
	if err := s3manager.NewBatchDeleteWithClient(svc).Delete(ctx, iter); err != nil {
		return false, err
	}

	return true, nil
}
