package metrics

import (
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type QueryRequest struct {
	Start        time.Time             `form:"start" url:"start" binding:"required"`
	End          time.Time             `form:"end" url:"end" binding:"required,gtefield=Start"`
	Step         string                `form:"step" url:"step" binding:"required"`
	QueryType    string                `form:"query_type" url:"query_type" binding:"required,oneof=range instant"`
	Promql       string                `form:"promql" url:"promql"`
	Pipeline     models.PromqlPipeline `form:"pipeline" url:"pipeline"`
	Filters      string                `form:"filters" url:"filters"`
	Conditions   []models.Condition    `form:"conditions" url:"conditions"`
	SubPipelines models.KnownPipelines `form:"subPipelines" url:"subPipelines"`
}
