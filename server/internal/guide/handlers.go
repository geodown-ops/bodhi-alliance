// Package guide is the 線上覺行小組 AI 組長: public chat and guided practice, plus the
// knowledge-base admin (upload, review, publish) and persona settings.
package guide

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/geodown-ops/bodhi-alliance/server/internal/auth"
	"github.com/geodown-ops/bodhi-alliance/server/internal/httpx"
)

const maxUpload = 20 << 20

type Handler struct {
	Store *Store
	Auth  *auth.Service
	LLM   LLM // nil when no API key is configured
}

func (h *Handler) Routes(r *gin.RouterGroup) {
	r.GET("/info", h.info)
	r.GET("/scripts", h.listScripts)
	r.POST("/chat", httpx.RateLimit(20, 30*time.Second), h.chat)

	admin := r.Group("/admin", h.Auth.Require(auth.RoleKnowledgeManager, auth.ScopeGuide))
	admin.GET("/documents", h.listDocuments)
	admin.POST("/documents", h.upload)
	admin.GET("/documents/:id", h.getDocument)
	admin.PUT("/documents/:id", h.updateDocument)
	admin.GET("/documents/:id/versions", h.versions)
	admin.POST("/documents/:id/status", h.setStatus)
	admin.POST("/try", h.try)
	admin.GET("/persona", h.persona)
	admin.PUT("/persona", h.savePersona)
	admin.GET("/usage", h.usage)
	admin.PUT("/settings", h.saveSettings)
}

func (h *Handler) info(c *gin.Context) {
	p, err := h.Store.Persona(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	u, err := h.Store.Usage(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": p.Name, "available": h.available(u)})
}

func (h *Handler) available(u MonthUsage) bool {
	return h.LLM != nil && u.Enabled && u.CostUSD < u.BudgetUSD
}

func (h *Handler) listScripts(c *gin.Context) {
	list, err := h.Store.Scripts(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, list)
}

type chatRequest struct {
	Messages []Message `json:"messages"`
	Mode     string    `json:"mode"`
	ScriptID string    `json:"script_id"`
}

func (h *Handler) chat(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "格式錯誤")
		return
	}
	u, err := h.Store.Usage(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "暫時無法回應")
		return
	}
	if !h.available(u) {
		httpx.Error(c, http.StatusServiceUnavailable, "線上組長目前休息中，請稍後再來，或到「覺行小組」頁面聯絡真人組長。")
		return
	}
	chunks, err := h.Store.PublishedChunks(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "暫時無法回應")
		return
	}
	h.respond(c, req, chunks, nil)
}

// respond runs one turn against the given knowledge chunks and records usage.
// script, when set, is used instead of looking up req.ScriptID (試問 on a draft script).
func (h *Handler) respond(c *gin.Context, req chatRequest, chunks []Chunk, script *Script) {
	msgs, err := CleanHistory(req.Messages)
	if err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	persona, err := h.Store.Persona(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "暫時無法回應")
		return
	}
	turn := Turn{
		Persona:   persona.Prompt,
		Knowledge: BuildKnowledge(SelectKnowledge(chunks, msgs[len(msgs)-1].Content)),
		Messages:  msgs,
	}
	if script != nil {
		turn.Script = script
	} else if req.Mode == modePractice {
		sc, err := h.Store.Script(c, req.ScriptID)
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "請先選一套共修流程")
			return
		}
		turn.Script = &sc
	}
	ctx, cancel := context.WithTimeout(c, 2*time.Minute)
	defer cancel()
	reply, usage, err := h.LLM.Complete(ctx, turn)
	if usage.Model != "" {
		if err := h.Store.RecordUsage(context.WithoutCancel(c), usage.Model, usage.Input, usage.Output, usage.CacheRead, usage.CacheW); err != nil {
			log.Printf("record usage: %v", err)
		}
	}
	if err != nil {
		log.Printf("guide chat: %v", err)
		httpx.Error(c, http.StatusBadGateway, "線上組長一時沒有回應，請再試一次。")
		return
	}
	c.JSON(http.StatusOK, reply)
}

