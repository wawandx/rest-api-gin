package routes

import (
	"os"
	"time"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/danilopolani/gocialite/structs"
	"github.com/dgrijalva/jwt-go"
	"github.com/wawandx/rest-api-gin/config"
	"github.com/wawandx/rest-api-gin/models"
)

// Redirect to correct oAuth URL
func RedirectHandler(context *gin.Context) {
	// Retrieve provider from route
	provider := context.Param("provider")

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("CLIENT_ID_GITHUB"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GITHUB"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/github/callback",
		},
		"google": {
			"clientID":     os.Getenv("CLIENT_ID_GOOGLE"),
			"clientSecret": os.Getenv("CLIENT_SECRET_GOOGLE"),
			"redirectURL":  os.Getenv("AUTH_REDIRECT_URL") + "/google/callback",
		},
	}

	providerScopes := map[string][]string{
		"github":   []string{},
		"google": []string{},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := config.Gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		context.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	context.Redirect(http.StatusFound, authURL)
}

// Handle callback of provider
func CallbackHandler(context *gin.Context) {
	// Retrieve query params for state and code
	state := context.Query("state")
	code := context.Query("code")
	provider := context.Param("provider")

	// Handle callback and check for errors
	user, _, err := config.Gocial.Handle(state, code)
	if err != nil {
		context.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	var newUser = getOrRegisterUser(provider, user)
	var jwtToken = createToken(&newUser)

	context.JSON(200, gin.H{
		"data": newUser,
		"token": jwtToken,
		"message": "Login Success",
	})
}

func getOrRegisterUser(provider string, user *structs.User) models.User{
	var userData models.User

	config.DB.Where("provider = ? AND social_id = ?", provider, user.ID).First(&userData)

	if userData.ID == 0 {
		newUser := models.User {
			FullName : user.FullName,
			Email    : user.Email,
			SocialId : user.ID,
			Provider : provider,
			Avatar   : user.Avatar,
		}
		config.DB.Create(&newUser)
		return newUser
	} else {
		return userData
	}
}

func createToken(user *models.User) string{
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "user_id": user.ID,
    "user_role": user.Role,
		"exp": time.Now().AddDate(0, 0, 7).Unix(),
		"iat": time.Now().Unix(),
	})

	jwtTokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		fmt.Println(err)
	}

	return jwtTokenString
}
