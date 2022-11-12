package wrapper

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strconv"
	"strings"
)

// boundary represents the boundary that is used in every multipart form.
var boundary = randomBoundary()

// randomBoundary generates a random value to be used as a boundary in multipart forms.
//
// Implemented from the mime/multipart package.
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", buf[:])
}

// quoteEscaper escapes quotes and backslashes in a multipart form.
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// createMultipartForm creates a multipart/form for Discord using a given JSON body and files.
func createMultipartForm(json []byte, files ...File) ([]byte, error) {
	b := bytes.NewBuffer(nil)

	// set the boundary.
	multipartWriter := multipart.NewWriter(b)
	err := multipartWriter.SetBoundary(boundary)
	if err != nil {
		return nil, fmt.Errorf("error setting multipart form boundary: %w", err)
	}

	// add the `payload_json` JSON to the form.
	multipartPayloadJSONPart, err := createPayloadJSONForm(multipartWriter)
	if err != nil {
		return nil, fmt.Errorf("error adding JSON payload header to multipart form: %w", err)
	}

	payloadJSON := bytes.NewReader(json)
	if _, err := io.Copy(multipartPayloadJSONPart, payloadJSON); err != nil {
		return nil, fmt.Errorf("error copying JSON payload data to multipart form: %w", err)
	}

	// add the remaining files to the form.
	for i, file := range files {
		name := strings.Join([]string{"file[", strconv.Itoa(i), "]"}, "")
		multipartFilePart, err := createFormFile(multipartWriter, name, file.Name, file.ContentType)
		if err != nil {
			return nil, fmt.Errorf("error adding a file %q to a multipart form: %w", file.Name, err)
		}

		payloadFile := bytes.NewReader(file.Data)
		if _, err := io.Copy(multipartFilePart, payloadFile); err != nil {
			return nil, fmt.Errorf("error copying file %q data to multipart form: %w", file.Name, err)
		}
	}

	// write the trailing boundary.
	if err := multipartWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing the multipart form: %w", err)
	}

	return b.Bytes(), nil
}

var (
	// contentTypeJSONString is an HTTP Header Content Type that indicates a payload with a JSON body.
	contentTypeJSONString = "application/json"

	// contentTypeOctetStreamString is an HTTP Header Content Type that indicates a payload with binary data.
	contentTypeOctetStreamString = "application/octet-stream"
)

// createPayloadJSONForm creates a form-data header for the `payload_json` file in a multipart form.
func createPayloadJSONForm(m *multipart.Writer) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `data; name="payload_json"`)
	h.Set("Content-Type", contentTypeJSONString)

	return m.CreatePart(h) // nolint:wrapcheck
}

// createFormFile creates a form-data header for file attachments in a multipart form.
func createFormFile(m *multipart.Writer, name, filename, contentType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, name, quoteEscaper.Replace(filename)))

	if "" == contentType {
		contentType = contentTypeOctetStreamString
	}

	h.Set("Content-Type", contentType)

	return m.CreatePart(h) // nolint:wrapcheck
}
