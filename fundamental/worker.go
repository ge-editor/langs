package fundamental

import (
	"github.com/ge-editor/gecore/lang"

	"github.com/ge-editor/langs/abstract"
)

// ------------------------------------------------------------------
// Fundamental implement gecore Mode interface
// ------------------------------------------------------------------

func NewFundamental() lang.Mode {
	return &Fundamental{
		Worker: abstract.Worker{
			Exts: []string{}, // Initialize the `Exts` field of abstract.Worker
		},
	}
}

type Fundamental struct {
	abstract.Worker
}

func (f *Fundamental) Name() string {
	return "Fundamental"
}
