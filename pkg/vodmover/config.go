package vodmover

type Config struct {
	OBS struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	} `yaml:"obs"`
	MoveVODs []MoveVODRule `yaml:"move_vods"`
}

type MoveVODRule struct {
	PatternWildcard string `yaml:"pattern_wildcard"`
	Destination     string `yaml:"destination"`
}
