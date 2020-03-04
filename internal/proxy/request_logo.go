package proxy

import "net/http"

func logoHandler(w http.ResponseWriter, r *http.Request) {
	sr, err := getStreamRequest(w, r, "/logo/")
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Find real URL of logo
	sr.Channel.LogoCacheMux.Lock()
	defer sr.Channel.LogoCacheMux.Unlock()
	if len(sr.Channel.LogoCache) == 0 {
		img, contentType, err := download(sr.Channel.Logo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
		sr.Channel.LogoCache = img
		sr.Channel.LogoCacheContentType = contentType
	}

	w.Header().Set("Content-Type", sr.Channel.LogoCacheContentType)
	w.Write(sr.Channel.LogoCache)
}
