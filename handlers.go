package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// JSON API Handlers
func indexHandler(c *gin.Context) {
	artists, err := loadArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"artists": artists,
	})
}

func artistsHandler(c *gin.Context) {
	artists, err := loadArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, artists)
}

func artistHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	artists, err := loadArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var artist Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}

	c.JSON(http.StatusOK, artist)
}

func locationsHandler(c *gin.Context) {
	locations, err := loadLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artists, err := loadArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artistNames := make(map[int]string)
	for _, artist := range artists {
		artistNames[artist.ID] = artist.Name
	}

	for i, location := range locations.Index {
		locations.Index[i].Name = artistNames[location.ID]
	}

	c.JSON(http.StatusOK, locations)
}

func datesHandler(c *gin.Context) {
	dates, err := loadDates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artists, err := loadArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artistNames := make(map[int]string)
	for _, artist := range artists {
		artistNames[artist.ID] = artist.Name
	}

	for i, date := range dates.Index {
		dates.Index[i].Name = artistNames[date.ID]
	}

	c.JSON(http.StatusOK, dates)
}

func relationsHandler(c *gin.Context) {
	relations, err := loadRelations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artists, err := loadArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	artistNames := make(map[int]string)
	for _, artist := range artists {
		artistNames[artist.ID] = artist.Name
	}

	for i, relation := range relations.Index {
		relations.Index[i].Name = artistNames[relation.ID]
	}

	c.JSON(http.StatusOK, relations)
}

// HTML Template Handlers
func artistLocationsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid artist ID"})
		return
	}

	artist, err := getArtistByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	locations, err := fetchLocations()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var artistLocations []string
	for _, loc := range locations.Index {
		if loc.ID == artist.ID {
			artistLocations = loc.Locations
			break
		}
	}

	c.HTML(http.StatusOK, "artist_locations.html", gin.H{
		"Name":      artist.Name,
		"Locations": artistLocations,
	})
}

func artistRelationsHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid artist ID"})
		return
	}

	artist, err := getArtistByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	relations, err := fetchRelations()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var artistRelations map[string][]string
	for _, rel := range relations.Index {
		if rel.ID == artist.ID {
			artistRelations = rel.DatesLocations
			break
		}
	}

	c.HTML(http.StatusOK, "artist_relations.html", gin.H{
		"Name":      artist.Name,
		"Relations": artistRelations,
	})
}

func artistDatesHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"error": "Invalid artist ID"})
		return
	}

	artist, err := getArtistByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	dates, err := fetchDates()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var artistDates []string
	for _, date := range dates.Index {
		if date.ID == artist.ID {
			artistDates = date.Dates
			break
		}
	}

	c.HTML(http.StatusOK, "artist_dates.html", gin.H{
		"Name":  artist.Name,
		"Dates": artistDates,
	})
}

// Uncomment and modify if you want to implement the search functionality

func searchHandler(c *gin.Context) {
	query := c.Query("query")
	artists, err := loadArtists()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	var searchResults []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			searchResults = append(searchResults, artist)
		}
	}

	c.HTML(http.StatusOK, "artists.html", gin.H{
		"Query":   query,
		"Artists": searchResults,
	})
}

func loadArtists() ([]Artist, error) {
	file, err := os.Open("artists.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var artists []Artist
	err = json.NewDecoder(file).Decode(&artists)
	return artists, err
}

func loadLocations() (LocationsData, error) {
	file, err := os.Open("locations.json")
	if err != nil {
		return LocationsData{}, err
	}
	defer file.Close()

	var locations LocationsData
	err = json.NewDecoder(file).Decode(&locations)
	return locations, err
}

func loadDates() (DatesData, error) {
	file, err := os.Open("dates.json")
	if err != nil {
		return DatesData{}, err
	}
	defer file.Close()

	var dates DatesData
	err = json.NewDecoder(file).Decode(&dates)
	return dates, err
}

func loadRelations() (RelationsData, error) {
	file, err := os.Open("relations.json")
	if err != nil {
		return RelationsData{}, err
	}
	defer file.Close()

	var relations RelationsData
	err = json.NewDecoder(file).Decode(&relations)
	return relations, err
}

func getArtistByID(id int) (Artist, error) {
	artists, err := loadArtists()
	if err != nil {
		return Artist{}, err
	}

	for _, artist := range artists {
		if artist.ID == id {
			return artist, nil
		}
	}

	return Artist{}, fmt.Errorf("artist not found")
}

