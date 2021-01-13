package main

import (
	"log"
	"strconv"
	"net/http"
	"grpc-mongodb-crud/client/handlers"
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	pb "grpc-mongodb-crud/proto"
)

func main(){
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil { log.Fatalf("Something went wrong: %s", err) }
	defer conn.Close()

	r := gin.Default()

	// Routes

	r.GET("/bookstore", func(c *gin.Context) {
	  if res, err := handlers.GetAllBooks(conn); err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  } else {
	    c.JSON(http.StatusOK, gin.H{ "books": res.Books })
	  }
	})

	r.GET("/bookstore/:id", func(c *gin.Context) {
	  id := c.Param("id")
	  if res, err := handlers.GetBookById(conn, id); err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  } else {
	    c.JSON(http.StatusOK, gin.H{ "book": res.Book })
	  }
	})

	r.POST("/bookstore/insert", func(c *gin.Context) {
	  release_year, err := strconv.Atoi(c.PostForm("release_year"))
	  if err != nil { 
	  	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  }

	  book := &pb.Book {
	  	Title: c.PostForm("title"),
	  	Description: c.PostForm("description"),
	  	Author: &pb.BookAuthor {
	  		Firstname:  c.PostForm("firstname"),
	  		Lastname:  c.PostForm("lastname"),
	  	},
	  	ReleaseYear: int32(release_year),
	  }

	  if res, err := handlers.InsertBook(conn, &pb.BookReq{ Book: book }); err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  } else {
	    c.JSON(http.StatusOK, gin.H{ "book": res.Book })
	  }
	})

	r.POST("/bookstore/update/:id", func(c *gin.Context) {
	  id := c.Param("id")
	  release_year, err := strconv.Atoi(c.PostForm("release_year"))
	  if err != nil { 
	  	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  }

	  book := &pb.Book {
	  	Title: c.PostForm("title"),
	  	Description: c.PostForm("description"),
	  	Author: &pb.BookAuthor {
	  		Firstname:  c.PostForm("firstname"),
	  		Lastname:  c.PostForm("lastname"),
	  	},
	  	ReleaseYear: int32(release_year),
	  }

	  if res, err := handlers.UpdateBook(conn, &pb.UpdateBookReq{ Id: id, Book: book }); err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  } else {
	    c.JSON(http.StatusOK, gin.H{ "book": res.Book })
	  }
	})

	r.POST("/bookstore/delete/:id", func(c *gin.Context) {
	  id := c.Param("id")

	  if res, err := handlers.DeleteBook(conn, id); err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	  } else {
	    c.JSON(http.StatusOK, gin.H{ "success": res.Success })
	  }
	})
	

	// Run server

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}