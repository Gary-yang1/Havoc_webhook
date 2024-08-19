package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		var message Message

		// 尝试将请求体绑定到结构体
		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse JSON"})
			return
		}

		// 检查 embeds 是否为空
		if len(message.Embeds) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No embeds found in the message"})
			return
		}
		dingtalk(message)
	})
	err := r.Run(":8080")
	if err != nil {
		return
	} // listen and serve on
}
