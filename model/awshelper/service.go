package awshelper

import (
	"ajebackend/helper"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"mime/multipart"
)

func UploadDocument(file *multipart.FileHeader, fileName string) (*s3manager.UploadOutput, error){
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
	contentDisposition := fmt.Sprintf("inline; filename=\"%s\"", fileName)

	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(fileName),
		Body:   fileBody,
		ContentType: &contentType,
		ContentDisposition: &contentDisposition,
	})

	return up, err
}
