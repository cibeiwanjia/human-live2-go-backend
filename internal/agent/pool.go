package agent

import (
	"errors"
	"sync"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/config"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
)

var (
	pool     *AgentPool
	poolOnce sync.Once
)

// AgentPool manages all registered agents
type AgentPool struct {
	mu          sync.RWMutex
	agents      map[string]Agent
	defaultName string
}

// GetPool returns the singleton agent pool
func GetPool() *AgentPool {
	poolOnce.Do(func() {
		pool = &AgentPool{
			agents: make(map[string]Agent),
		}
	})
	return pool
}

// Setup initializes the agent pool with configuration
func (p *AgentPool) Setup(cfg *config.AgentsConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, agentCfg := range cfg.SupportList {
		agent := NewAgent(agentCfg.Name, agentCfg.Desc, agentCfg.Config)
		if agent != nil {
			p.agents[agent.Name()] = agent
		}
	}
	p.defaultName = cfg.Default

	return nil
}

// Register adds an agent to the pool
func (p *AgentPool) Register(agent Agent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.agents[agent.Name()] = agent
}

// Get returns an agent by name
func (p *AgentPool) Get(name string) (Agent, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if name == "default" || name == "" {
		name = p.defaultName
	}

	agent, ok := p.agents[name]
	if !ok {
		return nil, ErrAgentNotFound
	}
	return agent, nil
}

// List returns all agent names
func (p *AgentPool) List() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	names := make([]string, 0, len(p.agents))
	for name := range p.agents {
		names = append(names, name)
	}
	return names
}

// Default returns the default agent name
func (p *AgentPool) Default() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.defaultName
}
