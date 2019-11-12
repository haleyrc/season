package season

type scanConfig struct {
	nested  bool
	garbage string
}

type Option func(cfg *scanConfig)

func WithNested(on bool) Option {
	return func(cfg *scanConfig) {
		cfg.nested = on
	}
}

func WithGarbage(g string) Option {
	return func(cfg *scanConfig) {
		cfg.garbage = g
	}
}
