package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/types"
)

// ConditionSet provides a fluent interface for building a list of conditions.
// It uses default values for Origin, Type, and Operator, which can be overridden.
type ConditionSet struct {
	conditions      []*models.Condition
	defaultOrigin   string
	defaultCondType string
	defaultOpStr    string
}

// NewConditionSet creates a new ConditionSet with default settings.
// Defaults are typically ConditionOriginRoot, ConditionTypeString, and OperatorEqual from the types package.
func NewConditionSet() *ConditionSet {
	return &ConditionSet{
		conditions:      []*models.Condition{},
		defaultOrigin:   types.ConditionOriginRoot,
		defaultCondType: types.ConditionTypeString,
		defaultOpStr:    types.OperatorEqual,
	}
}

// Add appends a new condition to the set. It infers the condition type
// (e.g., string, int64, float64, bool, datetime, string_array)
// from the provided Go type of the value and uses default origin and operator.
func (cs *ConditionSet) Add(key string, value interface{}) *ConditionSet {
	var valueStr string
	var condType string

	switch v := value.(type) {
	case string:
		valueStr = v
		condType = types.ConditionTypeString
	case []string:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			// Handle error: fallback to fmt.Sprintf or log, then use default string type
			// For simplicity in this example, we'll use Sprintf and default type
			valueStr = fmt.Sprintf("%v", v) // e.g., "[elem1 elem2]"
			condType = cs.defaultCondType   // Or a specific error type if defined
		} else {
			valueStr = string(jsonBytes)
			condType = types.ConditionTypeStringArray
		}
	case int, int8, int16, int32, int64:
		valueStr = strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
		condType = types.ConditionTypeInt64
	case uint, uint8, uint16, uint32, uint64:
		valueStr = strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
		condType = types.ConditionTypeInt64 // Backend might still expect "int64" for unsigned
	case float32:
		valueStr = strconv.FormatFloat(float64(v), 'f', -1, 32)
		condType = types.ConditionTypeFloat64
	case float64:
		valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		condType = types.ConditionTypeFloat64
	case time.Time:
		valueStr = v.Format(time.RFC3339)
		condType = types.ConditionTypeDatetime
	case bool:
		valueStr = strconv.FormatBool(v)
		condType = types.ConditionTypeBool
	default:
		valueStr = fmt.Sprintf("%v", v)
		condType = cs.defaultCondType // Fallback to default string type
	}

	return cs.addInternal(key, cs.defaultOrigin, condType, valueStr, cs.defaultOpStr)
}

// AddRawCondition appends a pre-constructed *models.Condition struct directly to the set.
// This is useful for complex conditions not easily built by other helper methods.
func (cs *ConditionSet) AddRawCondition(condition *models.Condition) *ConditionSet {
	if condition != nil {
		cs.conditions = append(cs.conditions, condition)
	}
	return cs
}

// AddFull appends a new condition to the set using explicitly provided origin, type, operator, and value.
func (cs *ConditionSet) AddFull(key, origin, condType, value, opStr string) *ConditionSet {
	return cs.addInternal(key, origin, condType, value, opStr)
}

// addInternal is a helper to construct and append the condition, returning the ConditionSet for chaining.
func (cs *ConditionSet) addInternal(key, origin, condType, value, opStr string) *ConditionSet {
	condition := &models.Condition{
		Key:    key,
		Origin: origin,
		Type:   condType,
		Filters: []*models.Filter{
			{Op: models.Op(opStr), Value: value},
		},
	}
	cs.conditions = append(cs.conditions, condition)
	return cs
}

// Build returns the final slice of *models.Condition.
func (cs *ConditionSet) Build() []*models.Condition {
	return cs.conditions
}

// AddOOMEventConditions appends a predefined set of conditions to identify OOM events.
// It adds a condition for reason=OOMKilled and type=container_crash.
func (cs *ConditionSet) AddOOMEventConditions() *ConditionSet {
	cs.Add(types.ConditionKeyReason, types.ConditionValueOOMKilled)
	cs.Add(types.ConditionKeyType, types.ConditionValueTypeContainerCrash)
	return cs
}
