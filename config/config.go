// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

// Config contains the DropbeatConfig
type Config struct {
	Dropbeat DropbeatConfig
}

// DropbeatConfig contains the dropbeat data
type DropbeatConfig struct {
	Period string `config:"period"`

	URLs []string

	Stats struct {
		Metrics *bool
		Health  *bool
	}
}

var DefaultConfig = Config{
	Dropbeat: DropbeatConfig{
		Period: "10ms",
	},
}
