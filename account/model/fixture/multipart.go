package fixture

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"runtime"
)

// MultipartImage is used for instantiating a test fixture
// for creating multipart file uploads containing an image.
type MultipartImage struct {
	ImagePath     string
	ImageFile     *os.File
	MultipartBody *bytes.Buffer
	ContentType   string
}

// NewMultipartImage creates an image file for testing
// and creates a Multipart form with this image file
// for testing.
func NewMultipartImage(fileName string, contentType string) *MultipartImage {
	// Create a test file in the same folder as this fixture.
	_, b, _, _ := runtime.Caller(0)
	directory := filepath.Dir(b)

	imagePath := filepath.Join(directory, fileName)
	f := createImage(imagePath)

	defer f.Close()

	// Create a multipart write onto which we
	// will write the image file.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Manually create a form file as CreateFormFile will
	// force file's content type to "application/octet-stream".
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "imageFile", fileName))
	header.Set("Content-Type", contentType)
	part, _ := writer.CreatePart(header)

	io.Copy(part, f)
	writer.Close()

	return &MultipartImage{
		ImagePath:     imagePath,
		ImageFile:     f,
		MultipartBody: body,
		ContentType:   writer.FormDataContentType(),
	}
}

// GetFormFile extracts a form file from a multipart body.
func (m *MultipartImage) GetFormFile() *multipart.FileHeader {
	_, params, _ := mime.ParseMediaType(m.ContentType)
	multipartReader := multipart.NewReader(m.MultipartBody, params["boundary"])
	form, _ := multipartReader.ReadForm(1024)
	files := form.File["imageFile"]

	return files[0]
}

// Close removes a created file for test.
func (m *MultipartImage) Close() {
	m.ImageFile.Close()
	os.Remove(m.ImagePath)
}

// createImage is used to create a quick example
// 1px x 1px image encoded as a PNG.
func createImage(imagePath string) *os.File {
	rectangle := image.Rect(0, 0, 1, 1)
	img := image.NewRGBA(rectangle)

	f, _ := os.Create(imagePath)
	png.Encode(f, img)

	return f
}
