package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shopping-mall/internal/models"
	"github.com/shopping-mall/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrInsufficientStock  = errors.New("insufficient stock")
)

type Service struct {
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
	jwtSecret   []byte
}

func NewService(userRepo repository.UserRepository, productRepo repository.ProductRepository, orderRepo repository.OrderRepository, jwtSecret string) *Service {
	return &Service{
		userRepo:    userRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
		jwtSecret:   []byte(jwtSecret),
	}
}

// User Service
func (s *Service) Register(ctx context.Context, req *models.RegisterRequest) error {
	existing, _ := s.userRepo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, req *models.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}

// Product Service
func (s *Service) CreateProduct(ctx context.Context, req *models.CreateProductRequest) error {
	product := &models.Product{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.productRepo.CreateProduct(ctx, product)
}

func (s *Service) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return s.productRepo.GetProductByID(ctx, id)
}

func (s *Service) ListProducts(ctx context.Context, page, pageSize int) ([]*models.Product, error) {
	skip := (page - 1) * pageSize
	return s.productRepo.ListProducts(ctx, skip, pageSize)
}

// Order Service
func (s *Service) CreateOrder(ctx context.Context, userID uuid.UUID, req *models.CreateOrderRequest) error {
	// Validate and update product stock
	var totalAmount float64
	for _, item := range req.Items {
		product, err := s.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return err
		}

		if product.Stock < item.Quantity {
			return ErrInsufficientStock
		}

		product.Stock -= item.Quantity
		if err := s.productRepo.UpdateProduct(ctx, product); err != nil {
			return err
		}

		totalAmount += product.Price * float64(item.Quantity)
	}

	order := &models.Order{
		ID:            uuid.New(),
		UserID:        userID,
		Items:         req.Items,
		TotalAmount:   totalAmount,
		Status:        models.OrderStatusPending,
		ShippingAddr:  req.ShippingAddr,
		PaymentMethod: req.PaymentMethod,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return s.orderRepo.CreateOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(ctx, id)
}

func (s *Service) ListUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	return s.orderRepo.ListOrdersByUserID(ctx, userID)
}

func (s *Service) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = status
	order.UpdatedAt = time.Now()
	return s.orderRepo.UpdateOrder(ctx, order)
}
