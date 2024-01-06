package handlers

import (
	"context"
	"fmt"
	"net/http"
	"test-app/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser model.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if !IsValidPhoneNumber(newUser.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	// Check if user already exists
	var existingUser model.User
	err := h.db.QueryRow(context.Background(), "SELECT id FROM users WHERE phone_number=$1", newUser.PhoneNumber).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this phone number already exists"})
		return
	}

	// Insert new user into the database
	_, err = h.db.Exec(context.Background(), "INSERT INTO users (name, phone_number) VALUES ($1, $2)", newUser.Name, newUser.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}
func (h *UserHandler) GenerateOTP(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if !IsValidPhoneNumber(req.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	var existingUser model.User
	errr := h.db.QueryRow(context.Background(), "SELECT id FROM users WHERE phone_number=$1", req.PhoneNumber).Scan(&existingUser.ID)
	fmt.Println(errr)
	if errr == pgx.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this phone number does not exists"})
		return
	} else if errr != nil {
		// For any other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	otp := GenerateOTP()
	expirationTime := time.Now().Add(time.Minute)

	// Update OTP and expiration time in the database
	// Replace this with actual database query
	_, err := h.db.Exec(context.Background(), "UPDATE users SET otp=$1, otp_expiration_time=$2 WHERE phone_number=$3", otp, expirationTime, req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"otp": otp}) // Consider not sending OTP back in response for real-world apps
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var request model.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if !IsValidPhoneNumber(request.PhoneNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone number format"})
		return
	}
	var user model.User
	err := h.db.QueryRow(context.Background(), "SELECT otp, otp_expiration_time FROM users WHERE phone_number=$1", request.PhoneNumber).Scan(&user.OTP, &user.OTPExpirationTime)

	// Check if the error is because the user was not found
	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		// For any other database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	fmt.Println(IsOTPExpired(user.OTPExpirationTime))
	if request.OTP == user.OTP && !IsOTPExpired(user.OTPExpirationTime) {
		h.db.Exec(context.Background(), "UPDATE users SET otp=$1, otp_expiration_time=$2 WHERE phone_number=$3", "", "", request.PhoneNumber)
		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
	}
}
