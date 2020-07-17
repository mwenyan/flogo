package azstorage

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"os"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var factory = &azStorageFactory{}

type Settings struct {
	Name           string `md:"name,required"`
	Description    string `md:"description"`
	ConnectionType string `md:"connectionType"`
	Accountname    string `md:"accountName"`
	SAS            string `md:"sas"`
	Accesskey      string `md:"accessKey"`
	ExpiryDate     string `md:"expiryDate"`
	CONNECTORINFO  string `md:"WI_STUDIO_OAUTH_CONNECTOR_INFO"`
	DocsMetadata   string `md:"DocsMetadata"`
}

// type MongodbClientConfig struct {
// 	Database    string
// 	MongoClient *mongo.Client
// }

func init() {
	if os.Getenv(log.EnvKeyLogLevel) == "DEBUG" {
		// Enable debug logs for sarama lib
		// sarama.Logger = dlog.New(os.Stderr, "[flogo-mongodb]", dlog.LstdFlags)
		// todo
	}
	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

type azStorageFactory struct {
}

func (*azStorageFactory) Type() string {
	return "azstorage"
}

func (*azStorageFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {
	sharedConn := &AzStorageSharedConfigManager{}
	var err error
	sharedConn.config, err = getazstorageConfig(settings)
	if err != nil {
		return nil, err
	}
	return sharedConn, nil
}

type AzStorageSharedConfigManager struct {
	config *Settings
	name   string
}

func (k *AzStorageSharedConfigManager) Type() string {
	return "azstorage"
}

func (k *AzStorageSharedConfigManager) GetConnection() interface{} {
	return k
}
func (k *AzStorageSharedConfigManager) GetClient() {
	return
}

func (k *AzStorageSharedConfigManager) GetClientConfiguration() *Settings {
	return k.config
}

func (k *AzStorageSharedConfigManager) ReleaseConnection(connection interface{}) {

}

func (k *AzStorageSharedConfigManager) Start() error {
	return nil
}

func (k *AzStorageSharedConfigManager) Stop() error {

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

func getazstorageConfig(settings map[string]interface{}) (*Settings, error) {
	connectionConfig := &Settings{}
	err := metadata.MapToStruct(settings, connectionConfig, false)
	if err != nil {
		return nil, err
	}
	if connectionConfig.ConnectionType == "Generate SAS" {
		stringToSign := connectionConfig.Accountname + "\n" + "rwdlacup" + "\n" + "bfqt" + "\n" + "sco" + "\n" + ""
		stringToSign = stringToSign + "\n" + connectionConfig.ExpiryDate
		stringToSign = stringToSign + "\n" + ""
		stringToSign = stringToSign + "" + "\n" + "https" + "\n" + "2017-11-09" + "\n"
		signingKey := connectionConfig.Accesskey
		decodekey, _ := base64.StdEncoding.DecodeString(signingKey)
		h := hmac.New(sha256.New, decodekey)
		h.Write([]byte(stringToSign))
		hashInBase64 := base64.StdEncoding.EncodeToString(h.Sum(nil))
		connectionConfig.SAS = "?sv=2017-11-09&ss=bfqt&srt=sco&sp=rwdlacup&se=" + connectionConfig.ExpiryDate + "&spr=https&sig=" + url.QueryEscape(hashInBase64)
	}
	return connectionConfig, nil
}
