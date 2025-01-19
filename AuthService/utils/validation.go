package utils

func Validate(username, password string) (message string) {
	if len(username) < 3 {
		return "Username must be 3 character"
	}
	if len(password) < 3 {
		return "Password must be atlest 5 character"
	}

	return "Passed"
}
