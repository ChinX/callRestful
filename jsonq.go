package matcher

import (
	"strconv"
	"strings"
)

type Matcher struct {
	mapper map[string]interface{}
	keyList []*part
}

func NewMatcher(mapper map[string]interface{}, keys []string) *Matcher {
	return &Matcher{
		mapper:mapper,
		keyList: NewPartList(keys),
	}
}

func (m *Matcher) Match() (interface{}, bool) {
	if m.keyList == nil{
		return "", false
	}
	var tempMap interface{} = m.mapper
	list := m.keyList
	ok := true
	for len(list) > 0{
		tempMap, list, ok = match(tempMap, list)
		if !ok {
			return "", false
		}
	}
	return tempMap, true
}

func match(values interface{}, list []*part) (interface{}, []*part, bool) {
		item := list[0]
		switch item.kind {
		case array:
			arr, ok := values.([]interface{})
			if ok {
				values, list, ok = matchInArray(arr, list)
			}
			if !ok {
				return nil, list, ok
			}
		case mapper:
			obj, ok := values.(map[string]interface{})
			if ok {
				values, list, ok = matchInMapper(obj, list)
			}
			if !ok {
				return nil, list, ok
			}
		}
		return values, list, true
}

func matchInArray(values []interface{}, list []*part) (interface{}, []*part, bool) {
	l := len(values)
	if l < 1 {
		return nil, list[1:], false
	}

	var val interface{}
	ok := true
	switch list[0].search {
	case normal:
		index, err := strconv.Atoi(list[0].key)
		if err != nil || index >= l {
			return nil, list[1:], false
		}
		return values[index], list[1:], true
	case between:
		val, list, ok = searchInArray(values, list)
	}
	return val, list, ok
}

func searchInArray(values []interface{}, list []*part) (interface{}, []*part, bool) {
	for _, v := range values {
		val, tmpList, ok := match(v, list[1:])
		if ok {
			return val, tmpList, ok
		}
	}
	return nil, list[1:], false
}

func matchInMapper(values map[string]interface{}, list []*part) (interface{}, []*part, bool) {
	var val interface{}
	ok := true
	switch list[0].search {
	case normal:
		val, list, ok = findInMapper(values, list)
	case prefix:
		val, list, ok = searchByHandle(values, list, strings.HasPrefix)
	case suffix:
		val, list, ok = searchByHandle(values, list, strings.HasSuffix)
	case between:
		val, list, ok = searchInMapper(values, list)
	}
	return val, list, ok
}

func searchByHandle(values map[string]interface{}, list []*part, handler func(s, substr string) bool) (interface{}, []*part, bool) {
	for k, v := range values {
		if handler(k, list[0].key) {
			return v, list[1:], true
		}
	}
	return nil, list[1:], false
}

func findInMapper(values map[string]interface{}, list []*part) (interface{}, []*part, bool) {
	for k, v := range values {
		if k == list[0].key {
			return v, list[1:], true
		}
	}
	return nil, list[1:], false
}

func searchInMapper(values map[string]interface{}, list []*part) (interface{}, []*part, bool) {
	if len(list[0].key) > 0 {
		return searchByHandle(values, list, strings.Contains)
	}
	for _, v := range values {
		val, tmpList, ok := match(v, list[1:])
		if ok {
			return val, tmpList, ok
		}
	}
	return nil, list[1:], false
}

package matcher

const (
	mapper keyKind = iota
	array
	normal
	prefix
	suffix
	between
)

type keyKind uint

type part struct {
	kind   keyKind
	key    string
	search keyKind
}

func NewPartList(keys []string) (list []*part) {
	if len(keys) == 0 {
		return nil
	}
	for i := range keys {
		p := matchPart(keys[i])
		if p == nil {
			return nil
		}
		list = append(list, p)
	}
	if len(list[len(list)-1].key) == 0{
		panic("the endpoint key must be not empty")
	}
	return
}

func matchPart(key string) *part {
	l := len(key)
	if l == 0 {
		return nil
	}
	if l > 1 && key[0] == '[' && key[l-1] == ']' {
		str := key[1 : l-1]
		p := &part{
			kind:   array,
			search: normal,
			key:    str,
		}
		if len(str) == 0 {
			p.search = between
		}
		return p
	}
	p := &part{
		kind:   mapper,
		search: normal,
		key:    key,
	}

	if suf, pre := key[0] == '%', key[l-1] == '%'; pre || suf {
		if pre && suf {
			p.search = between
			if l == 1{
				p.key = ""
			}else{
				p.key = key[1 : l-1]
			}
		} else if pre {
			p.search = prefix
			p.key = key[0:l-1]
		} else {
			p.search = suffix
			p.key = key[1:l]
		}
	}
	return p
}
