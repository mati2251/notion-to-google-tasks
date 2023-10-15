package sync

func updates() ([]string, error) {
	_, err := DB.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	return []string{}, nil
}
