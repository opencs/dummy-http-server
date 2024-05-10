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
package config

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	file := path.Join("..", "_samples", "config-empty.yaml")
	c, err := LoadConfig(file)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, ":8080", c.Address)
	assert.Equal(t, "var", c.CaptureDir)
	assert.Equal(t, 15, c.ReadTimeout)
	assert.Equal(t, 15, c.WriteTimeout)
	assert.Equal(t, 1024*1024, c.MaxRequestSize)
	assert.Nil(t, c.Responses)

	file = path.Join("..", "_samples", "config-simple.yaml")
	c, err = LoadConfig(file)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "localhost:8080", c.Address)
	assert.Equal(t, "capture2", c.CaptureDir)
	assert.Equal(t, 15, c.ReadTimeout)
	assert.Equal(t, 15, c.WriteTimeout)
	assert.Equal(t, 1024*1024, c.MaxRequestSize)
	assert.Nil(t, c.Responses)

	file = path.Join("..", "_samples", "config-full.yaml")
	c, err = LoadConfig(file)
	assert.Nil(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "localhost2:8080", c.Address)
	assert.Equal(t, "capture2", c.CaptureDir)
	assert.Equal(t, 123, c.ReadTimeout)
	assert.Equal(t, 456, c.WriteTimeout)
	assert.Equal(t, 789, c.MaxRequestSize)
	assert.Len(t, c.Responses, 2)

	assert.Equal(t, "\\/b.*", c.Responses[0].PathPattern)
	assert.Equal(t, []string{"GET", "POST"}, c.Responses[0].Methods)
	assert.Equal(t, "text/plain", c.Responses[0].ContentType)
	assert.Equal(t, "AAAA", c.Responses[0].Body)
	assert.True(t, c.Responses[0].SkipCapture)
	assert.Equal(t, 201, c.Responses[0].ReturnCode)

	assert.Equal(t, "\\/a.*", c.Responses[1].PathPattern)
	assert.Nil(t, c.Responses[1].Methods)
	assert.Equal(t, "text/html", c.Responses[1].ContentType)
	assert.Equal(t, "BBBB", c.Responses[1].Body)
	assert.False(t, c.Responses[1].SkipCapture)
	assert.Equal(t, 0, c.Responses[1].ReturnCode)
}
