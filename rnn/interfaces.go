package rnn

import (
	"github.com/unixpickle/autofunc"
	"github.com/unixpickle/num-analysis/linalg"
)

// UpstreamGradient stores the gradients of some
// output with respect to the outputs and output
// states of some Block.
// Either one of the slices (States or Outputs)
// may be nil, indicating that said gradient is
// completely 0.
type UpstreamGradient struct {
	States  []linalg.Vector
	Outputs []linalg.Vector
}

// UpstreamRGradient is like UpstreamGradient,
// but it stores the derivatives of all the
// partials with respect to some variable R.
//
// A slice (States or Outputs) can be nil if and
// only if its corresponding R slice is also nil.
type UpstreamRGradient struct {
	UpstreamGradient
	RStates  []linalg.Vector
	ROutputs []linalg.Vector
}

// A BlockInput stores a batch of states and inputs
// for a Block.
type BlockInput struct {
	States []*autofunc.Variable
	Inputs []*autofunc.Variable
}

// A BlockOutput represents a batch of outputs and new
// states from a Block.
type BlockOutput interface {
	States() []linalg.Vector
	Outputs() []linalg.Vector

	// Gradient updates the gradients in g given the
	// upstream gradient from this BlockOutput.
	Gradient(u *UpstreamGradient, g autofunc.Gradient)
}

// A BlockRInput is like a BlockInput, but includes
// derivatives of all the inputs and states with
// respect to some variable R.
type BlockRInput struct {
	States []*autofunc.RVariable
	Inputs []*autofunc.RVariable
}

// An BlockROutput is like a BlockOutput, but includes
// derivatives of the outputs and states with respect
// to some variable R.
type BlockROutput interface {
	States() []linalg.Vector
	Outputs() []linalg.Vector
	RStates() []linalg.Vector
	ROutputs() []linalg.Vector

	// RGradient updates the gradients in g and the
	// r-gradients in rg given the upstream gradient
	// u and the derivative of u with respect to R,
	// stored in ru.
	// The gradient g may be nil to indicate that only
	// the r-gradient is desired.
	RGradient(u *UpstreamRGradient, rg autofunc.RGradient, g autofunc.Gradient)
}

// A Block is a unit in a Recurrent Neural Network that
// transforms input-state pairs into output/state pairs.
type Block interface {
	// StateSize returns the number of values in each
	// state of the Block.
	StateSize() int

	// Batch applies forward propagation to a BlockInput.
	Batch(in *BlockInput) BlockOutput

	// BatchR is like Batch, but for an BlockRInput.
	// It is necessary to provide an RVector so that the
	// block knows how much each of its hidden parameters
	// changes with respect to R.
	BatchR(v autofunc.RVector, in *BlockRInput) BlockROutput
}

// A Learner is a Block with parameters that can be
// learned using some variant of gradient descent.
type Learner interface {
	Block
	Parameters() []*autofunc.Variable
}