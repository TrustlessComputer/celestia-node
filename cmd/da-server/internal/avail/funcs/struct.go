package funcs

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

type HeaderF struct {
	Root             types.Hash
	Proof            []types.Hash
	Number_Of_Leaves int
	Leaf_index       int
	Leaf             types.Hash
}
