package internal

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/pubgo/funk/proto/errorpb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const errorPbPkg = "github.com/pubgo/funk/proto/errorpb"

// GenerateFile generates a .errors.pb.go file containing service definitions.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".errors.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	genFile := jen.NewFile(string(file.GoPackageName))
	genFile.HeaderComment("Code generated by protoc-gen-lava-errors. DO NOT EDIT.")
	genFile.HeaderComment("versions:")
	genFile.HeaderComment(fmt.Sprintf("- protoc-gen-lava-errors %s", version))
	genFile.HeaderComment(fmt.Sprintf("- protoc                 %s", protocVersion(gen)))
	if file.Proto.GetOptions().GetDeprecated() {
		genFile.HeaderComment(fmt.Sprintf("%s is a deprecated file.", file.Desc.Path()))
	} else {
		genFile.HeaderComment(fmt.Sprintf("source: %s", file.Desc.Path()))
	}

	genFile.Comment("This is a compile-time assertion to ensure that this generated file")
	genFile.Comment("is compatible with the grpc package it is being compiled against.")
	genFile.Comment("Requires gRPC-Go v1.32.0 or later.")
	genFile.Id("const _ =").Qual("google.golang.org/grpc", "SupportPackageIsVersion7")
	g.Skip()

	for i := range file.Enums {
		m := file.Enums[i]
		tag, ok := proto.GetExtension(m.Desc.Options(), errorpb.E_Opts).(*errorpb.Options)
		if !ok || tag == nil || !tag.GetGen() {
			continue
		}

		g.Unskip()

		for _, codeName := range m.Values {
			var name = strings.ToLower(fmt.Sprintf("%s.%s.%s",
				file.Desc.Package(),
				"err_code",
				strcase.ToSnake(string(codeName.Desc.Name())),
			))

			// comment
			var rr = string(codeName.Desc.Name())
			if codeName.Comments.Leading.String() != "" {
				rr = codeName.Comments.Leading.String()
				rr = strings.Trim(strings.TrimSpace(rr), "/")
			}
			rr = strings.ToLower(strcase.ToSnake(rr))
			rr = strings.ReplaceAll(rr, "_", " ")
			rr = strings.TrimSpace(strings.ReplaceAll(rr, "  ", " "))

			var statusName = "OK"
			field, ok := proto.GetExtension(codeName.Desc.Options(), errorpb.E_Field).(*errorpb.Fields)
			if ok && field != nil {
				statusName = field.Code.String()
			}

			genFile.Var().
				Id("ErrCode"+string(codeName.Desc.Name())).
				Id("=").
				Op("&").Qual(errorPbPkg, "ErrCode").
				Values(jen.Dict{
					jen.Id("Code"):   jen.Qual(errorPbPkg, "Code_"+statusName),
					jen.Id("Status"): jen.Lit(name),
					jen.Id("Reason"): jen.Lit(rr),
				}).Line()
		}
	}

	g.P(genFile.GoString())
	return g
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}
