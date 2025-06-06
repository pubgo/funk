package main

import (
	"flag"

	"github.com/pubgo/funk/cmds/protoc-gen-go-cloudevent/internal"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var _ = flag.String("version", internal.Version, "version")

func main() {
	flag.Parse()

	protogen.Options{ParamFunc: flag.CommandLine.Set}.
		Run(func(plugin *protogen.Plugin) error {
			plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

			for _, name := range plugin.Request.FileToGenerate {
				internal.GenerateFile(plugin, plugin.FilesByPath[name])
			}

			return nil
		})
}
