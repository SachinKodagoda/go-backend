package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"go-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var database *mongo.Database

type SampleData struct {
	Categories []models.Category `json:"categories"`
	Products   []models.Product  `json:"products"`
	Users      []models.User     `json:"users"`
}

// Connect establishes a connection to MongoDB
func Connect() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mydb"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	database = client.Database(dbName)
	log.Println("Connected to MongoDB!")
	return nil
}

// Disconnect closes the MongoDB connection
func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}

// get the categories collection
func GetCategoriesCollection() *mongo.Collection {
	return database.Collection("categories")
}

// get the products collection
func GetProductsCollection() *mongo.Collection {
	return database.Collection("products")
}

// get the users collection
func GetUsersCollection() *mongo.Collection {
	return database.Collection("users")
}

// InitializeDatabase initializes the database with sample data if collections are empty
func InitializeDatabase() error {
	// Create indexes
	err := createIndexes()
	if err != nil {
		return err
	}

	// Check if collections already have data
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	categoriesCount, err := GetCategoriesCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	productsCount, err := GetProductsCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	usersCount, err := GetUsersCollection().CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}

	if categoriesCount > 0 && productsCount > 0 && usersCount > 0 {
		log.Println("Database already initialized")
		return nil
	}

	// Load sample data
	sampleData, err := loadSampleData()
	if err != nil {
		return err
	}

	// Insert categories
	if categoriesCount == 0 && len(sampleData.Categories) > 0 {
		var categoriesInterface []interface{}
		for _, category := range sampleData.Categories {
			categoriesInterface = append(categoriesInterface, category)
		}

		_, err = GetCategoriesCollection().InsertMany(ctx, categoriesInterface)
		if err != nil {
			return err
		}
		log.Println("Categories initialized successfully")
	}

	// Insert products
	if productsCount == 0 && len(sampleData.Products) > 0 {
		var productsInterface []interface{}
		for _, product := range sampleData.Products {
			productsInterface = append(productsInterface, product)
		}

		_, err = GetProductsCollection().InsertMany(ctx, productsInterface)
		if err != nil {
			return err
		}
		log.Println("Products initialized successfully")
	}

	// Insert users
	if usersCount == 0 && len(sampleData.Users) > 0 {
		var usersInterface []interface{}
		for _, user := range sampleData.Users {
			usersInterface = append(usersInterface, user)
		}

		_, err = GetUsersCollection().InsertMany(ctx, usersInterface)
		if err != nil {
			return err
		}
		log.Println("Users initialized successfully")
	}

	return nil
}

// Create indexes for better query performance
func createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create index on categories
	_, err := GetCategoriesCollection().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Create index on products
	_, err = GetProductsCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "category_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "category_group", Value: 1},
			},
		},
	})
	if err != nil {
		return err
	}

	// Create index on users
	_, err = GetUsersCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// Load sample data from file
