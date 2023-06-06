package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"todoGin/config"
	"todoGin/model/respErr"
)

// secret key untuk signing token
// middleware konsep nya adalah sesuatu yang ibaratnya intercept , request -> server,
func Authmiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// mengambil token dari header Authorization
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
				Message: "Unauthorized",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// split token dari header
		tokenString := authHeader[len("Bearer "):]

		// parsing token dengan secret key
		claims := &config.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
					Message: "Unauthorized",
					Status:  http.StatusUnauthorized,
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, &respErr.ErrorResponse{
				Message: "invalid or expired token",
				Status:  http.StatusBadRequest,
			})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, &respErr.ErrorResponse{
				Message: "Unauthorized (non Valid)",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// token valid, ambil username dari claim dan simpan ke dalam konteks
		ctx.Set("username", claims.Username)

		// token valid, melanjutkan ke handler
		ctx.Next()

	}
}
