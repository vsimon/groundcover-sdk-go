package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/groundcover-com/groundcover-sdk-go/pkg/models"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/types"
)

func TestNewConditionSet(t *testing.T) {
	cs := NewConditionSet()

	if cs == nil {
		t.Fatal("NewConditionSet() returned nil")
	}
	if cs.defaultOrigin != types.ConditionOriginRoot {
		t.Errorf("Expected defaultOrigin %s, got %s", types.ConditionOriginRoot, cs.defaultOrigin)
	}
	if cs.defaultCondType != types.ConditionTypeString { // Used as fallback type in Add
		t.Errorf("Expected defaultCondType %s, got %s", types.ConditionTypeString, cs.defaultCondType)
	}
	if cs.defaultOpStr != types.OperatorEqual {
		t.Errorf("Expected defaultOpStr %s, got %s", types.OperatorEqual, cs.defaultOpStr)
	}
	if len(cs.conditions) != 0 {
		t.Errorf("Expected initial conditions to be empty, got %d", len(cs.conditions))
	}
}

func TestConditionSet_Add(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)

	stringArray := []string{"apple", "banana"}
	stringArrayJSON, _ := json.Marshal(stringArray)

	testCases := []struct {
		name           string
		key            string
		value          interface{}
		expectedType   string
		expectedValStr string
	}{
		{"string value", "strKey", "testString", types.ConditionTypeString, "testString"},
		{"[]string value", "arrStrKey", stringArray, types.ConditionTypeStringArray, string(stringArrayJSON)},
		{"int value", "intKey", 123, types.ConditionTypeInt64, "123"},
		{"int8 value", "int8Key", int8(12), types.ConditionTypeInt64, "12"},
		{"int16 value", "int16Key", int16(1234), types.ConditionTypeInt64, "1234"},
		{"int32 value", "int32Key", int32(12345), types.ConditionTypeInt64, "12345"},
		{"int64 value", "int64Key", int64(123456), types.ConditionTypeInt64, "123456"},
		{"uint value", "uintKey", uint(789), types.ConditionTypeInt64, "789"},
		{"uint64 value", "uint64Key", uint64(789012), types.ConditionTypeInt64, "789012"},
		{"float32 value", "float32Key", float32(12.3), types.ConditionTypeFloat64, "12.3"},
		{"float64 value", "float64Key", 78.9, types.ConditionTypeFloat64, "78.9"},
		{"time.Time value", "timeKey", now, types.ConditionTypeDatetime, nowStr},
		{"bool true value", "boolKeyTrue", true, types.ConditionTypeBool, "true"},
		{"bool false value", "boolKeyFalse", false, types.ConditionTypeBool, "false"},
		{"struct value (fallback)", "structKey", struct{ A int }{5}, types.ConditionTypeString, "{5}"}, // fmt.Sprintf("%v")
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cs := NewConditionSet()
			cs.Add(tc.key, tc.value)

			conditions := cs.Build()
			if len(conditions) != 1 {
				t.Fatalf("Expected 1 condition, got %d", len(conditions))
			}

			cond := conditions[0]
			if cond.Key != tc.key {
				t.Errorf("Key: expected %s, got %s", tc.key, cond.Key)
			}
			if cond.Origin != types.ConditionOriginRoot {
				t.Errorf("Origin: expected %s, got %s", types.ConditionOriginRoot, cond.Origin)
			}
			if cond.Type != tc.expectedType {
				t.Errorf("Type: expected %s, got %s", tc.expectedType, cond.Type)
			}
			if len(cond.Filters) != 1 {
				t.Fatalf("Expected 1 filter, got %d", len(cond.Filters))
			}
			if string(cond.Filters[0].Op) != types.OperatorEqual {
				t.Errorf("Op: expected %s, got %s", types.OperatorEqual, cond.Filters[0].Op)
			}
			if cond.Filters[0].Value != tc.expectedValStr {
				t.Errorf("Value: expected %s, got %s", tc.expectedValStr, cond.Filters[0].Value)
			}
		})
	}

	// Test case for json.Marshal error in []string
	t.Run("[]string json.Marshal error", func(t *testing.T) {
		cs := NewConditionSet()
		// float64 cannot be marshalled into a JSON string array directly, causing json.Marshal to error
		// This is a bit of a contrived way to make json.Marshal fail for typical Add usage,
		// as []string should contain strings. A custom type that fails to marshal would be better.
		// However, we can test the fallback by passing a type that our Add method doesn't explicitly handle for StringArray
		// and which also is not one of the other explicitly handled types (int, float, bool, time.Time)
		// Let's simulate with a map, which hits the default fmt.Sprintf in Add and also the json.Marshal fallback in Add if it were []interface{}
		// The current code for []string branch will only take []string.
		// To test the json.Marshal error path within the []string case is hard without a custom unmarshallable type.
		// Instead, we test the general default fallback for an unhandled type.

		complexValue := map[string]interface{}{"fail": make(chan int)} // This will cause fmt.Sprintf to produce something like map[fail:0x...]
		cs.Add("fallbackKey", complexValue)
		conditions := cs.Build()
		cond := conditions[0]
		if cond.Type != types.ConditionTypeString { // Should use defaultCondType which is ConditionTypeString
			t.Errorf("Expected fallback type %s, got %s", types.ConditionTypeString, cond.Type)
		}
		expectedFallbackValue := fmt.Sprintf("%v", complexValue)
		if cond.Filters[0].Value != expectedFallbackValue {
			t.Errorf("Expected fallback value %s, got %s", expectedFallbackValue, cond.Filters[0].Value)
		}
	})
}

