package mock

func (s *mockServer) routes() {
	s.router.Get("/", s.handleRoot())
}
