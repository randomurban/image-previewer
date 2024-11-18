package fill

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/randomurban/image-previewer/internal/http/server/handle"
	"github.com/randomurban/image-previewer/internal/model"
	"github.com/randomurban/image-previewer/internal/service"
)

type Handle struct {
	previewer service.Previewer
}

func NewHandle(previewer service.Previewer) handle.Handler {
	return &Handle{
		previewer: previewer,
	}
}

func (h Handle) FillHandle(w http.ResponseWriter, r *http.Request) {
	log.Printf("url: %s", r.URL)
	width, err := strconv.Atoi(r.PathValue("width"))
	if err != nil {
		log.Printf("width: %v", err)
		http.Error(w, "Bad request: illegal width", http.StatusBadRequest)
		return
	}
	log.Printf("width: %v", width)
	if width <= 0 {
		http.Error(w, "Bad request: illegal width", http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(r.PathValue("height"))
	if err != nil {
		log.Printf("height: %v", err)
		http.Error(w, "Bad request: illegal height", http.StatusBadRequest)
		return
	}
	log.Printf("height: %v", height)
	if height <= 0 {
		http.Error(w, "Bad request: illegal height", http.StatusBadRequest)
		return
	}

	urlParam := strings.SplitN(r.URL.String(), "/", 5)
	if len(urlParam) != 5 {
		http.Error(w, "Bad request: illegal url", http.StatusBadRequest)
		return
	}
	url := urlParam[4]
	log.Printf("image url: %v", url)

	imgPreview, err := h.previewer.PreviewImage(width, height, url, r.Header)
	if err != nil {
		log.Printf("preview image: %v", err)
		switch {
		case errors.Is(err, model.ErrNotFound):
			http.Error(w, "Not found", http.StatusNotFound)
		case errors.Is(err, model.ErrTooLarge):
			http.Error(w, "Too Large", http.StatusRequestEntityTooLarge)
		case errors.Is(err, model.ErrBadGateway):
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		case errors.Is(err, model.ErrTimeout):
			http.Error(w, "Request Timeout", http.StatusRequestTimeout)
		case errors.Is(err, model.ErrUnsupportedMediaType):
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		default:
			http.Error(w, "preview: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if imgPreview.IsCacheHit {
		log.Printf("image preview cache HIT")
		w.Header().Set("X-Cache", "HIT")
	} else {
		log.Printf("image preview cache MISS")
		w.Header().Set("X-Cache", "MISS")
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgPreview.Buf)))

	_, err = w.Write(imgPreview.Buf)
	if err != nil {
		log.Printf("encode: %v", err)
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
