package hw04lrucache

// List - интерфейс для двусвязного списка
type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// ListItem - структура элемента списка
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

// list - реализация двусвязного списка.
type list struct {
	head *ListItem
	tail *ListItem
	size int
}

// Len - возвращает размер списка.
func (l *list) Len() int {
	return l.size
}

// Front - возвращает первый элемент списка.
func (l *list) Front() *ListItem {
	return l.head
}

// Back - возвращает последний элемент списка.
func (l *list) Back() *ListItem {
	return l.tail
}

// PushFront - добавляет элемент в начало списка.
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.head == nil {
		l.head = newItem
		l.tail = newItem
	} else {
		newItem.Next = l.head
		l.head.Prev = newItem
		l.head = newItem
	}
	l.size++
	return newItem
}

// PushBack - добавляет элемент в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.tail == nil {
		l.head = newItem
		l.tail = newItem
	} else {
		newItem.Prev = l.tail
		l.tail.Next = newItem
		l.tail = newItem
	}
	l.size++
	return newItem
}

// Remove - удаляет указанный элемент из списка.
func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}
	l.size--
}

// MoveToFront - перемещает указанный элемент в начало списка.
func (l *list) MoveToFront(i *ListItem) {
	if i == nil || i == l.head {
		return
	}
	// Удаляем элемент из текущей позиции.
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}
	// Перемещаем в начало.
	i.Next = l.head
	i.Prev = nil
	l.head.Prev = i
	l.head = i
}

// NewList - создает новый пустой список.
func NewList() List {
	return new(list)
}
