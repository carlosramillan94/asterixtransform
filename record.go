package goasterix

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"math/bits"
	"strings"

	"github.com/mokhtarimokhtar/goasterix/uap"
)

const (
	fixed      string = "fixed"
	extended   string = "extended"
	compound   string = "compound"
	repetitive string = "repetitive"
	explicit   string = "explicit"
	sp         string = "sp"
	re         string = "re"
	rfs        string = "rfs"
)

var (
	// ErrDatafieldUnknown reports which ErrDatafield Unknown.
	ErrDatafieldUnknown = errors.New("[Items] Type of Datafield Not found")
)

type Record struct {
	Fspec []byte
	Items []uap.DataField
}

// Decode extracts a Record of asterix data block (only one record).
// An asterix data block can contain a or more records.
// It returns the number of bytes unread and fills the Record Struct(Fspec, Items array) in byte.
func (rec *Record) Decode(data []byte, uap []uap.DataField) (unRead int, err error) {
	rb := bytes.NewReader(data)

	rec.Fspec, err = FspecReader(rb, 1)
	unRead = rb.Len()
	if err != nil {
		return unRead, err
	}

	frnIndex, _ := FspecIndex(rec.Fspec)

	for _, frn := range frnIndex {
		dataItem := uap[frn-1] // here the index corresponds to the FRN
		var tmp []byte

		switch strings.ToLower(dataItem.Type.Name) {
		case fixed:
			tmp, err = FixedDataFieldReader(rb, dataItem.Type.Size)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		case extended:
			tmp, err = ExtendedDataFieldReader(rb, dataItem.Type.Size)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		case compound:
			tmp, err = CompoundDataFieldReader(rb, dataItem.Type.Meta)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		case repetitive:
			tmp, err = RepetitiveDataFieldReader(rb, dataItem.Type.Size)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		case explicit:
			tmp, err = ExplicitDataFieldReader(rb)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		case sp, re:
			tmp, err = SPAndREDataFieldReader(rb)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		case rfs:
			tmp, err = RFSDataFieldReader(rb, uap)
			if err != nil {
				unRead = rb.Len()
				return unRead, err
			}

		default:
			err = ErrDatafieldUnknown
			return unRead, err
		}

		dataItem.Payload = tmp
		rec.Items = append(rec.Items, dataItem)
		unRead = rb.Len()
	}

	return unRead, nil
}

// Payload returns a slice of byte for one asterix record.
func (rec *Record) Payload() (b []byte) {
	b = append(b, rec.Fspec...)
	for _, item := range rec.Items {
		b = append(b, item.Payload...)
	}
	return b
}

// String returns a string(hex) representation of one asterix record (only existing items).
func (rec *Record) String() (items []string) {
	tmp := "FSPEC: " + hex.EncodeToString(rec.Fspec)
	items = append(items, tmp)

	for _, item := range rec.Items {
		tmp := item.DataItem + ": " + hex.EncodeToString(item.Payload)
		items = append(items, tmp)
	}
	return items
}

