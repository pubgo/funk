package langx

import (
	"github.com/pubgo/funk/v2/log"
	"golang.org/x/text/language"
)

func ParseTags(lang string) []language.Tag {
	tags, _, err := language.ParseAcceptLanguage(lang)
	if err != nil {
		log.Err(err).Str("lang", lang).Msg("parse accept-language failed")
		return []language.Tag{language.Chinese}
	}

	return tags
}