func TestConditionSet_AddFull(t *testing.T) {
	cs := NewConditionSet()
	key := "fullKey"
	origin := "customOrigin"
	condType := "customType"
	value := "fullValue"
	opStr := types.OperatorNotEqual

	cs.AddFull(key, origin, condType, value, opStr)

	conditions := cs.Build()
	if len(conditions) != 1 {
		t.Fatalf("Expected 1 condition, got %d", len(conditions))
	}

	cond := conditions[0]
	if cond.Key != key {
		t.Errorf("Expected Key %s, got %s", key, cond.Key)
	}
	if cond.Origin != origin {
		t.Errorf("Expected Origin %s, got %s", origin, cond.Origin)
	}
	if cond.Type != condType {
		t.Errorf("Expected Type %s, got %s", condType, cond.Type)
	}
	if len(cond.Filters) != 1 {
		t.Fatalf("Expected 1 filter, got %d", len(cond.Filters))
	}
	if string(cond.Filters[0].Op) != opStr {
		t.Errorf("Expected Op %s, got %s", opStr, cond.Filters[0].Op)
	}
	if cond.Filters[0].Value != value {
		t.Errorf("Expected Value %s, got %s", value, cond.Filters[0].Value)
	}
}

func TestConditionSet_AddOOMEventConditions(t *testing.T) {
	cs := NewConditionSet()
	cs.AddOOMEventConditions()

	conditions := cs.Build()
	if len(conditions) != 2 {
		t.Fatalf("Expected 2 conditions for OOM, got %d", len(conditions))
	}

	expectedOOMReason := &models.Condition{
		Key:    types.ConditionKeyReason,
		Origin: types.ConditionOriginRoot,
		Type:   types.ConditionTypeString,
		Filters: []*models.Filter{
			{Op: models.Op(types.OperatorEqual), Value: types.ConditionValueOOMKilled},
		},
	}
	expectedOOMType := &models.Condition{
		Key:    types.ConditionKeyType,
		Origin: types.ConditionOriginRoot,
		Type:   types.ConditionTypeString,
		Filters: []*models.Filter{
			{Op: models.Op(types.OperatorEqual), Value: types.ConditionValueTypeContainerCrash},
		},
	}

	if !reflect.DeepEqual(conditions[0], expectedOOMReason) {
		t.Errorf("First OOM condition (Reason) mismatch.\nExpected: %+v\nGot:      %+v", expectedOOMReason, conditions[0])
	}
	if !reflect.DeepEqual(conditions[1], expectedOOMType) {
		t.Errorf("Second OOM condition (Type) mismatch.\nExpected: %+v\nGot:      %+v", expectedOOMType, conditions[1])
	}
}

func TestConditionSet_Build_Empty(t *testing.T) {
	cs := NewConditionSet()
	conditions := cs.Build()
	if len(conditions) != 0 {
		t.Errorf("Expected 0 conditions for an empty set, got %d", len(conditions))
	}
}

