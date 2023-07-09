package api

import (
	"fmt"
	kafkaHelper "github.com/ShobenHou/monitor/internal/pkg/kafka"
	"github.com/gin-gonic/gin"
)

/*--------------old code-------------------------
func GetMonitorConfig(c *gin.Context) {
	agentID := c.Param("id")

	monitorConfig, err := getMonitorConfigFromDB(agentID) // TODO:Replace this with your own logic to get the configuration from the database.
	//the user get agent config(that is running) from database

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, monitorConfig)
}
*/

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
	fmt.Printf("MonitorConf from /config/update PUT: %+v\n", monitorConf)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = kafkaHelper.PublishConfigToKafka(monitorConf) //TODO: add agentID
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Configuration updated successfully."})
	fmt.Printf("Configuration updated to kafka successfully.")
}
