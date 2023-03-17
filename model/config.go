package model

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	Config = &AppConfig{}
)

type AppConfig struct {
	Version int `json:"version"`
	// Such as "300ms", "-1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// Values less than 1 minute are not allowed.
	Interval string `json:"interval"`
	Rules    []Rule `json:"rules"`
	Pushes   []Push `json:"pushes"`
}

func ReadAppConfig() error {
	if !utils.Exist(res.AppConfigPath) {
		configBytes, err := json.MarshalIndent(DefaultappConfig, "", "\t")
		if err != nil {
			utils.Error("[CONFIG] marshal default app config failed: %v", err)
			return err
		}
		err = os.WriteFile(res.AppConfigPath, configBytes, 0644)
		if err != nil {
			utils.Error("[CONFIG] write default app config failed: %v", err)
			return err
		}
		Config = DefaultappConfig
		return nil
	}

	configBytes, err := os.ReadFile(res.AppConfigPath)
	if err != nil {
		utils.Error("[CONFIG] read app config failed: %v", err)
		return err
	}
	err = json.Unmarshal(configBytes, Config)
	if err != nil {
		utils.Error("[CONFIG] unmarshal app config failed: %v", err)
	} else if Config.Version < DefaultappConfig.Version {
		utils.Warn("[CONFIG] app config version is too old, please update it")
	}
	return err
}

func GetInterval() time.Duration {
	ac := DefaultappConfig
	if Config != nil {
		ac = Config
	}
	d, err := time.ParseDuration(ac.Interval)
	if err == nil {
		if d < res.DefaultInterval {
			utils.Warn("[CONFIG] interval is too short, use default interval: 1m")
			return res.DefaultInterval
		}
		return d
	}
	utils.Warn("[CONFIG] parse interval failed: %v, use default interval: 1m", err)
	return res.DefaultInterval
}

var (
	defaultWekhookBody = map[string]interface{}{
		"action": "send_group_msg",
		"params": map[string]interface{}{
			"group_id": 123456789,
			"message":  "ServerBox Notification\n{{key}}: {{value}}",
		},
	}
	defaultWekhookBodyBytes, _ = json.Marshal(defaultWekhookBody)
	defaultWebhookIface        = PushWebhook{
		Name: "QQ Group",
		Url:  "http://localhost:5700",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_SECRET",
		},
		Method: "POST",
		Body:   defaultWekhookBodyBytes,
	}
	defaultWebhookIfaceBytes, _ = json.Marshal(defaultWebhookIface)

	defaultIOSIface = PushIOS{
		Name:    "My iPhone",
		Token:   "YOUR_TOKEN",
		Title:   "Server Notification",
		Content: "{{key}}: {{value}}",
	}
	defaultIOSIfaceBytes, _ = json.Marshal(defaultIOSIface)

	DefaultappConfig = &AppConfig{
		Version:  1,
		Interval: "1m",
		Rules: []Rule{
			{
				MonitorType: MonitorTypeCPU,
				Threshold:   ">=77%",
				Matcher:     "0",
			},
			{
				MonitorType: MonitorTypeNetwork,
				Threshold:   ">=7.7m/s",
				Matcher:     "eth0",
			},
		},
		Pushes: []Push{
			{
				Type:      PushTypeWebhook,
				Iface:     defaultWebhookIfaceBytes,
				BodyRegex: ".*",
				Code:      200,
			},
			{
				Type:      PushTypeIOS,
				Iface:     defaultIOSIfaceBytes,
				BodyRegex: ".*",
				Code:      200,
			},
		},
	}
)
