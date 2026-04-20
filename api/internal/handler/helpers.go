package handler

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
