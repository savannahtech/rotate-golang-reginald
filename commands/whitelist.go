package commands

func GetWhitelist() map[string]bool {
	return map[string]bool{
		"ls":     true,
		"pwd":    true,
		"whoami": true,
		"date":   true,
	}

}
