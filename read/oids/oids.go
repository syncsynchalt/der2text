package oids

var oids map[string]string = map[string]string{
	"2.5.4.3":               "CommonName",
	"2.5.4.4":               "SurName",
	"2.5.4.5":               "DeviceSerialNumber",
	"2.5.4.6":               "Country",
	"2.5.4.7":               "Locality",
	"2.5.4.8":               "State",
	"2.5.4.9":               "StreetAddress",
	"2.5.4.10":              "OrganizationalUnit",
	"2.5.4.11":              "Organization",
	"2.5.4.12":              "Title",
	"2.5.4.13":              "Description",
	"2.5.4.16":              "PostalAddress",
	"2.5.4.17":              "PostalCode",
	"2.5.4.20":              "TelephoneNumber",
	"2.5.4.42":              "GivenName",
	"2.5.4.43":              "Initials",
	"1.3.101.110":           "x25519",
	"1.2.840.113549.1.1.1":  "RSA Encryption",
	"1.2.840.113549.1.1.11": "Sha256WithRSAEncryption",
	"1.2.840.113549.1.9.1":  "Email",
}

func Name(oid string) string {
	return oids[oid]
}
