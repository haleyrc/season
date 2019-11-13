package season

type scanConfig struct {
	nested  bool
	garbage string
	debug   bool
}

type Option func(cfg *scanConfig)

func WithDebug(on bool) Option {
	return func(cfg *scanConfig) {
		cfg.debug = on
	}
}

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
