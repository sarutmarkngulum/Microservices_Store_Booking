package services

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
	"gitlab.com/final_project1240930/user_management_service/internal/logs"
	"gitlab.com/final_project1240930/user_management_service/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	userRepo repository.UserRepository
	UnimplementedUserServiceServer
}

func NewUserServer(userRepo repository.UserRepository) UserServiceServer {
	return userServer{userRepo: userRepo}
}

// func (userServer) mustEmbedUnimplementedUserServiceServer() {}

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Println("Warning: JWT_SECRET environment variable is not set")
	} else {
		log.Println("JWT_SECRET loaded successfully")
	}
}

func generateToken(username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s userServer) Register(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	// Validate both username and password
	if req.Username == "" || req.Password == "" {
		logs.Error("Username and password must be provided")
		return nil, status.Errorf(codes.InvalidArgument, "Username and password must be provided")
	}
	if len(req.Password) < 8 {
		logs.Error("Password must be at least 8 characters long")
		return nil, status.Errorf(codes.InvalidArgument, "Password must be at least 8 characters long")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error("Failed to hash password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	newUser := &repository.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		logs.Error("Failed to register user", zap.Error(err))
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, status.Errorf(codes.AlreadyExists, "Username already exists")
		}
		return nil, status.Errorf(codes.Internal, "Failed to register user")
	}

	res := &UserResponse{Username: req.Username}
	return res, nil
}

func (s userServer) Login(ctx context.Context, req *LoginRequest) (*Token, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		logs.Error("Username and password must be provided")
		return nil, status.Errorf(codes.InvalidArgument, "Username and password must be provided")
	}

	// Check if the user is valid
	user, err := s.userRepo.GetUser(req.Username)
	if err != nil {
		logs.Error("Failed to get user", zap.String("username", req.Username), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logs.Error("Invalid credentials for username", zap.String("username", req.Username))
		return nil, status.Errorf(codes.Unauthenticated, "Invalid credentials")
	}

	// Generate the token with role
	accessToken, err := generateToken(req.Username, user.Role)
	if err != nil {
		logs.Error("Failed to generate token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to generate token")
	}

	return &Token{AccessToken: accessToken}, nil
}

func (s userServer) GetAllUsers(ctx context.Context, req *Empty) (*UserList, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		logs.Error("Failed to get all users", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to get all users")
	}

	var userResponses []*UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &UserResponse{Username: user.Username, Role: user.Role})
	}

	return &UserList{Users: userResponses}, nil
}

func (s userServer) GetUser(ctx context.Context, req *UserIdentifier) (*UserResponse, error) {
	user, err := s.userRepo.GetUser(req.Identifier)
	if err != nil {
		logs.Error("Failed to get user", zap.String("identifier", req.Identifier), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	return &UserResponse{Username: user.Username, Role: user.Role}, nil
}

func (s userServer) UpdateUserRole(ctx context.Context, req *UpdateRoleRequest) (*UserResponse, error) {
	user, err := s.userRepo.GetUser(req.Identifier)
	if err != nil {
		logs.Error("Failed to get user", zap.String("identifier", req.Identifier), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	user.Role = req.Role
	if err := s.userRepo.UpdateUser(user.UUID, user); err != nil {
		logs.Error("Failed to update user role", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to update user role")
	}

	return &UserResponse{Username: user.Username, Role: user.Role}, nil
}

func (s userServer) UpdateUserPassword(ctx context.Context, req *UpdatePasswordRequest) (*UserResponse, error) {
	// Validate password length
	if len(req.Password) < 8 {
		logs.Error("Password must be at least 8 characters long")
		return nil, status.Errorf(codes.InvalidArgument, "Password must be at least 8 characters long")
	}

	user, err := s.userRepo.GetUser(req.Identifier)
	if err != nil {
		logs.Error("Failed to get user", zap.String("identifier", req.Identifier), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Hash the new password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		logs.Error("Failed to hash password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	user.Password = hashedPassword
	if err := s.userRepo.UpdateUser(user.UUID, user); err != nil {
		logs.Error("Failed to update user password", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to update user password")
	}

	return &UserResponse{Username: user.Username}, nil
}

// Utility function
func hashPassword(password string) (string, error) {
	// Example: Use bcrypt for hashing
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func (s userServer) DeleteUser(ctx context.Context, req *UserIdentifier) (*Empty, error) {
	if err := s.userRepo.DeleteUser(req.Identifier); err != nil {
		logs.Error("Failed to delete user", zap.String("identifier", req.Identifier), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Failed to delete user")
	}

	return &Empty{}, nil
}
