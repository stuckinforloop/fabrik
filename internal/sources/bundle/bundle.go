package bundle

import (
	_ "github.com/stuckinforloop/fabrik/internal/sources/google"
	_ "github.com/stuckinforloop/fabrik/internal/sources/linear"
)

// Import imports all the sources
func Import() {
	// this is required to register the data sources.
	// done here to avoid cyclic imports.
}
