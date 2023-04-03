package workspaceproxy

import (
	"context"
	"net"
	"net/http"

	"github.com/coder/coder/agent"
	"github.com/coder/coder/coderd/httpapi"
	"github.com/coder/coder/coderd/workspaceapps"
	"github.com/coder/coder/codersdk"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

func (proxy *Proxy) WorkspaceAgentPTY(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	proxy.WebsocketWaitMutex.Lock()
	proxy.WebsocketWaitGroup.Add(1)
	proxy.WebsocketWaitMutex.Unlock()
	defer proxy.WebsocketWaitGroup.Done()

	ticket, ok := proxy.WorkspaceAppsProvider.ResolveRequest(rw, r, workspaceapps.Request{
		AccessMethod:  workspaceapps.AccessMethodTerminal,
		BasePath:      r.URL.Path,
		AgentNameOrID: chi.URLParam(r, "workspaceagent"),
	})
	if !ok {
		return
	}

	values := r.URL.Query()
	parser := httpapi.NewQueryParamParser()
	reconnect := parser.UUID(values, uuid.New(), "reconnect")
	height := parser.UInt(values, 80, "height")
	width := parser.UInt(values, 80, "width")
	if len(parser.Errors) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, codersdk.Response{
			Message:     "Invalid query parameters.",
			Validations: parser.Errors,
		})
		return
	}

	conn, err := websocket.Accept(rw, r, &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionDisabled,
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, codersdk.Response{
			Message: "Failed to accept websocket.",
			Detail:  err.Error(),
		})
		return
	}

	ctx, wsNetConn := websocketNetConn(ctx, conn, websocket.MessageBinary)
	defer wsNetConn.Close() // Also closes conn.

	go httpapi.Heartbeat(ctx, conn)

	agentConn, release, err := proxy.WorkspaceAgentCache.Acquire(ticket.AgentID)
	if err != nil {
		_ = conn.Close(websocket.StatusInternalError, httpapi.WebsocketCloseSprintf("dial workspace agent: %s", err))
		return
	}
	defer release()
	ptNetConn, err := agentConn.ReconnectingPTY(ctx, reconnect, uint16(height), uint16(width), r.URL.Query().Get("command"))
	if err != nil {
		_ = conn.Close(websocket.StatusInternalError, httpapi.WebsocketCloseSprintf("dial: %s", err))
		return
	}
	defer ptNetConn.Close()
	agent.Bicopy(ctx, wsNetConn, ptNetConn)
}

// wsNetConn wraps net.Conn created by websocket.NetConn(). Cancel func
// is called if a read or write error is encountered.
type wsNetConn struct {
	cancel context.CancelFunc
	net.Conn
}

func (c *wsNetConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if err != nil {
		c.cancel()
	}
	return n, err
}

func (c *wsNetConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if err != nil {
		c.cancel()
	}
	return n, err
}

func (c *wsNetConn) Close() error {
	defer c.cancel()
	return c.Conn.Close()
}

// websocketNetConn wraps websocket.NetConn and returns a context that
// is tied to the parent context and the lifetime of the conn. Any error
// during read or write will cancel the context, but not close the
// conn. Close should be called to release context resources.
func websocketNetConn(ctx context.Context, conn *websocket.Conn, msgType websocket.MessageType) (context.Context, net.Conn) {
	ctx, cancel := context.WithCancel(ctx)
	nc := websocket.NetConn(ctx, conn, msgType)
	return ctx, &wsNetConn{
		cancel: cancel,
		Conn:   nc,
	}
}
