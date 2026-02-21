package rest

func (s *echoServer) setupRouter() {
	apiV1 := s.echo.Group("/api/v1")
	{
		productGroup := apiV1.Group("/products")
		{
			productGroup.POST("", s.handler.Product().Create)
			productGroup.GET("", s.handler.Product().List)
			productGroup.GET("/:id", s.handler.Product().Get)
			productGroup.PUT("/:id", s.handler.Product().Update)
		}
	}
}
