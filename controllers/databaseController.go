package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"Database_Analyzer/models"
	"Database_Analyzer/services"

	"github.com/gin-gonic/gin"
)

func SaveDatabaseConfiguration(c *gin.Context) {

	var dbConn models.DatabaseConfiguration

	if err := c.ShouldBindJSON(&dbConn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "Error": err.Error()})
		return
	}

	id, err := services.SaveDatabaseConfiguration(&dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": 201,
		"id":     id,
	})
}

func ScanDatabaseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "ID must be a Integer"})
		return
	}

	err = services.ScanDatabaseByID(id)
	if err != nil {
		if errors.Is(err, services.ErrConfigurationNotFound) || errors.Is(err, services.ErrUserWithoutPrivilegies) {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err.Error()})
			return
		}
		if errors.Is(err, services.ErrDatabaseFailed) || errors.Is(err, services.ErrEnconding) {
			c.JSON(http.StatusFailedDependency, gin.H{"status": http.StatusFailedDependency, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func GetReportByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "ID must be a Integer"})
		return
	}

	report, err := services.GetReportByID(id)
	if err != nil {
		if errors.Is(err, services.ErrReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Report": report})
}

func GetHTMLReportByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "error": "ID must be a Integer"})
		return
	}

	report, err := services.GetReportByID(id)
	if err != nil {
		if errors.Is(err, services.ErrReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}

	htmlReport, err := services.GenerateScanSummary(*report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlReport))
}
