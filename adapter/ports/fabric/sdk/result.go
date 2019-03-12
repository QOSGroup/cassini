package sdk

const (
	// DefaultResultJSON default result json string when json.Marshal error
	DefaultResultJSON string = `{"code": "500", "message": "unknown error"}`
)

// CallResult api call result
type CallResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
