package options

import (
	. "github.com/miekg/dns"
	_ "gopkg.in/dealancer/validate.v2"
)

// SVCBRecord holds a SVCB RR not ServiceParams because the dns library doesn't
// have an API capable of packing SvcParams directly. So we parse SVCB,
// pack it, and calculate an offset to get the SvcParams when Serializing
type SVCBRecord struct {
	Record *SVCB
}

func (d *SVCBRecord) UnmarshalText(text []byte) error {
	sbStr := ". 1 SVCB 10 . " + string(text)
	r, err := NewRR(sbStr)
	if err != nil {
		return err
	}
	d.Record = r.(*SVCB)
	return nil
}

func (d *SVCBRecord) Serialize() ([]byte, error) {
	var outBuf = make([]byte, Len(d.Record))

	_, err := PackRR(d.Record, outBuf, 0, nil, false)
	if err != nil {
		return nil, err
	}
	// SvcB starts at index 11 (Priority), SvcParams begin at index 14
	return outBuf[14:], nil
}
