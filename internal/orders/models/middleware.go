package models

type Logger struct {
	RequestMethod string `json:"request_method"`
	StatusCode    int    `json:"status_code"`
	ExecutonTime  string `json:"execution_time"`
	RequestID     string `json:"request_id"`
	RequestURI    string `json:"request_uri"`
	RequestSize   int    `json:"request_size"`
	ClientIP      string `json:"client_ip"`
	Errors        int    `json:"errors,omitempty"`
}
