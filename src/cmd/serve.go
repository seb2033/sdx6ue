/*
Copyright © 2023 Vinzenz Stadtmueller vinzenz.stadtmueller@fh-hagenberg.at
*/
package cmd

import (
	"fmt"

	"net/http"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// STRUCTS

type Recipe struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Ingredients pq.StringArray `gorm:"type:varchar(64)[]" json:"ingredients"`
}

type request struct {
	URL      string      `json:"url"`
	Method   string      `json:"method"`
	Headers  http.Header `json:"headers"`
	Body     []byte      `json:"body"`
	ClientIP string      `json:"client_ip"`
}

// GLOBAL VARS
var db *gorm.DB

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")
		gin.SetMode(gin.ReleaseMode)
		var err error

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			viper.GetString("db_host"),
			viper.GetString("db_user"),
			viper.GetString("db_password"),
			viper.GetString("db_name"),
			viper.GetString("db_port"),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		err = db.AutoMigrate(&Recipe{})
		if err != nil {
			panic("failed to migrate")
		}

		r := gin.Default()

		r.GET("/", func (c *gin.Context) {
			c.Data(http.StatusOK, "text/plain", []byte("Recipe service"))
		})

		r.GET("/recipes", func(c *gin.Context) {
			var recipes []Recipe
			result := db.Find(&recipes)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, recipes)
		})

		// Endpoint to return a specific recipe with the given ID
		r.GET("/recipes/:id", func(c *gin.Context) {
			var recipe Recipe
			result := db.First(&recipe, c.Param("id"))
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
				return
			}
			c.JSON(http.StatusOK, recipe)
		})

		// Endpoint to create a new recipe
		r.POST("/recipes", func(c *gin.Context) {
			var recipe Recipe
			err := c.ShouldBindJSON(&recipe)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			result := db.Create(&recipe)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, recipe)
		})

		// Endpoint to update a specific recipe with the given ID
		r.PUT("/recipes/:id", func(c *gin.Context) {
			var recipe Recipe
			result := db.First(&recipe, c.Param("id"))
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
				return
			}
			err := c.ShouldBindJSON(&recipe)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			result = db.Save(&recipe)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			c.JSON(http.StatusOK, recipe)
		})

		// Endpoint to delete a specific recipe with the given ID
		r.DELETE("/recipes/:id", func(c *gin.Context) {
			var recipe Recipe
			result := db.First(&recipe, c.Param("id"))
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
				return
			}
			result = db.Delete(&recipe)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
			c.Status(http.StatusNoContent)
		})

		// Echo the request
		r.GET("/debug", func(ctx *gin.Context) {
			var err error
			rr := &request{}
			rr.Method = ctx.Request.Method
			rr.Headers = ctx.Request.Header
			rr.URL = ctx.Request.URL.String()
			rr.ClientIP = ctx.ClientIP()
			if err != nil {
				return
			}

			if err != nil {
				return
			}
			ctx.JSON(http.StatusOK, rr)
		})

		// Check for database connection
		r.GET("/health", func(c *gin.Context) {
			if d, ok := db.DB(); ok == nil {
				if err = d.Ping(); err != nil {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusOK)
		})

		err = r.Run(":8080")
		if err != nil {
			panic("error starting the server")
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
