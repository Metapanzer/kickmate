package error

// ErrMapping checks if the given error exists in the predefined error lists.
func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(GeneralErrors[:], UserErrors[:]...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}
	return false
}
