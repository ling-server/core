package config

import "sync"

var metaDataOnce sync.Once
var metaDataInstance *ConfigMetaData

// Item - Configure item include default value, type, env name
type Item struct {
	// The Scope of this configuration item: eg: SystemScope, UserScope
	Scope string `json:"scope,omitempty"`
	// email, ldapbasic, ldapgroup, uaa settings, used to retieve configure items by group
	Group string `json:"group,omitempty"`
	// environment key to retrieves this value when initialize, for example: POSTGRESQL_HOST, only used for system settings, for user settings no EnvKey
	EnvironmentKey string `json:"environment_key,omitempty"`
	// The default string value for this key
	DefaultValue string `json:"default_value,omitempty"`
	// The key for current configure settings in database or rest api
	Name string `json:"name,omitempty"`
	// It can be &IntType{}, &StringType{}, &BoolType{}, &PasswordType{}, &MapType{} etc, any type interface implementation
	ItemType Type
	// Editable means it can updated by configure api, For system configure, the editable is always false, for user configure, it may depends
	Editable bool `json:"editable,omitempty"`
	// Description - Describle the usage of the configure item
	Description string
}

// Instance - Get Instance, make it singleton because there is only one copy of metadata in an env
func Instance() *ConfigMetaData {
	metaDataOnce.Do(func() {
		metaDataInstance = newConfigMetaData()
	})
	return metaDataInstance
}

func newConfigMetaData() *ConfigMetaData {
	return &ConfigMetaData{metaMap: make(map[string]Item)}
}

// ConfigMetaData ...
type ConfigMetaData struct {
	metaMap map[string]Item
}

// initFromArray - Initial metadata from an array
func (c *ConfigMetaData) InitFromArray(items []Item) {
	c.metaMap = make(map[string]Item)
	for _, item := range items {
		c.metaMap[item.Name] = item
	}
}

// Registe - append the item to metadata
func (c *ConfigMetaData) Register(name string, item Item) {
	c.metaMap[name] = item
}

// GetByName - Get current metadata of current name, if not defined, return false in second params
func (c *ConfigMetaData) GetByName(name string) (*Item, bool) {
	if item, ok := c.metaMap[name]; ok {
		return &item, true
	}
	return nil, false
}

// GetAll - Get all metadata in current env
func (c *ConfigMetaData) GetAll() []Item {
	metaDataList := make([]Item, 0)
	for _, value := range c.metaMap {
		metaDataList = append(metaDataList, value)
	}
	return metaDataList
}
