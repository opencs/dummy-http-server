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

import "github.com/spf13/viper"

func SetDefaults(v *viper.Viper) {

	v.SetDefault("address", ":8080")
	v.SetDefault("captureDir", "var")
	v.SetDefault("readTimeout", 15)
	v.SetDefault("writeTimeout", 15)
	v.SetDefault("maxRequestSize", 1024*1024)
}

type ResponseConfig struct {
	PathPattern string
	Methods     []string
	ContentType string
	Body        string
	SkipCapture bool
	ReturnCode  int
}

type Config struct {
	// Binding address.
	Address string
	// Capture directory.
	CaptureDir string
	// Read timeout in seconds.
	ReadTimeout int
	// Write timeout in seconds.
	WriteTimeout int
	// Maximum request size in bytes.
	MaxRequestSize int
	// Responses
	Responses []*ResponseConfig
	// Source configuration.
	source *viper.Viper
}

func (c *Config) GetSource() *viper.Viper {
	return c.source
}

func LoadConfig(file string) (*Config, error) {
	v := viper.New()
	SetDefaults(v)
	// Load
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	// Unmarshal
	c := new(Config)
	if err := v.Unmarshal(c); err != nil {
		return nil, err
	}
	c.source = v
	return c, nil
}
