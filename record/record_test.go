package record_test

import (
	"bytes"
	"testing"

	"github.com/decred/dcrlnd/lnwire"
	"github.com/decred/dcrlnd/record"
	"github.com/decred/dcrlnd/tlv"
)

type recordEncDecTest struct {
	name      string
	encRecord func() tlv.RecordProducer
	decRecord func() tlv.RecordProducer
	assert    func(*testing.T, interface{})
}

var (
	testTotal = lnwire.MilliAtom(45)
	testAddr  = [32]byte{0x01, 0x02}
)

var recordEncDecTests = []recordEncDecTest{
	{
		name: "mpp",
		encRecord: func() tlv.RecordProducer {
			return record.NewMPP(testTotal, testAddr)
		},
		decRecord: func() tlv.RecordProducer {
			return new(record.MPP)
		},
		assert: func(t *testing.T, r interface{}) {
			mpp := r.(*record.MPP)
			if mpp.TotalMAtoms() != testTotal {
				t.Fatal("incorrect total msat")
			}
			if mpp.PaymentAddr() != testAddr {
				t.Fatal("incorrect payment addr")
			}
		},
	},
}

// TestRecordEncodeDecode is a generic test framework for custom TLV records. It
// asserts that records can encode and decode themselves, and that the value of
// the original record matches the decoded record.
func TestRecordEncodeDecode(t *testing.T) {
	for _, test := range recordEncDecTests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			r := test.encRecord()
			r2 := test.decRecord()
			encStream := tlv.MustNewStream(r.Record())
			decStream := tlv.MustNewStream(r2.Record())

			test.assert(t, r)

			var b bytes.Buffer
			err := encStream.Encode(&b)
			if err != nil {
				t.Fatalf("unable to encode record: %v", err)
			}

			err = decStream.Decode(bytes.NewReader(b.Bytes()))
			if err != nil {
				t.Fatalf("unable to decode record: %v", err)
			}

			test.assert(t, r2)
		})
	}
}