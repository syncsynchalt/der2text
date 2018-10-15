package der_test

import (
	. "github.com/syncsynchalt/der2text/read/der"
	"github.com/syncsynchalt/der2text/read/indenter"
	"github.com/syncsynchalt/der2text/test"
	"strconv"
	"strings"
	"testing"
)

// helper function used by all tests below
func testDerOctets(tb testing.TB, inputOctets string, output string) {
	test.CallerDepth = 2
	defer func() { test.CallerDepth = 1 }()

	// strip spaces
	inputOctets = strings.Join(strings.Fields(inputOctets), "")

	// convert octets string to byte array
	bytes := make([]byte, 0)
	for len(inputOctets) > 0 {
		thisOctet := inputOctets[0:2]
		inputOctets = inputOctets[2:]
		byteVal, err := strconv.ParseInt(thisOctet, 16, 32)
		test.Ok(tb, err)
		bytes = append(bytes, byte(byteVal))
	}

	// run Parse, compare output
	parseOut := strings.Builder{}
	out := indenter.New(&parseOut)
	err := Parse(out, bytes)
	if err != nil && err.Error() == output {
		return
	}
	test.Ok(tb, err)
	test.Equals(tb, output, parseOut.String())
}

func TestEmpty(t *testing.T) {
	testDerOctets(t, "", "")
}

func TestShortRead(t *testing.T) {
	testDerOctets(t, "00", "short DER read, need at least two bytes, got 1")
	testDerOctets(t, "00 84", "Can't satisfy request to read 4 bytes to get length")
}

func TestEoc(t *testing.T) {
	testDerOctets(t, "00 00", `UNIVERSAL PRIMITIVE END-OF-CONTENT
`)
	testDerOctets(t, "00 01 00", `End-of-content had unexpected length 1`)
}

func TestBooleans(t *testing.T) {
	testDerOctets(t, "01 01 00", `UNIVERSAL PRIMITIVE BOOLEAN FALSE
`)
	testDerOctets(t, "01 01 01", `UNIVERSAL PRIMITIVE BOOLEAN TRUE
`)
	testDerOctets(t, "01 01 FF", `UNIVERSAL PRIMITIVE BOOLEAN TRUE
`)
	testDerOctets(t, "01 00", `Boolean had unexpected length 0`)
	testDerOctets(t, "01 02 00 00", `Boolean had unexpected length 2`)
}

func TestIntegerSmallPositive(t *testing.T) {
	testDerOctets(t, "02 00", `UNIVERSAL PRIMITIVE INTEGER :
`)
	testDerOctets(t, "02 01 00", `UNIVERSAL PRIMITIVE INTEGER 0
`)
	testDerOctets(t, "02 01 7f", `UNIVERSAL PRIMITIVE INTEGER 127
`)
	testDerOctets(t, "02 02 00 80", `UNIVERSAL PRIMITIVE INTEGER 128
`)
	testDerOctets(t, "02 02 01 00", `UNIVERSAL PRIMITIVE INTEGER 256
`)
}

func TestIntegerNegative(t *testing.T) {
	// -128
	testDerOctets(t, "02 01 80", `UNIVERSAL PRIMITIVE INTEGER :80
`)
	// -129
	testDerOctets(t, "02 02 ff 7f", `UNIVERSAL PRIMITIVE INTEGER :FF7F
`)
}

func TestIntegerLarge(t *testing.T) {
	// large number from a random RSA key
	testDerOctets(t, "02 82 01 01 00A6CEC888FFB8116BFF3562E3638D797BAABC39D9E227A1319C48EBCE60452ADEC36D66DE4F2FC68A4A54082A9CFE761EB6A59C185E42110B357FDD4C1E8DDC35308A7B72B732E7983DDA2493A0407253D29AD9E44F05C896C4DBED21CDF913E828F4C1E8E6D3F973E64F0050BDD70457796C278F445F4E0F006ADA91B2180ECD5B4E57EB272BD4D2C9855FEFF657D5DE01A261F46FEE26839E3AD1D52577CAA2BC08EE70559380BE7BCE6696A5EEAABEBB8C8367B772CA155BA2C258CCEC5B41EAF472D2D0E8AAA8F2A4985178A67F82076E3E9F2301AD6622F41F938F1D1F6F51A1CCEF50C91C06861487120C7C82A9E33FBDDF00F8DE6060EE248F4815E92B", `UNIVERSAL PRIMITIVE INTEGER :00A6CEC888FFB8116BFF3562E3638D797BAABC39D9E227A1319C48EBCE60452ADEC36D66DE4F2FC68A4A54082A9CFE761EB6A59C185E42110B357FDD4C1E8DDC35308A7B72B732E7983DDA2493A0407253D29AD9E44F05C896C4DBED21CDF913E828F4C1E8E6D3F973E64F0050BDD70457796C278F445F4E0F006ADA91B2180ECD5B4E57EB272BD4D2C9855FEFF657D5DE01A261F46FEE26839E3AD1D52577CAA2BC08EE70559380BE7BCE6696A5EEAABEBB8C8367B772CA155BA2C258CCEC5B41EAF472D2D0E8AAA8F2A4985178A67F82076E3E9F2301AD6622F41F938F1D1F6F51A1CCEF50C91C06861487120C7C82A9E33FBDDF00F8DE6060EE248F4815E92B
`)
}

