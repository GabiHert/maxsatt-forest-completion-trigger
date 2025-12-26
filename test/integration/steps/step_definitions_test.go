package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/dependency"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	aws2 "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/aws"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/timeutils"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/test/integration/mock"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go"
	"github.com/redis/go-redis/v9"
)

var tags string

func init() {
	flag.StringVar(&tags, "scenarios", "", "tags to run")
}

func TestFeatures(t *testing.T) {
	flag.Parse()

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			InitializeScenario(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			Tags:     tags,
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type testUtils struct {
	uri            string
	headers        map[string]string
	client         *http.Client
	response       *response
	time           *mock.Time
	api            *mock.ApiMock
	oauthApi       *mock.ApiMock
	redis          *redis.Client
	sns            *mock.SnsClient
	sqs            *mock.SqsClient
	secretsManager *mock.SecretsManagerClient
	dynamoDb       *mock.DynamoDbClient
	s3             *mock.S3Client
	rekognition    *mock.RekognitionClient
	ses            *mock.SesClient
	receiveCount   int
	createdFiles   []string
}

type response struct {
	status  int
	body    any
	err     error
	headers map[string]string
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	test := &testUtils{
		uri:            "",
		client:         &http.Client{},
		time:           mock.NewTime(),
		api:            mock.NewApiServer(),
		oauthApi:       mock.NewApiServer(),
		redis:          mock.NewRedis(),
		secretsManager: mock.NewSecretsManagerMock(),
		dynamoDb:       mock.NewDynamoDbClient(),
		sns:            mock.NewSnsClient(),
		sqs:            mock.NewSqsClient(),
		s3:             mock.NewS3Client(),
		rekognition:    mock.NewRekognitionClient(),
		ses:            mock.NewSesClient(),
	}

	ctx.Before(
		func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			test.before()
			// Set up mocks AFTER ResetInjector() is called in before()
			timeutils.GetTimeConfig().SetCurrentTimeProvider(test.time.Now)
			_ = setInjectorField("SecretsManager", test.secretsManager)
			_ = setInjectorField("Redis", test.redis)
			_ = setInjectorField("Sns", test.sns)
			_ = setInjectorField("Sqs", test.sqs)
			_ = setInjectorField("DynamoDb", test.dynamoDb)
			_ = setInjectorField("Time", test.time)
			_ = setInjectorField("S3", test.s3)
			_ = setInjectorField("Rekognition", test.rekognition)
			_ = setInjectorField("Ses", test.ses)
			test.startApp()
			return ctx, nil
		},
	)

	ctx.After(
		func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
			return nil, nil
		},
	)

	ctx.Given(`^the aws secret named "([^"]*)" exists with the following values$`, test.theSecretNamedExistsWithTheFollowingValues)
	ctx.Given(`^the "([^"]*)" env var is set to "([^"]*)"$`, test.theEnvVarIsSetTo)
	ctx.Given(`^the "([^"]*)" env var is set to the mocked api url$`, test.theEnvVarIsSetToTheMockedApiUrl)
	ctx.Given(`^the "([^"]*)" env var is set to the mocked oauth url$`, test.theEnvVarIsSetToTheMockedOAuthUrl)
	ctx.Given(`^the oauth server returns a valid token$`, test.theOAuthServerReturnsValidToken)
	ctx.Given(`^the oauth server returns an invalid token response$`, test.theOAuthServerReturnsInvalidToken)
	ctx.Given(`^the header is empty$`, test.theHeaderIsEmpty)
	ctx.Given(`^the header contains the key "([^"]*)" with "([^"]*)"$`, test.theHeaderContainsTheKeyWith)
	ctx.Given(`^the "([^"]*)" exists$`, test.theExists)
	ctx.Given(`^all the "([^"]*)" bellow exists$`, test.theExists)
	ctx.Given(`^the "([^"]*)" local file exists in s3 as "([^"]*)"$`, test.theS3FileExists)
	ctx.Given(`^the tables are empty$`, test.theTablesAreEmpty)
	ctx.Given(`^the "([^"]*)" table is empty$`, test.theTableIsEmpty)
	ctx.Given(`^the dynamodb table "([^"]*)" with key "([^"]*)" exists$`, test.theDynamoDbTableWithKeyExists)
	ctx.Given(`^the (-?\d+) "([^"]*)" request to "([^"]*)" returns status (\d+) with the following response$`, test.theRequestToReturnsStatusWithTheFollowingResponse)
	ctx.Given(`^the "([^"]*)" "([^"]*)" should not return any results$`, test.theShouldNotReturnAnyResults)
	ctx.Given(`^aws rekognition returns the following response$`, test.awsRekognitionRetunrsTheFollowingResponse)
	ctx.Step(`^I call "([^"]*)" "([^"]*)" with the following csv payload$`, test.iCallWithTheFollowingCsvPayload)
	ctx.Step(`^I call "([^"]*)" "([^"]*)" with the following payload$`, test.iCallWithTheFollowingPayload)
	ctx.Step(`^I call "([^"]*)" "([^"]*)" with multipart form data$`, test.iCallWithMultipartFormData)
	ctx.Step(`^the status returned should be (\d+)$`, test.theStatusReturnedShouldBe)
	ctx.Step(
		`^the response should contain the field "([^"]*)" equal to "([^"]*)"$`,
		test.theResponseShouldContainTheFieldEqualTo,
	)
	ctx.Step(
		`^the response should contain the field "([^"]*)" that contains "([^"]*)"$`,
		test.theResponseShouldContainTheFieldThatContains,
	)
	ctx.Step(
		`^the response should contain the field "([^"]*)" array length equal to (\d+)$`,
		test.theResponseShouldContainTheFieldArrayLengthEqualTo,
	)
	ctx.Step(
		`^the response should contain the text$`,
		test.theResponseShouldContainTheText,
	)
	ctx.Step(
		`^the response should not be nil$`,
		test.theResponseShouldNotBeNil,
	)
	ctx.Step(`^the response content type should be "([^"]*)"$`, test.theResponseContentTypeShouldBe)
	ctx.When(`^the following event is received via dynamodb stream`, test.theFollowingEventIsReceivedViaDynamoDbStream)
	ctx.When(`^the following event is received via sns$`, test.theFollowingEventIsReceivedViaSns)
	ctx.When(`^the following event is received via sqs$`, test.theFollowingEventIsReceivedViaSqs)
	ctx.Then(`^the lambda should finish without errors$`, test.theLambdaShouldFinishWithoutErrors)
	ctx.Then(`^the lambda should finish with "([^"]*)" error$`, test.theLambdaShouldFinishWithError)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have (\d+) messages published$`,
		test.theSqsQueueShouldHaveMessagesPublished,
	)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have the (\d+) message published with "([^"]*)" field equal to "([^"]*)"$`,
		test.theSqsQueueShouldHaveTheMessagePublishedWithFieldEqualTo,
	)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have the (\d+) message published with "([^"]*)" message attribute equal to "([^"]*)"$`,
		test.theSqsQueueShouldHaveTheMessagePublishedWithMessageAttributeEqualTo,
	)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have the (\d+) message published with message group id equal to "([^"]*)"$`,
		test.theSqsQueueShouldHaveTheMessagePublishedWithMessageGroupIdEqualTo,
	)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have the (\d+) message published with message deduplication id equal to "([^"]*)"$`,
		test.theSqsQueueShouldHaveTheMessagePublishedWithMessageDeduplicationIdEqualTo,
	)
	ctx.Then(
		`^the sns topic "([^"]*)" should have (\d+) messages published$`,
		test.theSnsTopicShouldHaveMessagesPublished,
	)
	ctx.Then(
		`^the sns topic "([^"]*)" should have a message published with "([^"]*)" field equal to "([^"]*)"$`,
		test.theSnsTopicShouldHaveAMessagePublishedWithFieldEqualTo,
	)
	ctx.Then(
		`^the sns topic "([^"]*)" should have a message published with "([^"]*)" message attribute equal to "([^"]*)"$`,
		test.theSnsTopicShouldHaveAMessagePublishedWithMessageAttributeEqualTo,
	)
	ctx.Then(`^the (\d+) "([^"]*)" request for "([^"]*)" headers should have the variable "([^"]*)" with value "([^"]*)"$`, test.theForRequestHeadersShouldHaveTheVariableWithValue)
	ctx.Then(`^the (\d+) "([^"]*)" request for "([^"]*)" queries should have the variable "([^"]*)" with value "([^"]*)"$`, test.theForRequestQueriesShouldHaveTheVariableWithValue)
	ctx.Then(`^the (\d+) "([^"]*)" request for "([^"]*)" json body should have the variable "([^"]*)" with value "([^"]*)"$`, test.theForRequestJsonBodyShouldHaveTheVariableWithValue)
	ctx.Then(`^rekognition should have disassociated (\d+) faces for collection "([^"]*)" and user "([^"]*)"$`, test.rekognitionShouldHaveDisassociatedFacesForCollectionAndUser)
	ctx.Then(`^rekognition should have deleted (\d+) faces for collection "([^"]*)"$`, test.rekognitionShouldHaveDeletedFacesForCollection)
	ctx.Then(`^rekognition should have indexed (\d+) faces for collection "([^"]*)"$`, test.rekognitionShouldHaveIndexedFacesForCollection)
	ctx.Then(`^rekognition should have associated (\d+) faces for collection "([^"]*)" and user "([^"]*)"$`, test.rekognitionShouldHaveAssociatedFacesForCollectionAndUser)
	ctx.Given(`^the file "([^"]*)" is duplicated to "([^"]*)"$`, test.theFileIsDuplicatedTo)
	ctx.Given(`^the file "([^"]*)" is deleted$`, test.theFileIsDeleted)
	ctx.Then(`^the file "([^"]*)" should exist$`, test.theFileShouldExist)
	ctx.Then(`^the file "([^"]*)" should be identical to "([^"]*)"$`, test.theFileShouldBeIdenticalTo)
	ctx.Step(`^the parquet file "([^"]*)" should have the field "([^"]*)" in all rows$`, test.theParquetFileShouldHaveFieldInAllRows)
	ctx.Step(`^the parquet file "([^"]*)" should have the field "([^"]*)" in row (\d+) with value "([^"]*)"$`, test.theParquetFileShouldHaveFieldInRowWithValue)
	ctx.Step(`^the parquet file "([^"]*)" should have the fields "([^"]*)" in all rows$`, test.theParquetFileShouldHaveFieldsInAllRows)
	ctx.Step(`^the parquet file "([^"]*)" should have (\d+) rows$`, test.theParquetFileShouldHaveRows)
	ctx.Then(`^the parquet file "([^"]*)" should have all fields from "([^"]*)"$`, test.theParquetFileShouldHaveAllFieldsFrom)

}

