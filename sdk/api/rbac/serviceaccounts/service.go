package serviceaccounts

import (
	"context"
	"net/http"
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

func (s *Service) CreateServiceAccount(ctx context.Context, request *CreateServiceAccountRequest, opts ...httpclient.RequestOption) (*CreateServiceAccountResponse, error) {
	urlPath := "/api/rbac/service-account/create"
	response := &CreateServiceAccountResponse{}
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

func (s *Service) DeleteServiceAccount(ctx context.Context, id string, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	urlPath := "/api/rbac/service-account/{id}"
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

func (s *Service) ListServiceAccounts(ctx context.Context, opts ...httpclient.RequestOption) ([]ListServiceAccountsResponseItem, error) {
	urlPath := "/api/rbac/service-accounts/list"
	response := []ListServiceAccountsResponseItem{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodGet,
		Path:     urlPath,
		Response: &response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) UpdateServiceAccount(ctx context.Context, request *UpdateServiceAccountRequest, opts ...httpclient.RequestOption) (*UpdateServiceAccountResponse, error) {
	urlPath := "/api/rbac/service-account/update"
	response := &UpdateServiceAccountResponse{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodPut,
		Path:     urlPath,
		Body:     request,
		Response: response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}
