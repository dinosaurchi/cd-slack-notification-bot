package notifier

func GetCDMessage(
	relatedThreadLink string,
	statuses []string,
) (string, error) {
	if isCDFailed(statuses) {
		return "" +
			"CD failed" +
			" - " + relatedThreadLink, nil
	}
	return "" +
		"CD succeeded" +
		" - " + relatedThreadLink, nil
}

func GetSuccessfulCDMessage(
	codeBuildURL string,
) (string, error) {
	return "" +
		"CD succeeded" +
		" - " + codeBuildURL, nil
}

func isCDFailed(
	statuses []string,
) bool {
	for _, status := range statuses {
		if status != "success" {
			return true
		}
	}
	return false
}
