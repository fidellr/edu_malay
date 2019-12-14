package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	teacherProfileHTTP "github.com/fidellr/edu_malay/teacher/delivery/http"
	teacherProfileRepo "github.com/fidellr/edu_malay/teacher/repository"
	teacherProfileServices "github.com/fidellr/edu_malay/teacher/usecase"

	clcProfileHTTP "github.com/fidellr/edu_malay/clc/delivery/http"
	clcProfileRepo "github.com/fidellr/edu_malay/clc/repository"
	clcProfileServices "github.com/fidellr/edu_malay/clc/usecase"
)

func init() {
	initConfig()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func initConfig() {
	configFile := ""
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	configFile = "config.json"
	viper.SetConfigType("json")

	if config := viper.GetString("config"); config != "" {
		configFile = config
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalln(err.Error())
	}

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Warn("edu_malay is Running in Debug Mode")
		return
	}

	logrus.SetLevel(logrus.InfoLevel)
	logrus.Warn("edu_malay is Running in Production Mode")

}

func initEduMalayApplication(e *echo.Echo) {
	var mongoDSN, mongoDatabase string
	contextTimeout := time.Duration(viper.GetInt("context.timeout")) * time.Second
	// timeout := time.Duration(viper.GetInt("http.timeout")) * time.Second
	// defaultTransport := &http.Transport{
	// 	MaxIdleConns:        viper.GetInt("http.max_idle_conns"),
	// 	MaxIdleConnsPerHost: viper.GetInt("http.max_idle_conns_per_host"),
	// 	IdleConnTimeout:     time.Duration(viper.GetInt("http.max_idle_conn_timeout")) * time.Second,
	// }

	if !viper.GetBool("debug") {
		mongoDSN = viper.GetString("mongo_prod.dsn")
		mongoDatabase = viper.GetString("mongo_prod.database")
	} else {
		mongoDSN = viper.GetString("mongo.dsn")
		mongoDatabase = viper.GetString("mongo.database")
	}

	masterSession, err := mgo.Dial(mongoDSN)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	masterSession.SetSafe(&mgo.Safe{})

	if mongoDatabase == "" {
		logrus.Fatalln(errors.New("Please provide a mongo database name"))
	}

	tpRepo := teacherProfileRepo.NewTeacherMongo(masterSession, mongoDatabase)
	tpServices := teacherProfileServices.NewTeacherProfileUsecase(tpRepo, contextTimeout)
	teacherProfileHTTP.NewTeacherProfileHandler(e, tpServices)

	clcRepo := clcProfileRepo.NewClcProfileMongo(masterSession, mongoDatabase)
	clcServices := clcProfileServices.NewClcProfileUsecase(clcRepo, contextTimeout)
	clcProfileHTTP.NewClcProfileHandler(e, clcServices)
	// userRepo := _mongoRepository.NewUserMongo(
	// 	_mongoRepository.UserSession(masterSession),
	// 	_mongoRepository.UserDBName(mongoDatabase),
	// )
	// uranusService := user.NewService(
	// 	user.Repository(userRepo),
	// 	user.Timeout(contextTimeout),
	// 	user.Validator(validator),
	// )

	// e.HTTPErrorHandler = delivery.HandleUncaughtHTTPError
	// _httpDelivery.NewUserHandler(e, _httpDelivery.UserService(uranusService))
}

func main() {
	echoInstance := echo.New()
	echoInstance.Server.ReadTimeout = time.Duration(viper.GetInt("http.server_read_timeout")) * time.Second
	echoInstance.Server.WriteTimeout = time.Duration(viper.GetInt("http.server_write_timeout")) * time.Second

	echoInstance.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong user")
	})

	initEduMalayApplication(echoInstance)

	address := viper.GetString("server.address")
	if err := echoInstance.Start(address); err != nil {
		logrus.Fatalln(err.Error())
	}

	logrus.Infof("Start listening on: %s", address)
}
