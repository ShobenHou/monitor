package api

import (
	"fmt"
	kafkaHelper "github.com/ShobenHou/monitor/internal/pkg/kafka"
	"github.com/gin-gonic/gin"
	"net/http"
)

// This is actually s.handleConfigUpdate()
func (s *Server) UpdateMonitorConfig(c *gin.Context) {
	//1. get config from user PUT method
	//2. store it in DB
	//3. update config to agents
	var monitorConf kafkaHelper.MonitorConf
	err := c.BindJSON(&monitorConf)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//TODO:err = updateMonitorConfigInDB(monitorConf) // TODO: add agentID

	//TODO: use logger instead of print
	fmt.Printf("UpdateMonitorConfig--Got MonitorConf from /config/update PUT: %+v\n", monitorConf)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("UpdateMonitorConfig--Publishing MonitorConf to Kafka \n")
	err = kafkaHelper.PublishConfigToKafka(monitorConf) //TODO: add agentID
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Configuration updated successfully."})
	fmt.Printf("UpdateMonitorConfig--Configuration updated to kafka successfully.")
}

func (s *Server) HandleHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello!")
}
