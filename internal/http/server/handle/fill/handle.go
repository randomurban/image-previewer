package fill

import (
	"log"
	"net/http"
	"strconv"

	"github.com/randomurban/image-previewer/internal/http/server/handle"
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
	width, err := strconv.Atoi(r.PathValue("width"))
	if err != nil {
		log.Printf("width: %v", err)
		http.Error(w, "Bad request: wrong width in url", http.StatusBadRequest)
		return
	}
	log.Printf("width: %v", width)

	height, err := strconv.Atoi(r.PathValue("height"))
	if err != nil {
		log.Printf("height: %v", err)
		http.Error(w, "Bad request: wrong height in url", http.StatusBadRequest)
		return
	}
	log.Printf("height: %v", height)

	url := r.PathValue("img")
	log.Printf("image url: %v", url)

	imgBuf, err := h.previewer.PreviewImage(width, height, url, r.Header)
	if err != nil {
		log.Printf("preview image: %v", err)
		http.Error(w, "preview: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(imgBuf)
	if err != nil {
		log.Printf("encode: %v", err)
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBuf)))
	_, err = w.Write(imgBuf)
	if err != nil {
		log.Printf("failed write: %s", err)
		return
	}
}
