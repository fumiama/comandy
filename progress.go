package main

import (
	"time"

	"github.com/FloatTech/ttl"
)

var pl = ttl.NewCache[string, *uintptr](time.Minute)

type progressLogger struct {
	key string
	sz  int64
	v   uintptr
	rd  uintptr
}

func newProgressLogger(key string, sz int64) *progressLogger {
	p := new(progressLogger)
	p.key = key
	p.sz = sz
	pl.Set(key, &p.v)
	return p
}

func (p *progressLogger) Write(b []byte) (int, error) {
	if p.sz == 0 {
		return len(b), nil
	}
	p.rd += uintptr(len(b))
	p.v = p.rd * 100 / uintptr(p.sz)
	return len(b), nil
}

func (p *progressLogger) remove() {
	pl.Delete(p.key)
}
