package protocol

import (
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var testHandler = func(_ sdk.Context, _ sdk.Msg) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func TestRouter(t *testing.T) {
	rtr := NewRouter()

	// require panic on invalid route
	require.Panics(t, func() {
		rtr.AddRoute("*", testHandler)
	})

	rtr.AddRoute("testRoute", testHandler)
	h := rtr.Route(sdk.Context{}, "testRoute")
	require.NotNil(t, h)

	// require panic on duplicate route
	require.Panics(t, func() {
		rtr.AddRoute("testRoute", testHandler)
	})
}
