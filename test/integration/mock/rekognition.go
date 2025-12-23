package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/google/uuid"
)

type RekognitionClient struct {
	createdUsers       map[string]bool
	matches            []types.UserMatch
	disassociatedFaces []DisassociatedFaceCall
	deletedFaces       []DeletedFaceCall
	indexedFaces       []IndexedFaceCall
	associatedFaces    []AssociatedFaceCall
	approve            bool
}

type DisassociatedFaceCall struct {
	CollectionId string
	UserId       string
	FaceIds      []string
}

type DeletedFaceCall struct {
	CollectionId string
	FaceIds      []string
}

type IndexedFaceCall struct {
	CollectionId string
	FaceIds      []string
}

type AssociatedFaceCall struct {
	CollectionId string
	UserId       string
	FaceIds      []string
}

var rekognitionInstance *RekognitionClient
var rekognitionInit sync.Once

func NewRekognitionClient() *RekognitionClient {
	rekognitionInit.Do(
		func() {
			rekognitionInstance = &RekognitionClient{
				approve:            true,
				matches:            []types.UserMatch{},
				createdUsers:       make(map[string]bool),
				disassociatedFaces: []DisassociatedFaceCall{},
				deletedFaces:       []DeletedFaceCall{},
				indexedFaces:       []IndexedFaceCall{},
				associatedFaces:    []AssociatedFaceCall{},
			}
		},
	)
	return rekognitionInstance
}

func (m *RekognitionClient) SetResponse(approve bool, similarityPercentage float64, userIdResponse string) {
	m.approve = approve
	if m.matches == nil {
		m.matches = []types.UserMatch{}
	}
	m.matches = append(m.matches, types.UserMatch{
		User: &types.MatchedUser{
			UserId: aws.String(userIdResponse),
		},
		Similarity: aws.Float32(float32(similarityPercentage)),
	})
}

func (m *RekognitionClient) CreateCollection(ctx context.Context, params *rekognition.CreateCollectionInput, optFns ...func(*rekognition.Options)) (*rekognition.CreateCollectionOutput, error) {
	if !m.approve {
		return nil, errors.New("CreateCollection denied by mock")
	}
	return &rekognition.CreateCollectionOutput{
		CollectionArn: params.CollectionId,
		StatusCode:    aws.Int32(200),
	}, nil
}

func (m *RekognitionClient) CreateUser(ctx context.Context, params *rekognition.CreateUserInput, optFns ...func(*rekognition.Options)) (*rekognition.CreateUserOutput, error) {
	if !m.approve {
		return nil, errors.New("CreateUser denied by mock")
	}
	if m.createdUsers == nil {
		m.createdUsers = make(map[string]bool)
	}
	m.createdUsers[*params.UserId] = true
	return &rekognition.CreateUserOutput{}, nil
}

func (m *RekognitionClient) IndexFaces(ctx context.Context, params *rekognition.IndexFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.IndexFacesOutput, error) {
	if !m.approve {
		return nil, errors.New("IndexFaces denied by mock")
	}

	faceId := uuid.NewString()

	m.indexedFaces = append(m.indexedFaces, IndexedFaceCall{
		CollectionId: *params.CollectionId,
		FaceIds:      []string{faceId},
	})

	return &rekognition.IndexFacesOutput{
		FaceRecords: []types.FaceRecord{
			{
				Face: &types.Face{
					FaceId: aws.String(faceId),
				},
			},
		},
	}, nil
}

func (m *RekognitionClient) AssociateFaces(ctx context.Context, params *rekognition.AssociateFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.AssociateFacesOutput, error) {
	if !m.approve {
		return nil, errors.New("AssociateFaces denied by mock")
	}
	if m.matches == nil {
		m.matches = []types.UserMatch{}
	}
	m.matches = append(m.matches, types.UserMatch{
		User: &types.MatchedUser{
			UserId: params.UserId,
		},
		Similarity: aws.Float32(99.9),
	})

	m.associatedFaces = append(m.associatedFaces, AssociatedFaceCall{
		CollectionId: *params.CollectionId,
		UserId:       *params.UserId,
		FaceIds:      params.FaceIds,
	})

	return &rekognition.AssociateFacesOutput{
		AssociatedFaces: []types.AssociatedFace{
			{
				FaceId: aws.String(uuid.NewString()),
			},
		},
	}, nil
}

