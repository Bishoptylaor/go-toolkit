package zcontainer

/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/26 -- 16:44
 @Author  : bishop ❤️ MONEY
 @Description: stack
*/

type StackElement struct {
	next  *StackElement
	stack *Stack
	Value any
}

type Stack struct {
	root StackElement
	len  int
}

func (s *Stack) init() *Stack {
	s.root.next = &s.root
	s.len = 0
	return s
}

func NewStack() *Stack {
	return new(Stack).init()
}

func (s *Stack) Len() int {
	return s.len
}

func (s *Stack) Push(v interface{}) {
	se := StackElement{s.root.next, s, v}
	s.root.next = &se
	s.len++
}

func (s *Stack) Pop() *StackElement {
	if s.len <= 0 {
		return nil
	}
	retVal := s.root.next
	s.root.next = retVal.next
	s.len--
	return retVal
}
