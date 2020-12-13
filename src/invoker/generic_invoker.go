package invoker

import (
	"config"
	"context"
	"errors"

	dubboConfig "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	log "github.com/sirupsen/logrus"
)

const (
	APP_NAME = "COMPARE-GO"
	CLUSTER  = "failover"
	REGISTRY = "hangzhouzk"
)

const DATASOURCE_METHOD_GET = "Get"
const EMPTY_STRING = ""

/**

 */
func Invoker(config *config.Config, param map[string]interface{}) (interface{}, error) {
	if !isConfigValid(config) {
		log.Warn("config is unvalid config:%v", config)
		return nil, errors.New("config unValid")
	}
	interfaceName := convertI2S(config.Reference["interfaceName"])
	var referenceConfig = dubboConfig.ReferenceConfig{
		InterfaceName:  interfaceName,
		Cluster:        CLUSTER,
		Registry:       REGISTRY,
		Protocol:       dubbo.DUBBO,
		Generic:        true,
		RequestTimeout: "1000",
	}

	//appName is the unique identification of RPCService
	referenceConfig.GenericLoad(APP_NAME)
	resut, err := referenceConfig.GetRPCService().(*dubboConfig.GenericService).Invoke(context.TODO(), []interface{}{DATASOURCE_METHOD_GET, param})
	if err != nil {
		log.Error("generic invoke error happend, referenceConfig:%v, param:%v", referenceConfig, param)
		return resut, err
	}
	return resut, nil
}

func isConfigValid(config *config.Config) bool {
	if config == nil || len(config.Reference) == 0 {
		return false
	}
	return true
}

//if val is nil and type don't match return ''
func convertI2S(val interface{}) string {
	if val == nil {
		return EMPTY_STRING
	}
	result, ok := val.(string)
	if ok {
		return result
	}
	return EMPTY_STRING
}