// FspecReader returns a slice of FSPEC asterix data.
// The step parameter defines the read jumps of the FX field.
func FspecReader(reader io.Reader, step uint8) (fspec []byte, err error) {
	for {
		tmp := make([]byte, step)
		err = binary.Read(reader, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		fspec = append(fspec, tmp[0])
		if tmp[0]&0x01 == 0 {
			break
		}
	}
	return fspec, err
}

// FspecIndex returns an array of uint8 corresponding to number FRN(Field Reference Number of Items).
// In other words, it converts a fspec bits to an array FRNs.
func FspecIndex(fspec []byte) (frnIndex []uint8, err error) {
	for j, val := range fspec {
		for i := 0; i < 7; i++ {
			frn := 7*j + i + 1
			tmp := bits.RotateLeft8(val, i)
			if tmp&0x80 != 0 {
				frnIndex = append(frnIndex, uint8(frn))
			}
		}
	}
	return frnIndex, nil
}

// FixedDataFieldReader extracts a number(nb) of bytes(size) and returns a slice of bytes(data of item).
// Fixed length Data Fields shall comprise a fixed number of octets.
func FixedDataFieldReader(rb *bytes.Reader, size uint8) (item []byte, err error) {
	for i := uint8(0); i < size; i++ {
		var tmp uint8
		err := binary.Read(rb, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		item = append(item, tmp)
	}
	return item, err
}

// ExtendedDataFieldReader extracts data item type Extended (FX: last bit = 1).
// Size parameter defines the size of extended field.
// Extended length Data Fields, being of a variable length, shall contain a primary part of predetermined length,
// immediately followed by a number of secondary parts, each of predetermined length.
// The presence of the next following secondary part shall be indicated by the setting to one of the
// Least Significant Bit (LSB) of the last octet of the preceding part (either the primary part or a secondary part).
// This bit which is reserved for that purpose is called the Field Extension Indicator (FX).
func ExtendedDataFieldReader(rb *bytes.Reader, size uint8) (item []byte, err error) {
	for {
		tmp := make([]byte, size)
		err = binary.Read(rb, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		item = append(item, tmp...)
		if tmp[size-1]&0x01 == 0 {
			break
		}
	}
	return item, err
}

// ExplicitDataFieldReader extracts a number of bytes define by the first byte.
// Explicit length Data Fields shall start with a one-octet length indicator giving
// the total field length in octets including the length indicator itself.
func ExplicitDataFieldReader(rb *bytes.Reader) (item []byte, err error) {
	l := make([]byte, 1)
	err = binary.Read(rb, binary.BigEndian, &l)
	if err != nil {
		return nil, err
	}

	tmp := make([]byte, l[0]-1)
	err = binary.Read(rb, binary.BigEndian, &tmp)
	if err != nil {
		return nil, err
	}

	item = append(item, l[0])
	item = append(item, tmp...)
	return item, err
}

// RepetitiveDataFieldReader extracts data item type Repetitive(1+rep*N byte).
// The first byte is REP(factor), nb is the size of bytes to repetition.
// Repetitive Data Fields, being of a variable length, shall comprise a one-octet Field Repetition Indicator (REP)
// signalling the presence of N consecutive sub-fields each of the same pre-determined length.
func RepetitiveDataFieldReader(rb *bytes.Reader, size uint8) (item []byte, err error) {
	rep := make([]byte, 1)
	err = binary.Read(rb, binary.BigEndian, &rep)
	if err != nil {
		return nil, err
	}

	tmp := make([]byte, rep[0]*size)
	err = binary.Read(rb, binary.BigEndian, &tmp)
	if err != nil {
		return nil, err
	}
	item = append(item, rep[0])
	item = append(item, tmp...)

	return item, err
}

// CompoundDataFieldReader
// Compound Data Fields, being of a variable length, shall comprise a primary subfield, followed by data subfields.
// The primary subfield determines the presence or absence of the subsequent data subfields. It comprises a first part
// of one octet extendable using the Field Extension (FX) mechanism.
// The definition, structure and format of the data subfields are part of the description of the relevant Compound Data
// Item. Data subfields shall be either fixed length, extended length, explicit length or repetitive, but not compound.
//func CompoundDataFieldReader(rb *bytes.Reader, prm uap.PrimaryField) (item []byte, err error) {
func CompoundDataFieldReader(rb *bytes.Reader, sub uap.MetaField) (item []byte, err error) {
	var primaries []byte
	for {
		tmp := make([]byte, 1)
		err = binary.Read(rb, binary.BigEndian, &tmp)
		if err != nil {
			return nil, err
		}
		primaries = append(primaries, tmp[0])
		if tmp[0]&0x01 == 0 {
			break
		}
	}
	item = append(item, primaries...)

	for _, primary := range primaries {
		var tmp []byte
		if primary&0x80 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[8])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
		if primary&0x40 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[7])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
		if primary&0x20 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[6])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
		if primary&0x10 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[5])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
		if primary&0x08 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[4])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
		if primary&0x04 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[3])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
		if primary&0x02 != 0 {
			tmp, err = SelectTypeFieldReader(rb, sub[2])
			if err != nil {
				return nil, err
			}
			item = append(item, tmp...)
		}
	}

	return item, err
}

func SelectTypeFieldReader(rb *bytes.Reader, sub uap.Subfield) (item []byte, err error) {

	typeOfField := strings.ToLower(sub.Name)
	switch typeOfField {
	case fixed:
		item, err = FixedDataFieldReader(rb, sub.Size)
		if err != nil {
			return nil, err
		}
	case repetitive:
		item, err = RepetitiveDataFieldReader(rb, sub.Size)
		if err != nil {
			return nil, err
		}
	case extended:
		item, err = ExtendedDataFieldReader(rb, sub.Size)
		if err != nil {
			return nil, err
		}
	case explicit:
		item, err = ExplicitDataFieldReader(rb)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrDatafieldUnknown
	}

	return item, err
}

// SPAndREDataFieldReader extracts returns a slice
// ref. EUROCONTROL-SPEC-0149 2.4
// 4.3.5 Non-Standard Data Fields:
// Reserved Expansion Data
// Field Special Purpose field
func SPAndREDataFieldReader(rb *bytes.Reader) (item []byte, err error) {
	l := make([]byte, 1)
	err = binary.Read(rb, binary.BigEndian, &l)
	if err != nil {
		return nil, err
	}

	tmp := make([]byte, l[0]-1)
	err = binary.Read(rb, binary.BigEndian, &tmp)
	if err != nil {
		return nil, err
	}

	item = append(item, l[0])
	item = append(item, tmp...)
	return item, err
}

// RFSDataFieldReader extracts Random Field Sequencing part and returns an array of byte(data item).
func RFSDataFieldReader(rb *bytes.Reader, uap []uap.DataField) (item []byte, err error) {
	// total is the number of datafields
	var n uint8
	err = binary.Read(rb, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}

	item = append(item, n)

	for i := uint8(0); i < n; i++ {
		// random FRN
		var frn uint8
		err := binary.Read(rb, binary.BigEndian, &frn)
		if err != nil {
			return nil, err
		}
		item = append(item, frn)

		for _, field := range uap {
			if frn == field.FRN {
				// todo: work just for Fixed datafield use case
				tmp := make([]byte, field.Type.Size)

				err := binary.Read(rb, binary.BigEndian, &tmp)
				if err != nil {
					return nil, err
				}

				item = append(item, tmp...)
			}
		}
	}

	return item, err
}
