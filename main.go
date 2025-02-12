package main // CompileDaemon -command="./hotell"

import (
	"github.com/NurymGM/hotell/controllers"
	"github.com/NurymGM/hotell/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.ConnectToRedis()
}

func main() {
	r := gin.Default()

	r.GET("/", controllers.RootRoute)

	r.Run()
}
