package lrucache

import (
	"container/list"
	"sync"
)

type LRUCache interface {
	// Добавляет новое значение с ключом в кеш (с наивысшим приоритетом), возвращает true, если все прошло успешно
	// В случае дублирования ключа вернуть false
	// В случае превышения размера - вытесняется наименее приоритетный элемент
	Add(key, value string) bool

	// Возвращает значение под ключом и флаг его наличия в кеше
	// В случае наличия в кеше элемента повышает его приоритет
	Get(key string) (value string, ok bool)

	// Удаляет элемент из кеша, в случае успеха возврашает true, в случае отсутствия элемента - false
	Remove(key string) (ok bool)

	// Удаление последнего элемента очереди и удаление этого элемента из мапы
	RemoveBack()

	// Удаление первого элемента очереди и удаление этого элемента из мапы
	RemoveFront()

	// Получение всех ключей, в таком порядке, в каком они в очереди
	GetKeys() []string
}

// Версия 2
// Использование связанного списка
type lrucache struct {
	mu     sync.RWMutex
	values map[string]*item
	queue  *list.List
	len    int
}

type item struct {
	value string
	ptr   *list.Element
}

func NewLRUCache(n int) LRUCache {
	return &lrucache{
		values: make(map[string]*item),
		queue:  list.New(),
		len:    n,
	}
}

func (l *lrucache) GetKeys() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	res := make([]string, 0, l.len)
	for e := l.queue.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		res = append(res, key)
	}
	return res
}

func (l *lrucache) Add(key, value string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Проверка наличия элемента
	if _, ok := l.values[key]; ok {
		return false
	}

	// Проверка выхода за пределы
	if l.len == len(l.values) {
		l.mu.Unlock()
		// Удаляем последний элемент очереди
		l.RemoveBack()
		l.mu.Lock()
	}

	// Добавление нового элемента
	l.values[key] = &item{value: value, ptr: l.queue.PushFront(key)}
	return true
}

func (l *lrucache) Get(key string) (value string, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Проверка наличия элемента
	if v, ok := l.values[key]; ok {
		// Переносим элемент в начало очереди
		l.queue.MoveToFront(v.ptr)
		return v.value, true
	}

	return "", false
}

func (l *lrucache) Remove(key string) (ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Проверка наличия элемента
	if v, ok := l.values[key]; ok {
		// Удаление элемента очереди
		l.queue.Remove(v.ptr)

		// Удаление элемента из мапы
		delete(l.values, key)

		return true
	}

	return false
}

func (l *lrucache) RemoveBack() {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Удаление из конца очереди
	endList := l.queue.Back()
	l.queue.Remove(endList)

	// Удаление из мапы
	delete(l.values, endList.Value.(string))
}

func (l *lrucache) RemoveFront() {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Удаление из конца очереди
	startList := l.queue.Front()
	l.queue.Remove(startList)

	// Удаление из мапы
	delete(l.values, startList.Value.(string))
}

// Версия 1
// Время не O(1)
/*type lrucache struct {
	mu     sync.RWMutex
	values map[string]*element
	len    int
}
type element struct {
	q     int
	value string
}

func NewLRUCache(n int) LRUCache {
	return &lrucache{
		values: make(map[string]*element, n),
		len:    n,
	}
}

func (l *lrucache) Add(key, value string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Проверяем наличие ключа
	if _, ok := l.values[key]; ok {
		return false
	}

	// Уменьшение приоретета всех ключей
	l.increment(false)

	// Проверка выхода за пределы
	l.test()

	// Установка нового значения
	l.values[key] = &element{q: 1, value: value}

	return true
}
func (l *lrucache) increment(f bool) {
	for _, e := range l.values {
		if e.q == l.len && f {
			return
		}
		e.q++
	}
}
func (l *lrucache) test() {
	for k, e := range l.values {
		if e.q == l.len+1 {
			delete(l.values, k)
		}
	}
}

func (l *lrucache) Get(key string) (value string, ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.values[key]; ok {
		// Уменьшение приоретета всех ключей
		l.increment(true)

		// Ставим на первое место в приоритете запрошенный ключ
		e.q = 1
		return e.value, true
	}
	return "", false
}

func (l *lrucache) Remove(key string) (ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.values, key)
	return true
}
*/
