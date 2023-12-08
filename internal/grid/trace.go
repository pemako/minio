// Copyright (c) 2015-2023 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package grid

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio/internal/pubsub"
)

// TraceParamsKey allows to pass trace parameters to the request via context.
// This is only needed when un-typed requests are used.
// MSS, map[string]string types are preferred, but any struct with exported fields will work.
type TraceParamsKey struct{}

// traceRequests adds request tracing to the connection.
func (c *Connection) traceRequests(p *pubsub.PubSub[madmin.TraceInfo, madmin.TraceType]) {
	c.trace = &tracer{
		Publisher: p,
		TraceType: madmin.TraceInternal,
		Prefix:    "grid",
		Local:     c.Local,
		Remote:    c.Remote,
		Subroute:  "",
	}
}

// subroute adds a specific subroute to the request.
func (c *tracer) subroute(subroute string) *tracer {
	if c == nil {
		return nil
	}
	c2 := *c
	c2.Subroute = subroute
	return &c2
}

type tracer struct {
	Publisher *pubsub.PubSub[madmin.TraceInfo, madmin.TraceType]
	TraceType madmin.TraceType
	Prefix    string
	Local     string
	Remote    string
	Subroute  string
}

func (c *muxClient) traceRoundtrip(ctx context.Context, t *tracer, h HandlerID, req []byte) ([]byte, error) {
	if t == nil || t.Publisher.NumSubscribers(t.TraceType) == 0 {
		return c.roundtrip(h, req)
	}
	start := time.Now()
	body := bytesOrLength(req)
	resp, err := c.roundtrip(h, req)
	end := time.Now()
	status := http.StatusOK
	errString := ""
	if err != nil {
		errString = err.Error()
		if IsRemoteErr(err) == nil {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
	}

	prefix := t.Prefix
	if p := handlerPrefixes[h]; p != "" {
		prefix = p
	}
	trace := madmin.TraceInfo{
		TraceType: t.TraceType,
		FuncName:  prefix + "." + h.String(),
		NodeName:  t.Remote,
		Time:      start,
		Duration:  end.Sub(start),
		Path:      t.Subroute,
		Error:     errString,
		HTTP: &madmin.TraceHTTPStats{
			ReqInfo: madmin.TraceRequestInfo{
				Time:    start,
				Proto:   "grid",
				Method:  "REQ",
				Client:  t.Local,
				Headers: nil,
				Path:    t.Subroute,
				Body:    []byte(body),
			},
			RespInfo: madmin.TraceResponseInfo{
				Time:       end,
				Headers:    nil,
				StatusCode: status,
				Body:       []byte(bytesOrLength(resp)),
			},
			CallStats: madmin.TraceCallStats{
				InputBytes:      len(req),
				OutputBytes:     len(resp),
				TimeToFirstByte: end.Sub(start),
			},
		},
	}
	// If the context contains a TraceParamsKey, add it to the trace path.
	v := ctx.Value(TraceParamsKey{})
	if p, ok := v.(*MSS); ok && p != nil {
		trace.Path += p.ToQuery()
		trace.HTTP.ReqInfo.Path = trace.Path
	} else if p, ok := v.(map[string]string); ok {
		m := MSS(p)
		trace.Path += m.ToQuery()
		trace.HTTP.ReqInfo.Path = trace.Path
	} else if v != nil {
		// Print exported fields as single request to path.
		obj := fmt.Sprintf("%+v", v)
		if len(obj) > 1024 {
			obj = obj[:1024] + "..."
		}
		trace.Path = fmt.Sprintf("%s?req=%s", trace.Path, url.QueryEscape(obj))
		trace.HTTP.ReqInfo.Path = trace.Path
	}
	t.Publisher.Publish(trace)
	return resp, err
}