func setInjectorField(fieldName string, value any) error {
	injector := dependency.Injector()
	v := reflect.ValueOf(injector).Elem()

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s doesn't exist", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	fieldValue := reflect.ValueOf(value)
	field.Set(fieldValue)
	return nil
}

var serverInit sync.Once

func (t *testUtils) startApp() {
	// For lambda functions, we don't start an HTTP server
	// The handler is invoked directly in the test steps
}

func (t *testUtils) before() {
	t.headers = make(map[string]string)
	t.createdFiles = []string{}

	err := mock.ClearRedis(t.redis)
	if err != nil {
		panic(err)
	}

	t.dynamoDb.Reset()
	t.receiveCount = 0
	t.secretsManager.Reset()
	t.rekognition.Reset()
	t.ses.Reset()
	t.s3.Reset()
	t.sns.Reset()
	t.sqs.Reset()

	// Reset properties and injector for each scenario
	properties.ResetProperties()
	dependency.ResetInjector()

	_ = os.Setenv("API_SECRET", "api-secret")
	_ = os.Setenv("AWS_S3_EMPLOYEES_BUCKET", "test-bucket")
	t.secretsManager.SetSecret("api-secret", map[string]any{
		"AUTH_CLIENT_ID":     "7073da1d-c57e-4fdc-a2e2-6208fee06ff1",
		"AUTH_CLIENT_SECRET": "b3d2c4f0-8c1b-4d2a-9e3f-5c6b7e8f9a0b",
	})

	t.api.Start()
	t.oauthApi.Start()
}

