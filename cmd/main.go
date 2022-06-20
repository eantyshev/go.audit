package main

import (
	"os"
	"strings"

	"github.com/eantyshev/go.audit/domain/event"
	"github.com/eantyshev/go.audit/domain/event/repository"
	"github.com/eantyshev/go.audit/domain/event/repository/mongodb"
	"github.com/eantyshev/go.audit/handlers"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfig(cfgPath string) {
	viper.SetEnvPrefix("audit_api")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigFile(cfgPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func chooseRepo(storageType, mongoUri string) (repo event.Repository) {
	switch storageType {
	case "inmem":
		repo = &repository.MemRepo{}
	case "mongodb":
		repo = mongodb.MakeMongoRepo(mongoUri)
	default:
		panic("unknown storage: " + storageType)
	}

	return repo
}

func main() {
	configPath := pflag.StringP("config", "c", "config.yaml", "configuration file")
	storageType := pflag.StringP("storage_type", "s", "mongodb", "storage: inmem or mongodb")
	pflag.Parse()

	LoadConfig(*configPath)
	addrPort := viper.GetString("http.listen")
	timeout := viper.GetDuration("http.timeout")
	mongoUri := viper.GetString("mongo.uri")

	api_key, ok := os.LookupEnv("AUDIT_API_KEY")
	if !ok {
		panic("No AUDIT_API_KEY env")
	}

	if storageType == nil {
		panic("storage is nil")
	}
	repo := chooseRepo(*storageType, mongoUri)
	defer repo.Close()

	srv := handlers.MakeAuditApi(addrPort, timeout, api_key, repo)
	srv.ServeForever()
}
