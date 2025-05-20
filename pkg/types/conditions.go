package types

// Condition fields
const (
	ConditionOriginRoot = "root"
)

// Condition keys
const (
	ConditionKeyNamespace = "namespace"
	ConditionKeyWorkload  = "workload"
	ConditionKeyPodName   = "podName"
	ConditionKeyReason    = "reason"
	ConditionKeyType      = "type"
	ConditionKeyEnv       = "env"
	ConditionKeyInstance  = "instance"
)

// Condition values
const (
	ConditionValueOOMKilled          = "OOMKilled"
	ConditionValueTypeContainerCrash = "container_crash"
)

// Filter operators
const (
	OperatorEqual                 = "eq"
	OperatorNotEqual              = "ne"
	OperatorContains              = "contains"
	OperatorNotContains           = "notcontains"
	OperatorContainsIgnoreCase    = "icontains"
	OperatorNotContainsIgnoreCase = "inotcontains"
	OperatorStartsWith            = "startswith"
	OperatorStartsWithIgnoreCase  = "istartswith"
)

// Condition Types exposed to the user, mapping to backend-expected type strings.
const (
	ConditionTypeString      = "string"       // For general string conditions
	ConditionTypeInt64       = "int64"        // For integer numbers
	ConditionTypeFloat64     = "float64"      // For floating-point numbers
	ConditionTypeBool        = "bool"         // For boolean values
	ConditionTypeDatetime    = "datetime"     // For time/date values
	ConditionTypeStringArray = "string_array" // For arrays of strings
)
