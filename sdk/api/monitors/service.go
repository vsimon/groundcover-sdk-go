package monitors

import (
	"context"
	"net/http"
	"net/url"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
	"gopkg.in/yaml.v2"
)

type Service struct {
	client httpclient.Doer
}

func New(client httpclient.Doer) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) CreateMonitorYaml(ctx context.Context, monitorYaml []byte, opts ...httpclient.RequestOption) (*CreateMonitorResponse, error) {
	urlPath := "/api/monitors"
	response := &CreateMonitorResponse{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:      http.MethodPost,
		Path:        urlPath,
		Body:        monitorYaml,
		Response:    response,
		ContentType: httpclient.ContentTypeYAML,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) CreateMonitor(ctx context.Context, request *CreateMonitorRequest, opts ...httpclient.RequestOption) (*CreateMonitorResponse, error) {
	monitorYaml, err := yaml.Marshal(request.MonitorModel)
	if err != nil {
		return nil, err
	}

	return s.CreateMonitorYaml(ctx, monitorYaml, opts...)
}

func (s *Service) GetMonitor(ctx context.Context, id string, opts ...httpclient.RequestOption) ([]byte, error) {
	urlPath, err := url.JoinPath("/api/monitors", id)
	if err != nil {
		return nil, err
	}
	response := []byte{}
	err = s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodGet,
		Path:     urlPath,
		Response: &response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) UpdateMonitor(ctx context.Context, id string, request *UpdateMonitorRequest, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	monitorYaml, err := yaml.Marshal(request)
	if err != nil {
		return nil, err
	}

	return s.UpdateMonitorYaml(ctx, id, monitorYaml, opts...)
}

func (s *Service) UpdateMonitorYaml(ctx context.Context, id string, monitorYaml []byte, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	urlPath, err := url.JoinPath("/api/monitors", id)
	if err != nil {
		return nil, err
	}
	response := &models.EmptyResponse{}
	err = s.client.Do(ctx, &httpclient.CallRequest{
		Method:      http.MethodPut,
		Path:        urlPath,
		Body:        monitorYaml,
		ContentType: httpclient.ContentTypeYAML,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) DeleteMonitor(ctx context.Context, id string, opts ...httpclient.RequestOption) (*models.EmptyResponse, error) {
	urlPath, err := url.JoinPath("/api/monitors", id)
	if err != nil {
		return nil, err
	}
	response := &models.EmptyResponse{}
	err = s.client.Do(ctx, &httpclient.CallRequest{
		Method: http.MethodDelete,
		Path:   urlPath,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}
