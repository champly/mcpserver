package server

type Handler interface {
	Start(stop <-chan struct{})
}

type Server struct {
	handlers []Handler
}

func New(opts ...Option) *Server {

	op := defaultOption()
	for _, f := range opts {
		f(op)
	}

	s := &Server{
		handlers: []Handler{},
	}
	s.handlers = append(s.handlers, newMCPServer(op))

	return s
}

func (s *Server) Start(stop <-chan struct{}) {
	for _, h := range s.handlers {
		go h.Start(stop)
	}
	<-stop
}
