package signify

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/agl/ed25519"
)

type readfiletest struct {
	file    string
	comment string
	content []byte
	parsed  interface{}
}

var readfiletests = []readfiletest{
	{
		file:    "_testdata/test.key",
		comment: "signify secret key",
		content: []byte{
			0x45, 0x64, 0x42, 0x4b, 0x00, 0x00, 0x00, 0x2a, 0xbb, 0x07, 0x17, 0x79, 0xb5, 0x84, 0x56, 0xe5,
			0xf6, 0x61, 0xdc, 0xe0, 0x44, 0x7b, 0x98, 0xd7, 0x42, 0x42, 0xc0, 0x8d, 0xc7, 0xc0, 0x52, 0x16,
			0xd3, 0xff, 0xb0, 0x73, 0xe8, 0x92, 0x09, 0x30, 0xd3, 0xdb, 0x4f, 0x63, 0xb2, 0x59, 0xa4, 0x78,
			0x26, 0x5a, 0x50, 0x04, 0xd3, 0x5a, 0xb5, 0xf8, 0x92, 0xb2, 0x75, 0x4c, 0x30, 0x12, 0x12, 0x63,
			0x6f, 0x15, 0x29, 0xd9, 0xdf, 0x41, 0x4c, 0xde, 0x4c, 0x14, 0x60, 0xb9, 0xb1, 0x14, 0x1c, 0xbc,
			0xc3, 0xde, 0xd1, 0xe7, 0x79, 0x6d, 0xd0, 0x12, 0xd7, 0xed, 0x92, 0x88, 0xf4, 0xf1, 0x6a, 0x2f,
			0x13, 0x38, 0x3d, 0x60, 0xb9, 0x35, 0x43, 0xd5},
		parsed: rawEncryptedKey{
			PKAlgo:    [2]byte{'E', 'd'},
			KDFAlgo:   [2]byte{'B', 'K'},
			KDFRounds: 42,
			Salt: [16]byte{
				0xbb, 0x07, 0x17, 0x79, 0xb5, 0x84, 0x56, 0xe5, 0xf6, 0x61, 0xdc, 0xe0, 0x44, 0x7b, 0x98, 0xd7,
			},
			Checksum:    [8]byte{0x42, 0x42, 0xc0, 0x8d, 0xc7, 0xc0, 0x52, 0x16},
			Fingerprint: [8]byte{0xd3, 0xff, 0xb0, 0x73, 0xe8, 0x92, 0x09, 0x30},
			PrivateKey: [ed25519.PrivateKeySize]byte{
				0xd3, 0xdb, 0x4f, 0x63, 0xb2, 0x59, 0xa4, 0x78,
				0x26, 0x5a, 0x50, 0x04, 0xd3, 0x5a, 0xb5, 0xf8, 0x92, 0xb2, 0x75, 0x4c, 0x30, 0x12, 0x12, 0x63,
				0x6f, 0x15, 0x29, 0xd9, 0xdf, 0x41, 0x4c, 0xde, 0x4c, 0x14, 0x60, 0xb9, 0xb1, 0x14, 0x1c, 0xbc,
				0xc3, 0xde, 0xd1, 0xe7, 0x79, 0x6d, 0xd0, 0x12, 0xd7, 0xed, 0x92, 0x88, 0xf4, 0xf1, 0x6a, 0x2f,
				0x13, 0x38, 0x3d, 0x60, 0xb9, 0x35, 0x43, 0xd5,
			},
		},
	}, {
		file:    "_testdata/test.pub",
		comment: "signify public key",
		content: []byte{
			0x45, 0x64, 0xd3, 0xff, 0xb0, 0x73, 0xe8, 0x92, 0x09, 0x30, 0xc8, 0x02, 0xe8, 0xf6, 0x4c, 0x35,
			0x63, 0xc2, 0x2e, 0xa3, 0x03, 0x56, 0xaf, 0x63, 0xf6, 0x92, 0xce, 0x2a, 0x63, 0x5c, 0xf6, 0x6e,
			0x7d, 0x48, 0x6c, 0xa8, 0x48, 0x8d, 0xe2, 0x04, 0xa6, 0x05},
		parsed: rawPublicKey{
			PKAlgo:      [2]byte{'E', 'd'},
			Fingerprint: [8]byte{0xd3, 0xff, 0xb0, 0x73, 0xe8, 0x92, 0x09, 0x30},
			PublicKey: [ed25519.PublicKeySize]byte{
				0xc8, 0x02, 0xe8, 0xf6, 0x4c, 0x35, 0x63, 0xc2, 0x2e, 0xa3, 0x03, 0x56, 0xaf, 0x63, 0xf6, 0x92,
				0xce, 0x2a, 0x63, 0x5c, 0xf6, 0x6e, 0x7d, 0x48, 0x6c, 0xa8, 0x48, 0x8d, 0xe2, 0x04, 0xa6, 0x05,
			},
		},
	}, {
		file:    "_testdata/test.msg.sig",
		comment: "signature from signify secret key",
		content: []byte{
			0x45, 0x64, 0xd3, 0xff, 0xb0, 0x73, 0xe8, 0x92, 0x09, 0x30, 0x9e, 0x9f, 0x91, 0x69, 0x08, 0x5d,
			0xa7, 0xb9, 0x1c, 0x82, 0x3c, 0x81, 0x69, 0x16, 0x16, 0x58, 0x7a, 0xd2, 0x53, 0xb4, 0xe9, 0x96,
			0x0b, 0x42, 0x3c, 0x8a, 0x40, 0x40, 0x47, 0x7e, 0xb0, 0x41, 0x74, 0x26, 0x47, 0x41, 0xa4, 0xe8,
			0x2f, 0xec, 0xfb, 0xde, 0xe2, 0x77, 0x58, 0x19, 0xca, 0xb0, 0x57, 0x5f, 0x73, 0x5f, 0x8b, 0xe2,
			0xac, 0x11, 0x00, 0x14, 0x55, 0xd6, 0xac, 0xd3, 0xd3, 0x03},
		parsed: rawSignature{
			PKAlgo:      [2]byte{'E', 'd'},
			Fingerprint: [8]byte{0xd3, 0xff, 0xb0, 0x73, 0xe8, 0x92, 0x09, 0x30},
			Signature: [ed25519.SignatureSize]byte{
				0x9e, 0x9f, 0x91, 0x69, 0x08, 0x5d, 0xa7, 0xb9, 0x1c, 0x82, 0x3c, 0x81, 0x69, 0x16, 0x16, 0x58,
				0x7a, 0xd2, 0x53, 0xb4, 0xe9, 0x96, 0x0b, 0x42, 0x3c, 0x8a, 0x40, 0x40, 0x47, 0x7e, 0xb0, 0x41,
				0x74, 0x26, 0x47, 0x41, 0xa4, 0xe8, 0x2f, 0xec, 0xfb, 0xde, 0xe2, 0x77, 0x58, 0x19, 0xca, 0xb0,
				0x57, 0x5f, 0x73, 0x5f, 0x8b, 0xe2, 0xac, 0x11, 0x00, 0x14, 0x55, 0xd6, 0xac, 0xd3, 0xd3, 0x03,
			},
		},
	}, {
		file:    "_testdata/test.nopass.key",
		comment: "nopass secret key",
		content: []byte{
			0x45, 0x64, 0x42, 0x4b, 0x00, 0x00, 0x00, 0x00, 0xf3, 0xff, 0x55, 0x86, 0xf2, 0x22, 0x74, 0xf4,
			0x35, 0x0f, 0xfc, 0x03, 0x2d, 0x31, 0x36, 0xb9, 0xb5, 0xc8, 0x61, 0xc5, 0xaf, 0x95, 0xb3, 0x8d,
			0x1f, 0x18, 0x5a, 0xca, 0x53, 0x7e, 0xd1, 0x45, 0xbb, 0x66, 0x94, 0xc0, 0x9f, 0x65, 0x2a, 0xae,
			0x0f, 0x24, 0x82, 0x9e, 0xbe, 0xae, 0xac, 0x9f, 0xec, 0x4c, 0x1c, 0xd9, 0x39, 0x1d, 0x3e, 0x4f,
			0x68, 0x61, 0x07, 0xf8, 0x50, 0x07, 0x1b, 0xc1, 0x94, 0xd0, 0x2a, 0x16, 0x22, 0xf2, 0x99, 0x28,
			0xea, 0x15, 0x9f, 0xbc, 0xd4, 0x2b, 0x69, 0x71, 0x03, 0xed, 0xa8, 0xd3, 0x7a, 0x59, 0x54, 0x82,
			0xeb, 0x2f, 0x87, 0xf4, 0x4d, 0x52, 0xb6, 0x49},
		parsed: rawEncryptedKey{
			PKAlgo:      [2]byte{'E', 'd'},
			KDFAlgo:     [2]byte{'B', 'K'},
			KDFRounds:   0,
			Salt:        [16]byte{0xf3, 0xff, 0x55, 0x86, 0xf2, 0x22, 0x74, 0xf4, 0x35, 0xf, 0xfc, 0x03, 0x2d, 0x31, 0x36, 0xb9},
			Checksum:    [8]byte{0xb5, 0xc8, 0x61, 0xc5, 0xaf, 0x95, 0xb3, 0x8d},
			Fingerprint: [8]byte{0x1f, 0x18, 0x5a, 0xca, 0x53, 0x7e, 0xd1, 0x45},
			PrivateKey: [ed25519.PrivateKeySize]byte{
				0xbb, 0x66, 0x94, 0xc0, 0x9f, 0x65, 0x2a, 0xae, 0x0f, 0x24, 0x82, 0x9e, 0xbe, 0xae, 0xac, 0x9f,
				0xec, 0x4c, 0x1c, 0xd9, 0x39, 0x1d, 0x3e, 0x4f, 0x68, 0x61, 0x07, 0xf8, 0x50, 0x07, 0x1b, 0xc1,
				0x94, 0xd0, 0x2a, 0x16, 0x22, 0xf2, 0x99, 0x28, 0xea, 0x15, 0x9f, 0xbc, 0xd4, 0x2b, 0x69, 0x71,
				0x03, 0xed, 0xa8, 0xd3, 0x7a, 0x59, 0x54, 0x82, 0xeb, 0x2f, 0x87, 0xf4, 0x4d, 0x52, 0xb6, 0x49},
		},
	}, {
		file:    "_testdata/test.nopass.pub",
		comment: "nopass public key",
		content: []byte{
			0x45, 0x64, 0x1f, 0x18, 0x5a, 0xca, 0x53, 0x7e, 0xd1, 0x45, 0x94, 0xd0, 0x2a, 0x16, 0x22, 0xf2,
			0x99, 0x28, 0xea, 0x15, 0x9f, 0xbc, 0xd4, 0x2b, 0x69, 0x71, 0x03, 0xed, 0xa8, 0xd3, 0x7a, 0x59,
			0x54, 0x82, 0xeb, 0x2f, 0x87, 0xf4, 0x4d, 0x52, 0xb6, 0x49},
		parsed: rawPublicKey{
			PKAlgo:      [2]byte{'E', 'd'},
			Fingerprint: [8]byte{0x1f, 0x18, 0x5a, 0xca, 0x53, 0x7e, 0xd1, 0x45},
			PublicKey: [ed25519.PublicKeySize]byte{
				0x94, 0xd0, 0x2a, 0x16, 0x22, 0xf2, 0x99, 0x28, 0xea, 0x15, 0x9f, 0xbc, 0xd4, 0x2b, 0x69, 0x71,
				0x03, 0xed, 0xa8, 0xd3, 0x7a, 0x59, 0x54, 0x82, 0xeb, 0x2f, 0x87, 0xf4, 0x4d, 0x52, 0xb6, 0x49},
		},
	}, {
		file:    "_testdata/test.nopass.msg.sig",
		comment: "signature from nopass secret key",
		content: []byte{
			0x45, 0x64, 0x1f, 0x18, 0x5a, 0xca, 0x53, 0x7e, 0xd1, 0x45, 0x12, 0xb1, 0xb1, 0xc8, 0xf5, 0xde,
			0xaf, 0xd4, 0xbb, 0x74, 0x14, 0x1d, 0x33, 0x7c, 0xbc, 0x0e, 0x83, 0xd6, 0x1f, 0x9e, 0x71, 0x9b,
			0x21, 0x0c, 0x1f, 0x79, 0x6e, 0x6d, 0x44, 0xb2, 0xe2, 0xfd, 0x1c, 0x35, 0x41, 0x1d, 0x89, 0xa8,
			0x34, 0x6f, 0x91, 0x11, 0x3e, 0xc3, 0xee, 0x28, 0x56, 0x59, 0x4a, 0xcb, 0x28, 0xa1, 0x14, 0xf8,
			0x77, 0x2f, 0x7a, 0x16, 0x58, 0x0e, 0x4f, 0x7e, 0xa6, 0x00},
		parsed: rawSignature{
			PKAlgo:      [2]byte{'E', 'd'},
			Fingerprint: [8]byte{0x1f, 0x18, 0x5a, 0xca, 0x53, 0x7e, 0xd1, 0x45},
			Signature: [ed25519.SignatureSize]byte{
				0x12, 0xb1, 0xb1, 0xc8, 0xf5, 0xde, 0xaf, 0xd4, 0xbb, 0x74, 0x14, 0x1d, 0x33, 0x7c, 0xbc, 0x0e,
				0x83, 0xd6, 0x1f, 0x9e, 0x71, 0x9b, 0x21, 0x0c, 0x1f, 0x79, 0x6e, 0x6d, 0x44, 0xb2, 0xe2, 0xfd,
				0x1c, 0x35, 0x41, 0x1d, 0x89, 0xa8, 0x34, 0x6f, 0x91, 0x11, 0x3e, 0xc3, 0xee, 0x28, 0x56, 0x59,
				0x4a, 0xcb, 0x28, 0xa1, 0x14, 0xf8, 0x77, 0x2f, 0x7a, 0x16, 0x58, 0x0e, 0x4f, 0x7e, 0xa6, 0x00},
		},
	},
}

