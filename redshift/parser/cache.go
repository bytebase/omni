package parser

import "sync"

// firstSetCache caches the candidates emitted by collect-mode parse functions.
type firstSetCache struct {
	mu sync.RWMutex
	m  map[string]*CandidateSet
}

var globalFirstSets = &firstSetCache{
	m: make(map[string]*CandidateSet),
}

func (c *firstSetCache) get(key string) *CandidateSet {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.m[key]
}

func (c *firstSetCache) set(key string, cs *CandidateSet) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = cs
}

// cachedCollect runs the given function to collect candidates.
// The cache is intentionally disabled because FIRST sets can depend on parser
// state that varies between Collect calls (e.g., cursor position relative to
// nested expressions), causing stale cache entries to return incomplete results.
func (p *Parser) cachedCollect(_ string, fn func()) bool {
	fn()
	return false
}
