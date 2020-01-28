package proxy

import (
	"bufio"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var m3u8TSCache *cache.Cache

var m3u8channels map[string]*M3U8Channel

// M3U8Channel stores information about m3u8 channel
type M3U8Channel struct {
	Channel *Channel

	linkOneAtOnceMux sync.Mutex // To ensure that only one routine at once is accessing & updating link

	link    string
	linkMux sync.RWMutex

	linkCache    []byte
	linkCacheMux sync.RWMutex

	linkCacheCreated    time.Time
	linkCacheCreatedMux sync.RWMutex

	linkRoot    string
	linkRootMux sync.RWMutex

	sessionUpdated    time.Time
	sessionUpdatedMux sync.RWMutex
}

// Link ...
func (c *M3U8Channel) Link() string {
	c.linkMux.RLock()
	defer c.linkMux.RUnlock()
	return c.link
}

// SetLink ...
func (c *M3U8Channel) SetLink(s string) {
	c.linkMux.Lock()
	defer c.linkMux.Unlock()
	c.link = s
}

// LinkCache ...
func (c *M3U8Channel) LinkCache() []byte {
	c.linkCacheMux.RLock()
	defer c.linkCacheMux.RUnlock()
	return c.linkCache
}

// SetLinkCache ...
func (c *M3U8Channel) SetLinkCache(b []byte) {
	c.linkCacheMux.Lock()
	defer c.linkCacheMux.Unlock()
	c.linkCache = b
}

// LinkCacheCreated ...
func (c *M3U8Channel) LinkCacheCreated() time.Time {
	c.linkCacheCreatedMux.RLock()
	defer c.linkCacheCreatedMux.RUnlock()
	return c.linkCacheCreated
}

// SetLinkCacheCreatedNow ...
func (c *M3U8Channel) SetLinkCacheCreatedNow() {
	c.linkCacheCreatedMux.Lock()
	defer c.linkCacheCreatedMux.Unlock()
	c.linkCacheCreated = time.Now()
}

// LinkRoot ...
func (c *M3U8Channel) LinkRoot() string {
	c.linkRootMux.RLock()
	defer c.linkRootMux.RUnlock()
	return c.linkRoot
}

// SetLinkRoot ...
func (c *M3U8Channel) SetLinkRoot(s string) {
	c.linkRootMux.Lock()
	defer c.linkRootMux.Unlock()
	c.linkRoot = s
}

// ----------

// LinkCacheValid ...
func (c *M3U8Channel) LinkCacheValid() bool {
	s := c.LinkCacheCreated()
	if time.Since(s).Seconds() > 1 || s.IsZero() {
		return false
	}
	return true
}

// Link ...
func (c *M3U8Channel) newRedirectedLink(s string) {
	c.SetLink(s)
	c.SetLinkRoot(deleteAfterLastSlash(s))
}

func deleteAfterLastSlash(str string) string {
	return str[0 : strings.LastIndex(str, "/")+1]
}

var reURILinkExtract = regexp.MustCompile(`URI="([^"]*)"`)

func rewriteLinks(origContent, prefix, linkRoot string) string {
	var sb strings.Builder

	linkRootURL, _ := url.Parse(linkRoot) // It will act as a base URL for full URLs

	modifyLink := func(link string) string {
		var l string

		switch {
		case strings.HasPrefix(link, "//"):
			tmpURL, _ := url.Parse(link)
			tmp2URL, _ := url.Parse(tmpURL.RequestURI())
			link = (linkRootURL.ResolveReference(tmp2URL)).String()
			l = strings.ReplaceAll(link, linkRoot, "")
		case strings.HasPrefix(link, "/"):
			tmp2URL, _ := url.Parse(link)
			link = (linkRootURL.ResolveReference(tmp2URL)).String()
			l = strings.ReplaceAll(link, linkRoot, "")
		default:
			l = link
		}

		return prefix + l
	}

	scanner := bufio.NewScanner(strings.NewReader(origContent))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			line = modifyLink(line)
		} else if strings.Contains(line, "URI=\"") && !strings.Contains(line, "URI=\"\"") {
			link := reURILinkExtract.FindStringSubmatch(line)[1]
			line = reURILinkExtract.ReplaceAllString(line, `URI="`+modifyLink(link)+`"`)
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}
