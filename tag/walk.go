package tag

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"runtime/debug"
)

/*
Walk performs a preorder traversal of s, applying walkfn to each node n that it visits
Walk returns an error if
- s is not a pointer to a struct type
- walkFn returns an error on any node
- Walk recovers from a panic, the recovered value is returned

*/
func Walk(walkfn WalkFn, s any) (err error) {
	return WalkContext(context.Background(), walkfn, s)
}

func WalkContext(ctx context.Context, walkfn WalkFn, s any) (err error) {
	if reflect.TypeOf(s).Kind() != reflect.Pointer || reflect.ValueOf(s).Elem().Kind() != reflect.Struct {
		return errors.New("value passed to Unmarshal must be a pointer to a struct type")
	}
	defer func() {
		if cause := recover(); cause != nil {
			slog.Error("Unmarshal recovered from panic", "cause", cause,)
			print(string(debug.Stack()))
			if cause_err, ok := cause.(error); ok {
				err = errors.Join(err, cause_err)
			}else {
				err = errors.Join(err, errors.New(cause.(string)))
			}
			
		}
	}()
	var root = NewRootNode(s)
	root.ctx = ctx
	var stack = []*Node{&root}
	for len(stack) > 0 {
		var curr *Node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		curr.ctx = ctx
		if curr.skip {
			continue
		}
		err = errors.Join(err, walkfn(curr))
		for next := range curr.descendants() {
			stack = append(stack, &next)
		}
	}
	return
}

var UnmarshalText WalkFn = func (n *Node) ( error) {
	return n.UnmarshalText(n.Bytes())
}

var UnmarshalTextContext WalkFn = func ( n *Node) ( error) {
	select {
	case<-n.Context().Done():
		return n.Context().Err()
	default:
		return n.UnmarshalText(n.Bytes())
	}
}

// func MarshalText(n *Node) ( err error) {
// 	var data []byte
// 	if data, err = n.MarshalText(); err != nil {
// 		return
// 	}
// 	n.Reset()
// 	_, err = n.Write(data)
// 	return
// }


// func ReverseLevelOrderWalk(walkfn WalkFn, s any) (err error) {
// 	return ReverseLevelOrderWalkContext(context.Background(), walkfn, s)
// }


/*

To get the reverse level order
Push the root node into a queue. While queue is not empty, perform steps 3,4.
Pop the front element of the queue. Instead of printing it (like in level order), push node's value into a stack. This way the elements present on upper levels will be printed later than elements on lower levels.
If the right child of current node exists, push it into queue before left child. This is because elements on the same level should be printed from left to right.
While stack is not empty, pop the top element and print it.

*/
func ReverseLevelOrderWalkContext(ctx context.Context, walkfn WalkFn, s any) (err error) {
	if reflect.TypeOf(s).Kind() != reflect.Pointer || reflect.ValueOf(s).Elem().Kind() != reflect.Struct {
		return errors.New("value passed to Unmarshal must be a pointer to a struct type")
	}
	defer func() {
		if cause := recover(); cause != nil {
			slog.Error("Unmarshal recovered from panic", "cause", cause)
			if cause_err, ok := cause.(error); ok {
				err = errors.Join(err, cause_err)
			}else {
				err = errors.Join(err, errors.New(cause.(string)))
			}
			
		}
	}()
	var root = NewRootNode( s)
	root.ctx = ctx
	var stack = []*Node{}
	var queue = []*Node{&root}
	for len(queue) > 0 {
		var curr *Node = queue[0]
		queue = queue[1:]
		stack = append(stack, curr)
		for next := range curr.descendants() {
			queue = append(queue, &next)
		}
	}
	for len(stack) > 0 {
		var curr *Node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		curr.ctx = ctx
		if curr.skip {
			continue
		}
		err = errors.Join(err, walkfn(curr))
		for next := range curr.descendants() {
			stack = append(stack, &next)
		}
	}
	return
}