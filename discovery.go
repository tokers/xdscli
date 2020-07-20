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
	gcontext "context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	apiv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	discoveryv2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
)

var (
	_xdsUserAgentName = "xdscli/" + _version
)

type mediateSuite struct {
	errc  chan error
	ackc  chan string
	stopc chan struct{}
	respc chan *apiv2.DiscoveryResponse
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func newGRPCConn(ctx *context) (*grpc.ClientConn, error) {
	dialCtx, dialCancel := gcontext.WithTimeout(ctx.rootCtx, ctx.flags.dialTimeout)
	defer dialCancel()

	dialOpts := grpc.WithContextDialer(
		func(ctx gcontext.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		},
	)

	kp := keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             2 * time.Second,
		PermitWithoutStream: true,
	}

	addr := ctx.endpoints[rand.Intn(len(ctx.endpoints))]
	// TODO TLS support
	conn, err := grpc.DialContext(dialCtx, addr,
		grpc.WithInsecure(),
		dialOpts,
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(kp),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ctx.flags.grpcMaxCallRecvSize)))

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func doDiscoveryService(ctx *context) error {
	conn, err := newGRPCConn(ctx)
	if err != nil {
		return err
	}

	adsClient, err := discoveryv2.NewAggregatedDiscoveryServiceClient(conn).StreamAggregatedResources(ctx.rootCtx)
	if err != nil {
		return err
	}

	suite := &mediateSuite{
		errc:  make(chan error, 1),
		stopc: make(chan struct{}),
		ackc:  make(chan string, 1),
		respc: make(chan *apiv2.DiscoveryResponse, 1),
	}

	ctx.wg.Add(2)

	go receiveThread(ctx, adsClient, suite)
	go sendThread(ctx, adsClient, suite)

	finalize := func() {
		close(suite.stopc)
		ctx.rootCancel()
		conn.Close()
		ctx.wg.Wait()
	}

	for {
		select {
		case <-ctx.interc:
			finalize()
			return nil
		case err := <-suite.errc:
			finalize()
			return err
		case resp := <-suite.respc:
			data, err := ctx.marshaller.marshal(resp)
			if err != nil {
				panic(err)
			}
			fmt.Println(data)
			if !ctx.flags.watch {
				finalize()
				return nil
			}
		}
	}

	return nil
}

func receiveThread(ctx *context, adsClient discoveryv2.AggregatedDiscoveryService_StreamAggregatedResourcesClient, suite *mediateSuite) {
	defer ctx.wg.Done()

	resp, err := adsClient.Recv()
	if err != nil {
		suite.errc <- err
		select {
		case <-suite.stopc:
			return
		}
	}

	nonce := resp.Nonce
	suite.ackc <- nonce
	suite.respc <- resp

	if !ctx.flags.watch {
		return
	}

	for {
		resp, err := adsClient.Recv()
		if err != nil {
			suite.errc <- err
			select {
			case <-suite.stopc:
				return
			}
		}

		nonce := resp.Nonce
		suite.ackc <- nonce
		suite.respc <- resp
	}
}

func sendThread(ctx *context, adsClient discoveryv2.AggregatedDiscoveryService_StreamAggregatedResourcesClient, suite *mediateSuite) {
	defer ctx.wg.Done()

	node := makeNode(ctx)
	// TODO Get ResourceName by spawning another CDS request when type url is
	// EDS and ResourceName is empty.
	discReq := makeDiscoveryRequest(ctx, node, ctx.flags.xds.resourceNames, "")
	if err := adsClient.Send(discReq); err != nil {
		suite.errc <- err
		select {
		case <-suite.stopc:
			return
		}
	}

	for {
		select {
		case <-suite.stopc:
			return
		case nonce := <-suite.ackc:
			// Send the ack.
			discReq = makeDiscoveryRequest(ctx, makeNode(ctx), nil, nonce)
			if err := adsClient.Send(discReq); err != nil {
				suite.errc <- err
				select {
				case <-suite.stopc:
					return
				}
			}
			if !ctx.flags.watch {
				return
			}
		}
	}
}

func makeDiscoveryRequest(ctx *context, node *core.Node, resourceNames []string, nonce string) *apiv2.DiscoveryRequest {
	discReq := &apiv2.DiscoveryRequest{
		VersionInfo:   ctx.flags.xds.initialVersionInfo,
		Node:          node,
		ResourceNames: resourceNames,
		TypeUrl:       ctx.typeUrl,
		// ErrorDetail: ctx.flags.xds.errorDetail,
		ResponseNonce: nonce,
	}
	return discReq
}

func makeNode(ctx *context) *core.Node {
	node := &core.Node{
		Id:            ctx.flags.xds.node,
		Metadata:      ctx.nodeMeta,
		UserAgentName: _xdsUserAgentName,
	}
	return node
}
