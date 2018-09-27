# der2text

Convert a PEM- or DER-encoded file to a human- and machine-readable form.  This form can be edited and fed into the reverse process to create an edited file.  If signed attributes are modified then the signature will be invalid.

This is useful for inspection and understanding of cryptographic data, and for setting up tests that need specific forms of cryptographic input (e.g. creating unusual or malformed certificates).

In all cases, assuming the file is in PEM- or DER- format, the following round-trip process should produce output identical to input:

```
cat input | der2text | text2der > output
```

### der2text utility

Reads PEM- or DER-encoded input and produces a readable and editable output.

Usage:
```
go get github.com/syncsynchalt/der2text/cmds/der2text
# add ~/go/bin/ to your $PATH
der2text /path/to/cert.pem
```

### interim format

`der2text` produces an interim text file for input to `text2der`.
This file is meant to be human readable and editable but also easily
machine-parseable.  The format is:

1. Lines consisting of zero or more spaces followed by "#" followed by data are ignored by machines
2. The first line may consist of the words `PEM ENCODED {FOO}` where `{FOO}` is one of `CERTIFICATE`, `CERTIFICATE REQUEST`, `PRIVATE KEY`, and so on.  This indicates a PEM wrapper of type `{FOO}`.
3. Indentation of more than zero spaces indicate the depth of the containing set or sequence, and are defined in terms of 2*level (e.g. 4 spaces indicates the data is in two levels of set or sequence)
4. After indentation, the first word `UNIVERSAL`, `APPLICATION`, `CONTEXT-SPECIFIC`, or `PRIVATE` indicates the ASN.1 type class.  As a rule this utility can only make types of class `UNIVERSAL` human-readable but will preserve data for all other types and classes found.
5. After class, the word `PRIMITIVE` or `CONSTRUCTED` indicates the ASN.1 type flag of primitive (content represents this single type) vs constructed (content contains multiple elements that can themselves be individual type-length-content datums).
    * In the case of types which can be either primitive or constructed this utility takes the opinionated stance of only representing the primitive type in human-readable form (this is also constrained by DER encoding rules in most cases).  In all cases data is preserved whether primitive or constructed.
6. After primitive/constructed flag, the ASN.1 type tag is listed as below:
   * `END-OF-CONTENT` followed by nothing
   * `INTEGER` followed by the number or content data
   * `BITSTRING1` followed by `PAD=n` of right padding amount followed by content data
   * `OCTETSTRING` followed by content data
   * `NULL` followed by nothing
   * `OID` followed by an ASN.1 object identifier
   * `OBJECTDESCRIPTION` followed by content data
   * `EXTERNAL` followed by content data
   * `REAL` followed by content data
   * `ENUMERATED` followed by the number or content data
   * `EMBEDDED-PDV` followed by content data
   * `UTF8STRING` followed by content data
   * `RELATIVEOID` followed by an ASN.1 object identifier
   * `NUMERICSTRING` followed by content data
   * `PRINTABLESTRING` followed by content data
   * `SET` followed by lines of higher indentation level that represent the data within this set
   * `SEQUENCE` followed by lines of higher indentation level that represent the data within this set
   * `T61STRING` followed by content data
   * `VIDEOTEXSTRING` followed by content data
   * `IA5STRING` (ASCII string) followed by content data
   * `UTCTIME` followed by content data
   * `GENERALIZEDTIME` followed by content data
   * `GRAPHICSTRING` followed by content data
   * `VISIBLESTRING` followed by content data
   * `GENERALSTRING` followed by content data
   * `UNIVERSALSTRING` (UTF-32BE string) followed by content data
   * `CHARACTERSTRING` followed by content data
   * `BMPSTRING` (UTF-16BE string) followed by content data
   * `UNHANDLED-TAG=xx` followed by content data
      * This represents a type that we can't show in a human-readable way without some knowledge of the ASN.1 schema.  This data is preserved but may be opaque to our desire to edit it.

In the above list, "content data" consists of either:

* the character `:` followed by the data converted to hexadecimal in pairs terminated by a newline
* the character `'` followed by the data terminated by a newline.  The data has been modified as below:
   * newlines are converted to `\n`
   * carriage returns are converted to `\r`

It is best to treat this data as ephemeral in case the format changes in the future.  In other words, keep data in PEM or DER form and convert it at the time that the changes should be made, then put it back in PEM or DER form for storage.

### text2der utility

Reads the "interim format" output of `der2text` and creates PEM- or DER-encoded data from it.

Usage:
```
go get github.com/syncsynchalt/der2text/cmds/text2der
# add ~/go/bin/ to your $PATH
text2der /path/to/der2text/output
```
