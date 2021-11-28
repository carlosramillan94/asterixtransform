package goasterix

import (
	"bytes"
	"github.com/mokhtarimokhtar/goasterix/uap"
	"testing"
)

func TestFixed_Payload(t *testing.T) {
	// Arrange
	fixed := new(Fixed)
	fixed.Data = []byte{0xff, 0xff, 0xff, 0xff}
	output := []byte{0xff, 0xff, 0xff, 0xff}

	// Act
	b := fixed.Payload()

	// Assert
	if len(b) != 4 {
		t.Errorf("FAIL: len(items) = %v; Expected: %v", len(b), 4)
	} else {
		t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(b), 4)
	}
	if bytes.Equal(b, output) == false {
		t.Errorf("FAIL: sp = % X; Expected: % X", b, output)
	} else {
		t.Logf("SUCCESS: sp = % X; Expected: % X", b, output)
	}
}

func TestExtended_Payload(t *testing.T) {
	// Arrange
	ext := new(Extended)
	ext.Primary = []byte{0xff}
	ext.Secondary = []byte{0xff, 0xff, 0xfe}
	output := []byte{0xff, 0xff, 0xff, 0xfe}

	// Act
	b := ext.Payload()

	// Assert
	if len(b) != 4 {
		t.Errorf("FAIL: len(items) = %v; Expected: %v", len(b), 4)
	} else {
		t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(b), 4)
	}
	if bytes.Equal(b, output) == false {
		t.Errorf("FAIL: sp = % X; Expected: % X", b, output)
	} else {
		t.Logf("SUCCESS: sp = % X; Expected: % X", b, output)
	}
}

func TestExplicit_Payload(t *testing.T) {
	// Arrange
	exp := new(Explicit)
	exp.Len = 0x04
	exp.Data = []byte{0xff, 0xff, 0xfe}
	output := []byte{0x04, 0xff, 0xff, 0xfe}

	// Act
	b := exp.Payload()

	// Assert
	if len(b) != 4 {
		t.Errorf("FAIL: len(items) = %v; Expected: %v", len(b), 4)
	} else {
		t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(b), 4)
	}
	if bytes.Equal(b, output) == false {
		t.Errorf("FAIL: sp = % X; Expected: % X", b, output)
	} else {
		t.Logf("SUCCESS: sp = % X; Expected: % X", b, output)
	}
}

func TestRepetitive_Payload(t *testing.T) {
	// Arrange
	rp := new(Repetitive)
	rp.Rep = 0x03
	rp.Data = []byte{0xff, 0xff, 0xfe}

	output := []byte{0x03, 0xff, 0xff, 0xfe}
	// Act
	b := rp.Payload()

	// Assert
	if len(b) != 4 {
		t.Errorf("FAIL: len(items) = %v; Expected: %v", len(b), 4)
	} else {
		t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(b), 4)
	}
	if bytes.Equal(b, output) == false {
		t.Errorf("FAIL: sp = % X; Expected: % X", b, output)
	} else {
		t.Logf("SUCCESS: sp = % X; Expected: % X", b, output)
	}
}

func TestCompound_Payload(t *testing.T) {
	// Arrange
	cp := new(Compound)
	cp.Primary = []byte{0xf0}
	cp.Secondary = []Item{
		{
			Meta: MetaItem{
				Type: uap.Fixed,
			},
			Fixed: &Fixed{
				Data: []byte{0xff},
			},
		},
		{
			Meta: MetaItem{
				Type: uap.Extended,
			},
			Extended: &Extended{
				Primary:   []byte{0xff},
				Secondary: []byte{0xff, 0xfe},
			},
		},
		{
			Meta: MetaItem{
				Type: uap.Explicit,
			},
			Explicit: &Explicit{
				Len:  0x02,
				Data: []byte{0xff},
			},
		},
		{
			Meta: MetaItem{
				Type: uap.Repetitive,
			},
			Repetitive: &Repetitive{
				Rep:  0x02,
				Data: []byte{0xff, 0xff},
			},
		},
	}
	output := []byte{0xf0, 0xff, 0xff, 0xff, 0xfe, 0x02, 0xff, 0x02, 0xff, 0xff}

	// Act
	b := cp.Payload()

	// Assert
	if len(b) != 10 {
		t.Errorf("FAIL: len(items) = %v; Expected: %v", len(b), 10)
	} else {
		t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(b), 10)
	}
	if bytes.Equal(b, output) == false {
		t.Errorf("FAIL: cp = % X; Expected: % X", b, output)
	} else {
		t.Logf("SUCCESS: cp = % X; Expected: % X", b, output)
	}
}

