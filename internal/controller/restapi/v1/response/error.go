package response

// Error is the REST error response body. It carries the HTTP status code and a
// human-readable message; except for `code` it mirrors nothing else from
// Ok, so callers can distinguish success and error bodies by shape
// (an error body never has a `data` field).
type Error struct {
	Code    int    `json:"code"    example:"400"`
	Message string `json:"message" example:"invalid request body"`
} // @name v1.Error