func (t *testUtils) execute(method, url string, request []byte, headers, queryParams map[string]string) (
	*http.Response,
	error,
) {
	var req *http.Request

	if request != nil {
		req, _ = http.NewRequest(method, t.uri+url, bytes.NewReader(request))
	} else {
		req, _ = http.NewRequest(method, t.uri+url, nil)
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return t.client.Do(req)
}

func (t *testUtils) theSecretNamedExistsWithTheFollowingValues(secret string, stringJson *godog.DocString) error {
	var secretValue map[string]any
	err := json.Unmarshal([]byte(stringJson.Content), &secretValue)
	if err != nil {
		return err
	}

	t.secretsManager.SetSecret(secret, secretValue)

	return nil
}

func (t *testUtils) theEnvVarIsSetTo(key, value string) error {
	_ = os.Setenv(key, value)

	return nil
}

func (t *testUtils) theEnvVarIsSetToTheMockedApiUrl(key string) error {
	mockUrl := t.api.GetUrl()
	_ = os.Setenv(key, mockUrl)

	return nil
}

func (t *testUtils) theEnvVarIsSetToTheMockedOAuthUrl(key string) error {
	mockUrl := t.oauthApi.GetUrl() + "/oauth/token"
	_ = os.Setenv(key, mockUrl)

	return nil
}

func (t *testUtils) theOAuthServerReturnsValidToken() error {
	t.oauthApi.SetResponse(-1, "POST", "/oauth/token", 200, map[string]any{
		"access_token": "mock-valid-access-token",
		"token_type":   "Bearer",
		"expires_in":   3600,
	})
	return nil
}

func (t *testUtils) theOAuthServerReturnsInvalidToken() error {
	t.oauthApi.SetResponse(-1, "POST", "/oauth/token", 401, map[string]any{
		"error":             "invalid_client",
		"error_description": "Invalid client credentials",
	})
	return nil
}

func (t *testUtils) theHeaderIsEmpty() error {
	t.headers = make(map[string]string)

	return nil
}

func (t *testUtils) theHeaderContainsTheKeyWith(key, value string) error {
	t.headers[key] = value

	return nil
}

func (t *testUtils) theExists(table string, jsonEntity *godog.DocString) error {
	entityContent := jsonEntity.Content
	var items []map[string]any
	err := json.Unmarshal([]byte(entityContent), &items)
	if err != nil {
		return err
	}

	// For HTTP-based tests, we add items to dynamoDB
	for _, item := range items {
		t.dynamoDb.AddItem(item, table)
	}

	return nil
}

func (t *testUtils) theS3FileExists(filePath, fileUrl string) error {
	var bucket, key string

	// Check if fileUrl is a full S3 URL or just a path
	if strings.HasPrefix(fileUrl, "http://") || strings.HasPrefix(fileUrl, "https://") || strings.HasPrefix(fileUrl, "s3://") {
		// Parse the S3 URL to extract bucket and key
		parsedURL, err := url.Parse(fileUrl)
		if err != nil {
			return err
		}

		if strings.HasPrefix(fileUrl, "s3://") {
			// s3://bucket/key format
			bucket = parsedURL.Host
			key = strings.TrimPrefix(parsedURL.Path, "/")
		} else {
			// Extract bucket from host (e.g., "test-bucket.s3.amazonaws.com" -> "test-bucket")
			parts := strings.Split(parsedURL.Host, ".")
			if len(parts) < 3 {
				return errors.New("invalid S3 URL format")
			}
			bucket = parts[0]
			// Extract key from path (remove leading slash)
			key = strings.TrimPrefix(parsedURL.Path, "/")
		}
	} else {
		// Simple path format - use the S3_BUCKET_NAME env var
		bucket = os.Getenv("S3_BUCKET_NAME")
		if bucket == "" {
			bucket = "maxsatt-data-bucket" // Default bucket name
		}
		key = fileUrl
	}

	// Get file content - check if filePath already includes 'test/integration/resources'
	var absPath string
	var err error
	if strings.HasPrefix(filePath, "test/integration/resources/") {
		// Remove the prefix since we're already in test/integration/steps
		filePath = strings.TrimPrefix(filePath, "test/integration/resources/")
		absPath, err = filepath.Abs("./../resources/" + filePath)
	} else {
		absPath, err = filepath.Abs("./../resources/" + filePath)
	}
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	// Add to mock S3
	if t.s3.Objects == nil {
		t.s3.Objects = make(map[string]map[string]mock.S3Object)
	}
	if t.s3.Objects[bucket] == nil {
		t.s3.Objects[bucket] = make(map[string]mock.S3Object)
	}

	contentType := "application/pdf"
	t.s3.Objects[bucket][key] = mock.S3Object{
		Data:        fileContent,
		ContentType: &contentType,
	}
	return nil
}

func (t *testUtils) theTablesAreEmpty() error {
	t.dynamoDb.Reset()
	err := mock.ClearRedis(t.redis)
	if err != nil {
		return err
	}
	return nil
}

func (t *testUtils) theTableIsEmpty(table string) error {
	err := t.dynamoDb.ResetTable(table)
	if err != nil {
		return err
	}

	return nil
}

func (t *testUtils) theDynamoDbTableWithKeyExists(table, key string) error {
	t.dynamoDb.AddTable(table, key)

	return nil
}

func (t *testUtils) theRequestToReturnsStatusWithTheFollowingResponse(index int, method, path string, status int, stringJson *godog.DocString) error {
	var apiResponse map[string]any
	err := json.Unmarshal([]byte(stringJson.Content), &apiResponse)
	if err != nil {
		return err
	}

	t.api.SetResponse(index, method, path, status, apiResponse)

	return nil
}

func (t *testUtils) theShouldNotReturnAnyResults(method, path string) error {
	t.api.ClearResponses(method, path)

	return nil
}

func (t *testUtils) awsRekognitionRetunrsTheFollowingResponse(content *godog.DocString) error {
	var awsResp map[string]any
	err := json.Unmarshal([]byte(content.Content), &awsResp)
	if err != nil {
		return err
	}

	approve := awsResp["approve"].(bool)

	for _, matches := range awsResp["matches"].([]any) {
		matchMap, ok := matches.(map[string]any)
		if !ok {
			return errors.New("invalid match format in AWS response")
		}

		t.rekognition.SetResponse(approve, matchMap["similarity"].(float64), matchMap["user_id"].(string))
	}

	return nil
}

func (t *testUtils) iCallWithTheFollowingCsvPayload(method, path string, body *godog.DocString) error {
	var payload string
	if body != nil && body.Content != "" {
		payload = body.Content
	}

	t.headers["Content-Type"] = "text/csv"

	resp, err := t.execute(method, path, []byte(payload), t.headers, nil)
	if err != nil {
		return err
	}

	t.response = &response{
		status: resp.StatusCode,
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var responseBody map[string]any
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		t.response.body = string(bodyBytes)
	} else {
		t.response.body = responseBody
	}

	return nil
}

func (t *testUtils) iCallWithTheFollowingPayload(method, path string, body *godog.DocString) error {
	var payload []byte
	if body != nil && body.Content != "" {
		payload = []byte(body.Content)
	}

	resp, err := t.execute(method, path, payload, t.headers, nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	for key, value := range resp.Header {
		headers[key] = value[0]
	}
	t.response = &response{
		status:  resp.StatusCode,
		headers: headers,
	}

	var responseBody map[string]any
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		t.response.body = string(bodyBytes)
	} else {
		t.response.body = responseBody
	}

	return nil
}

func (t *testUtils) iCallWithMultipartFormData(method, path string, body *godog.DocString) error {
	var payload map[string]any
	if body != nil && body.Content != "" {
		err := json.Unmarshal([]byte(body.Content), &payload)
		if err != nil {
			return err
		}
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	for key, value := range payload {
		switch value.(type) {
		case string:
			// value is a file name located in ./../resources/file_name
			isFile := true
			file, err := os.Open("./../resources/" + value.(string))
			if err != nil {
				isFile = false
			}
			defer file.Close()

			if isFile {
				part, err := writer.CreateFormFile(key, value.(string))
				if err != nil {
					return err
				}

				_, err = io.Copy(part, file)
				if err != nil {
					return err
				}
			} else {
				err = writer.WriteField(key, value.(string))
				if err != nil {
					return err
				}
			}
		default:
			// value is a json
			jsonValue, err := json.Marshal(value)
			if err != nil {
				return err
			}
			part, err := writer.CreateFormField(key)
			if err != nil {
				return err
			}
			_, err = part.Write(jsonValue)
		}
	}

	err := writer.Close()
	if err != nil {
		return err
	}

	t.headers["Content-Type"] = writer.FormDataContentType()

	resp, err := t.execute(method, path, buffer.Bytes(), t.headers, nil)
	if err != nil {
		return err
	}

	headers := make(map[string]string)
	for key, value := range resp.Header {
		headers[key] = value[0]
	}
	t.response = &response{
		status:  resp.StatusCode,
		headers: headers,
	}

	contentType := resp.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		var responseBody map[string]any
		err = json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			return err
		}
		t.response.body = responseBody
	case "application/octet-stream", "application/pdf", "application/zip":
		fileName := resp.Header.Get("Content-Disposition")
		if fileName == "" {
			if contentType == "application/pdf" {
				fileName = "response_file.pdf"
			} else if contentType == "application/zip" {
				fileName = "response_file.zip"
			} else {
				fileName = "response_file"
			}
		} else {
			fileName = strings.TrimPrefix(fileName, "inline; filename=")
			fileName = strings.TrimPrefix(fileName, "attachment; filename=")
			fileName = strings.Trim(fileName, "\"")
		}
		fileTmpPath := fmt.Sprintf("./../resources/tmp/%s", fileName)

		dir := filepath.Dir(fileTmpPath)
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		file, err := os.Create(fileTmpPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", fileTmpPath, err)
		}
		defer file.Close()
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to write response body to file %s: %w", fileTmpPath, err)
		}
		t.response.body = map[string]string{"file_path": fileTmpPath}
	default:
		return fmt.Errorf("unexpected response Content-Type: %s", resp.Header.Get("Content-Type"))
	}

	return nil
}