func TestItem_Payload(t *testing.T) {
	// setup
	type dataTest struct {
		TestCaseName string
		input        Item
		output       []byte
		len          int
	}
	dataSet := []dataTest{
		{
			TestCaseName: "testcase 1",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Fixed,
				},
				Fixed: &Fixed{
					Data: []byte{0xff, 0xff},
				},
			},
			output: []byte{0xff, 0xff},
			len:    2,
		},
		{
			TestCaseName: "testcase 2",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Extended,
				},
				Extended: &Extended{
					Primary:   []byte{0xff},
					Secondary: []byte{0xff, 0xfe},
				},
			},
			output: []byte{0xff, 0xff, 0xfe},
			len:    3,
		},
		{
			TestCaseName: "testcase 3",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Explicit,
				},
				Explicit: &Explicit{
					Len:  0x04,
					Data: []byte{0xff, 0xff, 0xff},
				},
			},
			output: []byte{0x04, 0xff, 0xff, 0xff},
			len:    4,
		},
		{
			TestCaseName: "testcase 4",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Repetitive,
				},
				Repetitive: &Repetitive{
					Rep:  0x02,
					Data: []byte{0xff, 0xff},
				},
			},
			output: []byte{0x02, 0xff, 0xff},
			len:    3,
		},
		{
			TestCaseName: "testcase 5",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Compound,
				},
				Compound: &Compound{
					Primary: []byte{0xc0},
					Secondary: []Item{
						{
							Meta: MetaItem{
								FRN:         1,
								DataItem:    "I000/010",
								Description: "Test item",
								Type:        uap.Fixed,
							},
							Fixed: &Fixed{
								Data: []byte{0xff, 0xff},
							},
						},
						{
							Meta: MetaItem{
								FRN:         1,
								DataItem:    "I000/010",
								Description: "Test item",
								Type:        uap.Fixed,
							},
							Fixed: &Fixed{
								Data: []byte{0xff, 0xff},
							},
						},
					},
				},
			},
			output: []byte{0xc0, 0xff, 0xff, 0xff, 0xff},
			len:    5,
		},
	}
	for _, row := range dataSet {
		// Arrange
		// Act
		b := row.input.Payload()

		// Assert
		if len(b) != row.len {
			t.Errorf("FAIL: len(items) = %v; Expected: %v", len(b), row.len)
		} else {
			t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(b), row.len)
		}
		if bytes.Equal(b, row.output) == false {
			t.Errorf("FAIL: item = % X; Expected: % X", b, row.output)
		} else {
			t.Logf("SUCCESS: item = % X; Expected: % X", b, row.output)
		}
	}

}

func TestItem_String(t *testing.T) {
	// setup
	type dataTest struct {
		TestCaseName string
		input        Item
		output       string
		len          int
	}
	dataSet := []dataTest{
		{
			TestCaseName: "testcase 1",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Fixed,
				},
				Fixed: &Fixed{
					Data: []byte{0xff, 0xff},
				},
			},
			output: "ffff",
			len:    10 + 4,
		},
		{
			TestCaseName: "testcase 2",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Extended,
				},
				Extended: &Extended{
					Primary:   []byte{0xff},
					Secondary: []byte{0xff, 0xfe},
				},
			},
			output: "fffffe",
			len:    10 + 6,
		},
		{
			TestCaseName: "testcase 3",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Explicit,
				},
				Explicit: &Explicit{
					Len:  0x04,
					Data: []byte{0xff, 0xff, 0xff},
				},
			},
			output: "04ffffff",
			len:    10 + 8,
		},
		{
			TestCaseName: "testcase 4",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Repetitive,
				},
				Repetitive: &Repetitive{
					Rep:  0x02,
					Data: []byte{0xff, 0xff},
				},
			},
			output: "02ffff",
			len:    10 + 6,
		},
		{
			TestCaseName: "testcase 5",
			input: Item{
				Meta: MetaItem{
					FRN:         1,
					DataItem:    "I000/010",
					Description: "Test item",
					Type:        uap.Compound,
				},
				Compound: &Compound{
					Primary: []byte{0xc0},
					Secondary: []Item{
						{
							Meta: MetaItem{
								FRN:         1,
								DataItem:    "I000/010",
								Description: "Test item",
								Type:        uap.Fixed,
							},
							Fixed: &Fixed{
								Data: []byte{0xff, 0xff},
							},
						},
						{
							Meta: MetaItem{
								FRN:         1,
								DataItem:    "I000/010",
								Description: "Test item",
								Type:        uap.Fixed,
							},
							Fixed: &Fixed{
								Data: []byte{0xff, 0xff},
							},
						},
					},
				},
			},
			output: "c0ffffffff",
			len:    10 + 10,
		},
	}
	for _, row := range dataSet {
		// Arrange
		// Act
		s := row.input.String()

		// Assert
		if len(s) != row.len {
			t.Errorf("FAIL: len(items) = %v; Expected: %v", len(s), row.len)
		} else {
			t.Logf("SUCCESS: len(items) = %v; Expected: %v", len(s), row.len)
		}
		if s == row.output {
			t.Errorf("FAIL: item = %s; Expected: %s", s, row.output)
		} else {
			t.Logf("SUCCESS: item = %s; Expected: %s", s, row.output)
		}
	}

}
