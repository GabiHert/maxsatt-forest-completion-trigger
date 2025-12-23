package webservice

import (
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
)

type analysisWebService struct {
	client      http.Client
	authService adapter.AuthWebService
}

func NewAnalysisWebService(httpClient http.Client, authService adapter.AuthWebService) adapter.AnalysisWebService {
	return &analysisWebService{
		client:      httpClient,
		authService: authService,
	}
}

func (a *analysisWebService) FetchByProcessingId(ctx context.Context, processingID string) (*entity.Analysis, error) {
	token, err := a.authService.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/v1/analysis/%s", properties.Properties().Services.ForestFieldAPIURL, processingID)

	req, err := http2.NewRequest(http2.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errs.InternalServerError(err, "ANL-00100", "Failed to create request for analysis data")
	}

	req.Header.Set("Authorization", *token)

	if properties.Properties().Services.AuthGtwId != "" {
		req.Header.Set("x-apigw-api-id", properties.Properties().Services.AuthGtwId)
	}

	response, err := a.client.Do(ctx, req)
	if err != nil {
		return nil, errs.InternalServerError(err, "ANL-00500", "Failed to fetch analysis data")
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return nil, domainError.NewProcessingNotFoundError(processingID)
	} else if response.StatusCode >= 400 && response.StatusCode < 500 {
		return nil, errs.ClientError("Error while trying to fetch analysis data", "ANL-00400", response.StatusCode)
	} else if response.StatusCode >= 500 {
		return nil, errs.ServerError("Error while trying to fetch analysis data", "ANL-01500", response.StatusCode)
	}

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errs.InternalServerError(err, "ANL-02500", "Failed to read response body")
	}

	var analysisDto dto.Analysis
	err = json.Unmarshal(responseBytes, &analysisDto)
	if err != nil {
		return nil, errs.InternalServerError(err, "ANL-03500", "Failed to unmarshal analysis response")
	}

	analysis := analysisDto.ToEntity()
	return &analysis, nil
}
