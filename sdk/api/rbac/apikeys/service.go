package apikeys

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type Service struct {
	client httpclient.Doer
}

func New(client httpclient.Doer) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) CreateApiKey(ctx context.Context, request *CreateApiKeyRequest, opts ...httpclient.RequestOption) (*CreateApiKeyResponse, error) {
	urlPath := "/api/rbac/apikey/create"
	response := &CreateApiKeyResponse{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodPost,
		Path:     urlPath,
		Body:     request,
		Response: response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) DeleteApiKey(ctx context.Context, id string, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	urlPath := "/api/rbac/apikey/{id}"
	urlPath = strings.Replace(urlPath, "{id}", id, 1)
	response := &models.EmptyResponse{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method: http.MethodDelete,
		Path:   urlPath,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) ListApiKeys(ctx context.Context, withRevoked *bool, withExpired *bool, opts ...httpclient.RequestOption) ([]ListApiKeysResponseItem, error) {
	urlPath := "/api/rbac/apikeys/list"

	queryParams := url.Values{}
	if withRevoked != nil {
		queryParams.Add("withRevoked", fmt.Sprint(*withRevoked))
	}
	if withExpired != nil {
		queryParams.Add("withExpired", fmt.Sprint(*withExpired))
	}

	response := []ListApiKeysResponseItem{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:      http.MethodGet,
		Path:        urlPath,
		QueryParams: queryParams,
		Response:    &response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}
