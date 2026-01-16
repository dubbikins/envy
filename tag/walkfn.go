package tag

//WalkFn is the type of function called by [Walk] to visit each field of a struct type.
//See the [Node] type for operations that can be performed on the node.
type WalkFn func(node *Node) (err error)
//func(node *Node) (err error)


func (fn WalkFn) Chain(fns ...ChainableWalkFn) WalkFn {
	var i = len(fns)-1
	for i >= 0 {
		if fns[i] == nil {
			break
		}
		fn = fns[i](fn)
		i--
	}

	return fn
}
//Generic WalkFn for a spefic kind of Walkfn
type WalkFnOf[Node any] func(node Node) (err error)

//ChainableWalkFn is a type of function that composes a series of [WalkFn]s,
//each performing actions on a [Node] in the order formed by the chain.
//A single TagWalkFunc usually is designed to handle the parsing for a specific struct tag,
//so ChainableWalkFn allows composing the parsing of multiple tags, in a specified order
//to acheive the desired effect.
type ChainableWalkFn func(next WalkFn) WalkFn


/*
Chained returns chained, the TagWalkFunc derived by either setting chained to walkFn, or applying each [ChainableWalkFn] in chain to the previous value of chained.
Each function in chain is applied in the reverse order that is provided, that is 
	Given chain [cw1,cw2,cw3] and walkFn
	Chained returns cw1(cw2(cw3(walkFn)))
	
*/
func Chained(fns ...ChainableWalkFn) (chained WalkFn) {
	chained = func(node *Node) (err error) {return nil} //do nothing to end the chain
	return chained.Chain(fns...)
}

