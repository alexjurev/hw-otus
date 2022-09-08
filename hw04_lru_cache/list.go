package hw04lrucache

import "fmt"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
	key   Key
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	n := ListItem{
		Value: v,
		Next:  l.front,
	}
	if l.len == 0 {
		l.back = &n
	}
	if l.front != nil {
		l.front.Prev = &n
	}
	l.len++
	l.front = &n
	return &n
}

func (l *list) PushBack(v interface{}) *ListItem {
	n := ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.len == 0 {
		l.front = &n
	}
	if l.back != nil {
		l.back.Next = &n
	}
	l.len++
	l.back = &n
	return &n
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		fmt.Println("can't delete nil")
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		if i.Next != nil {
			i.Next.Prev = nil
		}
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		if i.Prev != nil {
			i.Prev.Next = nil
		}
		l.back = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev == nil {
		return
	}

	if i.Next == nil {
		l.front.Prev = i
		i.Prev.Next = nil
		i.Next = l.front
		l.front = i
		l.back = i.Prev
		return
	}

	i.Next.Prev = i.Prev
	i.Prev.Next = i.Next
	l.front.Prev = i
	i.Next = l.front
	l.front = i
}
