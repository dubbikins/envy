package tag

//WalkFn is the type of function called by [Walk] to visit each field of a struct type.
//See the [Node] type for operations that can be performed on the node.
type WalkFn func(node *Node) (err error)
//func(node *Node) (err error)


//Generic WalkFn for a spefic kind of Walkfn
type WalkFnOf[Node any] func(node Node) (err error)

//ChainableWalkFn is a type of function that composes a series of [WalkFn]s,
//each performing actions on a [Node] in the order formed by the chain.
//A single TagWalkFunc usually is designed to handle the parsing for a specific struct tag,
//so ChainableWalkFn allows composing the parsing of multiple tags, in a specified order
//to acheive the desired effect.
type ChainableWalkFn func(next WalkFn) WalkFn


/*
Chain returns chained, the TagWalkFunc derived by either setting chained to walkFn, or applying each [ChainableWalkFn] in chain to the previous value of chained.
Each function in chain is applied in the reverse order that is provided, that is 
	Given chain [cw1,cw2,cw3] and walkFn
	Chain returns cw1(cw2(cw3(walkFn)))
	
*/
func Chain(walkFn WalkFn, chain ...ChainableWalkFn) (chained WalkFn) {
	chained = walkFn
	var i = len(chain)-1
	for i >= 0 {
		chained = chain[i](chained)
		i--
	}
	return chained
}

