module github.com/alis-exchange/cli

go 1.17

require (
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pterm/pterm v0.12.33
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	go.protobuf.alis.alis.exchange v0.0.0-20211201131156-ccfa4d80795e
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

require (
	cloud.google.com/go v0.93.3 // indirect
	github.com/atomicgo/cursor v0.0.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gookit/color v1.4.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Use the replace when developing locally
//replace go.protobuf.alis.dev => ../../alis/protobuf/go/alis

//replace go.lib.alis.dev => ../../alis/os1/os-go
//replace go.protobuf.alis.alis.exchange => ../../alis.exchange/alis/protobuf/go
