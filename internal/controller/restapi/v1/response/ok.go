package response

// Ok is the REST success response body. `code` mirrors the HTTP status code,
// `message` is a short human-readable summary (`"ok"`), and `data` holds the
// actual payload (null when there is nothing to return, e.g. delete).
type Ok[T any] struct {
	Code    int    `json:"code"    example:"200"`
	Message string `json:"message" example:"ok"`
	Data    T      `json:"data"`
} // @name v1.Ok