func (t *testUtils) theStatusReturnedShouldBe(status int) error {
	if t.response.status != status {
		// If there's an error response body, include it in the error message for debugging
		if t.response.status >= 400 && t.response.body != nil {
			bodyBytes, _ := json.Marshal(t.response.body)
			return fmt.Errorf("%d should be equal to %d. Response: %s", t.response.status, status, string(bodyBytes))
		}
		return fmt.Errorf("%d should be equal to %d", t.response.status, status)
	}
	return nil
}

func (t *testUtils) theResponseShouldContainTheFieldEqualTo(dotSeparatedField, value string) error {
	field := getFieldValue(t.response.body, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	} else {
		return assertEqual(field, value)
	}
}

func (t *testUtils) theResponseShouldContainTheFieldThatContains(dotSeparatedField, value string) error {
	field := getFieldValue(t.response.body, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	} else {
		return assertContains(field, value)
	}
}

func (t *testUtils) theResponseShouldContainTheFieldArrayLengthEqualTo(dotSeparatedField string, value int) error {
	field := getFieldValue(t.response.body, dotSeparatedField)

	err := assertNotNull(field)
	if err != nil {
		return err
	}

	var length int
	switch field.(type) {
	case []any:
		length = len(field.([]any))
	default:
		return errors.New("field is not an array")
	}

	return assertEqual(length, value)
}

func (t *testUtils) theResponseShouldContainTheText(content *godog.DocString) error {
	field := t.response.body

	contentString := strings.Replace(content.Content, "\\r", "\r", -1)

	return assertEqual(field, contentString)
}

func (t *testUtils) theResponseShouldNotBeNil() error {
	return assertNotNull(t.response.body)
}

