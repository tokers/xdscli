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
	// "fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	_gFlags = &globalFlags{}

	_rootCmd = &cobra.Command{
		Use:          "xdscli [options] <xds>",
		Short:        "xDS protocol client",
		Long:         "xDS protocol client to talk with management servers like Istio Pilot",
		SilenceUsage: true,
		Run:          rootCommandFunc,
	}
)

const (
	_defaultDialTimeout = 2 * time.Second
	_defaultReadTimeout = 5 * time.Second
	_defaultSendTimeout = 5 * time.Second
	_defaultTimeout     = 5 * time.Second
)

func init() {
	_rootCmd.PersistentFlags().BoolVar(&_gFlags.showVersion, "version", false, "show the version of xdscli")
	_rootCmd.PersistentFlags().StringSliceVar(&_gFlags.servers, "servers", nil, "xDS server addresses")
	_rootCmd.PersistentFlags().StringVar(&_gFlags.outputFormat, "write-out", "simple", "set the output format (json, yaml, protobuf)")
	_rootCmd.PersistentFlags().DurationVar(&_gFlags.dialTimeout, "dial-timeout", _defaultDialTimeout, "dial timeout for client connections")
	_rootCmd.PersistentFlags().DurationVar(&_gFlags.readTimeout, "read-timeout", _defaultReadTimeout, "read timeout for client connections")
	_rootCmd.PersistentFlags().DurationVar(&_gFlags.sendTimeout, "send-timeout", _defaultSendTimeout, "send timeout for client connections")
	_rootCmd.PersistentFlags().DurationVar(&_gFlags.timeout, "timeout", _defaultTimeout, "timeout for client connections ( dailing, reading and sending)")

	_rootCmd.PersistentFlags().StringVar(&_gFlags.xds.node, "node", "", "the node making the request")
	_rootCmd.PersistentFlags().StringVar(&_gFlags.xds.initialVersionInfo, "initial-version-info", "", "the version_info received with the most recent successfully processed response")
	_rootCmd.PersistentFlags().StringVar(&_gFlags.xds.errorDetail, "error-detail", "", "the error reason that update configuration cannot be applied, using non-empty string means the discovery response will be rejected by xdscli")
	_rootCmd.PersistentFlags().StringSliceVar(&_gFlags.xds.resourceNames, "resource-names", nil, "list of resources to subscribe to")
	_rootCmd.PersistentFlags().StringVar(&_gFlags.xds.apiVersion, "api-version", "v2", "version of xDS protocol")
	_rootCmd.PersistentFlags().StringVar(&_gFlags.xds.nodeMetadata, "node-metadata", "", "comma splitted key value pairs reresent node metadata")

	cobra.EnablePrefixMatching = true
}

func rootCommandFunc(cmd *cobra.Command, args []string) {
	if _gFlags.showVersion {
		showVersionAndQuit()
	}

	if len(args) != 1 {
		exitWithError(_exitBadArgs, errors.New("need exactly one argument as the discovery service type (like eds, cds and etc)."))
	}

	if err := validateAPIVersion(_gFlags.xds.apiVersion); err != nil {
		exitWithError(_exitBadArgs, err)
	}

	typeUrl, err := getDiscoveryServiceTypeUrl(_gFlags.xds.apiVersion, args[0])
	if err != nil {
		exitWithError(_exitBadArgs, err)
	}

	if len(_gFlags.servers) == 0 {
		exitWithError(_exitBadArgs, _errNoServers)
	}

	endpoints, err := validateAndResolveServers(_gFlags.servers)
	if err != nil {
		exitWithError(_exitError, err)
	}

	doDiscoveryService(endpoints, typeUrl)
}

func main() {
	// _rootCmd.SetUsageFunc(usageFunc)
	if err := _rootCmd.Execute(); err != nil {
		exitWithError(_exitSuccess, err)
	}
}
