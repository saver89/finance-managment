package server

func (s *Server) runHttpServer(httpAddress string) error {
	s.mapRoutes()

	if err := s.router.Run(httpAddress); err != nil {
		return err
	}

	return nil
}

func (s *Server) mapRoutes() {
	// fill default routes (healthcheck, metrics, etc)
}
