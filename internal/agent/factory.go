package agent

func NewAgent(name string, desc string, config map[string]interface{}) Agent {
	switch name {
	case "RepeaterAgent":
		return NewRepeaterAgent()
	case "OpenAIAgent":
		cfg := parseOpenAIConfig(config)
		return NewOpenAIAgent(cfg)
	case "DifyAgent":
		cfg := parseDifyConfig(config)
		return NewDifyAgent(cfg)
	case "CozeAgent":
		cfg := parseCozeConfig(config)
		return NewCozeAgent(cfg)
	case "FastGPTAgent":
		cfg := parseFastGPTConfig(config)
		return NewFastGPTAgent(cfg)
	default:
		return nil
	}
}

func parseOpenAIConfig(config map[string]interface{}) *OpenAIConfig {
	cfg := &OpenAIConfig{
		Model: "gpt-3.5-turbo",
	}
	if v, ok := config["base_url"].(string); ok {
		cfg.BaseURL = v
	}
	if v, ok := config["api_key"].(string); ok {
		cfg.APIKey = v
	}
	if v, ok := config["model"].(string); ok {
		cfg.Model = v
	}
	return cfg
}

func parseDifyConfig(config map[string]interface{}) *DifyConfig {
	cfg := &DifyConfig{}
	if v, ok := config["api_server"].(string); ok {
		cfg.APIServer = v
	}
	if v, ok := config["api_key"].(string); ok {
		cfg.APIKey = v
	}
	if v, ok := config["username"].(string); ok {
		cfg.Username = v
	}
	return cfg
}

func parseCozeConfig(config map[string]interface{}) *CozeConfig {
	cfg := &CozeConfig{}
	if v, ok := config["token"].(string); ok {
		cfg.Token = v
	}
	if v, ok := config["bot_id"].(string); ok {
		cfg.BotID = v
	}
	return cfg
}

func parseFastGPTConfig(config map[string]interface{}) *FastGPTConfig {
	cfg := &FastGPTConfig{}
	if v, ok := config["base_url"].(string); ok {
		cfg.BaseURL = v
	}
	if v, ok := config["api_key"].(string); ok {
		cfg.APIKey = v
	}
	if v, ok := config["uid"].(string); ok {
		cfg.UID = v
	}
	return cfg
}
