package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type TTSHandler struct {
	pool *engine.EnginePool
}

func NewTTSHandler() *TTSHandler {
	return &TTSHandler{
		pool: engine.GetPool(),
	}
}

func (h *TTSHandler) GetEngineList(c *gin.Context) {
	engines := h.pool.ListEngines(protocol.ENGINE_TYPE_TTS)
	resp := protocol.NewEngineListResp(engines)
	c.JSON(http.StatusOK, resp)
}

func (h *TTSHandler) GetDefaultEngine(c *gin.Context) {
	engineDesc := h.pool.GetDefaultEngine(protocol.ENGINE_TYPE_TTS)
	resp := protocol.NewEngineDefaultResp(engineDesc)
	c.JSON(http.StatusOK, resp)
}

func (h *TTSHandler) GetEngineParams(c *gin.Context) {
	engineName := c.Param("engine")

	ttsEngine, err := h.pool.GetTTS(engineName)
	if err != nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResponse("engine not found"))
		return
	}

	params := ttsEngine.(interface{ Parameters() []protocol.ParamDesc }).Parameters()
	resp := protocol.NewEngineParamResp(params)
	c.JSON(http.StatusOK, resp)
}

func (h *TTSHandler) GetVoiceList(c *gin.Context) {
	engineName := c.Param("engine")
	configStr := c.Query("config")

	ttsEngine, err := h.pool.GetTTS(engineName)
	if err != nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResponse("engine not found"))
		return
	}

	config := make(map[string]interface{})
	if configStr != "" {
		json.Unmarshal([]byte(configStr), &config)
	}

	voices, err := ttsEngine.Voices(c.Request.Context(), config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResponse(err.Error()))
		return
	}

	resp := protocol.NewVoiceListResp(voices)
	c.JSON(http.StatusOK, resp)
}

type TTSEngineInput struct {
	Engine string                 `json:"engine"`
	Config map[string]interface{} `json:"config"`
	Data   string                 `json:"data"`
}

func (h *TTSHandler) Infer(c *gin.Context) {
	var input TTSEngineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResponse("invalid request body"))
		return
	}

	ttsEngine, err := h.pool.GetTTS(input.Engine)
	if err != nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResponse("engine not found"))
		return
	}

	textMsg := &protocol.TextMessage{
		Data: input.Data,
	}

	audioMsg, err := ttsEngine.Run(c.Request.Context(), textMsg, input.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResponse(err.Error()))
		return
	}

	resp := &protocol.TTSEngineOutput{
		Code:        protocol.RESPONSE_CODE_OK,
		Message:     "success",
		Data:        string(audioMsg.Data),
		SampleRate:  audioMsg.SampleRate,
		SampleWidth: audioMsg.SampleWidth,
	}
	c.JSON(http.StatusOK, resp)
}
