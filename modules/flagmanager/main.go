package flagmanager

type BuaConfig struct {
	InFile          string
	OutFile         string
	Encode          bool
	Mute            bool
	Verbosity       int
	BestCompression bool
	V2              bool
}

type ReportConfig struct {
	OutFile    string
	Stdout     bool
	PrintLevel int
	Yaml       bool
	Json       bool
}

type TarmountConfig struct {
	MountPoint string
	Format     string
	TarPath    string
}