func (m *RekognitionClient) DisassociateFaces(ctx context.Context, params *rekognition.DisassociateFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.DisassociateFacesOutput, error) {
	if !m.approve {
		return nil, errors.New("AssociateFaces denied by mock")
	}

	m.disassociatedFaces = append(m.disassociatedFaces, DisassociatedFaceCall{
		CollectionId: *params.CollectionId,
		UserId:       *params.UserId,
		FaceIds:      params.FaceIds,
	})

	return &rekognition.DisassociateFacesOutput{
		DisassociatedFaces: []types.DisassociatedFace{
			{
				FaceId: aws.String(params.FaceIds[0]),
			},
		},
	}, nil
}

func (m *RekognitionClient) DeleteFaces(ctx context.Context, params *rekognition.DeleteFacesInput, optFns ...func(*rekognition.Options)) (*rekognition.DeleteFacesOutput, error) {
	if !m.approve {
		return nil, errors.New("DeleteFaces denied by mock")
	}

	m.deletedFaces = append(m.deletedFaces, DeletedFaceCall{
		CollectionId: *params.CollectionId,
		FaceIds:      params.FaceIds,
	})

	return &rekognition.DeleteFacesOutput{
		DeletedFaces: params.FaceIds,
	}, nil
}

func (m *RekognitionClient) SearchUsersByImage(ctx context.Context, params *rekognition.SearchUsersByImageInput, optFns ...func(*rekognition.Options)) (*rekognition.SearchUsersByImageOutput, error) {
	if !m.approve {
		return nil, errors.New("SearchUsersByImage denied by mock")
	}
	return &rekognition.SearchUsersByImageOutput{
		UserMatches: m.matches,
	}, nil
}

func (m *RekognitionClient) SearchUsers(ctx context.Context, params *rekognition.SearchUsersInput, optFns ...func(*rekognition.Options)) (*rekognition.SearchUsersOutput, error) {
	if !m.approve {
		return nil, errors.New("SearchUsersInput denied by mock")
	}
	if m.createdUsers == nil {
		m.createdUsers = make(map[string]bool)
	}
	if !m.createdUsers[*params.UserId] {
		return nil, &types.InvalidParameterException{
			Message: aws.String("User not found in collection"),
		}
	}
	return &rekognition.SearchUsersOutput{
		UserMatches: m.matches,
	}, nil
}

func (m *RekognitionClient) Reset() {
	m.approve = true
	m.matches = []types.UserMatch{}
	m.createdUsers = make(map[string]bool)
	m.disassociatedFaces = []DisassociatedFaceCall{}
	m.deletedFaces = []DeletedFaceCall{}
	m.indexedFaces = []IndexedFaceCall{}
	m.associatedFaces = []AssociatedFaceCall{}
}

func (m *RekognitionClient) GetDisassociatedFaces() []DisassociatedFaceCall {
	return m.disassociatedFaces
}

func (m *RekognitionClient) GetDeletedFaces() []DeletedFaceCall {
	return m.deletedFaces
}

func (m *RekognitionClient) GetIndexedFaces() []IndexedFaceCall {
	return m.indexedFaces
}

func (m *RekognitionClient) GetAssociatedFaces() []AssociatedFaceCall {
	return m.associatedFaces
}

func (m *RekognitionClient) CountDisassociatedFacesForUser(collectionId, userId string) int {
	count := 0
	for _, call := range m.disassociatedFaces {
		if call.CollectionId == collectionId && call.UserId == userId {
			count += len(call.FaceIds)
		}
	}
	return count
}

func (m *RekognitionClient) CountDeletedFacesForCollection(collectionId string) int {
	count := 0
	for _, call := range m.deletedFaces {
		if call.CollectionId == collectionId {
			count += len(call.FaceIds)
		}
	}
	return count
}

func (m *RekognitionClient) CountIndexedFacesForCollection(collectionId string) int {
	count := 0
	for _, call := range m.indexedFaces {
		if call.CollectionId == collectionId {
			count += len(call.FaceIds)
		}
	}
	return count
}

func (m *RekognitionClient) CountAssociatedFacesForUser(collectionId, userId string) int {
	count := 0
	for _, call := range m.associatedFaces {
		if call.CollectionId == collectionId && call.UserId == userId {
			count += len(call.FaceIds)
		}
	}
	return count
}
