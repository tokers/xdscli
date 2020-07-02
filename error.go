// Copyright 2020 xdscli Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	// http://tldp.org/LDP/abs/html/exitcodes.html
	_exitSuccess = iota
	_exitError

	_exitBadArgs = 128
)

var (
	_errNoServers                  = errors.New("no servers")
	_errInvalidDialTimeout         = errors.New("invalid --dial-timeout value")
	_errInvalidReadTimeout         = errors.New("invalid --read-timeout value")
	_errInvalidSendTimeout         = errors.New("invalid --send-timeout value")
	_errInvalidOutputFormat        = errors.New("invalid --write-out value")
	_errInvalidNode                = errors.New("invalid --node value")
	_errInvalidNodeMetaFormat      = errors.New("invalid --node-metadata value")
	_errInvalidGRPCMaxCallRecvSize = errors.New("invalid --grpc-max-call-recv-size")
	_errUnknownTypeUrl             = errors.New("server sent unknown resource type url")
)

func exitWithError(code int, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	os.Exit(code)
}
