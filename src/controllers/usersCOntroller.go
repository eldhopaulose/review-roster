package controllers

import (
	"os"
	"time"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"github.com/eldhopaulose/ReviewRoster/src/models"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/eldhopaulose/ReviewRoster/src/config/db"
	passwordvalidator "github.com/wagslane/go-password-validator"
)


var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var verifier = emailverifier.NewVerifier()

type AuthClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}


// Get all users
func GetAllUsersController(c *gin.Context) {
	var users []models.Users

	db := db.GetDB()

	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}


// Create a new user
func CreateUserController(c *gin.Context) {
	var user models.Users
	  const minEntropyBits = 60
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}


// Validate required fields
	if user.Email == "" || user.Name == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required: email, name, and password"})
		return
	}

	ret, err := verifier.Verify(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	if !ret.Syntax.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	//validate password

	if err := passwordvalidator.Validate(user.Password, minEntropyBits); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong enough"})
		return
	}


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	db := db.GetDB()
	if err := db.Create(&user).Error; err != nil {
		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "dtails": err.Error })
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "dtails": err.Error })
		}

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})
	  tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        // Handle the error (e.g., log it or return an error response)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"user": user, "token": tokenString})
}


// Login a user
func LoginUserController(c *gin.Context) {
    var user models.Users
    var loginRequest struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if loginRequest.Email == "" || loginRequest.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required: Email and Password"})
        return
    }

    ret, err := verifier.Verify(loginRequest.Email)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
        return
    }

    if !ret.Syntax.Valid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
        return
    }

    db := db.GetDB()
    if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
        }
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthClaims{
        UserID: user.ID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
        },
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
