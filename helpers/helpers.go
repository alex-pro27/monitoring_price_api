package helpers

func IndexOf(element interface{}, data []interface{}) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func Unique(intSlice []uint) []uint {
	keys := make(map[uint]bool)
	var list []uint
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
