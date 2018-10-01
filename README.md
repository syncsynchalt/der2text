# der2text

Convert a PEM- or DER-encoded file to a human- and machine-readable form.  This form can be edited and fed into the reverse process to create an edited PEM or DER file.

This is useful for inspection and understanding of cryptographic data and for setting up tests that need specific forms of cryptographic input (e.g. creating unusual or malformed certificates).

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

### text2der utility

Reads the "interim format" output of `der2text` and creates PEM- or DER-encoded data from it.

Usage:
```
go get github.com/syncsynchalt/der2text/cmds/text2der
# add ~/go/bin/ to your $PATH
text2der /path/to/der2text/output
```

### interim format

`der2text` produces an interim text file for input to `text2der`.
This file is meant to be human readable and editable but also easily
machine-parseable.  The format is:

1. Blank lines or lines consisting of zero or more spaces followed by "#" are ignored
2. The first line may consist of the words `PEM ENCODED {FOO}` where `{FOO}` is `CERTIFICATE`, `CERTIFICATE REQUEST`, `PRIVATE KEY`, and so on.  This indicates a PEM wrapper of type `{FOO}`.
3. Indentation with spaces indicate the items that are contained within a preceding `SET` or `SEQUENCE` or `PEM ENCODING`.  For example a `SET` that occurs with an indentation of two spaces will contain all immediately following lines indented by more than two spaces.
4. After indentation, the first word `UNIVERSAL`, `APPLICATION`, `CONTEXT-SPECIFIC`, or `PRIVATE` indicates the ASN.1 type class.  This utility can only make types of class `UNIVERSAL` human-readable but will preserve data for all other types found.
5. After type class the word `PRIMITIVE` or `CONSTRUCTED` indicates the ASN.1 type flag of primitive (content represents this single type) vs constructed (content contains multiple elements).
    * In the case of types which can be either primitive or constructed this utility only represents the primitive type in human-readable form.  This is also enforced by DER encoding rules in most cases.  In all cases the data is preserved whether primitive or constructed.
6. After primitive/constructed flag, the ASN.1 type tag and element content is as below:
   * `END-OF-CONTENT` followed by nothing
   * `INTEGER` followed by the number or by content data
   * `BITSTRING` followed by `PAD=n` of right padding amount (0-7 bits) followed by content data
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
   * `UNHANDLED-TAG=nn` followed by content data
      * This represents a type that we can't show in a human-readable way without some knowledge of the ASN.1 schema.  This data is preserved but may be opaque to our desire to edit it.

In the above list, "content data" consists of either:

* the character `:` followed by the data converted to hexadecimal in pairs terminated by a newline
* the character `'` followed by the data terminated by a newline.  The data has been modified as below:
   * newlines are converted to `\n`
   * carriage returns are converted to `\r`

It is best to treat this data as ephemeral in case the format changes in the future.  In other words, keep data in PEM or DER form and convert it at the time that the changes should be made, then put it back in PEM or DER form for storage.

#### Example output

