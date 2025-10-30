package examples_commongo

func WithConfigPath(configPath string) OptionsBuilder {
	return func(options *Options) {
		options.KV[ConfigPathOpt] = configPath
	}
}

func WithSMCABIPath(smcABIPath string) OptionsBuilder { 
		return func(options *Options) {
		options.KV[SMCABITPathOpt] = smcABIPath
	}
}