func testReadFile(t *testing.T, file, comment string, content []byte) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("%s: %s\n", file, err)
	}

	rcomment, rcontent, err := ReadFile(bytes.NewReader(buf))
	if err != nil {
		t.Fatal(err)
	}

	if rcomment != comment {
		t.Errorf("%s: comment\nexpected: %q\ngot %q\n", file, comment, rcomment)
	}

	if !bytes.Equal(rcontent, content) {
		t.Errorf("%s: content\nexpected: %x\ngot %x\n", file, content, rcontent)
	}
}

func TestReadFile(t *testing.T) {
	for _, tc := range readfiletests {
		testReadFile(t, tc.file, tc.comment, tc.content)
	}
}

func TestParsePrivateKey(t *testing.T) {
	for _, tc := range readfiletests {
		want, ok := tc.parsed.(rawEncryptedKey)
		if !ok {
			continue
		}

		ek, err := ParsePrivateKey(tc.content, "")
		if err != nil {
			t.Errorf("%s: %s\n", tc.file, err)
			continue
		}

		if want != *ek {
			t.Errorf("%s: expected: %+v got: %+v\n", tc.file, want, ek)
		}
	}
}

func TestParsePublicKey(t *testing.T) {
	for _, tc := range readfiletests {
		want, ok := tc.parsed.(rawPublicKey)
		if !ok {
			continue
		}

		pub, err := ParsePublicKey(tc.content)
		if err != nil {
			t.Errorf("%s: %s\n", tc.file, err)
			continue
		}

		if want != *pub {
			t.Errorf("%s: expected: %+v got: %+v\n", tc.file, want, pub)
		}
	}
}

func TestParseSignature(t *testing.T) {
	for _, tc := range readfiletests {
		want, ok := tc.parsed.(rawSignature)
		if !ok {
			continue
		}

		sig, err := ParseSignature(tc.content)
		if err != nil {
			t.Errorf("%s: %s\n", tc.file, err)
			continue
		}

		if want != *sig {
			t.Errorf("%s: expected: %+v got: %+v\n", tc.file, want, sig)
		}
	}
}
