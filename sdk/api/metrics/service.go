package metrics

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

func (s *Service) Query(ctx context.Context, request *QueryRequest, opts ...httpclient.RequestOption) ([]byte, error) {
	urlPath := "/api/metrics/query"
	response := []byte{}
	err := s.client.Do(ctx, &httpclient.CallRequest{
		Method:   http.MethodPost,
		Path:     urlPath,
		Body:     request,
		Response: &response,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return response, nil
}
