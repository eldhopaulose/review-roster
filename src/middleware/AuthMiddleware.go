package middleware

import (
    "net/http"
    "strings"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "github.com/google/uuid"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type AuthClaims struct {
    UserID uuid.UUID `json:"user_id"`
    jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
            c.Abort()
            return
        }

        token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
            c.Set("userID", claims.UserID)
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        c.Next()
    }
}
