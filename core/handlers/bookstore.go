package handlers

import (
	"fmt"
	"context"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	_ "google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"

	pb "grpc-mongodb-crud/proto"
	db "grpc-mongodb-crud/core/config"
	"grpc-mongodb-crud/core/models"
)

type Server struct {
	pb.UnimplementedBookstoreServer
}

func (s *Server) CreateBook(ctx context.Context, in *pb.BookReq) (*pb.BookRes, error) {
	book := in.GetBook()

	data := models.Book{
		// Id:    Empty, so it gets omitted and MongoDB generates a unique Object ID upon insertion.
		Title:    book.GetTitle(),
		Description:  book.GetDescription(),
		Author: models.BookAuthor{
			Firstname: book.GetAuthor().GetFirstname(),
			Lastname: book.GetAuthor().GetLastname(),
		},
		ReleaseYear: book.GetReleaseYear(),
	}

	// Insert the data into the database, result contains the newly generated Object ID for the new document
	result, err := db.Collection.InsertOne(db.Mongoctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v\n", err),
		)
	}

	oid := result.InsertedID.(primitive.ObjectID)
	book.Id = oid.Hex()

	return &pb.BookRes{ Book: book }, nil
}

func (s *Server) GetBook(ctx context.Context, in *pb.BookQuery) (*pb.BookRes, error) {

	var book *pb.Book

	oid, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid argument: %v\n", err),
		)
	}

	if err := db.Collection.FindOne(ctx, bson.M{"_id": oid }).Decode(&book); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v\n", err),
		)
	}

	return &pb.BookRes{ Book: book }, nil
}

func (s *Server) GetAllBooks(ctx context.Context, in *pb.EmptyQuery) (*pb.BookList, error) {
	var books []*pb.Book

	cur, err := db.Collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v\n", err),
		)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
	  var book *pb.Book

	  if err := cur.Decode(&book); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v\n", err),
		)
	  }
	  	
	  // insert into books
	  books = append(books, book)
	}

	if err := cur.Err(); err != nil {
	  return nil, status.Errorf(
	  	codes.Internal,
	  	fmt.Sprintf("Internal error: %s\n", err),
	  )
	}

	return &pb.BookList{ Books: books }, nil
}

func (s *Server) UpdateBook(ctx context.Context, in *pb.UpdateBookReq) (*pb.BookRes, error) {
	oid, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid argument error: %v", err),
		)
	}

	book := in.GetBook()

	updated := bson.M{
		"title": book.GetTitle(),
		"description": book.GetDescription(),
		"author": book.GetAuthor(),
		"release_year": book.GetReleaseYear(),
	}

	var result *pb.Book

	if err := db.Collection.FindOneAndUpdate(
		db.Mongoctx,
		bson.M{"_id": oid},
		bson.M{"$set": updated},
		options.FindOneAndUpdate().SetReturnDocument(1),
	).Decode(&result); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	return &pb.BookRes{ Book: result }, nil
}


func (s *Server) DeleteBook(ctx context.Context, in *pb.BookQuery) (*pb.DeleteBookRes, error) {
	id, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return &pb.DeleteBookRes{ Success: false }, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Invalid argument error: %s", err),
		)
	}

	if _, err := db.Collection.DeleteOne(
		db.Mongoctx,
		bson.M{ "_id": id },
	); err != nil {
		return &pb.DeleteBookRes{ Success: false }, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %s", err),
		)
	}

	return &pb.DeleteBookRes{ Success: true }, nil
}