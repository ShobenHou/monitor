package api

import (
	"github.com/ShobenHou/monitor/internal/agent"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Agent  *agent.Agent //Temp Okay. To update: a list of agents
	Router *gin.Engine
}

func NewServer(agent *agent.Agent) *Server {
	s := &Server{
		Agent:  agent,
		Router: gin.Default(),
	}
	s.Router.Use(s.logRequest)

	//router.GET("/agents/:id/config", GetMonitorConfig)
	s.Router.PUT("agents/:id/config/update", s.UpdateMonitorConfig)

	return s
}

func (s *Server) logRequest(c *gin.Context) {
	// TODO:log the incoming request
	c.Next()
}

//

/*-----------old code---------------------
func SetupRoutes(router *gin.Engine) {
	router.GET("/agents/:id/config", GetMonitorConfig)
	router.PUT("/agents/:id/config", UpdateMonitorConfig)
}

*/
