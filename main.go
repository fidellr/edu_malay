package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	pictHTTP "github.com/fidellr/edu_malay/pict/delivery/http"
	pictRepo "github.com/fidellr/edu_malay/pict/repository"
	pictServices "github.com/fidellr/edu_malay/pict/usecase"

	teacherProfileHTTP "github.com/fidellr/edu_malay/teacher/delivery/http"
	teacherProfileRepo "github.com/fidellr/edu_malay/teacher/repository"
	teacherProfileServices "github.com/fidellr/edu_malay/teacher/usecase"

	clcProfileHTTP "github.com/fidellr/edu_malay/clc/delivery/http"
	clcProfileRepo "github.com/fidellr/edu_malay/clc/repository"
	clcProfileServices "github.com/fidellr/edu_malay/clc/usecase"

	assemblerProfileHTTP "github.com/fidellr/edu_malay/assembler/profile/delivery/http"
	assemblerProfileRepo "github.com/fidellr/edu_malay/assembler/profile/repository"
	assemblerProfileServices "github.com/fidellr/edu_malay/assembler/profile/usecase"
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

	db := masterSession.DB(mongoDatabase)
	colTeacher := db.C("teacher")
	teacherIndex := mgo.Index{
		Key: []string{"$text:first_name", "$text:last_name", "$text:university"},
		Weights: map[string]int{
			"first_name": 9,
			"last_name":  8,
			"university": 3,
		},
		Name: "teacherIndex",
	}
	if idx, err := colTeacher.Indexes(); len(idx) > 0 {
		if err != nil {
			panic(err)
		}
		colTeacher.DropAllIndexes()
	}
	err = colTeacher.EnsureIndex(teacherIndex)

	colCLC := db.C("clc")
	clcIndex := mgo.Index{
		Key: []string{"$text:name", "$text:clc_level", "$text:status"},
		Weights: map[string]int{
			"name":      9,
			"clc_level": 8,
			"status":    3,
		},
		Name: "clcIndex",
	}

	if idx, err := colCLC.Indexes(); len(idx) > 0 {
		if err != nil {
			panic(err)
		}

		colCLC.DropAllIndexes()
	}
	err = colCLC.EnsureIndex(clcIndex)

	if err != nil {
		panic(err)
	}

	// masterS3Cfg := &aws.Config{
	// 	Region: aws.String("US"),
	// 	Credentials: credentials.NewStaticCredentials(
	// 		"GELA5TBFIZLBWZXT648T",
	// 		"yrj9EDq1oZPLQIrM04GsSdpyu7dB0ePkzcTuneBi",
	// 		"",
	// 	),
	// }

	// masterS3Session, err := session.NewSession(masterS3Cfg)
	// if err != nil {
	// 	panic(err)
	// }

	endpointHost := "cellar-c2.services.clever-cloud.com"
	awsAccessKeyID := "GELA5TBFIZLBWZXT648T"
	awsSecretKey := "yrj9EDq1oZPLQIrM04GsSdpyu7dB0ePkzcTuneBi"
	minioClient, err := minio.New(endpointHost, awsAccessKeyID, awsSecretKey, true)
	if err != nil {
		panic(err)
	}

	picRepo := pictRepo.NewS3Pict(minioClient)
	picServices := pictServices.NewPictUsecase(picRepo, contextTimeout)
	pictHTTP.NewPictHandler(e, picServices)

	tpRepo := teacherProfileRepo.NewTeacherMongo(masterSession, mongoDatabase)
	tpServices := teacherProfileServices.NewTeacherProfileUsecase(tpRepo, contextTimeout)
	teacherProfileHTTP.NewTeacherProfileHandler(e, tpServices)

	clcRepo := clcProfileRepo.NewClcProfileMongo(masterSession, mongoDatabase)
	clcServices := clcProfileServices.NewClcProfileUsecase(clcRepo, contextTimeout)
	clcProfileHTTP.NewClcProfileHandler(e, clcServices)

	assmblrProfileRepo := assemblerProfileRepo.NewProfileAssemblerMongo(masterSession, mongoDatabase)
	assmblrProfileServices := assemblerProfileServices.NewProfileAssemblerUsecase(assmblrProfileRepo, contextTimeout)
	assemblerProfileHTTP.NewProfileAssemblerHandler(e, assmblrProfileServices)

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

	echoInstance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://edu-malaysia.now.sh"},
	}))

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
