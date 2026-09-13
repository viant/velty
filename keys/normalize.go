package keys

func Normalize(key interface{}) interface{} {
	if key == nil {
		return nil
	}
	switch actual := key.(type) {
	case *uint64:
		if actual == nil {
			return nil
		}
		return *actual
	case *uint32:
		if actual == nil {
			return nil
		}
		return *actual
	case *uint16:
		if actual == nil {
			return nil
		}
		return *actual
	case *uint8:
		if actual == nil {
			return nil
		}
		return *actual
	case *uint:
		if actual == nil {
			return nil
		}
		return *actual
	case *uintptr:
		if actual == nil {
			return nil
		}
		return *actual
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
