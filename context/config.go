package context

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config for app
type Config struct {
	Daazweb struct {
		BaseURL    string `mapstructure:"base_url"`
		PaymentURL string `mapstructure:"payment_url"`
		Limit      string `mapstructure:"limit"`
		Operator   string `mapstructure:"operator"`
		Token      string `mapstructure:"token"`
	} `mapstructure:"daazweb"`
	SMTP struct {
		From     string `mapstructure:"from"`
		Password string `mapstructure:"password"`
		To       string `mapstructure:"to"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
	} `mapstructure:"smtp"`
	Telegram struct {
		Token       string  `mapstructure:"token"`
		AccessToken string  `mapstructure:"access_token"`
		Chats       []int64 `mapstructure:"chats"`
	} `mapstructure:"telegram"`
	SendSMTPTest bool `mapstructure:"send_smtp_test"`
	v            *viper.Viper
}

// InitConfig from file
func InitConfig() (*Config, error) {
	var cfg Config

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read the configuration file: %+v", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("cannot unmarshal config: %+v", err)
	}

	cfg.v = v

	return &cfg, nil
}

// Save new data to file
func (c *Config) Save() error {
	r := map[string]interface{}{}
	if err := mapstructure.Decode(c, &r); err != nil {
		return err
	}

	if err := c.v.MergeConfigMap(r); err != nil {
		return err
	}

	return c.v.WriteConfig()
}
