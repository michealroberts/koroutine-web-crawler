/*****************************************************************************************************************/

//	@author		Michael Roberts

/*****************************************************************************************************************/

package main

/*****************************************************************************************************************/

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	crawler "github.com/michealroberts/koroutine-web-crawler/pkg/crawler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*****************************************************************************************************************/

func setupRouter() *gin.Engine {
	// A new gin base router:
	router := gin.Default()

	config := cors.DefaultConfig()
	// This should be more restrictive based on the requirements:
	config.AllowOrigins = []string{"*"}

	// Setup the CORS middleware:
	router.Use(cors.New(config))

	// Setup the crawl endpoint
	// Setup the SSE route
	router.GET("/crawl", func(c *gin.Context) {
		// Retrieve query parameters or set defaults
		domain := c.DefaultQuery("domain", "https://example.com")

		depth := c.DefaultQuery("depth", "2")

		maxDepth, err := strconv.Atoi(depth)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid depth parameter"})
			return
		}

		// Setup headers for SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// Ensure that the Gin does not buffer the responses
		c.Writer.Flush()

		// Start the crawler in a new goroutine
		crawler := crawler.New()

		go crawler.Crawl(domain, maxDepth)

		// Create a ticker for keep-alive messages
		ticker := time.NewTicker(30 * time.Second)

		defer ticker.Stop()

		done := c.Request.Context().Done()

		for {
			select {
			case _, ok := <-crawler.Stream():
				if !ok {
					// Close the stream if channel is closed
					return
				}
				jsonNode, err := json.Marshal(crawler.Root)

				if err != nil {
					log.Printf("Failed to marshal node: %v", err)
					continue // Skip bad data
				}

				data := fmt.Sprintf("data: %s\n\n", string(jsonNode))

				_, writeErr := c.Writer.WriteString(data)

				if writeErr != nil {
					// Handle errors such as broken connections
					log.Println("Write error:", writeErr)
					return
				}

				// Flush the response
				c.Writer.Flush()
			case <-ticker.C:
				// Send a keep-alive message periodically
				_, writeErr := c.Writer.WriteString(": keep-alive\n\n")
				if writeErr != nil {
					log.Println("Write error on keep-alive:", writeErr)
					return
				}
				c.Writer.Flush()
			case <-done:
				// End the request when the client disconnects
				log.Println("Client has disconnected")
				return
			}
		}
	})

	return router
}

/*****************************************************************************************************************/

func main() {
	router := setupRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// Service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listen: %s\n", err) // Changed from log.Fatalf to log.Printf
			os.Exit(1)                      // Explicit exit call after log, which allows defers to execute
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a buffered channel
	quit := make(chan os.Signal, 1) // Buffer size of 1
	
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %s", err) // Changed from log.Fatal to log.Printf
	}

	log.Println("Server exiting")
}

/*****************************************************************************************************************/
