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
package capture

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

type CapturedRequest struct {
	Host      string              `json:"host,omitempty"`
	Remote    string              `json:"remote,omitempty"`
	URL       string              `json:"url"`
	Method    string              `json:"Method"`
	Timestamp time.Time           `json:"timestamp"`
	Headers   map[string][]string `json:"headers"`
	Body      []byte              `json:"body,omitempty"`
}

/*
Creates a new CapturedRequest from the given request.
*/
func NewFromRequest(request *http.Request, maxBody int64) (CapturedRequest, error) {
	// Read the body
	body, err := io.ReadAll(io.LimitReader(request.Body, maxBody))
	if err != nil {
		return CapturedRequest{}, nil
	}
	headers := make(map[string][]string)
	for k, v := range request.Header {
		headers[k] = v
	}
	return CapturedRequest{
		Host:      request.Host,
		Remote:    request.RemoteAddr,
		URL:       request.URL.String(),
		Method:    request.Method,
		Timestamp: time.Now().UTC(),
		Headers:   headers,
		Body:      body,
	}, nil
}

/*
Returns the file title for this instance.
*/
func (r *CapturedRequest) GetFileTitle() string {
	ts := r.Timestamp.UTC().Format("2006-01-02T150405")
	ms := fmt.Sprintf("%09d", r.Timestamp.Nanosecond())
	return ts + "." + ms + "." + r.Method
}

/*
Saves this request into a file.
*/
func (r *CapturedRequest) Save(writer io.Writer) error {
	// Convert to JSON
	data, err := json.MarshalIndent(r, " ", " ")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

/*
Saves this request into a file.
*/
func (r *CapturedRequest) SaveTo(parentDir string) error {
	// Save to file.
	writer, err := os.Create(path.Join(parentDir, r.GetFileTitle()))
	if err != nil {
		return err
	}
	err = r.Save(writer)
	if err != nil {
		writer.Close()
		return err
	}
	return writer.Close()
}
