package cmd

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	teacherProfileHTTP "github.com/fidellr/edu_malay/teacher/delivery/http"
	teacherProfileRepo "github.com/fidellr/edu_malay/teacher/repository"
	teacherProfileServices "github.com/fidellr/edu_malay/teacher/usecase"
)

var eduServerCMD = &cobra.Command{
	Use:   "http",
	Short: "Start http server for edu_malay",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

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
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCMD.AddCommand(eduServerCMD)
	eduServerCMD.PersistentFlags().String("config", "", "Set this flag to use a configuration file")
}

func initConfig() {
	configFile := ""
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if eduServerCMD.Flags().Lookup("config") != nil {
		configFile = "config.json"
		viper.BindPFlag("config", eduServerCMD.Flags().Lookup("config"))
		viper.SetConfigType("json")
	}

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
	// timeout := time.Duration(viper.GetInt("http.timeout")) * time.Second
	contextTimeout := time.Duration(viper.GetInt("context.timeout")) * time.Second
	// defaultTransport := &http.Transport{
	// 	MaxIdleConns:        viper.GetInt("http.max_idle_conns"),
	// 	MaxIdleConnsPerHost: viper.GetInt("http.max_idle_conns_per_host"),
	// 	IdleConnTimeout:     time.Duration(viper.GetInt("http.max_idle_conn_timeout")) * time.Second,
	// }

	mongoDSN := viper.GetString("mongo.dsn")

	masterSession, err := mgo.Dial(mongoDSN)
	if err != nil {
		logrus.Fatalln(err.Error())
	}

	masterSession.SetSafe(&mgo.Safe{})

	mongoDatabase := viper.GetString("mongo.database")
	if mongoDatabase == "" {
		logrus.Fatalln(errors.New("Please provide a mongo database name"))
	}

	tpRepo := teacherProfileRepo.NewTeacherMongo(masterSession, mongoDatabase)
	tpServices := teacherProfileServices.NewTeacherProfileUsecase(tpRepo, contextTimeout)
	teacherProfileHTTP.NewTeacherProfileHandler(e, tpServices)
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