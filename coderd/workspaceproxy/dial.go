package workspaceproxy

import (
	"bufio"
	"context"
	"net"
	"net/netip"

	"tailscale.com/derp"

	"cdr.dev/slog"
	"github.com/coder/coder/codersdk"
	"github.com/coder/coder/tailnet"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"tailscale.com/tailcfg"
)

func Dialer(ctx context.Context, coordinator tailnet.Coordinator, derpServer *derp.Server, derpMap *tailcfg.DERPMap, logger slog.Logger) func(id uuid.UUID) (*codersdk.WorkspaceAgentConn, error) {
	appCtx := ctx
	return func(agentID uuid.UUID) (*codersdk.WorkspaceAgentConn, error) {
		clientConn, serverConn := net.Pipe()
		conn, err := tailnet.NewConn(&tailnet.Options{
			Addresses: []netip.Prefix{netip.PrefixFrom(tailnet.IP(), 128)},
			DERPMap:   derpMap,
			Logger:    logger.Named("tailnet"),
		})
		if err != nil {
			_ = clientConn.Close()
			_ = serverConn.Close()
			return nil, xerrors.Errorf("create tailnet conn: %w", err)
		}
		ctx, cancel := context.WithCancel(appCtx)
		conn.SetDERPRegionDialer(func(_ context.Context, region *tailcfg.DERPRegion) net.Conn {
			if !region.EmbeddedRelay {
				return nil
			}
			left, right := net.Pipe()
			go func() {
				defer left.Close()
				defer right.Close()
				brw := bufio.NewReadWriter(bufio.NewReader(right), bufio.NewWriter(right))
				derpServer.Accept(ctx, right, brw, "internal")
			}()
			return left
		})

		sendNodes, _ := tailnet.ServeCoordinator(clientConn, func(node []*tailnet.Node) error {
			err = conn.UpdateNodes(node, true)
			if err != nil {
				return xerrors.Errorf("update nodes: %w", err)
			}
			return nil
		})
		conn.SetNodeCallback(sendNodes)
		agentConn := &codersdk.WorkspaceAgentConn{
			Conn: conn,
			CloseFunc: func() {
				cancel()
				_ = clientConn.Close()
				_ = serverConn.Close()
			},
		}
		go func() {
			err := coordinator.ServeClient(serverConn, uuid.New(), agentID)
			if err != nil {
				// Sometimes, we get benign closed pipe errors when the server is
				// shutting down.
				if appCtx.Err() == nil {
					logger.Warn(ctx, "tailnet coordinator client error", slog.Error(err))
				}
				_ = agentConn.Close()
			}
		}()
		if !agentConn.AwaitReachable(ctx) {
			_ = agentConn.Close()
			return nil, xerrors.Errorf("agent not reachable")
		}
		return agentConn, nil
	}
}
