package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
)

// Plant represents a plant entity
type Plant struct {
	ID               string  `json:"id" dynamodbav:"id"`
	Name             string  `json:"name" dynamodbav:"name"`
	ScientificName   string  `json:"scientific_name" dynamodbav:"scientific_name"`
	Family           string  `json:"family" dynamodbav:"family"`
	Type             string  `json:"type" dynamodbav:"type"`
	SunlightRequired string  `json:"sunlight_required" dynamodbav:"sunlight_required"`
	WaterInterval    int     `json:"water_interval" dynamodbav:"water_interval"`
	Height           float64 `json:"height" dynamodbav:"height"`
	Native           string  `json:"native" dynamodbav:"native"`
	Indoor           bool    `json:"indoor" dynamodbav:"indoor"`
}

// this will be changed in the main function
var (
	tableName = ""
	region    = "" // Change this to your AWS region
)

var dbClient *dynamodb.Client

func main() {
	err := godotenv.Load()

	tableName = os.Getenv("DYNAMO_DB_TABLE")
	region = os.Getenv("AWS_REGION")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create DynamoDB client
	dbClient = dynamodb.NewFromConfig(cfg)

	// Create DynamoDB table if it doesn't exist
	ensureTableExists()

	// Setup API routes
	router := mux.NewRouter()
	router.HandleFunc("/plants", getAllPlants).Methods("GET")
	router.HandleFunc("/plants/{id}", getPlant).Methods("GET")
	router.HandleFunc("/plants", createPlant).Methods("POST")
	router.HandleFunc("/plants/{id}", updatePlant).Methods("PUT")
	router.HandleFunc("/plants/{id}", deletePlant).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	fmt.Printf("Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

// Ensure the Plants table exists, create it if it doesn't
func ensureTableExists() {
	_, err := dbClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		// Create the table if it doesn't exist
		_, err = dbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("id"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("id"),
					KeyType:       types.KeyTypeHash,
				},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})

		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}

		// Wait for table to be created
		waiter := dynamodb.NewTableExistsWaiter(dbClient)
		if err := waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		}, 5*60); err != nil {
			log.Fatalf("Failed waiting for table to be active: %v", err)
		}

		fmt.Println("Plants table created successfully!")
	} else {
		fmt.Println("Plants table already exists.")
	}
}

// CRUD Operations

// GET all plants
func getAllPlants(w http.ResponseWriter, r *http.Request) {
	// Scan operation to get all items from the table
	result, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to scan items: %v", err))
		return
	}

	var plants []Plant
	err = attributevalue.UnmarshalListOfMaps(result.Items, &plants)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal items: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, plants)
}

// GET a specific plant
func getPlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get item operation to fetch a specific plant by ID
	result, err := dbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get item: %v", err))
		return
	}

	if result.Item == nil {
		respondWithError(w, http.StatusNotFound, "Plant not found")
		return
	}

	var plant Plant
	err = attributevalue.UnmarshalMap(result.Item, &plant)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to unmarshal item: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, plant)
}

// CREATE a new plant
func createPlant(w http.ResponseWriter, r *http.Request) {
	var plant Plant
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&plant); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Generate a UUID if ID is not provided
	if plant.ID == "" {
		plant.ID = uuid.New().String()
	}

	// Convert plant to DynamoDB attribute map
	av, err := attributevalue.MarshalMap(plant)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal plant: %v", err))
		return
	}

	// Put item operation to add a new plant
	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create plant: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, plant)
}

// UPDATE an existing plant
func updatePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var plant Plant
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&plant); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Ensure ID in path matches ID in body or set it
	plant.ID = id

	// Convert plant to DynamoDB attribute map
	av, err := attributevalue.MarshalMap(plant)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to marshal plant: %v", err))
		return
	}

	// Put item operation to update the plant
	_, err = dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update plant: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, plant)
}

// DELETE a plant
func deletePlant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete item operation
	_, err := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete plant: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success", "message": "Plant deleted"})
}

// Helper functions for HTTP responses
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
