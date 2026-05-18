package service

import (
	"net/url"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
)

type utmFields struct {
	source   string
	medium   string
	campaign string
	content  string
	term     string
}

func newUTMFields(values url.Values) utmFields {
	return utmFields{
		source:   textutil.TruncateRunes(values.Get("utm_source"), maxUTMValueLength),
		medium:   textutil.TruncateRunes(values.Get("utm_medium"), maxUTMValueLength),
		campaign: textutil.TruncateRunes(values.Get("utm_campaign"), maxUTMValueLength),
		content:  textutil.TruncateRunes(values.Get("utm_content"), maxUTMValueLength),
		term:     textutil.TruncateRunes(values.Get("utm_term"), maxUTMValueLength),
	}
}
