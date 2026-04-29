package profiles

import (
	"encoding/json"
	"net/http"
	"strings"
	httpx "job-crawler/apps/api/internal/platform/http"
)

type Handler struct{ Service Service }

func (h Handler) GetByID(w http.ResponseWriter, r *http.Request){ if r.Method!=http.MethodGet{httpx.WriteError(w,405,"method not allowed");return}; id:=strings.TrimPrefix(r.URL.Path,"/api/profiles/"); requester:=httpx.UserIDFromContext(r.Context()); p,ok:=h.Service.Get(requester,id); if !ok {httpx.WriteError(w,404,"not found");return}; if requester!=id && p.IsPrivate {httpx.WriteError(w,403,"private profile");return}; httpx.WriteJSON(w,200,p) }
func (h Handler) UpdateMe(w http.ResponseWriter, r *http.Request){ if r.Method!=http.MethodPatch{httpx.WriteError(w,405,"method not allowed");return}; var patch map[string]any; if err:=json.NewDecoder(r.Body).Decode(&patch); err!=nil {httpx.WriteError(w,400,"invalid payload");return}; uid:=httpx.UserIDFromContext(r.Context()); p,ok:=h.Service.UpdateMe(uid,patch); if !ok {httpx.WriteError(w,404,"not found");return}; httpx.WriteJSON(w,200,p)}