func loadSampleData() (*SampleData, error) {
	sampleDataJSON := `{
  "users": [
    {
      "id": "1",
      "email": "admin@gmail.com",
      "password": "lucytech@123",
      "name": "Admin User",
      "role": "admin"
    },
    {
      "id": "2",
      "email": "user@gmail.com",
      "password": "lucytech@123",
      "name": "Regular User",
      "role": "user"
    }
  ],
  "categories": [
    {
      "id": "1",
      "name": "Electronics",
      "parent_id": null
    },
    {
      "id": "2",
      "name": "Smartphones",
      "parent_id": "1"
    },
    {
      "id": "3",
      "name": "Laptops",
      "parent_id": "1"
    },
    {
      "id": "4",
      "name": "Clothing",
      "parent_id": null
    },
    {
      "id": "5",
      "name": "Men",
      "parent_id": "4"
    },
    {
      "id": "6",
      "name": "Women",
      "parent_id": "4"
    },
    {
      "id": "7",
      "name": "Home & Kitchen",
      "parent_id": null
    },
    {
      "id": "8",
      "name": "Appliances",
      "parent_id": "7"
    },
    {
      "id": "9",
      "name": "Furniture",
      "parent_id": "7"
    },
    {
      "id": "10",
      "name": "Books",
      "parent_id": null
    }
  ],
  "products": [
    {
      "id": "1",
      "name": "iPhone 13",
      "category_id": "2",
      "category_group": "1",
      "attributes": [
        {
          "code": "price",
          "value": 4,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Black",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "tags",
          "value": [
            "smartphone",
            "apple",
            "5G",
            "test"
          ],
          "type": "tags"
        }
      ]
    },
    {
      "id": "2",
      "name": "Samsung Galaxy S21",
      "category_id": "2",
      "category_group": "1",
      "attributes": [
        {
          "code": "price",
          "value": 899,
          "type": "number"
        },
        {
          "code": "color",
          "value": "White",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "tags",
          "value": [
            "smartphone",
            "samsung",
            "5G"
          ],
          "type": "tags"
        }
      ]
    },
    {
      "id": "3",
      "name": "MacBook Pro",
      "category_id": "3",
      "category_group": "1",
      "attributes": [
        {
          "code": "price",
          "value": 1999,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Silver",
          "type": "text"
        },
        {
          "code": "specs_url",
          "value": "https://www.apple.com/macbook-pro/specs/",
          "type": "url"
        },
        {
          "code": "in_stock",
          "value": false,
          "type": "boolean"
        }
      ]
    },
    {
      "id": "4",
      "name": "Dell XPS 13",
      "category_id": "3",
      "category_group": "1",
      "attributes": [
        {
          "code": "price",
          "value": 1299,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Black",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        }
      ]
    },
    {
      "id": "5",
      "name": "Men's T-Shirt",
      "category_id": "5",
      "category_group": "4",
      "attributes": [
        {
          "code": "price",
          "value": 35,
          "type": "number"
        },
        {
          "code": "sizes",
          "value": [
            "S",
            "M",
            "L",
            "XL"
          ],
          "type": "tags"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        }
      ]
    },
    {
      "id": "6",
      "name": "Women's Jeans",
      "category_id": "6",
      "category_group": "4",
      "attributes": [
        {
          "code": "price",
          "value": 49.99,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Blue",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        }
      ]
    },
    {
      "id": "7",
      "name": "Coffee Maker",
      "category_id": "8",
      "category_group": "7",
      "attributes": [
        {
          "code": "price",
          "value": 79.95,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Stainless Steel",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "tags",
          "value": [
            "kitchen",
            "appliance",
            "brewing"
          ],
          "type": "tags"
        }
      ]
    },
    {
      "id": "8",
      "name": "Microwave Oven",
      "category_id": "8",
      "category_group": "7",
      "attributes": [
        {
          "code": "price",
          "value": 149.99,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Black",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "wattage",
          "value": 1100,
          "type": "number"
        }
      ]
    },
    {
      "id": "9",
      "name": "Sofa Sectional",
      "category_id": "9",
      "category_group": "7",
      "attributes": [
        {
          "code": "price",
          "value": 1299.99,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Gray",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": false,
          "type": "boolean"
        },
        {
          "code": "dimensions",
          "value": "112\" x 84\" x 38\"",
          "type": "text"
        }
      ]
    },
    {
      "id": "10",
      "name": "The Great Gatsby",
      "category_id": "10",
      "category_group": "10",
      "attributes": [
        {
          "code": "price",
          "value": 14.99,
          "type": "number"
        },
        {
          "code": "author",
          "value": "F. Scott Fitzgerald",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "formats",
          "value": [
            "paperback",
            "hardcover",
            "ebook"
          ],
          "type": "tags"
        }
      ]
    },
    {
      "id": "11",
      "name": "Wireless Earbuds",
      "category_id": "1",
      "category_group": "1",
      "attributes": [
        {
          "code": "price",
          "value": 149.99,
          "type": "number"
        },
        {
          "code": "color",
          "value": "White",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "battery_life",
          "value": "6 hours",
          "type": "text"
        }
      ]
    },
    {
      "id": "12",
      "name": "Women's Sweater",
      "category_id": "6",
      "category_group": "4",
      "attributes": [
        {
          "code": "price",
          "value": 45.5,
          "type": "number"
        },
        {
          "code": "color",
          "value": "Red",
          "type": "text"
        },
        {
          "code": "in_stock",
          "value": true,
          "type": "boolean"
        },
        {
          "code": "sizes",
          "value": [
            "S",
            "M",
            "L"
          ],
          "type": "tags"
        },
        {
          "code": "material",
          "value": "Cotton blend",
          "type": "text"
        }
      ]
    }
  ]
}`

	var sampleData SampleData
	err := json.Unmarshal([]byte(sampleDataJSON), &sampleData)
	if err != nil {
		return nil, err
	}

	return &sampleData, nil
}
