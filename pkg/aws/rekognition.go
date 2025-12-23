package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

type RekognitionHelperAdapter interface {
	CreateCollection(ctx context.Context, clientId string) error
	CreateUser(ctx context.Context, collectionId, userId string) error
	AssociateFace(ctx context.Context, collectionId, userId string, minimumSimilarityPercentage float32, imageBytes []byte) (*string, error)
	SearchUsersByImage(ctx context.Context, collectionId string, minimumSimilarityPercentage float32, imageBytes []byte) ([]UserMatch, error)
	DisassociateFace(ctx context.Context, collectionId, userId, faceId string) error
}

type Rekognition interface {
	CreateCollection(ctx context.Context, params *rekognition.CreateCollectionInput, optFns ...func(*rekognition.Options)) (*rekognition.CreateCollectionOutput, error)
	CreateUser(ctx context.Context, params *rekognition.CreateUserInput, optFns ...func(*rekognition.Options)) (*rekognition.CreateUserOutput, error)
	IndexFaces(ctx context.Context, params *rekognition.IndexFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.IndexFacesOutput, error)
	AssociateFaces(ctx context.Context, params *rekognition.AssociateFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.AssociateFacesOutput, error)
	DisassociateFaces(ctx context.Context, params *rekognition.DisassociateFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.DisassociateFacesOutput, error)
	DeleteFaces(ctx context.Context, params *rekognition.DeleteFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.DeleteFacesOutput, error)
	SearchUsersByImage(ctx context.Context, params *rekognition.SearchUsersByImageInput, optFns ...func(*rekognition.Options)) (*rekognition.SearchUsersByImageOutput, error)
}

type UserMatch struct {
	Id         *string
	Similarity *float32
}

type rekognitionHelper struct {
	logger      loggerAdapter
	rekognition Rekognition
	instance    *rekognition.Client
}

func RekognitionClient(region string) Rekognition {
	cfg := getConfig(region)
	if awsUrl := os.Getenv("AWS_URL"); awsUrl != "" {
		return rekognition.NewFromConfig(cfg)
	}

	return rekognition.NewFromConfig(cfg)
}

func RekognitionHelper(rekognition Rekognition, logger loggerAdapter) RekognitionHelperAdapter {
	return &rekognitionHelper{
		logger:      logger,
		rekognition: rekognition,
	}
}

func (s *rekognitionHelper) CreateCollection(ctx context.Context, clientId string) error {
	s.logger.Debug(ctx, "Started CreateCollection", map[string]any{
		"clientId": clientId,
	})

	result, err := s.rekognition.CreateCollection(ctx, &rekognition.CreateCollectionInput{
		CollectionId: &clientId,
		Tags: map[string]string{
			"client_id": clientId,
		},
	})
	if err != nil {
		s.logger.Debug(ctx, "CreateCollection failed", err)
		return fmt.Errorf("failed to create collection %q: %w", clientId, err)
	}

	s.logger.Debug(ctx, "Finished CreateCollection", result.StatusCode, result.CollectionArn)
	return nil
}

func (s *rekognitionHelper) CreateUser(ctx context.Context, collectionId, userId string) error {
	s.logger.Debug(ctx, "Started CreateUser", map[string]any{
		"collectionId": collectionId,
		"userId":       userId,
	})

	_, err := s.rekognition.CreateUser(ctx, &rekognition.CreateUserInput{
		CollectionId: &collectionId,
		UserId:       &userId,
	})
	if err != nil {
		s.logger.Debug(ctx, "CreateUser failed", err)
		return fmt.Errorf("failed to create user %q in collection %q: %w", userId, collectionId, err)
	}

	s.logger.Debug(ctx, "Finished CreateUser", userId)
	return nil
}

