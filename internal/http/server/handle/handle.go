package handle

import "net/http"

type Handler interface {
	FillHandle(w http.ResponseWriter, r *http.Request)
}
