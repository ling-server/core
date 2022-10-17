package config

import (
	"context"
	"errors"
	"sync"

	"github.com/ling-server/core/config/metadata"
	"github.com/ling-server/core/log"
	"github.com/ling-server/core/orm"
)

var (
	// DefaultCfgManager the default change manager, default is DBCfgManager.
	// If InMemoryConfigManager is used, need to set to InMemoryCfgManager in
	// test code.
	DefaultConfigManager = DBConfigManager
	managersMU           sync.RWMutex
	managers             = make(map[string]Manager)
)

// Manager defines the operation for config
type Manager interface {
	Load(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{})
	Save(ctx context.Context) error
	Get(ctx context.Context, key string) *metadata.ConfigureValue
	UpdateConfig(ctx context.Context, cfgs map[string]interface{}) error
	GetUserConfigs(ctx context.Context) map[string]interface{}
	ValidateConfig(ctx context.Context, cfgs map[string]interface{}) error
	GetAll(ctx context.Context) map[string]interface{}
}

// Register  register the config manager
func Register(name string, mgr Manager) {
	managersMU.Lock()
	defer managersMU.Unlock()
	if mgr == nil {
		log.Error("Register manager is nil")
	}
	managers[name] = mgr
}

// GetManager get the configure manager by name
func GetManager(name string) (Manager, error) {
	mgr, ok := managers[name]
	if !ok {
		return nil, errors.New("config manager is not registered: " + name)
	}
	return mgr, nil
}

// DefaultMgr get default config manager
func DefaultManager() Manager {
	manager, err := GetManager(DefaultConfigManager)
	if err != nil {
		log.Error("failed to get config manager")
	}
	return manager
}

func GetConfigManager(ctx context.Context) Manager {
	return DefaultManager()
}

func Load(ctx context.Context) error {
	return DefaultManager().Load(ctx)
}

func Upload(cfg map[string]interface{}) error {
	return DefaultManager().UpdateConfig(orm.Context(), cfg)
}
