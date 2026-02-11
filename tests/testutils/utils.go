package testutils

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// IntPtr 返回整数指针
func IntPtr(i int) *int {
	return &i
}

// Float64Ptr 返回浮点数指针
func Float64Ptr(f float64) *float64 {
	return &f
}