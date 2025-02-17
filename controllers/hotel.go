package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/NurymGM/hotell/initializers"
	"github.com/NurymGM/hotell/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RootRoute(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello wws!"})
}

func CreateRoom(c *gin.Context) {
	input := models.Room{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	room := models.Room{Type: input.Type, Price: input.Price, Info: input.Info, IsAvailable: input.IsAvailable, Image: input.Image}
	result := initializers.DB.Create(&room)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error creating room"})
		return
	}
	c.IndentedJSON(http.StatusCreated, room)
}

func ReadRooms(c *gin.Context) {
	rooms := []models.Room{}
	result := initializers.DB.Select("id, type, price, info, is_available, image").Limit(10).Find(&rooms)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error loading rooms"})
		return
	}
	if len(rooms) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No rooms available"})
		return
	}
	c.IndentedJSON(http.StatusOK, rooms)
}

func ReadRoomByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid room ID"})
		return
	}

	// check Redis cache, if hit, respond with it
	val, err := initializers.RDB.Get(context.Background(), strconv.Itoa(id)).Result()
	if err == nil {
		// Deserialize room
		room := models.Room{}
		err2 := json.Unmarshal([]byte(val), &room)
		if err2 == nil {
			c.IndentedJSON(http.StatusOK, room)
			return
		}
	}

	// else get post from PostgreSQL
	room := models.Room{}
	result := initializers.DB.First(&room, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Room not found"})
		return
	}
	if result.Error != nil {
		log.Println("Database error:", result.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error loading room"})
		return
	}

	// Serialize room then add it to Redis cache
	jsonData, err := json.Marshal(room)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error marshalling room"})
		return
	}
	err = initializers.RDB.Set(context.Background(), strconv.Itoa(id), jsonData, 0).Err()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to set value in Redis"})
		return
	}

	c.IndentedJSON(http.StatusOK, room)
}

func UpdateRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid room ID"})
		return
	}
	// validate JSON input before DB query and parse json into a map
    var updateData map[string]interface{}
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }
	// now look for the room
	room := models.Room{}
	if err := initializers.DB.First(&room, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Room not found"})
			return
		}
		log.Println("Database error:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error finding room"})
		return
	}
    // Perform the update dynamically
    result := initializers.DB.Model(&room).Updates(updateData)
    if result.Error != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error updating room"})
        return
    }
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Room updated successfully"})
}

func DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid room ID"})
		return
	}
	// first look for the room
	room := models.Room{}
	if err := initializers.DB.First(&room, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Room not found"})
			return
		}
		log.Println("Database error:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error finding room"})
		return
	}
	// now delete
	result := initializers.DB.Delete(&room)
	if result.Error != nil {
		log.Println("Database error:", result.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error deleting room"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}