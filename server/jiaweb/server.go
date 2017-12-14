package jiaweb

import (
	"jiacrontab/server/view"
	"net/http"
	"sync"
)

type (
	HttpServer struct {
		stdServer *http.Server
		pool      *pool
		modelView *view.ModelView
		end       bool
	}

	pool struct {
		request  sync.Pool
		response sync.Pool
		context  sync.Pool
	}
)

func NewHttpServer() *HttpServer {
	s := &HttpServer{
		end: false,
		pool: &pool{

			context: sync.Pool{
				New: func() interface{} {
					return &HttpContext{}
				},
			},

			request: sync.Pool{
				New: func() interface{} {
					return &Request{}
				},
			},

			response: sync.Pool{
				New: func() interface{} {
					return &Response{}
				},
			},
		},
	}
	s.stdServer = &http.Server{
		Handler: s,
	}

	return s
}

func (s *HttpServer) ListenAndServe(addr string) error {
	s.stdServer.Addr = addr
	return s.stdServer.ListenAndServe()
}

func (s *HttpServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	httpctx := s.pool.context.Get().(*HttpContext)
	request := s.pool.request.Get().(*Request)
	response := s.pool.response.Get().(*Response)

	httpctx.reset(request, response, s)
	request.reset(req, httpctx)
	response.reset(rw)

}