package s3

// Borrowed from https://github.com/aws/aws-sdk-go-v2/blob/main/example/service/s3/listObjects/listObjects.go

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

// PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type putObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// PutFile uploads a file to an Amazon Simple Storage Service (Amazon S3) bucket
//
// Inputs
//
//		c is the context of the method call, which includes the AWS Region
//		api is the interface that defines the method call
//		input defines the input arguments to the service call.
//
// Output
//
// If success, a PutObjectOutput object containing the result of the service call and nil. Otherwise, nil and an error from the call to PutObject.
func PutFile(c context.Context, api putObjectAPI, input *s3.PutObjectInput) error {
	output, err := api.PutObject(c, input)
	if err != nil {
		log.Fatal().
			Str("service", "driftctl-s3").
			Msg("Error when writing file: ")
		return err
	}
	log.Info().
		Str("service", "driftctl-slack").
		Str("VersionId", *output.VersionId).
		Msg("Report uploaded to S3")

	return nil
}