```
PEM ENCODED CERTIFICATE
  UNIVERSAL CONSTRUCTED SEQUENCE
    UNIVERSAL CONSTRUCTED SEQUENCE
      UNIVERSAL PRIMITIVE INTEGER :00F4B0DA1F5D4A2788
      UNIVERSAL CONSTRUCTED SEQUENCE
        UNIVERSAL PRIMITIVE OID 1.2.840.113549.1.1.11
        # Sha256WithRSAEncryption
        UNIVERSAL PRIMITIVE NULL
      UNIVERSAL CONSTRUCTED SEQUENCE
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.6
            # Country
            UNIVERSAL PRIMITIVE PRINTABLESTRING 'US
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.8
            # State
            UNIVERSAL PRIMITIVE UTF8STRING 'Colorado
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.7
            # Locality
            UNIVERSAL PRIMITIVE UTF8STRING 'Parker
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.10
            # OrganizationalUnit
            UNIVERSAL PRIMITIVE UTF8STRING 'Ülfheim
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.3
            # CommonName
            UNIVERSAL PRIMITIVE UTF8STRING 'testcertificate.example.com
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 1.2.840.113549.1.9.1
            # Email
            UNIVERSAL PRIMITIVE IA5STRING 'fenris@ulfheim.net
      UNIVERSAL CONSTRUCTED SEQUENCE
        UNIVERSAL PRIMITIVE UTCTIME '180927171537Z
        # 2018-09-27 17:15:37 GMT
        UNIVERSAL PRIMITIVE UTCTIME '190927171537Z
        # 2019-09-27 17:15:37 GMT
      UNIVERSAL CONSTRUCTED SEQUENCE
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.6
            # Country
            UNIVERSAL PRIMITIVE PRINTABLESTRING 'US
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.8
            # State
            UNIVERSAL PRIMITIVE UTF8STRING 'Colorado
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.7
            # Locality
            UNIVERSAL PRIMITIVE UTF8STRING 'Parker
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.10
            # OrganizationalUnit
            UNIVERSAL PRIMITIVE UTF8STRING 'Ülfheim
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 2.5.4.3
            # CommonName
            UNIVERSAL PRIMITIVE UTF8STRING 'testcertificate.example.com
        UNIVERSAL CONSTRUCTED SET
          UNIVERSAL CONSTRUCTED SEQUENCE
            UNIVERSAL PRIMITIVE OID 1.2.840.113549.1.9.1
            # Email
            UNIVERSAL PRIMITIVE IA5STRING 'fenris@ulfheim.net
      UNIVERSAL CONSTRUCTED SEQUENCE
        UNIVERSAL CONSTRUCTED SEQUENCE
          UNIVERSAL PRIMITIVE OID 1.2.840.113549.1.1.1
          # RSA Encryption
          UNIVERSAL PRIMITIVE NULL
        UNIVERSAL PRIMITIVE BITSTRING PAD=0 :3082010A0282010100B772CB8D8B8E85F833350439EDC6C55E3AB4686A83A0DCF6C80C6DBBF10CCEF5AC799CBD0F5A62D2467AB708C2A2C34016B25C6C3057256328A0A4B7780FC08333E3238CB0E8E290144589EDA87AAB74CE6970DA1B29366B7A32E1CAF010EA4BCC8344774C896A04EDEA7B1A0937F952130706925AA6F42EDE577C081C825BE75232500F17077D8596D26B955456EC2A6CA01943A5FC328442B3C43BD540D8A40A8A8088A677E298ED3C8F16860B2007D17073241B5DFFA5835BDE5D200EC0BEEC798AEDAC587BF532073AC4664AEE6D091C3B92298A6BC461F3F5C16987980B5F449B0339B56F3F36353D4EE505728687E6E5000B70093BF64953C61AB253AF0203010001
    UNIVERSAL CONSTRUCTED SEQUENCE
      UNIVERSAL PRIMITIVE OID 1.2.840.113549.1.1.11
      # Sha256WithRSAEncryption
      UNIVERSAL PRIMITIVE NULL
    UNIVERSAL PRIMITIVE BITSTRING PAD=0 :AFC3B8595D59E2EDB8CBB31486DE45FF89769215D6CE9FC46E43962BA77E8628FC68568AF7F6FB349CDC56D35CB318C2AB005BF22E2B3BF8DE1A38F6030BAB71135D1BD2F88222D7E5342794263A7C416689E69B90FF39C1CC54E8BAF1B0CDB92EF3A6F2B74AFA1985EDF095AD96C3A218F6C4B1E3449AE8D756BF23B059DCD35AB013DFC098CA4F8F44FCB76BE6E5BA003A0FC4CF0EDF39F85F7F1CBAF7F87C8479A168D24085F8EA705F245D9A7F9F91F27A73A18658243C7CE1B3E3F68BFDEDB4722621CBB21AA3B69511F078E741EB933B6457C13F08748361374ACDDAE5B00F757D3DE0D27A93C23908E967C5AE71D279D77264371DDB18CB200E44DBD9
```
