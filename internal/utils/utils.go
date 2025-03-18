package utils

// EmptyStruct is good for coalescing variatic optional parameters in func calls
func StructEmpty(data ...interface{}) interface{} {
    if len(data) > 0 {
        if s, ok := data[0].([]interface{}); ok && len(s) > 0 {
            return s[0]
        }
        return data[0]
    }

    return struct{}{}
}

// Make sure struct is not nil where it is needed
func StructNil(data interface{}) interface{} {
    if data == nil {
        return struct{}{}
    }

    return data
}
