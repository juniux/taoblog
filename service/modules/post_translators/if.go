package post_translators

type PostTranslator interface {
	Translate(source string) (string, error)
}
