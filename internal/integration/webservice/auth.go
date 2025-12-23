package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	http2 "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
	redishelper "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/redis"
)

type authWebService struct {
	client http.Client
	redis  redishelper.RedisHelper
}

func NewAuthWebService(httpClient http.Client, redis redishelper.RedisHelper) adapter.AuthWebService {
	return &authWebService{
		client: httpClient,
		redis:  redis,
	}
}

func (a *authWebService) GetToken(ctx context.Context) (*string, error) {
	logger.Debug(ctx, "Started")

	var token string
	if err := a.redis.Get(ctx, "api-token", &token); err != nil || token == "" {
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		data.Set("client_id", properties.Properties().Services.AuthClientId)
		data.Set("client_secret", properties.Properties().Services.AuthClientSecret)

		payload := strings.NewReader(data.Encode())

		req, err := http2.NewRequest(http2.MethodPost, properties.Properties().Services.AuthUrl+"/oauth2/token", payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		if properties.Properties().Services.AuthGtwId != "" {
			req.Header.Set("x-apigw-api-id", properties.Properties().Services.AuthGtwId)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response, err := a.client.Do(ctx, req)
		if err != nil {
			return nil, errs.InternalServerError(err, "TKN-00500", "Failed to find JWKs")
		}
		defer response.Body.Close()

		if response.StatusCode >= 400 && response.StatusCode < 500 {
			return nil, errs.ClientError("Error while trying to get token", "TKN-00400", response.StatusCode)
		} else if response.StatusCode >= 500 {
			return nil, errs.ServerError("Error while trying to get token", "TKN-01500", response.StatusCode)
		}

		accessTokenBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, errs.ServerError("Error while trying to parse token: "+err.Error(), "TKN-02500", http2.StatusInternalServerError)
		}

		var accessToken dto.AccessToken
		err = json.Unmarshal(accessTokenBytes, &accessToken)
		if err != nil {
			return nil, errs.ServerError("Error while trying to parse token: "+err.Error(), "TKN-03500", http2.StatusInternalServerError)
		}

		token = accessToken.AccessToken
		_ = a.redis.Put(ctx, "api-token", token, time.Duration(float64(accessToken.ExpiresIn)*0.8)*time.Second)
	}

	logger.Debug(ctx, "Finished")
	return &token, nil
}
