package users

import (
	"net/http"
	httpx "job-crawler/apps/api/internal/platform/http"
)

type Handler struct{ Service Service }

func (h Handler) Search(w http.ResponseWriter, r *http.Request){ if r.Method!=http.MethodGet{httpx.WriteError(w,405,"method not allowed");return}; httpx.WriteJSON(w,200,h.Service.Search(r.URL.Query().Get("q"))) }