func (s *rekognitionHelper) AssociateFace(ctx context.Context, collectionId, userId string, minimumSimilarityPercentage float32, imageBytes []byte) (*string, error) {
	s.logger.Debug(ctx, "Started AssociateFaces (via image bytes)", map[string]any{
		"collectionId":                collectionId,
		"userId":                      userId,
		"minimumSimilarityPercentage": minimumSimilarityPercentage,
	})

	indexResult, err := s.rekognition.IndexFaces(ctx,
		&rekognition.IndexFacesInput{
			CollectionId: &collectionId,
			Image: &types.Image{
				Bytes: imageBytes,
			},
			ExternalImageId:     &userId,
			QualityFilter:       types.QualityFilterHigh,
			MaxFaces:            aws.Int32(1),
			DetectionAttributes: []types.Attribute{types.AttributeDefault},
		},
	)
	if err != nil {
		s.logger.Debug(ctx, "IndexFaces (for association) failed", err)
		return nil, fmt.Errorf("failed to index face for association: %w", err)
	}
	if len(indexResult.FaceRecords) == 0 || indexResult.FaceRecords[0].Face == nil || indexResult.FaceRecords[0].Face.FaceId == nil {
		return nil, fmt.Errorf("no face detected for association")
	}

	var faceId string
	if indexResult.FaceRecords[0].Face != nil && indexResult.FaceRecords[0].Face.FaceId != nil {
		faceId = *indexResult.FaceRecords[0].Face.FaceId
	}

	result, err := s.rekognition.AssociateFaces(ctx,
		&rekognition.AssociateFacesInput{
			CollectionId:       &collectionId,
			UserId:             &userId,
			ClientRequestToken: aws.String(s.logger.GetTransactionID(ctx)),
			UserMatchThreshold: &minimumSimilarityPercentage,
			FaceIds:            []string{faceId},
		},
	)
	if err != nil {
		s.logger.Debug(ctx, "AssociateFaces failed", err)
		return nil, fmt.Errorf("failed to associate faces: %w", err)
	}

	s.logger.Debug(ctx, "Finished AssociateFaces", result.AssociatedFaces)
	return &faceId, nil
}

func (s *rekognitionHelper) DisassociateFace(ctx context.Context, collectionId, userId, faceId string) error {
	s.logger.Debug(ctx, "Started DeleteFaces", map[string]any{
		"collectionId": collectionId,
		"faceId":       faceId,
	})

	_, err := s.rekognition.DisassociateFaces(ctx, &rekognition.DisassociateFacesInput{
		CollectionId:       &collectionId,
		UserId:             &userId,
		FaceIds:            []string{faceId},
		ClientRequestToken: aws.String(s.logger.GetTransactionID(ctx)),
	})
	if err != nil {
		s.logger.Debug(ctx, "DisassociateFaces failed", err)
		return fmt.Errorf("failed to disassociate face: %w", err)
	}

	result, err := s.rekognition.DeleteFaces(ctx, &rekognition.DeleteFacesInput{
		CollectionId: &collectionId,
		FaceIds:      []string{faceId},
	})
	if err != nil {
		s.logger.Debug(ctx, "DeleteFaces failed", err)
		return fmt.Errorf("failed to delete faces: %w", err)
	}

	s.logger.Debug(ctx, "Finished DeleteFaces", result.DeletedFaces)
	return nil
}

func (s *rekognitionHelper) SearchUsersByImage(ctx context.Context, collectionId string, minimumSimilarityPercentage float32, imageBytes []byte) ([]UserMatch, error) {
	s.logger.Debug(ctx, "Started SearchUsersByImage (via image bytes)", map[string]any{
		"collectionId": collectionId,
	})

	result, err := s.rekognition.SearchUsersByImage(ctx, &rekognition.SearchUsersByImageInput{
		CollectionId: &collectionId,
		Image: &types.Image{
			Bytes: imageBytes,
		},
		UserMatchThreshold: &minimumSimilarityPercentage,
		QualityFilter:      types.QualityFilterHigh,
	})
	if err != nil {
		s.logger.Debug(ctx, "SearchUsersByImage failed", err)
		return nil, fmt.Errorf("failed to search users by image: %w", err)
	}

	userMatches := make([]UserMatch, 0)
	for _, userMatch := range result.UserMatches {
		userMatches = append(userMatches, UserMatch{
			Id:         userMatch.User.UserId,
			Similarity: userMatch.Similarity,
		})
	}

	s.logger.Debug(ctx, "Finished SearchUsersByImage", userMatches)
	return userMatches, nil
}
