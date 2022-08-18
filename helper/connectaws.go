package helper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectAws() (*session.Session, error) {
	AccessKeyID := GetEnvWithKey("AWS_ID")
	SecretAccessKey := GetEnvWithKey("AWS_SECRET_KEY")
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String("ap-southeast-1"),
			Credentials: credentials.NewStaticCredentials(
				AccessKeyID,
				SecretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		return sess, err
	}
	return sess, nil
}