func (t *testUtils) theResponseContentTypeShouldBe(expectedContentType string) error {
	if t.response == nil {
		return fmt.Errorf("no response available")
	}

	actualContentType := ""
	for key, value := range t.response.headers {
		if strings.ToLower(key) == "content-type" {
			actualContentType = strings.ToLower(value)
			break
		}
	}

	return assertEqual(actualContentType, expectedContentType)
}

func (t *testUtils) theFollowingEventIsReceivedViaDynamoDbStream(stringJson *godog.DocString) error {
	var content map[string]any
	err := json.Unmarshal([]byte(stringJson.Content), &content)
	if err != nil {
		return err
	}

	oldImage, err := aws2.MarshalToStreamImage(content)
	if err != nil {
		panic(err)
	}

	dynamodbStreamEvent := events.DynamoDBStreamRecord{
		ApproximateCreationDateTime: events.SecondsEpochTime{},
		Keys:                        nil,
		NewImage:                    nil,
		OldImage:                    oldImage,
		SequenceNumber:              "",
		SizeBytes:                   0,
		StreamViewType:              "",
	}

	awsCtx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	})

	responseData, err := dependency.Injector().Wire(awsCtx).Handler.Handle(awsCtx, dynamodbStreamEvent)

	t.response = &response{
		status: 0,
		body:   responseData,
		err:    err,
	}
	return nil
}

func (t *testUtils) theFollowingEventIsReceivedViaSns(stringJson *godog.DocString) error {
	snsEvent := events.SNSEntity{
		MessageAttributes: map[string]any{
			"ApproximateReceiveCount": map[string]any{
				"Type":  "String",
				"Value": strconv.Itoa(t.receiveCount),
			},
			"correlationId": map[string]any{
				"Type":  "String",
				"Value": "d003c2bc-9d14-44ab-9245-3206ac723b8c",
			},
		},
		Message: stringJson.Content,
	}

	snsEventByte, err := json.Marshal(snsEvent)
	if err != nil {
		return err
	}
	stringSnsEvent := string(snsEventByte)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Attributes: map[string]string{
					"ApproximateReceiveCount": strconv.Itoa(t.receiveCount),
				},
				Body: stringSnsEvent,
				MessageAttributes: map[string]events.SQSMessageAttribute{
					"correlationId": {
						DataType:    "String",
						StringValue: aws.String("d003c2bc-9d14-44ab-9245-3206ac723b8c"),
					},
				},
			},
		},
	}

	awsCtx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	})

	responseData, err := dependency.Injector().Wire(awsCtx).Handler.Handle(awsCtx, sqsEvent)

	t.response = &response{
		status: 0,
		body:   responseData,
		err:    err,
	}
	return nil
}

func (t *testUtils) theFollowingEventIsReceivedViaSqs(stringJson *godog.DocString) error {
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: stringJson.Content,
			},
		},
	}

	awsCtx := lambdacontext.NewContext(context.Background(), &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	})

	responseData, err := dependency.Injector().Wire(awsCtx).Handler.Handle(awsCtx, sqsEvent)

	t.response = &response{
		status: 0,
		body:   responseData,
		err:    err,
	}
	return nil
}

func (t *testUtils) theLambdaShouldFinishWithoutErrors() error {
	if t.response.err == nil {
		return nil
	}
	return fmt.Errorf("expected no error but got: %v", t.response.err)
}

func (t *testUtils) theLambdaShouldFinishWithError(message string) error {
	if err := assertNotNull(t.response.err); err == nil {
		if strings.Contains(t.response.err.Error(), message) {
			return nil
		}
		return fmt.Errorf("expected error containing '%s' but got: '%s'", message, t.response.err.Error())
	} else {
		return fmt.Errorf("expected error but got none")
	}
}

func (t *testUtils) theSqsQueueShouldHaveMessagesPublished(queue string, quantity int) error {
	messagesSent := t.sqs.Messages[queue]
	err := assertNotNull(messagesSent)
	if err != nil {
		return err
	}

	err = assertEqual(len(messagesSent), quantity)
	if err != nil {
		return err
	}
	return nil
}