func TestConditionSet_Chaining(t *testing.T) {
	cs := NewConditionSet()
	tm := time.Now()
	returnedCs := cs.Add("key1", "value1").
		Add("keyInt", 42).
		Add("keyFloat", 3.14).
		Add("keyBool", true).
		Add("keyTime", tm).
		Add("keyStrArr", []string{"a", "b"}).
		AddFull("key2", "o2", "customTypeForFull", "v2", types.OperatorContains).
		AddOOMEventConditions()

	if cs != returnedCs {
		t.Error("Chaining broken: methods did not return the original ConditionSet pointer")
	}

	conditions := cs.Build()
	// 1+1+1+1+1+1 (Add) + 1 (AddFull) + 2 (OOM) = 9
	if len(conditions) != 9 {
		t.Errorf("Expected 9 conditions after chaining, got %d. Conditions: %+v", len(conditions), conditions)
	}
}

func TestConditionSet_MixedOperations(t *testing.T) {
	cs := NewConditionSet()
	tm := time.Now().Add(time.Hour)
	strArr := []string{"x", "y", "z"}

	cs.Add(types.ConditionKeyNamespace, "ns1")
	cs.Add(types.ConditionKeyEnv, 123) // Int -> Int64 type
	cs.AddOOMEventConditions()
	cs.AddFull(types.ConditionKeyPodName, types.ConditionOriginRoot, types.ConditionTypeString, "pod-abc", types.OperatorEqual)
	cs.Add("metricValue", 99.9) // Float64 -> Float64 type
	cs.Add("eventTime", tm)
	cs.Add("tags", strArr)

	conditions := cs.Build()

	// 1 + 1 + 2 (OOM) + 1 + 1 + 1 + 1 = 8
	if len(conditions) != 8 {
		t.Fatalf("Expected 8 conditions in mixed operations, got %d. Conditions: %+v", len(conditions), conditions)
	}

	// Spot check a few conditions for type and value
	if conditions[0].Key != types.ConditionKeyNamespace || conditions[0].Type != types.ConditionTypeString || conditions[0].Filters[0].Value != "ns1" {
		t.Errorf("Mismatch in first condition (ns1): %+v", conditions[0])
	}
	if conditions[1].Key != types.ConditionKeyEnv || conditions[1].Type != types.ConditionTypeInt64 || conditions[1].Filters[0].Value != "123" {
		t.Errorf("Mismatch in second condition (env 123): %+v", conditions[1])
	}
	// OOM conditions are [2] and [3]
	if conditions[2].Type != types.ConditionTypeString {
		t.Errorf("OOM1 type error: %+v", conditions[2])
	}
	if conditions[3].Type != types.ConditionTypeString {
		t.Errorf("OOM2 type error: %+v", conditions[3])
	}

	if conditions[4].Key != types.ConditionKeyPodName || conditions[4].Type != types.ConditionTypeString || conditions[4].Filters[0].Value != "pod-abc" {
		t.Errorf("Mismatch in AddFull condition (pod-abc): %+v", conditions[4])
	}
	if conditions[5].Key != "metricValue" || conditions[5].Type != types.ConditionTypeFloat64 || conditions[5].Filters[0].Value != "99.9" {
		t.Errorf("Mismatch in float condition (metricValue): %+v", conditions[5])
	}

	if conditions[6].Key != "eventTime" || conditions[6].Type != types.ConditionTypeDatetime || conditions[6].Filters[0].Value != tm.Format(time.RFC3339) {
		t.Errorf("Mismatch in time condition (eventTime): %+v", conditions[6])
	}
	strArrJSON, _ := json.Marshal(strArr)
	if conditions[7].Key != "tags" || conditions[7].Type != types.ConditionTypeStringArray || conditions[7].Filters[0].Value != string(strArrJSON) {
		t.Errorf("Mismatch in string array condition (tags): %+v", conditions[7])
	}
}

func TestConditionSet_AddRawCondition(t *testing.T) {
	cs := NewConditionSet()

	rawCondition := &models.Condition{
		Key:    "rawKey",
		Origin: "rawOrigin",
		Type:   "rawType",
		Filters: []*models.Filter{
			{Op: models.Op("rawOp"), Value: "rawValue"},
		},
	}

	cs.AddRawCondition(rawCondition)

	conditions := cs.Build()
	if len(conditions) != 1 {
		t.Fatalf("Expected 1 condition after AddRawCondition, got %d", len(conditions))
	}

	if !reflect.DeepEqual(conditions[0], rawCondition) {
		t.Errorf("AddRawCondition did not append the condition correctly.\nExpected: %+v\nGot:      %+v", rawCondition, conditions[0])
	}

	// Test adding nil condition
	cs.AddRawCondition(nil)
	if len(cs.Build()) != 1 {
		t.Errorf("Expected condition count to remain 1 after adding nil, got %d", len(cs.Build()))
	}
}
