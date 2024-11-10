package main

import (
	"encoding/json"
	"html/template"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var templates *template.Template

func main() {
	// Parse templates
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Fetch and store data
	artists, err := fetchArtists()
	if err != nil {
		log.Fatalf("Error fetching artists: %v", err)
	}
	saveJSON("artists.json", artists)

	locations, err := fetchLocations()
	if err != nil {
		log.Fatalf("Error fetching locations: %v", err)
	}
	saveJSON("locations.json", locations)

	dates, err := fetchDates()
	if err != nil {
		log.Fatalf("Error fetching dates: %v", err)
	}
	saveJSON("dates.json", dates)

	relations, err := fetchRelations()
	if err != nil {
		log.Fatalf("Error fetching relations: %v", err)
	}
	saveJSON("relations.json", relations)

	// Initialize Gin router
	router := gin.Default()

	// Load HTML templates
	router.SetHTMLTemplate(templates)
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:4173", "http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}

	router.Use(cors.New(config))
	// Serve static files
	router.Static("/static", "./static")

	// Routes
	router.GET("/", indexHandler)
	router.GET("/artists", artistsHandler)
	router.GET("/artist/:id", artistHandler)
	router.GET("/locations", locationsHandler)
	router.GET("/dates", datesHandler)
	router.GET("/relations", relationsHandler)

	// Artist-specific routes
	router.GET("/artist/locations/:id", artistLocationsHandler)
	router.GET("/artist/relations/:id", artistRelationsHandler)
	router.GET("/artist/dates/:id", artistDatesHandler)

	// Start server
	log.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}

func saveJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(data)
}
