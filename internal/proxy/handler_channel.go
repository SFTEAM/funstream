package proxy

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
)

func quickWrite(ctx *fasthttp.RequestCtx, content []byte, contentType string, httpStatus int) {
	ctx.SetContentType(contentType)
	ctx.SetStatusCode(httpStatus)
	ctx.SetBody(content)
}

func channelHandler(ctx *fasthttp.RequestCtx) {
	reqPath := strings.Replace(string(ctx.RequestURI()), "/iptv/", "", 1)
	reqPathParts := strings.SplitN(reqPath, "/", 2)
	if len(reqPathParts) == 0 {
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	// Debug
	if len(reqPathParts) > 1 {
		log.Println("Received channel-only request")
	} else {
		log.Println("Received channel-with-data request")
	}

	// Unescape title
	unescapedTitle, err := url.QueryUnescape(reqPathParts[0])
	if err != nil {
		ctx.Error("invalid request", http.StatusBadRequest)
		return
	}

	// Find channel
	c, ok := playlist.Channels[reqPathParts[0]]
	if !ok {
		ctx.Error("channel not found", http.StatusNotFound)
		return
	}

	// Find link reference
	c.ActiveLinkMux.RLock()
	l := c.ActiveLink
	c.ActiveLinkMux.RUnlock()

	// Find channel type
	l.Mux.RLock()
	link := l.Link
	linkType := l.LinkType
	l.Mux.RUnlock()

	// Error if channel type is unknown and request URL contains additional path
	if linkType == linkTypeUnknown && len(reqPathParts) == 2 {
		ctx.Error("invalid request", http.StatusBadRequest)
		return
	}

	// Lock mutex if channel's type is unknown, so no other routine tries to identify it at the same time
	c.ActiveLink.Mux.Lock()
	if linkType == linkTypeUnknown {
		handleLinkUnknown(ctx, &reqPathParts[0], &unescapedTitle, link, c, c.ActiveLink)
		c.ActiveLink.Mux.Unlock()
		return
	}
	c.ActiveLink.Mux.Unlock()

	// Understand what do we need to do with this link
	switch linkType {
	case linkTypeMedia:
		log.Println("Processing type: Media")
		handleStream(ctx, &reqPathParts[0], &unescapedTitle, link, c, c.ActiveLink)
	case linkTypeStream:
		log.Println("Processing type: Stream")
		handleStream(ctx, &reqPathParts[0], &unescapedTitle, link, c, c.ActiveLink)
	case linkTypeM3U8:
		log.Println("Processing type: M3U8")
		ctx.Error("not yet supported", http.StatusNotImplemented) // TODO
	case linkTypeUnsupported:
		ctx.Error("unsupported channel format", http.StatusServiceUnavailable)
	default:
		ctx.Error("internal server error", http.StatusInternalServerError)
	}
}
