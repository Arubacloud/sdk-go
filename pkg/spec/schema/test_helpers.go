package schema

// Test helper functions shared across all domain test files

// StringPtr returns a pointer to the given string value
func StringPtr(s string) *string {
	return &s
}

// Int32Ptr returns a pointer to the given int32 value
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr returns a pointer to the given int64 value
func Int64Ptr(i int64) *int64 {
	return &i
}

// BoolPtr returns a pointer to the given bool value
func BoolPtr(b bool) *bool {
	return &b
}

// Float64Ptr returns a pointer to the given float64 value
func Float64Ptr(f float64) *float64 {
	return &f
}
