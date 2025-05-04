package k8s

import (
	"context"
	"net/http"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
)

type Service struct {
	client httpclient.Doer
}

func New(client httpclient.Doer) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) ClustersList(ctx context.Context, request *ClustersListRequest, opts ...httpclient.RequestOption) (*ClustersListResponse, error) {
	urlPath := "/api/k8s/v3/clusters/list"
	response := &ClustersListResponse{}
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

func (s *Service) WorkloadsList(ctx context.Context, request *WorkloadsListRequest, opts ...httpclient.RequestOption) (*WorkloadsListResponse, error) {
	urlPath := "/api/k8s/v3/workloads/list"
	response := &WorkloadsListResponse{}
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
