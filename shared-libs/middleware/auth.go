package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	userv1 "pb_schemas/user/v1"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthConfig holds configuration for the authentication middleware
type AuthConfig struct {
	UserServiceURL string
	Timeout        time.Duration
}

// UserInfo contains authenticated user information
type UserInfo struct {
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// JWTAuthMiddleware creates a Gin middleware for JWT authentication
func JWTAuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Extract JWT token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header is required",
				"message": "Please provide a valid JWT token in the Authorization header",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "JWT token is required",
				"message": "Please provide a valid JWT token",
			})
			c.Abort()
			return
		}

		// Validate token with user service
		userInfo, err := validateTokenWithUserService(token, config)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid or expired token",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user", userInfo)
		c.Set("user_email", userInfo.Email)
		c.Set("user_roles", userInfo.Roles)

		// Continue to next handler
		c.Next()
	})
}

// validateTokenWithUserService validates JWT token by calling user service gRPC endpoint
func validateTokenWithUserService(token string, config AuthConfig) (*UserInfo, error) {
	// Set default timeout if not provided
	timeout := config.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Connect to user service
	conn, err := grpc.DialContext(ctx, config.UserServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create user service client
	client := userv1.NewUserServiceClient(conn)

	// Validate token
	resp, err := client.ValidateToken(ctx, &userv1.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	// Check if token is valid
	if !resp.Valid {
		return nil, gin.Error{Err: err, Type: gin.ErrorTypePublic}
	}

	return &UserInfo{
		Email: resp.UserEmail,
		Roles: resp.Roles,
	}, nil
}

// GetUserFromContext extracts user information from Gin context
func GetUserFromContext(c *gin.Context) (*UserInfo, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userInfo, ok := user.(*UserInfo)
	return userInfo, ok
}

// RequireRole creates a middleware that requires specific roles
func RequireRole(roles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		user, exists := GetUserFromContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "User information not found in context",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range user.Roles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Insufficient permissions",
				"message": "User does not have required role(s)",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}
