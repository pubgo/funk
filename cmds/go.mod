module github.com/pubgo/funk/v2/cmds

go 1.23.0

replace github.com/pubgo/funk/v2 => ../

replace github.com/pubgo/funk/v2/component => ../component

require (
	entgo.io/ent v0.13.1
	github.com/dave/jennifer v1.7.0
	github.com/iancoleman/strcase v0.2.0
	github.com/moby/term v0.5.0
	github.com/samber/lo v1.51.0
	github.com/urfave/cli/v3 v3.4.1
	google.golang.org/protobuf v1.34.3-0.20240816073751-94ecbc261689
	gopkg.in/yaml.v3 v3.0.1
)

require (
	ariga.io/atlas v0.21.1 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/DATA-DOG/go-sqlmock v1.5.2 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/go-openapi/inflect v0.19.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.13.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/zclconf/go-cty v1.8.0 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
