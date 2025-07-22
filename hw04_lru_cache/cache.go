package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

// Set - добавляет или обновляет элемент в кэше.
func (c *lruCache) Set(key Key, value interface{}) bool {
	// Если ключ уже существует, обновляем значение и перемещаем в начало.
	if item, exists := c.items[key]; exists {
		item.Value = value
		c.queue.MoveToFront(item)
		return true
	}

	// Если кэш полон, удаляем последний элемент.
	if c.queue.Len() >= c.capacity {
		last := c.queue.Back()
		if last != nil {
			// Находим ключ, соответствующий последнему элементу.
			for k, v := range c.items {
				if v == last {
					delete(c.items, k)
					break
				}
			}
			c.queue.Remove(last)
		}
	}

	// Добавляем новый элемент в начало.
	listItem := c.queue.PushFront(value)
	c.items[key] = listItem
	return false
}

// Get - получает значение по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	// Перемещаем элемент в начало.
	c.queue.MoveToFront(item)
	return item.Value, true
}

// Clear - очищает кэш.
func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