func (t *testUtils) theSqsQueueShouldHaveTheMessagePublishedWithFieldEqualTo(queue string, index int, dotSeparatedField, value string) error {
	messages := t.sqs.Messages[queue]

	err := assertNotNull(messages)
	if err != nil {
		return err
	}

	message := messages[index]
	err = assertNotNull(message)
	if err != nil {
		return err
	}

	var entity map[string]interface{}
	err = json.Unmarshal([]byte(*message), &entity)
	if err != nil {
		return err
	}

	field := getFieldValue(entity, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func (t *testUtils) theSqsQueueShouldHaveTheMessagePublishedWithMessageAttributeEqualTo(queue string, index int, dotSeparatedField, value string) error {
	messages := t.sqs.MessageAttributes[queue]

	err := assertNotNull(messages)
	if err != nil {
		return err
	}

	message := messages[index]
	err = assertNotNull(message)
	if err != nil {
		return err
	}

	var entity map[string]interface{}
	err = json.Unmarshal([]byte(*message), &entity)
	if err != nil {
		return err
	}

	field := getFieldValue(entity, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func (t *testUtils) theSqsQueueShouldHaveTheMessagePublishedWithMessageGroupIdEqualTo(queue string, index int, value string) error {
	messages := t.sqs.MessageGroupId[queue]

	err := assertNotNull(messages)
	if err != nil {
		return err
	}

	message := messages[index]
	err = assertNotNull(message)
	if err != nil {
		return err
	}

	if value == "nil" {
		return assertNull(*message)
	} else if value == "not nil" {
		return assertNotNull(*message)
	}

	return assertEqual(*message, value)
}

func (t *testUtils) theSqsQueueShouldHaveTheMessagePublishedWithMessageDeduplicationIdEqualTo(queue string, index int, value string) error {
	messages := t.sqs.MessageDeduplicationId[queue]

	err := assertNotNull(messages)
	if err != nil {
		return err
	}

	message := messages[index]
	err = assertNotNull(message)
	if err != nil {
		return err
	}

	if value == "nil" {
		return assertNull(*message)
	} else if value == "not nil" {
		return assertNotNull(*message)
	}

	return assertEqual(*message, value)
}

func (t *testUtils) theSnsTopicShouldHaveMessagesPublished(topic string, quantity int) error {
	err := assertEqual(t.sns.NumberOfMessagesSent[topic], quantity)
	if err != nil {
		return err
	}
	return nil
}

func (t *testUtils) theSnsTopicShouldHaveAMessagePublishedWithFieldEqualTo(topic, dotSeparatedField, value string) error {
	message := t.sns.Messages[topic]

	if message == nil {
		return errors.New("message is nil")
	}

	var entity map[string]interface{}
	err := json.Unmarshal([]byte(*message), &entity)
	if err != nil {
		return err
	}

	field := getFieldValue(entity, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func (t *testUtils) theSnsTopicShouldHaveAMessagePublishedWithMessageAttributeEqualTo(topic, dotSeparatedField, value string) error {
	message := t.sns.MessageAttributes[topic]

	if message == nil {
		return errors.New("message attribute is nil")
	}

	var entity map[string]interface{}
	err := json.Unmarshal([]byte(*message), &entity)
	if err != nil {
		return err
	}

	field := getFieldValue(entity, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func (t *testUtils) theForRequestHeadersShouldHaveTheVariableWithValue(index int, method, path string, header, value string) error {
	headers := t.api.GetRequestHeaders(method, path, index)
	if headers == nil {
		return errors.New("headers is nil")
	}

	var field any
	for headerKey, headerValue := range headers {
		if strings.ToLower(headerKey) == strings.ToLower(header) {
			field = &headerValue
		}
	}

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	} else {
		return assertEqual(*field.(*string), value)
	}
}

func (t *testUtils) theForRequestQueriesShouldHaveTheVariableWithValue(index int, method, path string, header, value string) error {
	queries := t.api.GetRequestQueries(method, path, index)
	if queries == nil {
		return errors.New("queries is nil")
	}

	var field *string
	for headerKey, headerValue := range queries {
		if strings.ToLower(headerKey) == strings.ToLower(header) {
			field = &headerValue
		}
	}

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	} else {
		return assertEqual(*field, value)
	}
}

func (t *testUtils) theForRequestJsonBodyShouldHaveTheVariableWithValue(index int, method, service string, dotSeparatedField, value string) error {
	requestBody := t.api.GetRequestBody(method, service, index)
	if requestBody == nil {
		return errors.New("requestBody is nil")
	}

	field := getFieldValue(requestBody, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func getFieldValue(object any, dotSeparatedField string) any {
	var field any

	var objectMap map[string]any
	objectJson, _ := json.Marshal(object)
	if err := json.Unmarshal(objectJson, &objectMap); err != nil {
		var listMap []any
		if err := json.Unmarshal(objectJson, &listMap); err != nil {
			panic(err)
		}
		field = listMap
	} else {
		field = objectMap
	}

	fields := strings.Split(dotSeparatedField, ".")

	for _, currentField := range fields {
		if i, err := strconv.Atoi(currentField); err == nil {
			field = field.([]any)[i]
		} else {
			field = field.(map[string]any)[currentField]
		}
	}

	switch v := field.(type) {
	case int:
		field = strconv.Itoa(v)
	case float64:
		field = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		field = strconv.FormatBool(v)
	}

	return field
}

func assertNull(val1 any) error {
	if val1 == nil {
		return nil
	}
	val1String, _ := json.Marshal(val1)
	return errors.New(string(val1String) + " should be nil")
}

func assertNotNull(val1 any) error {
	if val1 != nil {
		return nil
	}
	return errors.New("value should not be nil")
}

func assertEqual(current, expected any) error {
	if current == expected {
		return nil
	}
	currentString, _ := json.Marshal(current)
	expectedString, _ := json.Marshal(expected)
	if bytes.Equal(currentString, expectedString) {
		return nil
	}
	if bytes.Equal(bytes.Trim(expectedString, "\""), bytes.Trim(currentString, "\"")) {
		return nil
	}
	return errors.New(string(currentString) + " should be equal to " + string(expectedString))
}

func assertContains(current, expected any) error {
	currentString := fmt.Sprintf("%v", current)
	expectedString := fmt.Sprintf("%v", expected)

	if strings.Contains(currentString, expectedString) {
		return nil
	}

	// Try with JSON marshaled versions if direct string conversion doesn't work
	currentBytes, err1 := json.Marshal(current)
	expectedBytes, err2 := json.Marshal(expected)

	if err1 == nil && err2 == nil {
		currentJSONString := string(bytes.Trim(currentBytes, "\""))
		expectedJSONString := string(bytes.Trim(expectedBytes, "\""))

		if strings.Contains(currentJSONString, expectedJSONString) {
			return nil
		}
	}

	return fmt.Errorf("%s should contain %s", currentString, expectedString)
}

func (t *testUtils) rekognitionShouldHaveDisassociatedFacesForCollectionAndUser(expectedCount int, collectionId, userId string) error {
	actualCount := t.rekognition.CountDisassociatedFacesForUser(collectionId, userId)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d disassociated faces for collection %s and user %s, but got %d", expectedCount, collectionId, userId, actualCount)
	}
	return nil
}

func (t *testUtils) rekognitionShouldHaveDeletedFacesForCollection(expectedCount int, collectionId string) error {
	actualCount := t.rekognition.CountDeletedFacesForCollection(collectionId)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d deleted faces for collection %s, but got %d", expectedCount, collectionId, actualCount)
	}
	return nil
}

func (t *testUtils) rekognitionShouldHaveIndexedFacesForCollection(expectedCount int, collectionId string) error {
	actualCount := t.rekognition.CountIndexedFacesForCollection(collectionId)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d indexed faces for collection %s, but got %d", expectedCount, collectionId, actualCount)
	}
	return nil
}

func (t *testUtils) rekognitionShouldHaveAssociatedFacesForCollectionAndUser(expectedCount int, collectionId, userId string) error {
	actualCount := t.rekognition.CountAssociatedFacesForUser(collectionId, userId)
	if actualCount != expectedCount {
		return fmt.Errorf("expected %d associated faces for collection %s and user %s, but got %d", expectedCount, collectionId, userId, actualCount)
	}
	return nil
}

// theFileIsDuplicatedTo duplicates a file from source to destination path.
// The destination file is tracked for cleanup and destination directories are created if they don't exist.
func (t *testUtils) theFileIsDuplicatedTo(sourcePath, destinationPath string) error {
	var absSourcePath string
	var err error
	if filepath.IsAbs(sourcePath) {
		absSourcePath = sourcePath
	} else if strings.HasPrefix(sourcePath, "test/integration/resources/") {
		sourcePath = strings.TrimPrefix(sourcePath, "test/integration/resources/")
		absSourcePath, err = filepath.Abs("./../resources/" + sourcePath)
		if err != nil {
			return fmt.Errorf("failed to resolve source file path: %w", err)
		}
	} else {
		absSourcePath, err = filepath.Abs("./../resources/" + sourcePath)
		if err != nil {
			return fmt.Errorf("failed to resolve source file path: %w", err)
		}
	}

	var absDestPath string
	if filepath.IsAbs(destinationPath) {
		absDestPath = destinationPath
	} else if strings.HasPrefix(destinationPath, "test/integration/resources/") {
		destinationPath = strings.TrimPrefix(destinationPath, "test/integration/resources/")
		absDestPath, err = filepath.Abs("./../resources/" + destinationPath)
		if err != nil {
			return fmt.Errorf("failed to resolve destination file path: %w", err)
		}
	} else {
		absDestPath, err = filepath.Abs("./../resources/" + destinationPath)
		if err != nil {
			return fmt.Errorf("failed to resolve destination file path: %w", err)
		}
	}

	sourceData, err := os.ReadFile(absSourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	destDir := filepath.Dir(absDestPath)
	if err = os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory %s: %w", destDir, err)
	}

	err = os.WriteFile(absDestPath, sourceData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}

	t.createdFiles = append(t.createdFiles, absDestPath)

	return nil
}

// theFileIsDeleted deletes a file at the specified path.
func (t *testUtils) theFileIsDeleted(filePath string) error {
	var absPath string
	var err error
	if filepath.IsAbs(filePath) {
		absPath = filePath
	} else if strings.HasPrefix(filePath, "test/integration/resources/") {
		filePath = strings.TrimPrefix(filePath, "test/integration/resources/")
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
	} else {
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
	}

	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}

	err = os.Remove(absPath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// theFileShouldExist verifies that a file exists at the specified path.
func (t *testUtils) theFileShouldExist(filePath string) error {
	var absPath string
	var err error
	if filepath.IsAbs(filePath) {
		absPath = filePath
	} else if strings.HasPrefix(filePath, "test/integration/resources/") {
		filePath = strings.TrimPrefix(filePath, "test/integration/resources/")
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
	} else {
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
	}

	if _, err = os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", absPath)
	} else if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}

	return nil
}

// theFileShouldBeIdenticalTo compares two files and verifies they have identical content.
func (t *testUtils) theFileShouldBeIdenticalTo(filePath, otherFilePath string) error {
	var absPath string
	var err error
	if filepath.IsAbs(filePath) {
		absPath = filePath
	} else if strings.HasPrefix(filePath, "test/integration/resources/") {
		filePath = strings.TrimPrefix(filePath, "test/integration/resources/")
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve first file path: %w", err)
		}
	} else {
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return fmt.Errorf("failed to resolve first file path: %w", err)
		}
	}

	var absOtherPath string
	if filepath.IsAbs(otherFilePath) {
		absOtherPath = otherFilePath
	} else if strings.HasPrefix(otherFilePath, "test/integration/resources/") {
		otherFilePath = strings.TrimPrefix(otherFilePath, "test/integration/resources/")
		absOtherPath, err = filepath.Abs("./../resources/" + otherFilePath)
		if err != nil {
			return fmt.Errorf("failed to resolve second file path: %w", err)
		}
	} else {
		absOtherPath, err = filepath.Abs("./../resources/" + otherFilePath)
		if err != nil {
			return fmt.Errorf("failed to resolve second file path: %w", err)
		}
	}

	fileData, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read first file: %w", err)
	}

	otherFileData, err := os.ReadFile(absOtherPath)
	if err != nil {
		return fmt.Errorf("failed to read second file: %w", err)
	}

	if !bytes.Equal(fileData, otherFileData) {
		return fmt.Errorf("files are not identical: %s and %s have different content", absPath, absOtherPath)
	}

	return nil
}

// theParquetFileShouldHaveFieldInAllRows validates that a specific field exists and has a value (not null) in all rows.
// Reads the parquet file in chunks to avoid OOM issues.
func (t *testUtils) theParquetFileShouldHaveFieldInAllRows(filePath, fieldName string) error {
	parquetData, err := t.readParquetFile(filePath)
	if err != nil {
		return err
	}

	parquetFile, err := parquet.OpenFile(bytes.NewReader(parquetData), int64(len(parquetData)))
	if err != nil {
		return fmt.Errorf("failed to open parquet file: %w", err)
	}

	rowGroups := parquetFile.RowGroups()
	totalRows := int64(0)

	for _, rowGroup := range rowGroups {
		numRowsInGroup := rowGroup.NumRows()
		rowGroupReader := parquet.NewRowGroupReader(rowGroup)

		for i := int64(0); i < numRowsInGroup; i++ {
			row := make(map[string]any)
			err = rowGroupReader.Read(&row)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read row %d: %w", totalRows+i, err)
			}

			if _, exists := row[fieldName]; !exists {
				return fmt.Errorf("row %d is missing field '%s'", totalRows+i, fieldName)
			}

			totalRows++
		}
	}

	if totalRows == 0 {
		return fmt.Errorf("parquet file is empty")
	}

	return nil
}

