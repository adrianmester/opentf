// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/placeholderplaceholderplaceholder/opentf/internal/opentf"
)

// UIInput is an implementation of terraform.UIInput that communicates
// over RPC.
type UIInput struct {
	Client *rpc.Client
}

func (i *UIInput) Input(ctx context.Context, opts *opentf.InputOpts) (string, error) {
	var resp UIInputInputResponse
	err := i.Client.Call("Plugin.Input", opts, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		err = resp.Error
		return "", err
	}

	return resp.Value, nil
}

type UIInputInputResponse struct {
	Value string
	Error *plugin.BasicError
}

// UIInputServer is a net/rpc compatible structure for serving
// a UIInputServer. This should not be used directly.
type UIInputServer struct {
	UIInput opentf.UIInput
}

func (s *UIInputServer) Input(
	opts *opentf.InputOpts,
	reply *UIInputInputResponse) error {
	value, err := s.UIInput.Input(context.Background(), opts)
	*reply = UIInputInputResponse{
		Value: value,
		Error: plugin.NewBasicError(err),
	}

	return nil
}
