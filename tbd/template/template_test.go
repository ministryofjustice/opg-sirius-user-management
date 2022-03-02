package template

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert := assert.New(t)

	tmpls, err := Parse("testdata/template/layout/*.gohtml", "testdata/template/*.gohtml", nil)
	assert.Nil(err)

	tmpl := tmpls.Get("test.gohtml")

	var buf bytes.Buffer
	err = tmpl(&buf, nil)
	assert.Nil(err)
	assert.Equal("<body>\n    <h1>Hello</h1>\n  </body>", buf.String())
}

func TestGetWithBadName(t *testing.T) {
	assert := assert.New(t)

	tmpls, err := Parse("testdata/template/layout/*.gohtml", "testdata/template/*.gohtml", nil)
	assert.Nil(err)

	tmpl := tmpls.Get("error.gohtml")
	err = tmpl(io.Discard, nil)
	assert.Equal(MissingError("error.gohtml"), err)
}
