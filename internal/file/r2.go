package image

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"strings"
	"time"

	"image"
	"image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

const BUCKET_NAME = "gear-ecommerce"
const IMAGE_LOCATION = "image-gear-ecommerce.goldenfealla.dev"

func UploadImageJpeg(client *s3.Client, base64 string, fileName string) (*string, error) {
	// convert base64 to image
	imgReader := b64.NewDecoder(b64.StdEncoding, strings.NewReader(base64))
	img, _, err := image.Decode(imgReader)
	if err != nil {
		return nil, err
	}

	resizeImage := imaging.Resize(img, 200, 0, imaging.Lanczos)

	var jpegImage bytes.Buffer

	err = jpeg.Encode(&jpegImage, resizeImage, &jpeg.Options{
		Quality: 85,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(BUCKET_NAME),
		Body:        &jpegImage,
		Key:         aws.String(fmt.Sprintf("image/%v", fileName)),
		ContentType: aws.String("image"),
	})

	if err != nil {
		return nil, err
	}

	publicURL := fmt.Sprintf("%v/%v", IMAGE_LOCATION, fileName)

	return &publicURL, nil
}
