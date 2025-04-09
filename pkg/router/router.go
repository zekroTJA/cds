package router

// A lot of this code is more or less directly copied
// from Go's net/http package.
// See: https://github.com/golang/go/blob/master/src/net/http/routing_tree.go
// License: https://github.com/golang/go/blob/master/LICENSE

import (
	"bytes"
	"strings"
)

type Router[H any] struct {
	root routeNode[H]
}

type routeNode[H any] struct {
	pattern  []string
	handlers []H

	children mapping[string, *routeNode[H]]
}

func (t *routeNode[H]) add(sub []string, pattern []string, handler H) {
	if len(sub) == 0 {
		t.pattern = pattern
		t.handlers = append(t.handlers, handler)
		return
	}

	key := sub[0]
	node, ok := t.children.get(key)
	if ok {
		node.add(sub[1:], pattern, handler)
		return
	}

	node = &routeNode[H]{}
	node.add(sub[1:], pattern, handler)
	t.children.add(key, node)
}

func (t *routeNode[H]) match(pattern []string) (handlers []H, sub string) {
	if len(pattern) == 0 {
		return t.handlers, ""
	}

	key := pattern[0]
	node, ok := t.children.get(key)
	if !ok {
		return t.handlers, strings.Join(pattern, "/")
	}

	return node.match(pattern[1:])
}

func (t *Router[H]) Add(pattern string, handler H) {
	patternSplit := splitPattern(pattern)
	t.root.add(patternSplit, patternSplit, handler)
}

func (t *Router[H]) Match(pattern string) (handlers []H, sub string, ok bool) {
	patternSplit := splitPattern(pattern)
	matchedHandlers, sub := t.root.match(patternSplit)
	if len(matchedHandlers) == 0 {
		return nil, "", false
	}
	return matchedHandlers, sub, true
}

func splitPattern(pattern string) []string {
	if pattern == "" || pattern == "/" {
		return nil
	}

	if pattern[0] == '/' {
		pattern = pattern[1:]
	}

	n := strings.Count(pattern, "/")
	split := make([]string, 0, n)

	var buf bytes.Buffer
	for _, c := range pattern {
		if c == '/' {
			if buf.Len() == 0 {
				continue
			}
			split = append(split, buf.String())
			buf.Reset()
		} else {
			buf.WriteRune(c)
		}
	}

	if buf.Len() > 0 {
		split = append(split, buf.String())
	}

	return split
}