func TestBitstring(t *testing.T) {
	testDerOctets(t, "03 04 06 6e 5d c0", `UNIVERSAL PRIMITIVE BITSTRING PAD=6 :6E5DC0
`)
	testDerOctets(t, "03 01 00", `UNIVERSAL PRIMITIVE BITSTRING PAD=0 :
`)
	testDerOctets(t, "03 07 00 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE BITSTRING PAD=0 :6869206D6F6D
`)
	// BER, not DER (long-form length coding)
	testDerOctets(t, "03 81 04 06 6e 5d c0", `UNIVERSAL PRIMITIVE BITSTRING PAD=6 :6E5DC0
`)
}

func TestBitstringBadInputs(t *testing.T) {
	testDerOctets(t, "03 00", `BitString had no padding byte`)
	testDerOctets(t, "03 02 FF FF", `BitString padding has illegal value 255`)
}

func TestBitstringComposed(t *testing.T) {
	testDerOctets(t, "23 09  03 03 00 6e 5d  03 02 06 c0", `UNIVERSAL CONSTRUCTED UNHANDLED-TAG=3
  UNIVERSAL PRIMITIVE BITSTRING PAD=0 :6E5D
  UNIVERSAL PRIMITIVE BITSTRING PAD=6 :C0
`)
}

func TestOctetString(t *testing.T) {
	testDerOctets(t, "04 08 01 23 45 67 89 ab cd ef", `UNIVERSAL PRIMITIVE OCTETSTRING :0123456789ABCDEF
`)
}

func TestOctetStringHinted(t *testing.T) {
	testDerOctets(t, "04 08 6869206D6F6D FFFF", `UNIVERSAL PRIMITIVE OCTETSTRING :6869206D6F6DFFFF
# strings: "hi mom.."
`)
}

func TestNull(t *testing.T) {
	testDerOctets(t, "05 00", `UNIVERSAL PRIMITIVE NULL
`)
	testDerOctets(t, "05 01 00", `Null has non-zero content`)
}

func TestOID(t *testing.T) {
	testDerOctets(t, "06 00", `OID doesn't have content`)
	testDerOctets(t, "06 06 2a 86 48 86 f7 0d", `UNIVERSAL PRIMITIVE OID 1.2.840.113549
`)
	testDerOctets(t, "06 09 2a 86 48 86 f7 0d 01 01 01", `UNIVERSAL PRIMITIVE OID 1.2.840.113549.1.1.1
# RSA Encryption
`)
	testDerOctets(t, "06 03 55 04 03", `UNIVERSAL PRIMITIVE OID 2.5.4.3
# CommonName
`)
}

func TestObjectDescription(t *testing.T) {
	// I don't have a good example for this
	testDerOctets(t, "07 03 01 02 03", `UNIVERSAL PRIMITIVE OBJECTDESCRIPTION :010203
`)
}

func TestExternal(t *testing.T) {
	testDerOctets(t, "28 0A 020107820500 31323334", `UNIVERSAL CONSTRUCTED EXTERNAL :02010782050031323334
`)
}

func TestReal(t *testing.T) {
	// 0
	testDerOctets(t, "09 00", `UNIVERSAL PRIMITIVE REAL :
`)
	// -0 base 10
	testDerOctets(t, "09 03 01 2D 30", `UNIVERSAL PRIMITIVE REAL :012D30
# strings: ".-0"
`)
	// 3.14 base 10
	testDerOctets(t, "09 08 03 33 31 34 2E 45 2D 32", `UNIVERSAL PRIMITIVE REAL :033331342E452D32
# strings: ".314.E-2"
`)
}

func TestEnumerated(t *testing.T) {
	// same as Integer
	testDerOctets(t, "0a 01 00", `UNIVERSAL PRIMITIVE ENUMERATED 0
`)
}

func TestEmbeddedPDV(t *testing.T) {
	// you are in a maze of twisty unused types, all alike
	testDerOctets(t, "2b 03 01 02 03", `UNIVERSAL CONSTRUCTED EMBEDDED-PDV :010203
`)
}

func TestUtf8String(t *testing.T) {
	testDerOctets(t, "0c 00", `UNIVERSAL PRIMITIVE UTF8STRING '
`)
	testDerOctets(t, "0c 06 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE UTF8STRING 'hi mom
`)
	testDerOctets(t, "0c 07 68 69 20 6d 6f 6d 0a", `UNIVERSAL PRIMITIVE UTF8STRING 'hi mom\n
`)
	testDerOctets(t, "0c 07 68 c3 ad 20 6d 6f 6d", `UNIVERSAL PRIMITIVE UTF8STRING 'hí mom
`)
}

func TestUtfStringBadUtf8(t *testing.T) {
	testDerOctets(t, "0c 06 68 ed 20 6d 6f 6d", "UNIVERSAL PRIMITIVE UTF8STRING 'h\xed mom\n")
}

func TestRelativeOID(t *testing.T) {
	testDerOctets(t, "0d 00", `Relative OID doesn't have content`)
	testDerOctets(t, "0d 02 04 03", `UNIVERSAL PRIMITIVE RELATIVEOID 4.3
`)
}

func TestNumericString(t *testing.T) {
	testDerOctets(t, "12 07 31203220332034", `UNIVERSAL PRIMITIVE NUMERICSTRING '1 2 3 4
`)
	testDerOctets(t, "12 00", `UNIVERSAL PRIMITIVE NUMERICSTRING '
`)
}

// strictly speaking, only [0-9 ] are allowed in numericstring
func TestNumericStringInvalid(t *testing.T) {
	testDerOctets(t, "12 06 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE NUMERICSTRING 'hi mom
`)
}

func TestPrintableString(t *testing.T) {
	testDerOctets(t, "13 00", `UNIVERSAL PRIMITIVE PRINTABLESTRING '
`)
	testDerOctets(t, "13 06 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE PRINTABLESTRING 'hi mom
`)
}

func TestT61String(t *testing.T) {
	testDerOctets(t, "14 00", `UNIVERSAL PRIMITIVE T61STRING :
`)
	testDerOctets(t, "14 06 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE T61STRING :6869206D6F6D
# strings: "hi mom"
`)
}

func TestVideotextString(t *testing.T) {
	testDerOctets(t, "15 00", `UNIVERSAL PRIMITIVE VIDEOTEXSTRING :
`)
	testDerOctets(t, "15 06 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE VIDEOTEXSTRING :6869206D6F6D
# strings: "hi mom"
`)
}

func TestIA5String(t *testing.T) {
	testDerOctets(t, "16 00", `UNIVERSAL PRIMITIVE IA5STRING '
`)
	testDerOctets(t, "16 06 68 69 20 6d 6f 6d", `UNIVERSAL PRIMITIVE IA5STRING 'hi mom
`)
}

func TestUTCTime(t *testing.T) {
	testDerOctets(t, "17 00", `UNIVERSAL PRIMITIVE UTCTIME '
`)
	testDerOctets(t, "17 0B 313830393130 303130325A", `UNIVERSAL PRIMITIVE UTCTIME '1809100102Z
`)
	testDerOctets(t, "17 0D 313830393130 3031303230305A", `UNIVERSAL PRIMITIVE UTCTIME '180910010200Z
# 2018-09-10 01:02:00 GMT
`)
}

func TestGeneralizedTime(t *testing.T) {
	testDerOctets(t, "18 00", `UNIVERSAL PRIMITIVE GENERALIZEDTIME '
`)
	testDerOctets(t, "18 0D 313830393130 3031303230305A", `UNIVERSAL PRIMITIVE GENERALIZEDTIME '180910010200Z
# 2018-09-10 01:02:00 GMT
`)
	testDerOctets(t, "18 13 3230303031323331 323335393539 2E 393939 5A",
		`UNIVERSAL PRIMITIVE GENERALIZEDTIME '20001231235959.999Z
# 2000-12-31 23:59:59.999 GMT
`)
}

func TestGraphicString(t *testing.T) {
	testDerOctets(t, "19 00", `UNIVERSAL PRIMITIVE GRAPHICSTRING :
`)
	testDerOctets(t, "19 06 6869206D6F6D", `UNIVERSAL PRIMITIVE GRAPHICSTRING :6869206D6F6D
# strings: "hi mom"
`)
}

func TestVisibleString(t *testing.T) {
	testDerOctets(t, "1a 00", `UNIVERSAL PRIMITIVE VISIBLESTRING '
`)
	testDerOctets(t, "1a 06 6869206D6F6D", `UNIVERSAL PRIMITIVE VISIBLESTRING 'hi mom
`)
}

// non-ascii not allowed but make sure we can round-trip it
func TestVisibleStringIllegal(t *testing.T) {
	testDerOctets(t, "1a 06 68ed206D6F6D", "UNIVERSAL PRIMITIVE VISIBLESTRING 'h\xed mom\n")
}

func TestGeneralString(t *testing.T) {
	testDerOctets(t, "1b 00", `UNIVERSAL PRIMITIVE GENERALSTRING :
`)
	testDerOctets(t, "1b 06 6869206D6F6D", `UNIVERSAL PRIMITIVE GENERALSTRING :6869206D6F6D
# strings: "hi mom"
`)
}

func TestUniversalString(t *testing.T) {
	testDerOctets(t, "1c 00", `UNIVERSAL PRIMITIVE UNIVERSALSTRING '
`)
	testDerOctets(t, "1C 18 00000068 00000069 00000020 0000006D 0000006F 0000006D",
		`UNIVERSAL PRIMITIVE UNIVERSALSTRING 'hi mom
`)
}

// lock behavior down
func TestUniversalStringIllegal(t *testing.T) {
	testDerOctets(t, "1C 17 00000068 00000069 00000020 0000006D 0000006F 000000",
		"UNIVERSAL PRIMITIVE UNIVERSALSTRING 'hi mo\uFFFD\n")
	testDerOctets(t, "1C 04 0000ffff", "UNIVERSAL PRIMITIVE UNIVERSALSTRING '\uFFFF\n")
	testDerOctets(t, "1C 04 0000fffe", "UNIVERSAL PRIMITIVE UNIVERSALSTRING '\uFFFE\n")
	testDerOctets(t, "1C 04 0000fffc", "UNIVERSAL PRIMITIVE UNIVERSALSTRING '\uFFFC\n")
}

func TestCharacterString(t *testing.T) {
	testDerOctets(t, "1D 06 6869206D6F6D", `UNIVERSAL PRIMITIVE CHARACTERSTRING :6869206D6F6D
# strings: "hi mom"
`)
}

func TestBMPString(t *testing.T) {
	testDerOctets(t, "1e 00", `UNIVERSAL PRIMITIVE BMPSTRING '
`)
	testDerOctets(t, "1e 0c 0068 0069 0020 006D 006F 006D",
		`UNIVERSAL PRIMITIVE BMPSTRING 'hi mom
`)
}

// lock behavior down
func TestBMPStringIllegal(t *testing.T) {
	testDerOctets(t, "1e 0b 0068 0069 0020 006D 006F 00",
		"UNIVERSAL PRIMITIVE BMPSTRING 'hi mo\uFFFD\n")
	testDerOctets(t, "1e 02 ffff", "UNIVERSAL PRIMITIVE BMPSTRING '\uFFFF\n")
	testDerOctets(t, "1e 02 fffe", "UNIVERSAL PRIMITIVE BMPSTRING '\uFFFE\n")
	testDerOctets(t, "1e 02 fffc", "UNIVERSAL PRIMITIVE BMPSTRING '\uFFFC\n")
	testDerOctets(t, "1e 02 d800", "UNIVERSAL PRIMITIVE BMPSTRING '\uFFFD\n")
}

func TestSequence(t *testing.T) {
	testDerOctets(t, "30 00", `UNIVERSAL CONSTRUCTED SEQUENCE
`)
	testDerOctets(t, "30 08"+"02017B"+"0C03616263", `UNIVERSAL CONSTRUCTED SEQUENCE
  UNIVERSAL PRIMITIVE INTEGER 123
  UNIVERSAL PRIMITIVE UTF8STRING 'abc
`)
}

func TestSet(t *testing.T) {
	testDerOctets(t, "31 00", `UNIVERSAL CONSTRUCTED SET
`)
	testDerOctets(t, "31 08"+"02017B"+"0C03616263", `UNIVERSAL CONSTRUCTED SET
  UNIVERSAL PRIMITIVE INTEGER 123
  UNIVERSAL PRIMITIVE UTF8STRING 'abc
`)
}

func TestLongTypeTag(t *testing.T) {
	testDerOctets(t, "5f 87 67 03 01 02 03", `APPLICATION PRIMITIVE UNHANDLED-TAG=999 :010203
`)
}
