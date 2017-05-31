package i18n

// LocalizedString represents a piece of textual information in multiple language.
// This must be kept in sync with the localized_string_schema.go.
type LocalizedString struct {
	// This is the primary value.
	EnUS string `json:"en-us"`
	ZhCN string `json:"zh-cn"`

	// The derived values only available from graphql result.
	DerivedCN string `json:"cn"`
	DerivedEN string `json:"en"`
}

func (s LocalizedString) String(intl *Context) string {
	switch intl {
	case CnContext:
		if s.ZhCN == "" {
			return s.EnUS
		}
		return s.ZhCN
	default:
		return s.EnUS
	}
}

// L constructs LocalizedString for deferred translation.
func L(key string) LocalizedString {
	return LocalizedString{
		ZhCN: CnContext.S(key),
		EnUS: EnContext.S(key),
	}
}

// SwitchLocalized is a healper function that is used in *_models.go.
func SwitchLocalized(intl *Context, en, cn string) string {
	if cn == "" {
		return en
	}
	switch intl.Language() {
	case EnUS:
		return en
	case ZhCN:
		return cn
	default:
		panic("Language Not supported!" + intl.LangCode())
	}
}

func (s LocalizedString) Merge(to LocalizedString) LocalizedString {
	l := s
	if l.ZhCN == "" {
		l.ZhCN = to.ZhCN
	}
	if l.EnUS == "" {
		l.EnUS = to.EnUS
	}
	return l
}
