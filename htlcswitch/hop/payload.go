package hop

import (
	"encoding/binary"
	"io"

	"github.com/decred/dcrlnd/lnwire"
	"github.com/decred/dcrlnd/record"
	"github.com/decred/dcrlnd/tlv"
	"github.com/decred/lightning-onion/v2"
)

// Payload encapsulates all information delivered to a hop in an onion payload.
// A Hop can represent either a TLV or legacy payload. The primary forwarding
// instruction can be accessed via ForwardingInfo, and additional records can be
// accessed by other member functions.
type Payload struct {
	// FwdInfo holds the basic parameters required for HTLC forwarding, e.g.
	// amount, cltv, and next hop.
	FwdInfo ForwardingInfo
}

// NewLegacyPayload builds a Payload from the amount, cltv, and next hop
// parameters provided by leegacy onion payloads.
func NewLegacyPayload(f *sphinx.HopData) *Payload {
	nextHop := binary.BigEndian.Uint64(f.NextAddress[:])

	return &Payload{
		FwdInfo: ForwardingInfo{
			Network:         DecredNetwork,
			NextHop:         lnwire.NewShortChanIDFromInt(nextHop),
			AmountToForward: lnwire.MilliAtom(f.ForwardAmount),
			OutgoingCTLV:    f.OutgoingCltv,
		},
	}
}

// NewPayloadFromReader builds a new Hop from the passed io.Reader. The reader
// should correspond to the bytes encapsulated in a TLV onion payload.
func NewPayloadFromReader(r io.Reader) (*Payload, error) {
	var (
		cid  uint64
		amt  uint64
		cltv uint32
	)

	tlvStream, err := tlv.NewStream(
		record.NewAmtToFwdRecord(&amt),
		record.NewLockTimeRecord(&cltv),
		record.NewNextHopIDRecord(&cid),
	)
	if err != nil {
		return nil, err
	}

	_, err = tlvStream.DecodeWithParsedTypes(r)
	if err != nil {
		return nil, err
	}

	return &Payload{
		FwdInfo: ForwardingInfo{
			Network:         DecredNetwork,
			NextHop:         lnwire.NewShortChanIDFromInt(cid),
			AmountToForward: lnwire.MilliAtom(amt),
			OutgoingCTLV:    cltv,
		},
	}, nil
}

// ForwardingInfo returns the basic parameters required for HTLC forwarding,
// e.g. amount, cltv, and next hop.
func (h *Payload) ForwardingInfo() ForwardingInfo {
	return h.FwdInfo
}
