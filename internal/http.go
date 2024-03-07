package internal

import (
	"github.com/tinyfluffs/http/memcache/internal/cache"
	"io"
	"log"
	"net/http"
	"strings"
)

func Run(server *http.Server, mc *cache.MemoryCache) {
	mc.Run()
	http.Handle("/{key}", &HttpCacheHandler{mc: mc})

	log.Printf("listening on: %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		if strings.Contains(err.Error(), "Server closed") {
			return
		}
		log.Panicln(err)
	}
}

type HttpCacheHandler struct {
	mc *cache.MemoryCache
}

func (h *HttpCacheHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.handleGet(w, req)
		break
	case http.MethodPost:
		h.handlePost(w, req)
		break
	}
}

func (h *HttpCacheHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	key := req.PathValue("key")
	val, expireOn, ok := h.mc.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Add("X-Expires-On", expireOn.String())
	_, _ = w.Write(val)
}

func (h *HttpCacheHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	key := req.PathValue("key")

	val, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	h.mc.Store(key, val)
	w.WriteHeader(http.StatusOK)
}
