// ... existing code ...

// EvalCondition evaluates a condition string or complex condition object against an unstructured object
func EvalCondition(obj *unstructured.Unstructured, rawCondition interface{}) (bool, error) {
	switch condition := rawCondition.(type) {
	case string:
		// Handle string condition (backward compatibility)
		return evalStringCondition(obj, condition)
	case []interface{}:
		// Handle array of conditions (AND logic, backward compatibility)
		return evalAndConditions(obj, condition)
	case map[string]interface{}:
		// Handle complex condition with operator
		return evalComplexCondition(obj, condition)
	default:
		return false, fmt.Errorf("unsupported condition type: %T", rawCondition)
	}
}

// evalStringCondition processes a single string condition (either jsonpath or condition)
func evalStringCondition(obj *unstructured.Unstructured, condition string) (bool, error) {
	// Handle jsonpath condition
	if strings.HasPrefix(condition, "jsonpath=") {
		return EvalJSONPath(obj, condition[len("jsonpath="):])
	}
	
	// Handle condition=Type[=Status] format
	if strings.HasPrefix(condition, "condition=") {
		return EvalStatusCondition(obj, condition[len("condition="):])
	}
	
	return false, fmt.Errorf("unknown condition format: %s", condition)
}

// evalAndConditions evaluates an array of conditions with AND logic
func evalAndConditions(obj *unstructured.Unstructured, conditions []interface{}) (bool, error) {
	for _, condition := range conditions {
		satisfied, err := EvalCondition(obj, condition)
		if err != nil {
			return false, err
		}
		if !satisfied {
			return false, nil
		}
	}
	return true, nil
}

// evalComplexCondition evaluates a complex condition with logical operators
func evalComplexCondition(obj *unstructured.Unstructured, condition map[string]interface{}) (bool, error) {
	operator, ok := condition["operator"].(string)
	if !ok {
		return false, fmt.Errorf("missing or invalid 'operator' field in complex condition")
	}
	
	conditionsRaw, ok := condition["conditions"].([]interface{})
	if !ok {
		return false, fmt.Errorf("missing or invalid 'conditions' field in complex condition")
	}
	
	switch strings.ToLower(operator) {
	case "and":
		return evalAndConditions(obj, conditionsRaw)
	case "or":
		return evalOrConditions(obj, conditionsRaw)
	default:
		return false, fmt.Errorf("unsupported operator: %s", operator)
	}
}

// evalOrConditions evaluates an array of conditions with OR logic
func evalOrConditions(obj *unstructured.Unstructured, conditions []interface{}) (bool, error) {
	for _, condition := range conditions {
		satisfied, err := EvalCondition(obj, condition)
		if err != nil {
			return false, err
		}
		if satisfied {
			return true, nil
		}
	}
	return false, nil
}

// ... existing code ...
