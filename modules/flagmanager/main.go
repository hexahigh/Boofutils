package flagmanager

type BuaConfig struct {
	InFile          string
	OutFile         string
	Encode          bool
	Mute            bool
	Verbosity       int
	BestCompression bool
}

type ReportConfig struct {
	OutFile    string
	Stdout     bool
	PrintLevel int
}
