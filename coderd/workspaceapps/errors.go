package workspaceapps

import (
	"context"
	"net/http"
	"net/url"

	"cdr.dev/slog"
	"github.com/coder/coder/site"
)

// WriteWorkspaceApp404 writes a HTML 404 error page for a workspace app. If
// appReq is not nil, it will be used to log the request details at debug level.
func WriteWorkspaceApp404(ctx context.Context, log slog.Logger, accessURL *url.URL, rw http.ResponseWriter, appReq *Request, msg string) {
	if appReq != nil {
		slog.Helper()
		log.Debug(ctx,
			"workspace app 404: "+msg,
			slog.F("username_or_id", appReq.UsernameOrID),
			slog.F("workspace_and_agent", appReq.WorkspaceAndAgent),
			slog.F("workspace_name_or_id", appReq.WorkspaceNameOrID),
			slog.F("agent_name_or_id", appReq.AgentNameOrID),
			slog.F("app_slug_or_port", appReq.AppSlugOrPort),
		)
	}

	site.RenderStaticErrorPage(ctx, rw, site.ErrorPageData{
		Status:       http.StatusNotFound,
		Title:        "Application Not Found",
		Description:  "The application or workspace you are trying to access does not exist or you do not have permission to access it.",
		RetryEnabled: false,
		DashboardURL: accessURL.String(),
	})
}

// WriteWorkspaceApp500 writes a HTML 500 error page for a workspace app. If
// appReq is not nil, it's fields will be added to the logged error message.
func WriteWorkspaceApp500(ctx context.Context, log slog.Logger, accessURL *url.URL, rw http.ResponseWriter, appReq *Request, err error, msg string) {
	if appReq != nil {
		slog.Helper()
		ctx = slog.With(ctx,
			slog.F("username_or_id", appReq.UsernameOrID),
			slog.F("workspace_and_agent", appReq.WorkspaceAndAgent),
			slog.F("workspace_name_or_id", appReq.WorkspaceNameOrID),
			slog.F("agent_name_or_id", appReq.AgentNameOrID),
			slog.F("app_name_or_port", appReq.AppSlugOrPort),
		)
	}
	log.Warn(ctx,
		"workspace app auth server error: "+msg,
		slog.Error(err),
	)

	site.RenderStaticErrorPage(ctx, rw, site.ErrorPageData{
		Status:       http.StatusInternalServerError,
		Title:        "Internal Server Error",
		Description:  "An internal server error occurred.",
		RetryEnabled: false,
		DashboardURL: accessURL.String(),
	})
}

// WriteWorkspaceAppOffline writes a HTML 502 error page for a workspace app. If
// appReq is not nil, it will be used to log the request details at debug level.
func WriteWorkspaceAppOffline(ctx context.Context, log slog.Logger, accessURL *url.URL, rw http.ResponseWriter, appReq *Request, msg string) {
	if appReq != nil {
		slog.Helper()
		log.Debug(ctx,
			"workspace app unavailable: "+msg,
			slog.F("username_or_id", appReq.UsernameOrID),
			slog.F("workspace_and_agent", appReq.WorkspaceAndAgent),
			slog.F("workspace_name_or_id", appReq.WorkspaceNameOrID),
			slog.F("agent_name_or_id", appReq.AgentNameOrID),
			slog.F("app_slug_or_port", appReq.AppSlugOrPort),
		)
	}

	site.RenderStaticErrorPage(ctx, rw, site.ErrorPageData{
		Status:       http.StatusBadGateway,
		Title:        "Application Unavailable",
		Description:  msg,
		RetryEnabled: true,
		DashboardURL: accessURL.String(),
	})
}
