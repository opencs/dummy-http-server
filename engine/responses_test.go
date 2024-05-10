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
	"bytes"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//------------------------------------------------------------------------------

var _ Response = (*DefaultResponse)(nil)

func TestDefaultResponse(t *testing.T) {

	assert.Equal(t, int(200), DEFAULT_RESPONSE.ResponseCode())
	assert.Equal(t, "application/json", DEFAULT_RESPONSE.ContentType())

	assert.True(t, DEFAULT_RESPONSE.Match("", ""))
	assert.True(t, DEFAULT_RESPONSE.Match("12312312 13", "13 1313123"))

	actual := bytes.NewBuffer(nil)
	assert.Nil(t, DEFAULT_RESPONSE.WriteBody(actual))
	assert.Equal(t, "{}", actual.String())
}

//------------------------------------------------------------------------------

var _ Response = (*responseImpl)(nil)

func TestNewResponseImpl(t *testing.T) {
	r := newResponseImpl()

	assert.Nil(t, r.pathPattern)
	assert.NotNil(t, r.methods)
	assert.Len(t, r.methods, 0)
	assert.Equal(t, 200, r.responseCode)
	assert.Equal(t, "", r.contentType)
	assert.Nil(t, r.body)
}

func TestResponseImpl_MatchPath(t *testing.T) {
	r := responseImpl{}

	assert.True(t, r.MatchPath(""))
	assert.True(t, r.MatchPath("whatever"))

	r.pathPattern = regexp.MustCompile("^/a$")
	assert.False(t, r.MatchPath(""))
	assert.False(t, r.MatchPath("/ab"))
	assert.False(t, r.MatchPath("/ba"))
	assert.True(t, r.MatchPath("/a"))

	r.pathPattern = regexp.MustCompile("/a")
	assert.False(t, r.MatchPath(""))
	assert.True(t, r.MatchPath("/ab"))
	assert.False(t, r.MatchPath("/ba"))
	assert.True(t, r.MatchPath("/a"))
}

func TestResponseImpl_MatchMethods(t *testing.T) {
	r := responseImpl{}

	assert.True(t, r.MatchMethods(""))
	assert.True(t, r.MatchMethods("123123"))
	assert.True(t, r.MatchMethods("POST"))

	r.methods = make(map[string]bool)
	r.methods["POST"] = false
	assert.False(t, r.MatchMethods(""))
	assert.False(t, r.MatchMethods("123123"))
	assert.True(t, r.MatchMethods("POST"))
}

func TestResponseImpl_Match(t *testing.T) {
	r := responseImpl{}

	assert.True(t, r.Match("", ""))
	assert.True(t, r.Match("fdasfd", "1 1231231323"))

	r.pathPattern = regexp.MustCompile("^/a$")
	assert.False(t, r.Match("", ""))
	assert.True(t, r.Match("", "/a"))
	assert.True(t, r.Match("fdasfd", "/a"))
	assert.True(t, r.Match("POST", "/a"))

	r.methods = make(map[string]bool)
	r.methods["POST"] = false
	assert.False(t, r.Match("", ""))
	assert.False(t, r.Match("", "/a"))
	assert.False(t, r.Match("fdasfd", "/a"))
	assert.True(t, r.Match("POST", "/a"))
}

func TestResponseImpl_ResponseCode(t *testing.T) {
	r := responseImpl{}

	assert.Equal(t, 0, r.ResponseCode())
	r.responseCode = 123
	assert.Equal(t, 123, r.ResponseCode())
}

func TestResponseImpl_ContentType(t *testing.T) {
	r := responseImpl{}

	assert.Equal(t, "", r.ContentType())
	r.contentType = "123"
	assert.Equal(t, "123", r.ContentType())
}

func TestResponseImpl_WriteBody(t *testing.T) {
	r := responseImpl{}

	b := bytes.NewBuffer(nil)
	assert.Nil(t, r.WriteBody(b))
	assert.Equal(t, "", b.String())

	r.body = []byte("123")
	b = bytes.NewBuffer(nil)
	assert.Nil(t, r.WriteBody(b))
	assert.Equal(t, "123", b.String())
}

// ------------------------------------------------------------------------------

func TestResponseBuilder_SetPathPatternStr(t *testing.T) {
	b := ResponseBuilder{}

	assert.Nil(t, b.SetPathPatternStr("a"))
	assert.Equal(t, "a", b.pathPattern.String())

	b = ResponseBuilder{}
	assert.NotNil(t, b.SetPathPatternStr("["))
	assert.Nil(t, b.pathPattern)
}

func TestResponseBuilder_SetPathPattern(t *testing.T) {
	b := ResponseBuilder{}

	p := regexp.MustCompile("a")
	b2 := b.SetPathPattern(p)
	assert.Same(t, &b, b2)
	assert.Same(t, p, b.pathPattern)

	b2 = b.SetPathPattern(nil)
	assert.Same(t, &b, b2)
	assert.Nil(t, b.pathPattern)
}

func TestResponseBuilder_AddMethod(t *testing.T) {
	b := ResponseBuilder{}

	assert.Nil(t, b.methods)
	b2 := b.AddMethod("A")
	assert.Same(t, &b, b2)
	assert.Equal(t, []string{"A"}, b.methods)

	b2 = b.AddMethod("B")
	assert.Same(t, &b, b2)
	assert.Equal(t, []string{"A", "B"}, b.methods)

	b2 = b.AddMethod("C", "D")
	assert.Same(t, &b, b2)
	assert.Equal(t, []string{"A", "B", "C", "D"}, b.methods)
}

