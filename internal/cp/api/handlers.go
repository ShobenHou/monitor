package api

import (
	"github.com/ShobenHou/monitor/pkg/kafka"
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
	agentID := c.Param("id")

	var monitorConf MonitorConf
	err := c.BindJSON(&monitorConf)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = updateMonitorConfigInDB(agentID, monitorConf) // TODO: Replace this with your own logic to update the configuration in the database.

	//1. get config from user PUT method
	//2. store it in DB
	//3. update config to agents

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	err = kafka.PublishConfigToKafka(agentID, monitorConf)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Configuration updated successfully."})
}
