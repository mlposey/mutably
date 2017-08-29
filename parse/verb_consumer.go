package parse

import "anvil/db"

type VerbConsumer struct {
	Key db.KeyRing
}

func (consumer *VerbConsumer) Consume(page Page) (bool, error) {
	// Divide page into language sections
	// For each section
		// If verb
			// Get language
			// Get title (the verb)
			// Get verb type
			// Get inflection rule
			// Add above to VerbDefinition struct
			// Send definition to database method
	return false, nil
}
