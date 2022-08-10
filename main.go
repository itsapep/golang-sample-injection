package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/itsapep/golang-sample-injection/config"
	"github.com/itsapep/golang-sample-injection/model"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.NewConfig()
	db := cfg.DbConn()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	routeEngine := gin.Default()
	routerGroup := routeEngine.Group("/api")
	routerGroup.POST("/auth/login", func(ctx *gin.Context) {
		var login model.Login
		if err := ctx.ShouldBindJSON(&login); err != nil {
			ctx.JSONP(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		var userCred = model.UserCredential{}
		// sql := fmt.Sprintf("select * from user_credential where user_name=%s and user_password=%s", login.User, login.Password)
		sql := "select * from user_credential where user_name=$1 and user_password=$2"
		log.Println("sql: ", sql)

		err := db.Get(&userCred, sql, login.User, login.Password)
		if err != nil {
			ctx.JSONP(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "PONG",
		})
	})

	var apiPort = config.APIConfig{}
	apiPort.APIHost = os.Getenv("API_HOST")
	apiPort.APIPort = os.Getenv("API_PORT")
	listenAddress := fmt.Sprintf("%s:%s", apiPort.APIHost, apiPort.APIPort)

	err := routeEngine.Run(listenAddress)
	if err != nil {
		panic(err)
	}
}
