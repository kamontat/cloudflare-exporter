package cloudflare

func filterT[T any](data []T, whitelist []string, blacklist []string, comparison func(data T) string) (output []T) {
	output = make([]T, 0)

	if len(whitelist) < 1 && len(blacklist) < 1 {
		return data
	}

	whitelistMap := make(map[string]bool)
	for _, w := range whitelist {
		whitelistMap[w] = true
	}
	blacklistMap := make(map[string]bool)
	for _, w := range blacklist {
		blacklistMap[w] = true
	}

	for _, v := range data {
		var actual = comparison(v)
		isWhite := whitelistMap[actual]
		isBlack := blacklistMap[actual]
		if isWhite && !isBlack {
			output = append(output, v)
		}
	}

	return
}
