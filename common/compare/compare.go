package compare

// InList If x is one of list
func InList[T comparable](x T, list []T) bool {
	for _, now := range list {
		if x == now {
			return true
		}
	}
	return false
}
