package main

var trayTranslations = map[string]map[string]string{
	"en": {
		"title":   "AICoder Dashboard",
		"show":    "Show AICoder",
		"launch":  "Launch Claude Code",
		"quit":    "Quit AICoder",
		"models":  "Providers",
		"actions": "Actions",
	},
	"zh-Hans": {
		"title":   "AICoder 控制台",
		"show":    "显示主窗口",
		"launch":  "启动 Claude Code",
		"quit":    "退出程序",
		"models":  "服务商选择",
		"actions": "操作",
	},
	"zh-Hant": {
		"title":   "AICoder 控制台",
		"show":    "顯示主視窗",
		"launch":  "啟動 Claude Code",
		"quit":    "退出程式",
		"models":  "服務商選擇",
		"actions": "操作",
	},
}

const RequiredNodeVersion = "22.14.0"