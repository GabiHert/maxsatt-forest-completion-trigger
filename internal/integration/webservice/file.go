package webservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	http2 "net/http"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	domainError "github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/error"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

type fileWebService struct {
	client      http.Client
	authService adapter.AuthWebService
}

func NewFileWebService(httpClient http.Client, authService adapter.AuthWebService) adapter.FileWebService {
	return &fileWebService{
		client:      httpClient,
		authService: authService,
	}
}

func (f *fileWebService) CreateFile(ctx context.Context, file *entity.FileRecord) (*entity.FileRecord, error) {
	token, err := f.authService.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	apiURL := properties.Properties().Services.ForestFieldAPIURL + "/v1/files"

	requestBody := dto.NewFileRecord(file)
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		logger.Error(ctx, err, "Failed to marshal file record request", "imageID", file.ImageID, "errorCode", "FIL-00100")
		return nil, errs.InternalServerError(err, "FIL-00100", "Failed to marshal file record request")
	}

	req, err := http2.NewRequest(http2.MethodPost, apiURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-00200", "Failed to create request for file record creation")
	}

	req.Header.Set("Authorization", *token)
	req.Header.Set("Content-Type", "application/json")

	if properties.Properties().Services.AuthGtwId != "" {
		req.Header.Set("x-apigw-api-id", properties.Properties().Services.AuthGtwId)
	}

	response, err := f.client.Do(ctx, req)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-00500", "Failed to create file record")
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil, errs.ClientError("Error while trying to create file record", "FIL-00400", response.StatusCode)
	} else if response.StatusCode >= 500 {
		return nil, errs.ServerError("Error while trying to create file record", "FIL-01500", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-02500", "Failed to read response body")
	}

	var fileRecordResponse dto.FileRecord
	err = json.Unmarshal(responseBytes, &fileRecordResponse)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-03500", "Failed to unmarshal file record response")
	}

	result := fileRecordResponse.ToEntity()
	return result, nil
}

// GetByProcessingID fetches the DELTA file for a processing ID to retrieve the image_id.
// Returns the file record or an error if not found or API fails.
func (f *fileWebService) GetByProcessingID(ctx context.Context, processingID string) (*entity.FileRecord, error) {
	token, err := f.authService.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/v1/files?processing_id=%s&type=DELTA&limit=1", properties.Properties().Services.ForestFieldAPIURL, processingID)

	req, err := http2.NewRequest(http2.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-04100", "Failed to create request for DELTA file lookup")
	}

	req.Header.Set("Authorization", *token)
	req.Header.Set("Content-Type", "application/json")

	if properties.Properties().Services.AuthGtwId != "" {
		req.Header.Set("x-apigw-api-id", properties.Properties().Services.AuthGtwId)
	}

	response, err := f.client.Do(ctx, req)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-04500", "Failed to fetch DELTA file")
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return nil, domainError.NewDeltaFileNotFoundError(processingID)
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil, errs.ClientError("Error while trying to fetch DELTA file", "FIL-04400", response.StatusCode)
	} else if response.StatusCode >= 500 {
		return nil, errs.ServerError("Error while trying to fetch DELTA file", "FIL-05500", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-06500", "Failed to read response body")
	}

	var listResponse dto.ListFilesResponse
	err = json.Unmarshal(responseBytes, &listResponse)
	if err != nil {
		return nil, errs.InternalServerError(err, "FIL-07500", "Failed to unmarshal files list response")
	}

	files := listResponse.GetFiles()
	if len(files) == 0 {
		return nil, domainError.NewDeltaFileNotFoundError(processingID)
	}

	deltaFile := files[0]
	if deltaFile.ImageID == "" {
		return nil, domainError.NewMissingImageIDError(processingID)
	}

	result := deltaFile.ToEntity()
	return result, nil
}
