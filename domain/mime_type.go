package domain

const (
	HtmlMimeType MimeTypeAlias = "html"
	TextMimeType MimeTypeAlias = "text"
)

type MimeTypeAlias string

func (mt MimeTypeAlias) IsValid() bool {
	return mt == HtmlMimeType || mt == TextMimeType
}

func (mt MimeTypeAlias) Header() string {
	switch mt {
	case HtmlMimeType:
		return "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	case TextMimeType:
		return "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	}

	return ""
}

func (mt MimeTypeAlias) String() string {
	return string(mt)
}