// theParquetFileShouldHaveFieldInRowWithValue validates that a specific field in a specific row has an expected value.
// Row numbers are 0-indexed.
func (t *testUtils) theParquetFileShouldHaveFieldInRowWithValue(filePath, fieldName string, rowNumber int, expectedValue string) error {
	parquetData, err := t.readParquetFile(filePath)
	if err != nil {
		return err
	}

	parquetFile, err := parquet.OpenFile(bytes.NewReader(parquetData), int64(len(parquetData)))
	if err != nil {
		return fmt.Errorf("failed to open parquet file: %w", err)
	}

	rowGroups := parquetFile.RowGroups()
	currentRow := 0

	for _, rowGroup := range rowGroups {
		numRowsInGroup := rowGroup.NumRows()
		rowGroupReader := parquet.NewRowGroupReader(rowGroup)

		for i := int64(0); i < numRowsInGroup; i++ {
			if currentRow == rowNumber {
				row := make(map[string]any)
				err = rowGroupReader.Read(&row)
				if err != nil {
					return fmt.Errorf("failed to read row %d: %w", rowNumber, err)
				}

				value, exists := row[fieldName]
				if !exists {
					return fmt.Errorf("row %d does not have field '%s'", rowNumber, fieldName)
				}

				actualValue := fmt.Sprintf("%v", value)
				if actualValue != expectedValue {
					return fmt.Errorf("row %d field '%s' has value '%s', expected '%s'", rowNumber, fieldName, actualValue, expectedValue)
				}

				return nil
			}

			row := make(map[string]any)
			err = rowGroupReader.Read(&row)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read row %d: %w", currentRow, err)
			}

			currentRow++
		}
	}

	return fmt.Errorf("row %d not found in parquet file (total rows: %d)", rowNumber, currentRow)
}

