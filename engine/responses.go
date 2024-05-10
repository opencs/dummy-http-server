// Copyright (c) 2023-2024, Open Communications Security
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//  3. Neither the name of the copyright holder nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package engine

import (
	"encoding/base64"
	"io"
	"net/http"
	"regexp"

	"gitlab.opencs.dev.br/opencs-commons/dummy-http-server/config"
)

var (
	// An empty JSON object.
	EMPTY_OBJECT string = "{}"
	// The default content-type. It is always "application/json".
	DEFAULT_CONTENT_TYPE string = "application/json"
	// The default response. It is a pointer to an instance of DefaultResponse.
	DEFAULT_RESPONSE = new(DefaultResponse)
)

// This is the interface for all responses.
type Response interface {
	// Checks if the given request matches with this response based on the
	// method and path.
	Match(method string, path string) bool
	// Returns the return code.
	ResponseCode() int
	// Returns the content type.
	ContentType() string
	// Writes the body of the response to a writer.
	WriteBody(writer io.Writer) error
	// If true, prevents the request from being captured.
	SkipCapture() bool
}

//------------------------------------------------------------------------------

// This type implements the response interface. It will always match a request and
// will reply with 200 and an empty JSON object.
type DefaultResponse struct{}

// Always return true.
func (r *DefaultResponse) Match(method string, path string) bool {
	return true
}

// Always return 200.
func (r *DefaultResponse) ResponseCode() int {
	return 200
}

// Always return DEFAULT_CONTENT_TYPE.
func (r *DefaultResponse) ContentType() string {
	return DEFAULT_CONTENT_TYPE
}

func (r *DefaultResponse) WriteBody(writer io.Writer) error {
	_, err := writer.Write([]byte(EMPTY_OBJECT))
	return err
}

func (r *DefaultResponse) SkipCapture() bool {
	return false
}

//------------------------------------------------------------------------------

/*
Implementation of a Response that uses regular expressions to decide if it
matches with the given request.
*/
type responseImpl struct {
	pathPattern  *regexp.Regexp
	methods      map[string]bool
	responseCode int
	contentType  string
	body         []byte
	skipCapture  bool
}

// Creates a new responseImpl and initializes some fields with default values.
// methods will be allocated but empty and responseCode will be set to 200.
func newResponseImpl() *responseImpl {
	return &responseImpl{
		pathPattern:  nil,
		methods:      make(map[string]bool),
		responseCode: 200,
		contentType:  "",
		body:         nil,
	}
}

// Checks if the given path matches this response.
func (r *responseImpl) MatchPath(path string) bool {
	if r.pathPattern == nil {
		return true
	}
	return r.pathPattern.Match([]byte(path))
}

func (r *responseImpl) MatchMethods(method string) bool {
	if len(r.methods) == 0 {
		return true
	}
	_, found := r.methods[method]
	return found
}

func (r *responseImpl) Match(method string, path string) bool {
	return r.MatchMethods(method) && r.MatchPath(path)
}

func (r *responseImpl) ResponseCode() int {
	return r.responseCode
}

func (r *responseImpl) ContentType() string {
	return r.contentType
}

func (r *responseImpl) WriteBody(writer io.Writer) error {
	if r.body == nil {
		return nil
	}
	_, err := writer.Write(r.body)
	return err
}

func (r *responseImpl) SkipCapture() bool {
	return r.skipCapture
}

// ------------------------------------------------------------------------------

// This builder is used to create responses.
type ResponseBuilder struct {
	pathPattern  *regexp.Regexp
	methods      []string
	responseCode int
	contentType  string
	body         []byte
	skipCapture  bool
}

// Sets the path pattern from a regex string.
func (b *ResponseBuilder) SetPathPatternStr(pattern string) error {
	p, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	b.pathPattern = p
	return nil
}

