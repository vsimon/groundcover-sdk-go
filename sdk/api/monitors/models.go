package monitors

import (
	"github.com/google/uuid"
	"github.com/groundcover-com/groundcover-sdk-go/sdk/models"
)

type CreateMonitorRequest struct {
	MonitorModel
}
type CreateMonitorResponse struct {
	MonitorID string `json:"monitorId"`
}

type UpdateMonitorRequest struct {
	MonitorModel
	ID uuid.UUID
}

type UpdateMonitorYamlRequest struct {
	MonitorModel
}

type MonitorModel struct {
	UUID                *uuid.UUID          `json:"-" yaml:"uuid,omitempty" hash:"ignore"`
	GrafanaMonitorID    string              `json:"-" yaml:"grafanaMonitorId,omitempty" hash:"ignore"`
	Title               string              `json:"title" yaml:"title" validate:"required"`
	Display             DisplayModel        `json:"display,omitempty" yaml:"display,omitempty"`
	Severity            string              `json:"severity,omitempty" yaml:"severity,omitempty" `
	MeasurementType     string              `json:"measurementType" yaml:"measurementType" validate:"omitempty,oneof=state event" `
	Model               Model               `json:"model" yaml:"model"`
	Labels              map[string]string   `json:"labels,omitempty" yaml:"labels,omitempty"`
	Category            string              `json:"category,omitempty" yaml:"category,omitempty"`
	Annotations         map[string]string   `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ExecutionErrorState string              `json:"executionErrorState,omitempty" yaml:"executionErrorState,omitempty" validate:"omitempty,oneof=OK Error Alerting"`
	NoDataState         string              `json:"noDataState,omitempty" yaml:"noDataState,omitempty" validate:"omitempty,oneof=OK NoData Alerting"`
	EvaluationInterval  *EvaluationInterval `json:"evaluationInterval,omitempty" yaml:"evaluationInterval,omitempty"`
	Catalog             CatalogModel        `json:"catalog,omitempty" yaml:"catalog,omitempty"`
	AutoResolve         *bool               `json:"autoResolve,omitempty" yaml:"autoResolve,omitempty" hash:"ignore"`
	Team                string              `json:"team,omitempty" yaml:"team,omitempty" hash:"ignore"`
	Routing             []string            `json:"routing,omitempty" yaml:"routing,omitempty" hash:"ignore"`
	IsPaused            *bool               `json:"isPaused,omitempty" yaml:"isPaused,omitempty" hash:"ignore"`
}

type DisplayModel struct {
	Header               *string  `json:"header,omitempty" yaml:"header,omitempty"`
	ResourceHeaderLabels []string `json:"resourceHeaderLabels,omitempty" yaml:"resourceHeaderLabels,omitempty"`
	ContextHeaderLabels  []string `json:"contextHeaderLabels,omitempty" yaml:"contextHeaderLabels,omitempty"`
	Description          string   `json:"description,omitempty" yaml:"description,omitempty"`
}

type Model struct {
	Queries    []Query        `json:"queries,omitempty" yaml:"queries,omitempty" validate:"required"`
	Reducers   []ReducerModel `json:"reducers,omitempty" yaml:"reducers,omitempty"`
	Thresholds []Threshold    `json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
}

type Query struct {
	DataType          string                 `json:"dataType,omitempty" yaml:"dataType,omitempty"`
	RelativeTimerange *RelativeTimerange     `json:"relativeTimerange,omitempty" yaml:"relativeTimerange,omitempty"`
	Name              string                 `json:"name" yaml:"name" validate:"required"`
	Expression        string                 `json:"expression,omitempty" yaml:"expression,omitempty" validate:"required_without_all=Pipeline SqlPipeline"`
	DatasourceType    string                 `json:"datasourceType,omitempty" yaml:"datasourceType,omitempty"`
	DatasourceID      string                 `json:"datasourceID,omitempty" yaml:"datasourceId,omitempty"`
	QueryType         string                 `json:"queryType,omitempty" yaml:"queryType,omitempty"`
	Pipeline          *models.PromqlPipeline `json:"pipeline,omitempty" yaml:"pipeline,omitempty" validate:"required_without_all=Expression SqlPipeline"`
	SqlPipeline       *models.Pipeline       `json:"sqlPipeline,omitempty" yaml:"sqlPipeline,omitempty" validate:"required_without_all=Expression Pipeline SqlPipeline"`
	Filters           string                 `json:"filters,omitempty" yaml:"filters,omitempty"`
	Conditions        []models.Condition     `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	InstantRollup     string                 `json:"instantRollup,omitempty" yaml:"instantRollup,omitempty"`
}

type EvaluationInterval struct {
	Interval   Duration  `json:"interval,omitempty" yaml:"interval,omitempty"`
	PendingFor *Duration `json:"pendingFor,omitempty" yaml:"pendingFor,omitempty" validation:"gtfield=Interval"`
}

type CatalogModel struct {
	CatalogId       string   `json:"id,omitempty" yaml:"id,omitempty"`
	CatalogVersion  int      `json:"version,omitempty" yaml:"version,omitempty"`
	CatalogTags     []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	CatalogCategory string   `json:"category,omitempty" yaml:"category,omitempty"`
}

type ReducerModel struct {
	Name              string             `json:"name" yaml:"name"`
	InputName         string             `json:"inputName,omitempty" yaml:"inputName,omitempty" validate:"required_unless=Type math"`
	Expression        string             `json:"expression,omitempty" yaml:"expression,omitempty" validate:"required_if=Type math"`
	Type              string             `json:"type" yaml:"type"`
	RelativeTimerange *RelativeTimerange `json:"relativeTimerange,omitempty" yaml:"relativeTimerange,omitempty"`
}

type Threshold struct {
	Name              string             `json:"name" yaml:"name" gorm:"not null" validate:"required"`
	InputName         string             `json:"inputName" yaml:"inputName" gorm:"not null" validate:"required"`
	Operator          string             `json:"operator" yaml:"operator" gorm:"not null" validate:"oneof=gt lt within_range outside_range"`
	Values            []float64          `json:"values" yaml:"values" gorm:"not null" validate:"required"`
	RelativeTimerange *RelativeTimerange `json:"relativeTimerange,omitempty" yaml:"relativeTimerange,omitempty" gorm:"embedded"`
}

type RelativeTimerange struct {
	From *Duration `json:"from,omitempty" yaml:"from,omitempty" gorm:"column:relative_timerange_from"`
	To   *Duration `json:"to,omitempty" yaml:"to,omitempty" gorm:"column:relative_timerange_to"`
}

type Duration string
