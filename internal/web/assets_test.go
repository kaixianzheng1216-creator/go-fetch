package web

import "testing"

func TestIndexHTMLIsEmbedded(t *testing.T) {
	html, err := IndexHTML()
	if err != nil {
		t.Fatal(err)
	}
	if len(html) == 0 {
		t.Fatal("expected embedded index.html")
	}
}