func TestResponseBuilder_SetResponseCode(t *testing.T) {
	b := ResponseBuilder{}

	b2 := b.SetResponseCode(123)
	assert.Same(t, &b, b2)
	assert.Equal(t, 123, b.responseCode)
}

func TestResponseBuilder_SetContentType(t *testing.T) {
	b := ResponseBuilder{}

	b2 := b.SetContentType("123")
	assert.Same(t, &b, b2)
	assert.Equal(t, "123", b.contentType)
}

func TestResponseBuilder_SetBody(t *testing.T) {
	b := ResponseBuilder{}

	exp := []byte("12345")
	b2 := b.SetBody(exp)
	assert.Same(t, &b, b2)
	assert.Equal(t, exp, b.body)
	assert.NotSame(t, &exp[0], &b.body[0])
}

func TestResponseBuilder_Build(t *testing.T) {
	b := ResponseBuilder{}

	r := b.Build()
	imp := r.(*responseImpl)
	require.NotNil(t, imp)
	assert.Nil(t, imp.pathPattern)
	assert.Empty(t, imp.methods)
	assert.Equal(t, 200, imp.responseCode)
	assert.Equal(t, "", imp.contentType)
	assert.Nil(t, imp.body)

	pattern := regexp.MustCompile("a")
	body := []byte("12345")

	b = ResponseBuilder{}
	b.SetPathPattern(pattern).AddMethod("PUT", "POST", "PUT").SetResponseCode(123).SetContentType("type1").SetBody(body)
	r = b.Build()
	imp = r.(*responseImpl)
	require.NotNil(t, imp)
	assert.Same(t, pattern, imp.pathPattern)
	assert.Len(t, imp.methods, 2)
	assert.Contains(t, imp.methods, "PUT")
	assert.Contains(t, imp.methods, "POST")
	assert.Equal(t, 123, imp.responseCode)
	assert.Equal(t, "type1", imp.contentType)
	assert.Equal(t, body, imp.body)
	assert.NotSame(t, &b.body[0], &imp.body[0])

	b = ResponseBuilder{}
	b.SetPathPattern(pattern).AddMethod("PUT", "POST", "PUT").SetResponseCode(-1).SetContentType("type1").SetBody(body)
	r = b.Build()
	imp = r.(*responseImpl)
	require.NotNil(t, imp)
	assert.Same(t, pattern, imp.pathPattern)
	assert.Len(t, imp.methods, 2)
	assert.Contains(t, imp.methods, "PUT")
	assert.Contains(t, imp.methods, "POST")
	assert.Equal(t, 200, imp.responseCode)
	assert.Equal(t, "type1", imp.contentType)
	assert.Equal(t, body, imp.body)
	assert.NotSame(t, &b.body[0], &imp.body[0])
}

//------------------------------------------------------------------------------

func TestResponseSet_AddResponse(t *testing.T) {
	s := ResponseSet{}

	r1 := new(DefaultResponse)
	r2 := new(DefaultResponse)

	s.AddResponse(r1)
	assert.Equal(t, []Response{r1}, s.responses)

	s.AddResponse(r2)
	assert.Equal(t, []Response{r1, r2}, s.responses)
}

func TestResponseSet_Find(t *testing.T) {
	s := ResponseSet{}

	r := s.Find("", "")
	assert.Same(t, DEFAULT_RESPONSE, r)

	b := ResponseBuilder{}
	b.SetPathPattern(regexp.MustCompile("^/a$")).AddMethod("PUT", "POST", "PUT")
	r1 := b.Build()

	b = ResponseBuilder{}
	b.SetPathPattern(regexp.MustCompile("^/b$"))
	r2 := b.Build()

	b = ResponseBuilder{}
	b.SetPathPattern(regexp.MustCompile("^/a$")).AddMethod("WHATEVER")
	r3 := b.Build()

	s.AddResponse(r1)
	s.AddResponse(r2)
	s.AddResponse(r3)

	r = s.Find("", "")
	assert.Same(t, DEFAULT_RESPONSE, r)

	r = s.Find("", "/a")
	assert.Same(t, DEFAULT_RESPONSE, r)

	r = s.Find("PUT", "/a")
	assert.Same(t, r1, r)

	r = s.Find("PUT", "/b")
	assert.Same(t, r2, r)

	r = s.Find("WHATEVER", "/a")
	assert.Same(t, r3, r)
}

//------------------------------------------------------------------------------

func TestWriteResponse(t *testing.T) {

	resp := httptest.NewRecorder()
	assert.Nil(t, WriteResponse(DEFAULT_RESPONSE, resp))
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("Content-Type"))
	assert.Equal(t, "{}", resp.Body.String())

	b := ResponseBuilder{}
	b.SetResponseCode(123)
	r := b.Build()
	resp = httptest.NewRecorder()
	assert.Nil(t, WriteResponse(r, resp))
	assert.Equal(t, 123, resp.Code)
	assert.Equal(t, "", resp.Header().Get("Content-Type"))
	assert.Equal(t, "", resp.Body.String())
}
