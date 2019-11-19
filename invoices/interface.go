package invoices

import (
	"github.com/decred/dcrlnd/htlcswitch/hop"
	"github.com/decred/dcrlnd/record"
)

// Payload abstracts access to any additional fields provided in the final hop's
// TLV onion payload.
type Payload interface {
	// MultiPath returns the record corresponding the option_mpp parsed from
	// the onion payload.
	MultiPath() *record.MPP

	// CustomRecords returns the custom tlv type records that were parsed
	// from the payload.
	CustomRecords() hop.CustomRecordSet
}
