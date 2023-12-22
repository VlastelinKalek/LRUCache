package lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type caseStruct struct {
	fn    string
	key   string
	value string
	want  bool
}

type cmd struct {
	c     LRUCache
	cases []caseStruct
}

// Проверка добавление, получение и удаление из кэша
func TestAddGetRemove(t *testing.T) {
	req := require.New(t)

	testsCase := map[string]cmd{
		"AddCache":        {NewLRUCache(1), []caseStruct{{"Add", "One", "1", true}}},
		"GetCache":        {NewLRUCache(1), []caseStruct{{"Add", "One", "1", true}, {"Get", "One", "1", true}}},
		"RemoveCache":     {NewLRUCache(1), []caseStruct{{"Add", "One", "1", true}, {"Remove", "One", "1", true}}},
		"AddDublicate":    {NewLRUCache(1), []caseStruct{{"Add", "One", "1", true}, {"Add", "One", "1", false}}},
		"RemoveNullCache": {NewLRUCache(1), []caseStruct{{"Remove", "One", "1", false}}},
		"GetNullCache":    {NewLRUCache(1), []caseStruct{{"Get", "One", "1", false}}},
	}

	for name, cmd := range testsCase {
		t.Run(name, func(t *testing.T) {
			for _, cases := range cmd.cases {
				var res bool
				switch cases.fn {
				case "Add":
					res = cmd.c.Add(cases.key, cases.value)
				case "AddDublicate":
					res = cmd.c.Add(cases.key, cases.value)
				case "Get":
					_, res = cmd.c.Get(cases.key)
				case "Remove":
					res = cmd.c.Remove(cases.key)
				}
				req.Equal(res, cases.want, name)
			}
		})
	}
}

// Проверка превышение лимита кэша
func TestOverLimit(t *testing.T) {
	req := require.New(t)
	c := NewLRUCache(2)
	want := []string{"Three", "Two"}

	t.Run("OverLimit", func(t *testing.T) {
		c.Add("One", "1")
		c.Add("Two", "2")
		c.Add("Three", "3")
		res := c.GetKeys()
		req.Equal(res, want)
	})
}

// Проверка удаления первых и последних элементов
func TestRemoveFrontBack(t *testing.T) {
	req := require.New(t)
	c := NewLRUCache(3)
	want := []string{"Two"}

	t.Run("OverLimit", func(t *testing.T) {
		c.Add("One", "1")
		c.Add("Two", "2")
		c.Add("Three", "3")
		c.RemoveFront()
		c.RemoveBack()
		res := c.GetKeys()
		req.Equal(res, want)
	})
}
