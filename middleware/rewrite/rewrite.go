package rewrite

import (
	"net/http"

	config "github.com/go-kratos/gateway/api/gateway/config/v1"
	v1 "github.com/go-kratos/gateway/api/gateway/middleware/rewrite/v1"

	"github.com/go-kratos/gateway/middleware"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	middleware.Register("rewrite", Middleware)
}

func Middleware(c *config.Middleware) (middleware.Middleware, error) {
	options := &v1.Rewrite{}
	if c.Options != nil {
		if err := anypb.UnmarshalTo(c.Options, options, proto.UnmarshalOptions{Merge: true}); err != nil {
			return nil, err
		}
	}
	requestHeadersRewrite := options.RequestHeadersRewrite
	respondHeadersRewrite := options.ReponseHeadersRewrite
	return func(next http.RoundTripper) http.RoundTripper {
		return middleware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if options.HostRewrite != nil {
				req.Host = *options.HostRewrite
			}
			if options.PathRewrite != nil {
				req.URL.Path = *options.PathRewrite
			}
			if requestHeadersRewrite != nil {
				for key, value := range options.GetRequestHeadersRewrite().Set {
					req.Header.Set(key, value)
				}
				for key, value := range options.GetRequestHeadersRewrite().Add {
					if req.Header.Get(key) == "" {
						req.Header.Add(key, value)
					} else {
						req.Header.Set(key, value)
					}
				}
				for _, value := range options.GetRequestHeadersRewrite().Remove {
					if req.Header.Get(value) != "" {
						req.Header.Del(value)
					}
				}
			}
			resp, err := next.RoundTrip(req)
			if err != nil {
				return nil, err
			}

			if respondHeadersRewrite != nil {
				for key, value := range options.GetRequestHeadersRewrite().Set {
					resp.Header.Set(key, value)
				}
				for key, value := range options.GetRequestHeadersRewrite().Add {
					if resp.Header.Get(key) == "" {
						req.Header.Add(key, value)
					} else {
						resp.Header.Set(key, value)
					}
				}
				for _, value := range options.GetRequestHeadersRewrite().Remove {
					if resp.Header.Get(value) != "" {
						resp.Header.Del(value)
					}
				}
			}
			return resp, nil
		})
	}, nil
}
