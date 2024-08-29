package main

import (
	"fmt"
	"golang/models"
	"os"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func GetConnection() *gorm.DB {
	db.LogMode(true)
	return db
}

func ConnectDB() (err error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASS"), os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_SCHEMA"))
	db, err = gorm.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	// Migration database
	err = db.AutoMigrate(
		&models.UserData{},
	).Error

	if err != nil {
		fmt.Println("ERROR")
	}
	return
}

// In-memory data store
var userDatas []models.UserData

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("Error load .env")
	}
	ConnectDB()

	r := gin.Default()
	// Lấy tất cả user
	r.GET("/users", func(c *gin.Context) {
		var users []models.UserData
		db.Find(&users)
		c.JSON(200, users)
	})

	// Tạo user mới
	r.POST("/users", func(c *gin.Context) {
		var user models.UserData
		if err := c.ShouldBindJSON(&user); err == nil {
			db.Create(&user)
			c.JSON(201, user)
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}
	})

	// Cập nhật user
	r.PUT("/users/:id", func(c *gin.Context) {
		var user models.UserData
		id := c.Param("id")
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		if err := c.ShouldBindJSON(&user); err == nil {
			db.Save(&user)
			c.JSON(200, user)
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}
	})

	// Xóa user
	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&models.UserData{}, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}
		c.Status(204)
	})
	r.Run("127.0.0.1:3000")
}


//http

// func usersHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		// Get all products
// 		json.NewEncoder(w).Encode(users)
// 	case http.MethodPost:
// 		// Create a new product
// 		var newUser User
// 		json.NewDecoder(r.Body).Decode(&newUser)
// 		newUser.ID = len(users) + 1
// 		users = append(users, newUser)
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(newUser)
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

// func userHandler(w http.ResponseWriter, r *http.Request) {
// 	userID := r.URL.Path[len("/user/"):]
// 	switch r.Method {
// 	case http.MethodGet:
// 		// Get a single user
// 		for _, user := range users {
// 			if fmt.Sprintf("%d", user.ID) == userID {
// 				json.NewEncoder(w).Encode(user)
// 				return
// 			}
// 		}
// 		w.WriteHeader(http.StatusNotFound)
// 	case http.MethodPut:
// 		// Update a user
// 		var updatedUser User
// 		json.NewDecoder(r.Body).Decode(&updatedUser)
// 		for i, user := range users {
// 			if fmt.Sprintf("%d", user.ID) == userID {
// 				updatedUser.ID = user.ID
// 				users[i] = updatedUser
// 				json.NewEncoder(w).Encode(updatedUser)
// 				return
// 			}
// 		}
// 		w.WriteHeader(http.StatusNotFound)
// 	case http.MethodDelete:
// 		// Delete a product
// 		for i, user := range users {
// 			if fmt.Sprintf("%d", user.ID) == userID {
// 				users = append(users[:i], users[i+1:]...)
// 				w.WriteHeader(http.StatusNoContent)
// 				json.NewEncoder(w).Encode(true)
// 				return
// 			}
// 		}
// 		w.WriteHeader(http.StatusNotFound)
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }
