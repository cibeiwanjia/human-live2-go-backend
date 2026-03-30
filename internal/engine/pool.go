package engine

import (
	"errors"
	"sync"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/config"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/tts"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

var (
	ErrEngineNotFound = errors.New("engine not found")
)

var (
	pool     *EnginePool
	poolOnce sync.Once
)

type EnginePool struct {
	mu         sync.RWMutex
	ttsEngines map[string]TTSEngine
	asrEngines map[string]ASREngine
	ttsDefault string
	asrDefault string
}

func GetPool() *EnginePool {
	poolOnce.Do(func() {
		pool = &EnginePool{
			ttsEngines: make(map[string]TTSEngine),
			asrEngines: make(map[string]ASREngine),
		}
	})
	return pool
}

func (p *EnginePool) Setup(cfg *config.EnginesConfig) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, engineCfg := range cfg.TTS.SupportList {
		var engine TTSEngine
		switch engineCfg.Name {
		case "EdgeTTS":
			engine = tts.NewEdgeTTS()
		case "TencentTTS":
			engine = tts.NewTencentTTS(nil)
		case "DifyTTS":
			engine = tts.NewDifyTTS(nil)
		case "CozeTTS":
			engine = tts.NewCozeTTS(nil)
		}
		if engine != nil {
			p.ttsEngines[engine.Name()] = engine
		}
	}
	p.ttsDefault = cfg.TTS.Default

	p.asrDefault = cfg.ASR.Default

	return nil
}

func (p *EnginePool) GetTTS(name string) (TTSEngine, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if name == "default" || name == "" {
		name = p.ttsDefault
	}

	engine, ok := p.ttsEngines[name]
	if !ok {
		return nil, ErrEngineNotFound
	}
	return engine, nil
}

func (p *EnginePool) GetASR(name string) (ASREngine, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if name == "default" || name == "" {
		name = p.asrDefault
	}

	engine, ok := p.asrEngines[name]
	if !ok {
		return nil, ErrEngineNotFound
	}
	return engine, nil
}

func (p *EnginePool) ListTTS() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	names := make([]string, 0, len(p.ttsEngines))
	for name := range p.ttsEngines {
		names = append(names, name)
	}
	return names
}

func (p *EnginePool) ListASR() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	names := make([]string, 0, len(p.asrEngines))
	for name := range p.asrEngines {
		names = append(names, name)
	}
	return names
}

func (p *EnginePool) TTSDefault() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ttsDefault
}

func (p *EnginePool) ASRDefault() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.asrDefault
}

func (p *EnginePool) ListEngines(engineType protocol.ENGINE_TYPE) []protocol.EngineDesc {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var engines []protocol.EngineDesc

	switch engineType {
	case protocol.ENGINE_TYPE_TTS:
		for _, engine := range p.ttsEngines {
			engines = append(engines, engine.Desc())
		}
	case protocol.ENGINE_TYPE_ASR:
		for _, engine := range p.asrEngines {
			engines = append(engines, engine.Desc())
		}
	}

	return engines
}

func (p *EnginePool) GetDefaultEngine(engineType protocol.ENGINE_TYPE) protocol.EngineDesc {
	p.mu.RLock()
	defer p.mu.RUnlock()

	switch engineType {
	case protocol.ENGINE_TYPE_TTS:
		if engine, ok := p.ttsEngines[p.ttsDefault]; ok {
			return engine.Desc()
		}
	case protocol.ENGINE_TYPE_ASR:
		if engine, ok := p.asrEngines[p.asrDefault]; ok {
			return engine.Desc()
		}
	}

	return protocol.EngineDesc{}
}
