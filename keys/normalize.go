package keys

func Normalize(key interface{}) interface{} {
	if key == nil {
		return nil
	}
	switch actual := key.(type) {
	case *int64:
		if actual == nil {
			return nil
		}
		return int(*actual)
	case *int32:
		if actual == nil {
			return nil
		}
		return int(*actual)
	case *float64:
		if actual == nil {
			return nil
		}
		return int(*actual)
	case *float32:
		if actual == nil {
			return nil
		}
		return int(*actual)
	case *int16:
		if actual == nil {
			return nil
		}
		return int(*actual)
	case int32:
		return int(actual)
	case int64:
		return int(actual)
	case int16:
		return int(actual)
	case *int:
		if actual == nil {
			return nil
		}
		return *actual
	case []byte:
		if len(actual) == 0 {
			return ""
		}
		return string(actual)

	case *[]byte:
		if actual == nil {
			return nil
		}
		return string(*actual)
	case *string:
		if actual == nil {
			return nil
		}
		return *actual
	}
	return key
}
