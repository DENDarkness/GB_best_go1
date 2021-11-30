package page

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLinks(t *testing.T) {

	htmlTest := `<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Testing</title>
</head>
<body>
<h1>Test</h1>
<a href="https://yandex.ru">url</a>
<a href="/testing">path</a>
<a href="mailto">mail</a>
<a href="#">#</a>
<a href="./">./</a>
<a href="">zero</a>
</body>
</html>`

	testPage, _ := NewPage(strings.NewReader(htmlTest), "https://example.com")
	links := testPage.GetLinks()
	want := []string{
		"https://example.com/testing",
	}
	r := assert.ElementsMatch(t, want, links)
	t.Logf("MSG: %v", r)
}
