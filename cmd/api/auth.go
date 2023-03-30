package api

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nuriansyah/services-logbook-mbkm/utils"
	"net/http"
	"time"
)

type LoginDosenReqBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginMahasiswaReqBody struct {
	Nrp      string `json:"nrp" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginSuccessResponse struct {
	Token string `json:"token"`
}

type RegisterMahasiswaReqBody struct {
	Name     string `json:"name" binding:"required"`
	Nrp      string `json:"nrp" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterDosenReqBody struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RegisterSuccessResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
type RequestPassword struct {
	password string `json:"password"`
}
type DosenResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var jwtKey = []byte("key")

type Claims struct {
	Id int
	jwt.StandardClaims
}

func (api *API) getUserIdFromToken(c *gin.Context) (int, error) {
	tokenString := c.GetHeader("Authorization")[(len("Bearer ")):]
	claim := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return -1, err
	}

	if token.Valid {
		claim := token.Claims.(*Claims)
		return claim.Id, nil

	} else {
		return -1, errors.New("invalid token")
	}
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	claim := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, err
}

func (api API) generateJWT(userId *int) (string, error) {
	expTime := time.Now().Add(60 * time.Minute)

	claims := &Claims{
		Id: *userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func (api *API) registerMahasiswa(c *gin.Context) {
	var input RegisterMahasiswaReqBody
	err := c.BindJSON(&input)
	var ve validator.ValidationErrors

	if err != nil {
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"errors": utils.GetErrorMessage(ve)},
			)
			return
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	userId, responseCode, err := api.userRepo.InsertUserMahasiswa(input.Name, input.Nrp, input.Password)
	if err != nil {
		c.AbortWithStatusJSON(responseCode, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := api.generateJWT(&userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RegisterSuccessResponse{Message: "success", Token: tokenString})
}

func (api *API) registerDosen(c *gin.Context) {
	var input RegisterDosenReqBody
	err := c.BindJSON(&input)
	var ve validator.ValidationErrors

	if err != nil {
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"errors": utils.GetErrorMessage(ve)},
			)
			return
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	userId, responseCode, err := api.userRepo.InsertUserDosen(input.Name, input.Email, input.Password)
	if err != nil {
		c.AbortWithStatusJSON(responseCode, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := api.generateJWT(&userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, RegisterSuccessResponse{Message: "success", Token: tokenString})
}

func (api *API) loginDosen(c *gin.Context) {
	var loginReq LoginDosenReqBody
	err := c.BindJSON(&loginReq)
	var ve validator.ValidationErrors

	if err != nil {
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"errors": utils.GetErrorMessage(ve)},
			)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	userId, err := api.userRepo.LoginDosen(loginReq.Email, loginReq.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := api.generateJWT(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginSuccessResponse{Token: tokenString})
}
func (api *API) loginMahasiswa(c *gin.Context) {
	var loginReq LoginMahasiswaReqBody
	err := c.BindJSON(&loginReq)
	var ve validator.ValidationErrors

	if err != nil {
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				gin.H{"errors": utils.GetErrorMessage(ve)},
			)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	userId, err := api.userRepo.LoginMahasiswa(loginReq.Nrp, loginReq.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := api.generateJWT(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, LoginSuccessResponse{Token: tokenString})
}

func (api *API) fetchDataDosen(c *gin.Context) {
	dosenID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	dsn, err := api.userRepo.FetchDataDosen(dosenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMahasiswaResponse{Message: err.Error()})
		return
	}
	if len(dsn) == 0 {
		c.JSON(http.StatusNotFound, ErrorMahasiswaResponse{Message: "Student not found"})
		return
	}

	c.JSON(http.StatusOK, DosenResponse{
		Id:    dosenID,
		Name:  dsn[0].Name,
		Email: dsn[0].Email,
	})
}

func (api *API) changePasswordDosen(c *gin.Context) {
	var req RequestPassword
	userID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Status Unaouthorized"})
		return
	}
	err = api.userRepo.ChangePasswordDosen(userID, req.password)
	c.JSON(http.StatusOK, gin.H{"message": "Change Password Successfully"})
}
func (api *API) changePasswordMahasiswa(c *gin.Context) {
	var req RequestPassword
	userID, err := api.getUserIdFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Status Unaouthorized"})
		return
	}
	err = api.userRepo.ChangePasswordMahasiswa(userID, req.password)
	c.JSON(http.StatusOK, gin.H{"message": "Change Password Successfully"})
}
