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
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"gitlab.opencs.dev.br/opencs-commons/dummy-http-server/capture"
	"gitlab.opencs.dev.br/opencs-commons/dummy-http-server/config"
	"go.uber.org/zap"
)

type Engine struct {
	Responses ResponseSet
	Config    *config.Config
	Logger    *zap.Logger
}

func NewEngine(config *config.Config) (*Engine, error) {
	ret := &Engine{
		Config: config,
	}
	if err := ret.initLogger(); err != nil {
		return nil, err
	}
	if err := ret.initResponses(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *Engine) initLogger() error {
	// Configure
	logFile := path.Join(e.Config.CaptureDir, "log.log")
	cfg := zap.NewProductionConfig()
	cfg.ErrorOutputPaths = []string{logFile}
	cfg.OutputPaths = []string{logFile}
	// Build
	logger, err := cfg.Build()
	if err != nil {
		return nil
	}
	e.Logger = logger
	return nil
}

func (e *Engine) initResponses() error {
	for i, cfg := range e.Config.Responses {
		r, err := NewResponseFromConfig(cfg)
		if err != nil {
			e.Logger.Error("Bad response definition.", zap.Int("index", i), zap.Error(err))
		} else {
			e.Responses.AddResponse(r)
		}
	}
	return nil
}

func (e *Engine) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// Select the response first
	method := request.Method
	path := request.URL.Path
	resp := e.Responses.Find(method, path)

	// Capture the request
	cap, err := capture.NewFromRequest(request, int64(e.Config.MaxRequestSize))
	if err != nil {
		e.Logger.Error("Unable to capture the request.", zap.Error(err))
	} else {
		if !resp.SkipCapture() {
			err := cap.SaveTo(e.Config.CaptureDir)
			if err != nil {
				e.Logger.Error("Unable to save the captured request.", zap.Error(err))
			}
		} else {
			e.Logger.Info("Capture skipped.", zap.String("URL", request.URL.String()),
				zap.String("host", request.Host), zap.String("remote", request.RemoteAddr))
		}
	}

	// Send the response
	err = WriteResponse(resp, response)
	if err != nil {
		e.Logger.Error("Unable to send the response.", zap.Error(err))
	}
}

func (e *Engine) StartServer() error {
	// Configure the server
	srv := &http.Server{
		Addr:           e.Config.Address,
		Handler:        e,
		ReadTimeout:    time.Duration(e.Config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(e.Config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 0,
	}

	// Wait for the kill signal
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		e.Logger.Info("Stopping the server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// Start the server
	e.Logger.Info("Server stopped.")
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		e.Logger.Error("Erro", zap.Error(err))
		close(idleConnsClosed)
	} else {
		err = nil
	}
	<-idleConnsClosed
	e.Logger.Info("Server stopped.")
	return err
}
