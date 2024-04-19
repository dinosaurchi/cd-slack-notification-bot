package slack

func GetThreadTimestamp(
	threadTimestamp string,
	timeStamp string,
) string {
	if threadTimestamp != "" {
		return threadTimestamp
	} else if timeStamp != "" {
		return timeStamp
	}
	panic("No thread timestamp or timestamp found in the message: ")
}
