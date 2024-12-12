package langs

import "github.com/ge-editor/gecore/lang"

// ------------------------------------------------------------------
// Fundamental implement gecore Mode interface
// ------------------------------------------------------------------

func NewFundamental() lang.Mode {
	return &Fundamental{
		AbstractWorker: AbstractWorker{
			exts: []string{}, // Initialize the `exts` field of AbstractWorker
		},
	}
}

type Fundamental struct {
	AbstractWorker
}

func (f *Fundamental) Name() string {
	return "Fundamental"
}
