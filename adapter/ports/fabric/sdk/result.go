package sdk

const (
	// errUnsuportedToken error response json
	errUnsuportedToken = `{"code": 404, "message": "unsupported chain's network or token"}`
	// defaultResultJSON success response json
	defaultResultJSON string = `{"code": 500, "message": "unknown error"}`
)

// CallResult api call result
type CallResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
