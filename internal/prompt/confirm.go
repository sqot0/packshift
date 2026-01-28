package prompt

func Confirm(label string) (bool, error) {
	result, err := Select(label, []string{"Yes", "No"})
	if err != nil {
		return false, err
	}
	return result == 0, nil
}
