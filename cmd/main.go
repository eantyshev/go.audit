package main

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.audit/internal/api"
	"go.audit/internal/repository"
	"go.audit/internal/repository/mongodb"
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

func chooseRepo(storageType, mongoUri string) (repo repository.RepoIface) {
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

	srv := api.MakeAuditApi(addrPort, timeout, api_key, repo)
	srv.ServeForever()
}
