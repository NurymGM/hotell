package main // CompileDaemon -command="./hotell"

import (
	"github.com/NurymGM/hotell/controllers"
	"github.com/NurymGM/hotell/initializers"
	"github.com/NurymGM/hotell/migrations"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.ConnectToRedis()
	migrations.Migrate()
}

func main() {
	r := gin.Default()

	r.GET("/", controllers.RootRoute)
	r.POST("/rooms", controllers.CreateRoom)
	r.GET("/rooms", controllers.ReadRooms)
	r.GET("/rooms/:id", controllers.ReadRoomByID)
	r.PUT("/rooms/:id", controllers.UpdateRoom)
	r.DELETE("/rooms/:id", controllers.DeleteRoom)

	r.Run()
}
