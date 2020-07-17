package cosmosdb

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logCache = log.ChildLogger(log.RootLogger(), "CosmosDB-connection")
var factory = &CosmosDBConnectorFactory{}

type Settings struct {
	Account    string `md:"account"`
	MasterKey  string `md:"key"`
	APIVersion string `md:"api"`
	Timeout    int    `md:"timeout"`
	Name       string `md:"name"`
}

type CosmosDBConfig struct {
	Account    string
	MasterKey  string
	APIVersion string
	Timeout    int
	Name       string
}

func init() {

	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

type CosmosDBConnectorFactory struct {
}

type CosmosDBSharedConfigManager struct {
	config *CosmosDBConfig
}

func (*CosmosDBConnectorFactory) Type() string {
	return "CosmosDBConnectorFactory"
}

func (*CosmosDBConnectorFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {

	sharedConn := &CosmosDBSharedConfigManager{}
	var err error
	sharedConn.config, err = getClientConfig(settings)
	if err != nil {
		return nil, err
	}
	return sharedConn, nil
}

func (k *CosmosDBSharedConfigManager) Type() string {
	return "CosmosDBServiceConnector"
}

func (k *CosmosDBSharedConfigManager) GetConnection() interface{} {
	return k
}

func (k *CosmosDBSharedConfigManager) GetClientConfiguration() *CosmosDBConfig {
	return k.config
}

func (k *CosmosDBSharedConfigManager) ReleaseConnection(connection interface{}) {

}

func (k *CosmosDBSharedConfigManager) Start() error {
	return nil
}

func (k *CosmosDBSharedConfigManager) Stop() error {

	return nil
}

func GetSharedConfiguration(conn interface{}) (connection.Manager, error) {

	var cManager connection.Manager
	var err error

	cManager, err = coerce.ToConnection(conn)

	if err != nil {
		return nil, err
	}

	return cManager, nil
}

func getClientConfig(settings map[string]interface{}) (*CosmosDBConfig, error) {
	connectionConfig := &CosmosDBConfig{}

	s := &Settings{}

	err := metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}

	connectionConfig.Name = s.Name
	connectionConfig.Account = s.Account
	connectionConfig.MasterKey = s.MasterKey
	connectionConfig.APIVersion = s.APIVersion
	connectionConfig.Timeout = s.Timeout
	return connectionConfig, nil
}
