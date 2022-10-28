package main

import (
	"bytes"
	"encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tacg0909/meshitero-put-post/resize"
)

func putImage(base64Image string, postId string) error {
    imageBinary, err := base64.StdEncoding.DecodeString(base64Image)
    if err != nil {
        return err
    }
    largeImageBuf := bytes.NewBuffer(imageBinary)
    smallImageBuf := *largeImageBuf
    largeImage, err := resize.Resize(largeImageBuf, 1000)
    if err != nil {
        return err
    }
    err = putImageToBucket("large/" + postId + ".jpg", largeImage)
    if err != nil {
        return err
    }
    smallImage, err := resize.Resize(&smallImageBuf, 300)
    if err != nil {
        return err
    }
    err = putImageToBucket("small/" + postId + ".jpg", smallImage)
    return err
}

func putImageToBucket(key string, image bytes.Buffer) error {
    client := s3.New(session.New(), &aws.Config{
        Region: aws.String(os.Getenv("S3_BUCKET_REGION")),
    })
    _, err := client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
        Key: aws.String(key),
        Body: bytes.NewReader(image.Bytes()),
        ContentType: aws.String("image/jpeg"),
    })
    return err
}
