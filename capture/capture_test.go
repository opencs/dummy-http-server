// Copyright (c) 2023-2024, Open Communications Security
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
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

package capture

import (
	"bytes"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromRequest(t *testing.T) {

	r := httptest.NewRequest("PUT", "http://host1/path1", bytes.NewReader([]byte("12345")))
	r.Header["a"] = []string{"b"}

	c, err := NewFromRequest(r, 10000)
	assert.Nil(t, err)
	assert.Equal(t, "PUT", c.Method)
	assert.Equal(t, []string{"b"}, c.Headers["a"])
	assert.Equal(t, []byte("12345"), c.Body)
	assert.Greater(t, time.Millisecond, time.Since(c.Timestamp))

	r = httptest.NewRequest("PUT", "http://host1/path1", bytes.NewReader([]byte("12345")))
	r.Header["a"] = []string{"b"}

	c, err = NewFromRequest(r, 2)
	assert.Nil(t, err)
	assert.Equal(t, "PUT", c.Method)
	assert.Equal(t, []string{"b"}, c.Headers["a"])
	assert.Equal(t, []byte("12"), c.Body)
	assert.Greater(t, time.Millisecond, time.Since(c.Timestamp))
}

func TestCapturedRequest_GetFileTitle(t *testing.T) {
	c := CapturedRequest{
		Method:    "M123",
		Timestamp: time.UnixMilli(12344310923).UTC(),
	}
	assert.Equal(t, "1970-05-23T205830.923000000.M123", c.GetFileTitle())
}

func TestCapturedRequest_Save(t *testing.T) {
	r := httptest.NewRequest("PUT", "http://host1/path1", bytes.NewReader([]byte("12345")))
	r.Header["a"] = []string{"b"}

	c, err := NewFromRequest(r, 10000)
	require.Nil(t, err)
	c.Timestamp = time.UnixMilli(12344310923).UTC()

	buff := bytes.NewBuffer(nil)
	err = c.Save(buff)
	assert.Nil(t, err)
	assert.Equal(t, "{\n  \"host\": \"host1\",\n  \"remote\": \"192.0.2.1:1234\",\n  \"url\": \"http://host1/path1\",\n  \"Method\": \"PUT\",\n  \"timestamp\": \"1970-05-23T20:58:30.923Z\",\n  \"headers\": {\n   \"a\": [\n    \"b\"\n   ]\n  },\n  \"body\": \"MTIzNDU=\"\n }",
		buff.String())
}

func TestCapturedRequest_SaveTo(t *testing.T) {

	r := httptest.NewRequest("PUT", "http://host1/path1", bytes.NewReader([]byte("12345")))
	r.Header["a"] = []string{"b"}

	c, err := NewFromRequest(r, 10000)
	require.Nil(t, err)
	c.Timestamp = time.UnixMilli(12344310923).UTC()

	root := os.TempDir()
	err = c.SaveTo(root)
	assert.Nil(t, err)

	actual, err := os.ReadFile(path.Join(root, c.GetFileTitle()))
	assert.Nil(t, err)
	assert.Equal(t, "{\n  \"host\": \"host1\",\n  \"remote\": \"192.0.2.1:1234\",\n  \"url\": \"http://host1/path1\",\n  \"Method\": \"PUT\",\n  \"timestamp\": \"1970-05-23T20:58:30.923Z\",\n  \"headers\": {\n   \"a\": [\n    \"b\"\n   ]\n  },\n  \"body\": \"MTIzNDU=\"\n }",
		string(actual))
}
