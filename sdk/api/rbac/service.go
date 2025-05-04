package rbac

import (
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/rbac/apikeys"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/rbac/policies"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/api/rbac/serviceaccounts"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/httpclient"
)

type Service struct {
	Apikeys         *apikeys.Service
	Policies        *policies.Service
	Serviceaccounts *serviceaccounts.Service
}

func New(client httpclient.Doer) *Service {
	return &Service{
		Apikeys:         apikeys.New(client),
		Policies:        policies.New(client),
		Serviceaccounts: serviceaccounts.New(client),
	}
}