// try answers a test question as the public chat would, optionally including one
// unpublished document so it can be checked before going live. A 共修腳本 is tried
// in practice mode instead of being added to the general knowledge.
func (h *Handler) try(c *gin.Context) {
	var req struct {
		chatRequest
		IncludeDocumentID string `json:"include_document_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "格式錯誤")
		return
	}
	if h.LLM == nil {
		httpx.Error(c, http.StatusServiceUnavailable, "伺服器還沒有設定 ANTHROPIC_API_KEY")
		return
	}
	chunks, err := h.Store.PublishedChunks(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	var script *Script
	if id := req.IncludeDocumentID; id != "" {
		doc, err := h.Store.GetDocument(c, id)
		if err != nil {
			httpx.Error(c, http.StatusNotFound, err.Error())
			return
		}
		kept := chunks[:0]
		for _, ch := range chunks {
			if ch.DocumentID != id {
				kept = append(kept, ch)
			}
		}
		chunks = kept
		if doc.Category == "script" {
			script = &Script{ID: doc.ID, Title: doc.Title, Body: doc.Body}
		} else {
			for _, ch := range SplitChunks(doc.Body) {
				ch.DocumentID, ch.Title, ch.Category = doc.ID, doc.Title, doc.Category
				chunks = append(chunks, ch)
			}
		}
	}
	h.respond(c, req.chatRequest, chunks, script)
}

func (h *Handler) listDocuments(c *gin.Context) {
	docs, err := h.Store.ListDocuments(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, gin.H{"documents": docs, "categories": categories})
}

func validDoc(title, category string) error {
	if strings.TrimSpace(title) == "" || utf8.RuneCountInString(title) > 200 {
		return errors.New("請填寫標題（200 字以內）")
	}
	if _, ok := categories[category]; !ok {
		return errors.New("分類不正確")
	}
	return nil
}

// upload accepts a file (multipart "file") or pasted text ("body"), stored as a draft.
func (h *Handler) upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUpload+1<<20)
	title, category := strings.TrimSpace(c.PostForm("title")), c.PostForm("category")
	body, source := c.PostForm("body"), "貼上的文字"
	if fh, err := c.FormFile("file"); err == nil {
		if fh.Size > maxUpload {
			httpx.Error(c, http.StatusRequestEntityTooLarge, "檔案超過 20 MB")
			return
		}
		f, err := fh.Open()
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "讀取檔案失敗")
			return
		}
		data, err := io.ReadAll(io.LimitReader(f, maxUpload))
		f.Close()
		if err != nil {
			httpx.Error(c, http.StatusBadRequest, "讀取檔案失敗")
			return
		}
		if body, err = ExtractText(fh.Filename, data); err != nil {
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		source = fh.Filename
		if title == "" {
			title = strings.TrimSuffix(fh.Filename, fileExt(fh.Filename))
		}
	} else {
		body = normalize(body)
	}
	if strings.TrimSpace(body) == "" {
		httpx.Error(c, http.StatusBadRequest, "請上傳檔案或貼上文字")
		return
	}
	if err := validDoc(title, category); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.Store.CreateDocument(c, title, category, source, body, auth.CurrentUser(c).ID)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func fileExt(name string) string {
	if i := strings.LastIndex(name, "."); i > 0 {
		return name[i:]
	}
	return ""
}

func (h *Handler) getDocument(c *gin.Context) {
	d, err := h.Store.GetDocument(c, c.Param("id"))
	if errors.Is(err, errNotFound) {
		httpx.Error(c, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *Handler) updateDocument(c *gin.Context) {
	var req struct {
		Title    string `json:"title"`
		Category string `json:"category"`
		Body     string `json:"body"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "格式錯誤")
		return
	}
	if err := validDoc(req.Title, req.Category); err != nil {
		httpx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.Store.UpdateDocument(c, c.Param("id"), strings.TrimSpace(req.Title), req.Category,
		normalize(req.Body), "後台編輯", auth.CurrentUser(c).ID)
	if errors.Is(err, errNotFound) {
		httpx.Error(c, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) versions(c *gin.Context) {
	v, err := h.Store.Versions(c, c.Param("id"))
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *Handler) setStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required,oneof=draft published archived"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "狀態不正確")
		return
	}
	err := h.Store.SetStatus(c, c.Param("id"), req.Status)
	if errors.Is(err, errNotFound) {
		httpx.Error(c, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "更新失敗")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) persona(c *gin.Context) {
	list, err := h.Store.PersonaHistory(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) savePersona(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required,max=50"`
		Prompt string `json:"prompt" binding:"required,max=20000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "請填寫名字與角色設定")
		return
	}
	id := auth.CurrentUser(c).ID
	if err := h.Store.SavePersona(c, strings.TrimSpace(req.Name), strings.TrimSpace(req.Prompt), &id); err != nil {
		httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) usage(c *gin.Context) {
	u, err := h.Store.Usage(c)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "讀取失敗")
		return
	}
	c.JSON(http.StatusOK, gin.H{"usage": u, "configured": h.LLM != nil, "available": h.available(u)})
}

func (h *Handler) saveSettings(c *gin.Context) {
	var req struct {
		MonthlyBudgetUSD *float64 `json:"monthly_budget_usd"`
		Enabled          *bool    `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Error(c, http.StatusBadRequest, "格式錯誤")
		return
	}
	if req.MonthlyBudgetUSD != nil {
		if *req.MonthlyBudgetUSD < 0 {
			httpx.Error(c, http.StatusBadRequest, "預算不能是負數")
			return
		}
		if err := h.Store.SetSetting(c, "monthly_budget_usd", strconv.FormatFloat(*req.MonthlyBudgetUSD, 'f', 2, 64)); err != nil {
			httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
			return
		}
	}
	if req.Enabled != nil {
		if err := h.Store.SetSetting(c, "chat_enabled", strconv.FormatBool(*req.Enabled)); err != nil {
			httpx.Error(c, http.StatusInternalServerError, "儲存失敗")
			return
		}
	}
	c.Status(http.StatusNoContent)
}