// theParquetFileShouldHaveFieldsInAllRows validates that multiple fields exist in all rows.
// Field names should be comma-separated.
func (t *testUtils) theParquetFileShouldHaveFieldsInAllRows(filePath, fieldNamesStr string) error {
	fieldNames := strings.Split(fieldNamesStr, ",")
	for i, name := range fieldNames {
		fieldNames[i] = strings.TrimSpace(name)
	}

	parquetData, err := t.readParquetFile(filePath)
	if err != nil {
		return err
	}

	parquetFile, err := parquet.OpenFile(bytes.NewReader(parquetData), int64(len(parquetData)))
	if err != nil {
		return fmt.Errorf("failed to open parquet file: %w", err)
	}

	rowGroups := parquetFile.RowGroups()
	totalRows := int64(0)

	for _, rowGroup := range rowGroups {
		numRowsInGroup := rowGroup.NumRows()
		rowGroupReader := parquet.NewRowGroupReader(rowGroup)

		for i := int64(0); i < numRowsInGroup; i++ {
			row := make(map[string]any)
			err = rowGroupReader.Read(&row)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to read row %d: %w", totalRows+i, err)
			}

			for _, fieldName := range fieldNames {
				if _, exists := row[fieldName]; !exists {
					return fmt.Errorf("row %d is missing field '%s'", totalRows+i, fieldName)
				}
			}

			totalRows++
		}
	}

	if totalRows == 0 {
		return fmt.Errorf("parquet file is empty")
	}

	return nil
}

// theParquetFileShouldHaveRows validates that the parquet file has exactly the expected number of rows.
func (t *testUtils) theParquetFileShouldHaveRows(filePath string, expectedRows int) error {
	parquetData, err := t.readParquetFile(filePath)
	if err != nil {
		return err
	}

	parquetFile, err := parquet.OpenFile(bytes.NewReader(parquetData), int64(len(parquetData)))
	if err != nil {
		return fmt.Errorf("failed to open parquet file: %w", err)
	}

	rowGroups := parquetFile.RowGroups()
	totalRows := int64(0)

	for _, rowGroup := range rowGroups {
		totalRows += rowGroup.NumRows()
	}

	if int(totalRows) != expectedRows {
		return fmt.Errorf("parquet file has %d rows, expected %d", totalRows, expectedRows)
	}

	return nil
}

// theParquetFileShouldHaveAllFieldsFrom validates that the target parquet file contains all fields from the source parquet file.
func (t *testUtils) theParquetFileShouldHaveAllFieldsFrom(targetPath, sourcePath string) error {
	sourceData, err := t.readParquetFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source parquet: %w", err)
	}

	sourceFile, err := parquet.OpenFile(bytes.NewReader(sourceData), int64(len(sourceData)))
	if err != nil {
		return fmt.Errorf("failed to open source parquet: %w", err)
	}

	sourceSchema := sourceFile.Schema()
	sourceFields := make(map[string]bool)
	for _, field := range sourceSchema.Fields() {
		sourceFields[field.Name()] = true
	}

	targetData, err := t.readParquetFile(targetPath)
	if err != nil {
		return fmt.Errorf("failed to read target parquet: %w", err)
	}

	targetFile, err := parquet.OpenFile(bytes.NewReader(targetData), int64(len(targetData)))
	if err != nil {
		return fmt.Errorf("failed to open target parquet: %w", err)
	}

	targetSchema := targetFile.Schema()
	targetFields := make(map[string]bool)
	for _, field := range targetSchema.Fields() {
		targetFields[field.Name()] = true
	}

	for fieldName := range sourceFields {
		if !targetFields[fieldName] {
			return fmt.Errorf("target parquet is missing field '%s' from source", fieldName)
		}
	}

	return nil
}

// readParquetFile is a helper function to read parquet files from either local filesystem or S3 mock.
func (t *testUtils) readParquetFile(filePath string) ([]byte, error) {
	var absPath string
	var err error

	if filepath.IsAbs(filePath) {
		absPath = filePath
	} else if strings.HasPrefix(filePath, "test/integration/resources/") {
		filePath = strings.TrimPrefix(filePath, "test/integration/resources/")
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve file path: %w", err)
		}
	} else {
		absPath, err = filepath.Abs("./../resources/" + filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve file path: %w", err)
		}
	}

	data, err := os.ReadFile(absPath)
	if err == nil {
		return data, nil
	}

	return nil, fmt.Errorf("file not found: %s", absPath)
}
