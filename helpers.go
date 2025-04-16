package main

import "strconv"

// Helper for MAX function
func max(a, b interface{}) interface{} {
	af := toFloat64(a)
	bf := toFloat64(b)
	if af > bf {
		return a
	}
	return b
}

// Helper for MIN function
func min(a, b interface{}) interface{} {
	af := toFloat64(a)
	bf := toFloat64(b)
	if af < bf {
		return a
	}
	return b
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}
