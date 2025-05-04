package policies

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

func (s *Service) ApplyPolicy(ctx context.Context, request *ApplyPolicyRequest, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	urlPath := "/api/rbac/policy/apply"
	response := &models.EmptyResponse{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method: http.MethodPost,
		Path:   urlPath,
		Body:   request,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) CreatePolicy(ctx context.Context, request *CreatePolicyRequest, opts ...httpclient.RequestOption) (*Policy, error) {
	urlPath := "/api/rbac/policy/create"
	response := &Policy{}
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

func (s *Service) DeletePolicy(ctx context.Context, id string, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	urlPath := "/api/rbac/policy/{id}"
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

func (s *Service) GetPolicies(ctx context.Context, opts ...httpclient.RequestOption) ([]PolicyWithEntityCount, error) {
	urlPath := "/api/rbac/policies/list"
	response := []PolicyWithEntityCount{}
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

func (s *Service) GetPolicy(ctx context.Context, id string, opts ...httpclient.RequestOption) (*Policy, error) {
	urlPath := "/api/rbac/policy/{id}"
	urlPath = strings.Replace(urlPath, "{id}", id, 1)
	response := &Policy{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodGet,
		Path:     urlPath,
		Response: response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) GetPolicyAuditTrail(ctx context.Context, id string, opts ...httpclient.RequestOption) ([]Policy, error) {
	urlPath := "/api/rbac/policy/{id}/auditTrail"
	urlPath = strings.Replace(urlPath, "{id}", id, 1)
	response := []Policy{}
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

func (s *Service) UpdatePolicy(ctx context.Context, id string, request *UpdatePolicyRequest, opts ...httpclient.RequestOption) (*Policy, error) {
	urlPath := "/api/rbac/policy/{id}"
	urlPath = strings.Replace(urlPath, "{id}", id, 1)
	response := &Policy{}
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
