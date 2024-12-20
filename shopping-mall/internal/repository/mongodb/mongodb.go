package mongodb

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopping-mall/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	db *mongo.Database
}
// Initialize MongoDB repository
func NewMongoRepository(uri string) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, &readpref.ReadPref{}); err != nil {
		return nil, err
	}
	return &MongoRepository{
		db: client.Database("shopping_mall"),
	}, nil
}

// User Repository Implementation
func (r *MongoRepository) CreateUser(ctx context.Context, user *models.User) error {
	collection := r.db.Collection("users")
	_, err := collection.InsertOne(ctx, user)
	return err
}

// User Repository Implementation
func (r *MongoRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	collection := r.db.Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// User Repository Implementation
func (r *MongoRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	collection := r.db.Collection("users")
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// User Repository Implementation
func (r *MongoRepository) UpdateUser(ctx context.Context, user *models.User) error {
	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

// Product Repository Implementation
func (r *MongoRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	collection := r.db.Collection("products")
	_, err := collection.InsertOne(ctx, product)
	return err
}
// Product Repository Implementation
func (r *MongoRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	collection := r.db.Collection("products")
	var product models.Product
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Product Repository Implementation
func (r *MongoRepository) ListProducts(ctx context.Context, skip, limit int) ([]*models.Product, error) {
	collection := r.db.Collection("products")
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *MongoRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	collection := r.db.Collection("products")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": product.ID}, bson.M{"$set": product})
	return err
}

func (r *MongoRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	collection := r.db.Collection("products")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Order Repository Implementation
func (r *MongoRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	collection := r.db.Collection("orders")
	_, err := collection.InsertOne(ctx, order)
	return err
}

func (r *MongoRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	collection := r.db.Collection("orders")
	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *MongoRepository) ListOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	collection := r.db.Collection("orders")
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*models.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *MongoRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	collection := r.db.Collection("orders")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": order.ID}, bson.M{"$set": order})
	return err
}
