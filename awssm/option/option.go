package option

import (
	"github.com/hypcode/go-awssm-env/awssm/internal"
)

type ClientOption interface {
	Apply(*internal.Settings)
}

func WithReferencePrefix(prefix string) ClientOption {
	return withReferencePrefix(prefix)
}

type withReferencePrefix string

func (w withReferencePrefix) Apply(s *internal.Settings) {
	s.ReferencePrefix = string(w)
}
