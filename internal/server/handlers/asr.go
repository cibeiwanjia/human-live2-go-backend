package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type ASRHandler struct {
	pool       *engine.EnginePool
	wsUpgrader websocket.Upgrader
}

func NewASRHandler() *ASRHandler {
	return &ASRHandler{
		pool: engine.GetPool(),
		wsUpgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *ASRHandler) GetEngineList(c *gin.Context) {
	engines := h.pool.ListEngines(protocol.ENGINE_TYPE_ASR)
	resp := protocol.NewEngineListResp(engines)
	c.JSON(http.StatusOK, resp)
}

func (h *ASRHandler) GetDefaultEngine(c *gin.Context) {
	engineDesc := h.pool.GetDefaultEngine(protocol.ENGINE_TYPE_ASR)
	resp := protocol.NewEngineDefaultResp(engineDesc)
	c.JSON(http.StatusOK, resp)
}

func (h *ASRHandler) GetEngineParams(c *gin.Context) {
	engineName := c.Param("engine")

	asrEngine, err := h.pool.GetASR(engineName)
	if err != nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResponse("engine not found"))
		return
	}

	params := asrEngine.(interface{ Parameters() []protocol.ParamDesc }).Parameters()
	resp := protocol.NewEngineParamResp(params)
	c.JSON(http.StatusOK, resp)
}

type ASREngineInput struct {
	Engine      string                 `json:"engine"`
	Config      map[string]interface{} `json:"config"`
	Data        string                 `json:"data"`
	Type        protocol.AUDIO_TYPE    `json:"type"`
	SampleRate  int                    `json:"sampleRate"`
	SampleWidth int                    `json:"sampleWidth"`
}

func (h *ASRHandler) Infer(c *gin.Context) {
	var input ASREngineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResponse("invalid request body"))
		return
	}

	asrEngine, err := h.pool.GetASR(input.Engine)
	if err != nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResponse("engine not found"))
		return
	}

	audioData, err := base64.StdEncoding.DecodeString(input.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResponse("invalid audio data"))
		return
	}

	audioMsg := &protocol.AudioMessage{
		Data:        audioData,
		Type:        input.Type,
		SampleRate:  input.SampleRate,
		SampleWidth: input.SampleWidth,
	}

	textMsg, err := asrEngine.Run(c.Request.Context(), audioMsg, input.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResponse(err.Error()))
		return
	}

	resp := protocol.NewASREngineOutput(string(textMsg.Data))
	c.JSON(http.StatusOK, resp)
}

func (h *ASRHandler) InferFile(c *gin.Context) {
	engine := c.PostForm("engine")
	audioType := c.PostForm("type")
	configStr := c.PostForm("config")
	sampleRate := c.PostForm("sampleRate")
	sampleWidth := c.PostForm("sampleWidth")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.NewErrorResponse("file is required"))
		return
	}
	defer file.Close()

	fileData := make([]byte, 0)
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		fileData = append(fileData, buf[:n]...)
	}

	asrEngine, err := h.pool.GetASR(engine)
	if err != nil {
		c.JSON(http.StatusNotFound, protocol.NewErrorResponse("engine not found"))
		return
	}

	config := make(map[string]interface{})
	if configStr != "" {
		json.Unmarshal([]byte(configStr), &config)
	}

	sr := 16000
	sw := 2
	if sampleRate != "" {
		json.Unmarshal([]byte(sampleRate), &sr)
	}
	if sampleWidth != "" {
		json.Unmarshal([]byte(sampleWidth), &sw)
	}

	audioMsg := &protocol.AudioMessage{
		Data:        fileData,
		Type:        protocol.AUDIO_TYPE(audioType),
		SampleRate:  sr,
		SampleWidth: sw,
	}

	textMsg, err := asrEngine.Run(c.Request.Context(), audioMsg, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.NewErrorResponse(err.Error()))
		return
	}

	resp := protocol.NewASREngineOutput(string(textMsg.Data))
	c.JSON(http.StatusOK, resp)
}

type EngineInput struct {
	Engine string                 `json:"engine"`
	Config map[string]interface{} `json:"config"`
}

func (h *ASRHandler) StreamInfer(c *gin.Context) {
	conn, err := h.wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var engineInput EngineInput
	engineStarted := false

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		action, payload, err := protocol.ParseMessage(data)
		if err != nil {
			h.sendError(conn, "invalid message format")
			continue
		}

		switch protocol.WS_RECV_ACTION_TYPE(action) {
		case protocol.WS_RECV_ACTION_PING:
			h.sendPong(conn)

		case protocol.WS_RECV_ACTION_ENGINE_START:
			if err := json.Unmarshal(payload, &engineInput); err != nil {
				h.sendError(conn, "invalid ENGINE_START payload")
				return
			}
			engineStarted = true
			h.sendAction(conn, string(protocol.WS_SEND_ACTION_ENGINE_STARTED), "")

		case protocol.WS_RECV_ACTION_PARTIAL_INPUT:
			if !engineStarted {
				h.sendError(conn, "engine not started")
				continue
			}

		case protocol.WS_RECV_ACTION_FINAL_INPUT:
			if !engineStarted {
				h.sendError(conn, "engine not started")
				continue
			}
			result := string(payload)
			h.sendAction(conn, string(protocol.WS_SEND_ACTION_FINAL_OUTPUT), result)

		case protocol.WS_RECV_ACTION_ENGINE_STOP:
			h.sendAction(conn, string(protocol.WS_SEND_ACTION_ENGINE_STOPPED), "")
			return
		}
	}
}

func (h *ASRHandler) sendPong(conn *websocket.Conn) {
	msg, _ := protocol.StructMessage(string(protocol.WS_SEND_ACTION_PONG), []byte{})
	conn.WriteMessage(websocket.BinaryMessage, msg)
}

func (h *ASRHandler) sendError(conn *websocket.Conn, errMsg string) {
	msg, _ := protocol.StructMessage(string(protocol.WS_SEND_ACTION_ERROR), []byte(errMsg))
	conn.WriteMessage(websocket.BinaryMessage, msg)
}

func (h *ASRHandler) sendAction(conn *websocket.Conn, action string, payload string) {
	msg, _ := protocol.StructMessage(action, []byte(payload))
	conn.WriteMessage(websocket.BinaryMessage, msg)
}
