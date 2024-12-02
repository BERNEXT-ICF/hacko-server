package adapter

import (
	"context"
	"hacko-app/internal/infrastructure/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

func WithDigihubStorage() Option {
	return func(a *Adapter) {
		env := config.Envs.HackoStorage

		a.HackoStorage = s3.New(s3.Options{
			BaseEndpoint: aws.String(env.Endpoint),
			Region:       env.Region,
			Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(env.Key, env.Secret, "")),
		})

		_, err := a.HackoStorage.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil {
			log.Fatal().Err(err).Msgf("Error while connecting to crowners storage: %s", env.Endpoint)
		}

		log.Info().Msg("Digihub storage connected")
	}
}