// Sets the path pattern.
//
// It always returns itself.
func (b *ResponseBuilder) SetPathPattern(pathPattern *regexp.Regexp) *ResponseBuilder {
	b.pathPattern = pathPattern
	return b
}

// Adds a method to this builder. If no method is set, the resulting response will
// match all methods.
//
// It always returns itself.
func (b *ResponseBuilder) AddMethod(method ...string) *ResponseBuilder {
	b.methods = append(b.methods, method...)
	return b
}

// Sets the response code. If not set, it defaults to 200.
//
// It always returns itself.
func (b *ResponseBuilder) SetResponseCode(code int) *ResponseBuilder {
	b.responseCode = code
	return b
}

// Sets the content type. If not set, defaults to "".
//
// It always returns itself.
func (b *ResponseBuilder) SetContentType(contentType string) *ResponseBuilder {
	b.contentType = contentType
	return b
}

// Sets the body of the response. If not set, defaults to no body. This version
// clones the provided body to prevent further changes in its contents.
//
// It always returns itself.
func (b *ResponseBuilder) SetBody(body []byte) *ResponseBuilder {
	b.body = append([]byte(nil), body...)
	return b
}

// Sets the body of the response. If not set, defaults to no body. This version
// claims ownership of the provided body.
//
// It always returns itself.
func (b *ResponseBuilder) SetBodyNoClone(body []byte) *ResponseBuilder {
	b.body = body
	return b
}

// Builds a new response based on the current builder state.
func (b *ResponseBuilder) Build() Response {
	r := newResponseImpl()

	if b.pathPattern != nil {
		r.pathPattern = b.pathPattern
	}
	for _, m := range b.methods {
		r.methods[m] = true
	}
	if b.responseCode > 0 {
		r.responseCode = b.responseCode
	}
	r.contentType = b.contentType
	if b.body != nil {
		r.body = append([]byte(nil), b.body...)
	}
	r.skipCapture = b.skipCapture
	return r
}

// Sets the skip capture. Defaults to false.
//
// It always returns itself.
func (b *ResponseBuilder) SkipCapture(skipCapture bool) *ResponseBuilder {
	b.skipCapture = skipCapture
	return b
}

//------------------------------------------------------------------------------

// This struct implements a response set. It stores a list of responses and implements
// the response matching mechanism.
type ResponseSet struct {
	responses []Response
}

// Adds a response to this list. The first added
func (s *ResponseSet) AddResponse(response Response) {
	s.responses = append(s.responses, response)
}

// Finds a response that matches the request. If no registered response matches it
// returns DEFAULT_RESPONSE.
func (s *ResponseSet) Find(method string, path string) Response {
	for _, r := range s.responses {
		if r.Match(method, path) {
			return r
		}
	}
	return DEFAULT_RESPONSE
}

//------------------------------------------------------------------------------

// Writes a response to a ResponseWriter.
func WriteResponse(resp Response, response http.ResponseWriter) error {
	if resp.ContentType() != "" {
		response.Header().Add("Content-Type", resp.ContentType())
	}
	response.WriteHeader(resp.ResponseCode())
	return resp.WriteBody(response)
}

// Creates a new response from the configuration.
func NewResponseFromConfig(config *config.ResponseConfig) (Response, error) {
	b := ResponseBuilder{}
	if config.PathPattern != "" {
		p, err := regexp.Compile(config.PathPattern)
		if err != nil {
			return nil, err
		}
		b.SetPathPattern(p)
	}
	if len(config.Methods) > 0 {
		b.AddMethod(config.Methods...)
	}
	b.SetContentType(config.ContentType)
	if config.Body != "" {
		body, err := base64.StdEncoding.DecodeString(config.Body)
		if err != nil {
			return nil, err
		}
		b.SetBody(body)
	}
	b.SkipCapture(config.SkipCapture)
	if config.ReturnCode != 0 {
		b.SetResponseCode(config.ReturnCode)
	}
	return b.Build(), nil
}
