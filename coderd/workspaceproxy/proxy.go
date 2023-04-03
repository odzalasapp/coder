package workspaceproxy

import (
	"sync"

	"github.com/coder/coder/coderd/wsconncache"

	"github.com/coder/coder/coderd/workspaceapps"
)

type Options struct {
	WebsocketWaitMutex    *sync.Mutex
	WebsocketWaitGroup    *sync.WaitGroup
	WorkspaceAppsProvider *workspaceapps.Provider
	WorkspaceAgentCache   *wsconncache.Cache
}

type Proxy struct {
	WorkspaceAppsProvider *workspaceapps.Provider
	WorkspaceAgentCache   *wsconncache.Cache

	// TODO: This is a dumb way to pass these imo
	WebsocketWaitMutex *sync.Mutex
	WebsocketWaitGroup *sync.WaitGroup
}

func New(opts *Options) *Proxy {
	if opts == nil {
		opts = &Options{}
	}
	if opts.WebsocketWaitMutex == nil {
		opts.WebsocketWaitMutex = &sync.Mutex{}
	}
	if opts.WebsocketWaitGroup == nil {
		opts.WebsocketWaitGroup = &sync.WaitGroup{}
	}

	return &Proxy{
		WorkspaceAppsProvider: opts.WorkspaceAppsProvider,
		WorkspaceAgentCache:   opts.WorkspaceAgentCache,
		WebsocketWaitMutex:    opts.WebsocketWaitMutex,
		WebsocketWaitGroup:    opts.WebsocketWaitGroup,
	}
}
