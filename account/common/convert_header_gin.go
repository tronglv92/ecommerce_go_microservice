package common

import "net/http"

func ConvertHeaders(ginHeaders http.Header) map[string]string {
	restyHeaders := make(map[string]string)
	for key, values := range ginHeaders {
		// Resty uses a map of string to string for headers
		// If there are multiple values for a header, you might want to concatenate them
		// or use a comma-separated string, depending on your use case
		restyHeaders[key] = values[0]
	}
	return restyHeaders
}
