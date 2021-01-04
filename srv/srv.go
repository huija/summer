package srv

import (
	"context"
	"github.com/huija/summer/utils"
	"net/http"
)

const (
	GIN = "gin"
)

const (
	defaultFramework = GIN
	defaultHost      = "localhost"
	defaultListen    = "127.0.0.1"
	defaultPort      = "9090"
	defaultBaseURL   = "/"
)

// Srv config
type Srv struct {
	// non-zero attr
	Name      string `json:",omitempty"`
	Tag       string `json:",omitempty"`
	Author    string `json:",omitempty"`
	Host      string `json:",omitempty"`
	IP        string `json:",omitempty"`
	FrameWork string `json:",omitempty"`
	Listen    string `json:",omitempty"`
	Port      string `json:",omitempty"`
	BaseURL   string `json:",omitempty"`

	// zero val is ok
	Release  bool
	Trace    bool
	Pprof    bool
	Cors     bool
	Swag     bool
	Static   string
	HtmlGlob string
}

// Handler http handler
var Handler http.Handler

// Server http server
var Server http.Server

// GetTraceId customize func to get request trace id
var GetTraceId func(c context.Context, wrap ...string) string

// Defaults srv
func Defaults(srv *Srv) (*Srv, error) {
	if srv == nil {
		return defaultsSrv(), nil
	}
	err := utils.MergeStructByMarshal(srv, defaultsSrv())
	return srv, err
}

func defaultsSrv() *Srv {
	return &Srv{
		Name:      "summer",
		Tag:       "v0.0.1",
		Author:    "huija",
		Host:      defaultHost,
		IP:        utils.GetOutboundIP(),
		FrameWork: defaultFramework,
		Listen:    defaultListen,
		Port:      defaultPort,
		BaseURL:   defaultBaseURL,
	}
}
