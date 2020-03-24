package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ErrCheck - obsługa błedów
// ================================================================================================
func ErrCheck(errNr error) bool {
	if errNr != nil {
		fmt.Println(errNr)
		return false
	}
	return true
}

// ErrCheck2 - obsługa błedów - wersja bez wypisywania błędu
// ================================================================================================
func ErrCheck2(errNr error) bool {
	if errNr != nil {
		return false
	}
	return true
}

//
// GetHostName - Wysłanie HOSTNAME
// ================================================================================================
func GetHostName(c *gin.Context) {
	log.Println("GetHostName()")
	c.Header("Access-Control-Allow-Origin", "*")

	hostname, err := os.Hostname()
	ErrCheck(err)
	c.JSON(http.StatusOK, hostname)
}

//
// Options - Obsługa request'u OPTIONS (CORS)
// ================================================================================================
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// main - program główny
// ================================================================================================
func main() {

	// Logowanie do pliku
	//
	logFileApp, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	ErrCheck(err)
	log.SetOutput(io.MultiWriter(os.Stdout, logFileApp))

	gin.DisableConsoleColor()
	logFileGin, err := os.OpenFile("gin.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	ErrCheck(err)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logFileGin)

	log.Println("=========================================")
	log.Println("======          APP START        ========")
	log.Println("=========================================")

	// SERVER HTTP
	// =======================================

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(Options)

	r.LoadHTMLGlob("./dist/*.html")
	r.StaticFS("/fonts", http.Dir("./dist/fonts"))
	r.StaticFS("/css", http.Dir("./dist/css"))
	r.StaticFS("/js", http.Dir("./dist/js"))
	r.StaticFile("/", "./dist/index.html")
	r.StaticFile("favicon.ico", "./dist/favicon.ico")
	r.StaticFile("samar.png", "./dist/samar.png")
	r.StaticFile("sign.png", "./dist/sign.png")

	api := r.Group("/api/v1")
	{
		api.GET("/gethostname", GetHostName)
	}
	// Start serwera WWW
	// =======================================
	r.Run(":8081")
}
