package dns

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	errNeedRefresh        = errors.New("cache need to refresh")
	errCacheEntryNotFound = errors.New("cache enntry not found")
)

type entry struct {
	ip        []string
	timepoint time.Time
	ttl       time.Duration
}

type recordCache struct {
	m *sync.RWMutex
	c map[string]entry
}

func newDnsCache() recordCache {
	return recordCache{
		m: &sync.RWMutex{},
		c: make(map[string]entry),
	}
}

func (cache recordCache) get(domain string) (string, error) {
	cache.m.RLock()

	e, ok := cache.c[domain]

	cache.m.RUnlock()

	if !ok {
		return "", errCacheEntryNotFound
	}

	if time.Since(e.timepoint) > e.ttl {
		cache.m.Lock()
		defer cache.m.Unlock()

		log.Printf("doh: dns cache entry expired %s -> %s \n", domain, e.ip[0])

		delete(cache.c, domain)

		return "", errNeedRefresh
	}

	return e.get(), nil
}

func (cache recordCache) put(domain string, ip []string, ttl time.Duration) string {
	if len(ip) == 0 {
		log.Printf("doh: no ip provided to put into cache for %s\n", domain)
		return ""
	}

	cache.m.Lock()
	defer cache.m.Unlock()

	e, ok := cache.c[domain]
	if ok {
		return e.get()
	}

	e = entry{
		ip:        ip,
		timepoint: time.Now(),
		ttl:       ttl,
	}

	cache.c[domain] = e

	log.Printf("doh: added cache record %s -> %s\n", domain, e.string())

	return e.get()
}

func (e entry) get() string {
	return e.ip[0]
}

func (e entry) string() string {
	return fmt.Sprintf("[ ttl=%s, ip=%q ]", e.ttl, e.ip)
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}

	return b
}
