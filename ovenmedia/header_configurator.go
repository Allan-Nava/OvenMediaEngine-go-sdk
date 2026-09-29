package ovenmedia

import "encoding/base64"

// HeaderConfigurator holds headers sent with every request; pass it to
// [WithHeaders].
type HeaderConfigurator struct {
	Headers map[string]string
}

func InitHeaderConfigurator() *HeaderConfigurator {
	return &HeaderConfigurator{
		Headers: map[string]string{},
	}
}

func (h *HeaderConfigurator) SetHeader(key string, value string) {
	h.Headers[key] = value
}

func (h *HeaderConfigurator) GetHeader(key string) *string {
	v, ok := h.Headers[key]
	if !ok {
		return nil
	}
	return &v
}

func (h *HeaderConfigurator) GetHeaders() map[string]string {
	return h.Headers
}

func (h *HeaderConfigurator) SetHeaders(headers map[string]string) {
	h.Headers = headers
}

func (h *HeaderConfigurator) DeleteHeader(key string) {
	delete(h.Headers, key)
}

func (h *HeaderConfigurator) DeleteHeaders() {
	h.Headers = make(map[string]string)
}

func (h *HeaderConfigurator) HasHeader(key string) bool {
	_, ok := h.Headers[key]
	return ok
}

func (h *HeaderConfigurator) HasHeaders() bool {
	return len(h.Headers) > 0
}

func (h *HeaderConfigurator) GetHeaderKeys() []string {
	keys := make([]string, len(h.Headers))
	i := 0
	for k := range h.Headers {
		keys[i] = k
		i++
	}
	return keys
}

func (h *HeaderConfigurator) GetHeaderValues() []string {
	values := make([]string, len(h.Headers))
	i := 0
	for _, v := range h.Headers {
		values[i] = v
		i++
	}
	return values
}

func (h *HeaderConfigurator) GetHeaderKeyValuePairs() map[string]string {
	return h.Headers
}

// CreateBasicAuthHeader sets Authorization: Basic base64(username:password),
// for an OME AccessToken of the form "username:password".
func (h *HeaderConfigurator) CreateBasicAuthHeader(username string, password string) {
	h.Headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

// CreateBasicAuthHeaderEncoded sets Authorization: Basic <base64EncodedToken>.
func (h *HeaderConfigurator) CreateBasicAuthHeaderEncoded(base64EncodedToken string) {
	h.Headers["Authorization"] = "Basic " + base64EncodedToken
}

// CreateOmeBasicAuthHeaderWord sets Authorization: Basic base64(word), for an
// AccessToken of any form.
func (h *HeaderConfigurator) CreateOmeBasicAuthHeaderWord(word string) {
	h.Headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(word))
}

// CreateOmeBasicAuthHeader used to set an "ome-access-token" header, which
// OME doesn't read. It now sets Authorization, like CreateBasicAuthHeader.
//
// Deprecated: use [CreateBasicAuthHeader] or [WithAccessToken].
func (h *HeaderConfigurator) CreateOmeBasicAuthHeader(username string, password string) {
	h.CreateBasicAuthHeader(username, password)
}

// CreateOmeBasicAuthHeaderEncoded used to set an "ome-access-token" header,
// which OME doesn't read. It now sets Authorization.
//
// Deprecated: use [CreateBasicAuthHeaderEncoded].
func (h *HeaderConfigurator) CreateOmeBasicAuthHeaderEncoded(base64EncodedToken string) {
	h.CreateBasicAuthHeaderEncoded(base64EncodedToken)
}
