package response

// Envelope is the unified REST response wrapper for both success and error
// payloads. `code` mirrors the HTTP status code, `message` carries a
// human-readable summary, and `data` holds the actual payload (null on error
// or when there is nothing to return, e.g. delete).
type Envelope struct {
	Code    int    `json:"code"    example:"200"`
	Message string `json:"message" example:"ok"`
	Data    any    `json:"data"`
} // @name v1.Envelope
