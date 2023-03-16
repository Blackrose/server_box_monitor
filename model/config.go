package model

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

type AppConfig struct {
	// such as "300ms", "-1.5h" or "2h45m". 
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	ScriptRunInterval string `json:"script_run_interval"`
	Rules			 []Rule `json:"rules"`
	Pushes []Push `json:"pushes"`
}

var (
	DefaultappConfig = &AppConfig{
		ScriptRunInterval: "3m",
		Rules: []Rule{
			{
				MonitorType: MonitorTypeCPU,
				Threshold: ">=80%",
				Matcher: "0",
			},
			{
				MonitorType: MonitorTypeNetwork,
				Threshold: ">=17.7m/s",
				Matcher: "eth0",
			},
			{
				MonitorType: MonitorTypeDisk,
				Threshold: ">=95.2%",
				Matcher: "sda1",
			},
		},
		Pushes: []Push{
			{
				PushType: PushTypeWebhook,
				PushIface: &PushWebhook{
					Url: "http://httpbin.org/post",
					Headers: map[string]string{"Content-Type": "application/json"},
					Method: "POST",
				},
				TitleFormat: "[ServerBox] Notification",
				ContentFormat: "{{key}}: {{value}}",
			},
		},
	}
)

func ReadAppConfig() (*AppConfig, error) {
	if !utils.Exist(res.AppConfigPath) {
		configBytes, err := json.Marshal(DefaultappConfig)
		if err != nil {
			utils.Error("[CONFIG] marshal default app config failed: %v", err)
			return nil, err
		}
		err = os.WriteFile(res.AppConfigPath, configBytes, 0644)
		if err != nil {
			utils.Error("[CONFIG] write default app config failed: %v", err)
			return nil, err
		}
		return DefaultappConfig, nil
	}
	
	appConfig := &AppConfig{}
	configBytes, err := os.ReadFile(res.AppConfigPath)
	if err != nil {
		utils.Error("[CONFIG] read app config failed: %v", err)
		return nil, err
	}
	err = json.Unmarshal(configBytes, appConfig)
	if err != nil {
		utils.Error("[CONFIG] unmarshal app config failed: %v", err)
	}
	return appConfig, err
}

func (ac *AppConfig) GetRunInterval() time.Duration {
	sri := ac.ScriptRunInterval
	d, err := time.ParseDuration(sri)
	if err == nil {
		return d
	}
	sri = "3m"
	d, err = time.ParseDuration(sri)
	if err == nil {
		return d
	}
	return time.Minute * 3
}