package notifications

import (
	"encoding/json"
	"net/http"
	httpx "job-crawler/apps/api/internal/platform/http"
)

type Handler struct{ Service Service }

func (h Handler) Update(w http.ResponseWriter, r *http.Request){ if r.Method!=http.MethodPatch{httpx.WriteError(w,405,"method not allowed");return}; var payload map[string]any; if err:=json.NewDecoder(r.Body).Decode(&payload); err!=nil{httpx.WriteError(w,400,"invalid payload");return}; uid:=httpx.UserIDFromContext(r.Context()); httpx.WriteJSON(w,200,h.Service.Update(uid,payload)) }
