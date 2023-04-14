package lrucache

import (
	"container/list"
	"sync"
)

type LRU[KeyT comparable, ValueT any] struct {
	// Для реализации LRU кэша используется хэш таблица для хранения значений,
	// а для хранения информации о последнем использовании - очередь, реализованная на списке.
	// Выбор списка обусловлен тем, что удаление из него имеет константную сложность.
	// Линейная сложность поиска по нему нивелируется тем, что в мапе хранится не только значение,
	// но и указатель на ноду из очереди.
	store   cacheStore[KeyT, ValueT]
	q       *list.List
	mx      *sync.Mutex
	maxSize int
}

type cacheStore[KeyT comparable, ValueT any] map[KeyT]*storeItem[ValueT]

type storeItem[ValueT any] struct {
	value   ValueT
	keyNode *list.Element
}

func New[KeyT comparable, ValueT any](maxSize int) *LRU[KeyT, ValueT] {
	if maxSize < 1 {
		panic("Max size of LRU cache must be great then zero.")
	}
	return &LRU[KeyT, ValueT]{
		store:   make(cacheStore[KeyT, ValueT], maxSize),
		q:       list.New(),
		mx:      new(sync.Mutex),
		maxSize: maxSize,
	}
}

func (c *LRU[KeyT, ValueT]) Get(key KeyT) (ValueT, bool) {
	var value ValueT

	// Для того, чтобы обеспечить консистентное добавление в очередь, ставим лок.
	c.mx.Lock()
	defer c.mx.Unlock()

	item := c.store[key]

	// Если значение в кэше, то перемещаем ключ в конец очереди,
	// таким образом пометив его как последнее используемое.
	if item != nil {
		c.q.MoveToBack(item.keyNode)
		value = item.value
	}

	return value, item != nil
}

func (c *LRU[KeyT, ValueT]) Set(key KeyT, value ValueT) {
	c.mx.Lock()
	defer c.mx.Unlock()

	item := c.store[key]
	if item == nil {
		// Если значение не в кэше, то добавляем его в стор,
		// не забыв положить ключ в очередь.
		c.store[key] = &storeItem[ValueT]{
			value:   value,
			keyNode: c.q.PushBack(key),
		}
	} else {
		// Если в кэше, то обновляем значение и перемещаем ключ
		// в конец очереди.
		item.value = value
		c.q.MoveToBack(item.keyNode)
	}

	// Также не забываем вытеснить лишнее значение при переполнении.
	if c.q.Len() > c.maxSize {
		front := c.q.Front() // Неравенство фронта списка nil'у обеспечивается положительным размером кэша.
		delete(c.store, front.Value.(KeyT))
		c.q.Remove(front)
	}
}

func (c *LRU[KeyT, ValueT]) Clear() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.store = make(cacheStore[KeyT, ValueT], c.maxSize)
	c.q = list.New()
}
