package connection

import (
	"fmt"
	"github.com/project-flogo/core/support"

	"github.com/project-flogo/core/support/log"
)

var (
	managerFactories = make(map[string]ManagerFactory)
	managers = make(map[string]Manager)
)

func RegisterManagerFactory(factory ManagerFactory) error {

	if factory == nil {
		return fmt.Errorf("cannot register with 'nil' manager factory")
	}

	ref := support.GetRef(factory)

	if _, dup := managerFactories[ref]; dup {
		return fmt.Errorf("manager factory already registered: %s", ref)
	}

	managerFactories[ref] = factory

	log.RootLogger().Debugf("Registering '%s' manager factory: %s", factory.Type, ref )

	return nil
}

func GetManagerFactory(ref string) ManagerFactory {
	return managerFactories[ref]
}

func RegisterManager(connectionId string, manager Manager) error {

	if connectionId == "" {
		return fmt.Errorf("'id' must be specified when registering")
	}

	if manager == nil {
		return fmt.Errorf("cannot register with 'nil' manager")
	}

	if _, dup := managers[connectionId]; dup {
		return fmt.Errorf("manager already registered: %s", connectionId)
	}

	log.RootLogger().Debugf("Registering manager: %s", connectionId)

	managers[connectionId] = manager

	return nil
}

func GetManager(id string) Manager {
	return managers[id]
}

func Managers() map[string]Manager {
	ret := make(map[string]Manager,len(managers) )
	for id, manager := range managers {
		ret[id] = manager
	}

	return ret
}
