package handlers

import (
	"backend-gin-gonic/database"
	_ "backend-gin-gonic/docs"
	"backend-gin-gonic/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// Login handles user login
// @Summary User Login
// @Description Authenticates a user and returns a JWT token.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param UserCred body models.LoginRequest true "User credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /login [post]
func Login(c *gin.Context) {

	var loginRequest models.LoginRequest

	if c.ShouldBindJSON(&loginRequest) != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to read body",
		})
		return
	}

	sqlStatement := `SELECT email, password FROM public.sk_user WHERE email = $1`
	var email, hashedPassword string

	fmt.Println(loginRequest.Email)
	fmt.Println(loginRequest.Password)
	fmt.Println(sqlStatement)
	result := database.DB.QueryRow(sqlStatement, loginRequest.Email).Scan(&email, &hashedPassword)
	fmt.Println(result)

	if result != nil {
		if result == sql.ErrNoRows {
			// If no user is found
			fmt.Println("No user found")
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid email or password",
			})
			return
		}
		// If some other error occurs
		fmt.Println("Query error:", result)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Something went wrong",
		})
		return
	}

	//fmt.Println(email, hashedPassword)

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": loginRequest.Email,
		"exp": time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Failed to create token",
		})
		return
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": loginRequest.Email,
		"exp": time.Now().Add(8 * time.Hour).Unix(), // Refresh token valid for 7 days
	})

	fmt.Println(os.Getenv("SECRET"))
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create refresh token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 900, "", "", false, true)
	c.SetCookie("RefreshToken", refreshTokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, models.LoginResponse{
		User:         "Login success for email " + loginRequest.Email,
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
		Email:        email,
	})
}

// @Summary Refresh Access Token
// @Description Refreshes the access token using the refresh token cookie.
// @Tags Authentication
// @Produce json
// @Success 200 {object} models.RefreshToken "Token refreshed successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Refresh token is missing or invalid"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error - Failed to create new access token"
// @Router /refresh [get]
func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("RefreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Refresh token is missing",
		})
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": claims["sub"],
			"exp": time.Now().Add(15 * time.Minute).Unix(),
		})

		accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to create new access token",
			})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", accessTokenString, 900, "", "", false, true) // 15 minute

		c.JSON(http.StatusOK, models.RefreshToken{
			Message:     "Token refreshed succecfully",
			AccessToken: accessTokenString,
		})
	} else {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid refresh token",
		})
	}
}
