package examples_commongo

type OptionsBuilder func(options *Options)

type Options struct {
	KV map[string]interface{}
}

func BuildOptions(opts ...OptionsBuilder) *Options {
	options := &Options{
		KV: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
