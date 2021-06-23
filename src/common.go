package cfzap

func StringInArray(value string, list []string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}

	return false
}

func CompareStringArray(array1 []string, array2 []string) bool {
	if array1 == nil && array2 == nil {
		return true
	}
	if array1 == nil || array2 == nil {
		return false
	}
	if len(array1) != len(array2) {
		return false
	}

	for _, s := range array1 {
		if !StringInArray(s, array2) {
			return false
		}
	}

	return true
}
