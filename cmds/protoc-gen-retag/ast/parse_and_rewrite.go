// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	retagpb "github.com/pubgo/funk/proto/retag"
	"github.com/samber/lo"
	"github.com/searKing/golang/go/reflect"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type FieldInfo struct {
	FieldNameInProto string
	FieldNameInGo    string
	FieldTag         []reflect.StructTag
}

type StructInfo struct {
	StructNameInProto string
	StructNameInGo    string
	FieldInfos        []FieldInfo
}

type FileInfo struct {
	FileName    string
	StructInfos []StructInfo
}

func (si *StructInfo) FindField(name string) (FieldInfo, bool) {
	for _, f := range si.FieldInfos {
		if f.FieldNameInGo == name {
			return f, true
		}
	}
	return FieldInfo{}, false
}

// WalkDescriptorProto returns all struct infos of dp， which contains FieldTag.
func WalkDescriptorProto(g *protogen.Plugin, dp *descriptorpb.DescriptorProto, typeNames []string) []StructInfo {
	var ss []StructInfo

	s := StructInfo{}
	s.StructNameInProto = dp.GetName()
	s.StructNameInGo = CamelCaseSlice(append(typeNames, CamelCase(dp.GetName())))

	for _, field := range dp.GetField() {
		var oneOfS *StructInfo
		if field.OneofIndex != nil { // Special Case: oneof
			oneOfS = &StructInfo{}
			oneOfS.StructNameInProto = field.GetName()
			oneOfS.StructNameInGo = CamelCaseSlice(append(typeNames, CamelCase(dp.GetName()), CamelCase(field.GetName())))
		}

		f := HandleFieldDescriptorProto(field)
		if f == nil {
			continue
		}

		if oneOfS != nil {
			oneOfS.FieldInfos = append(oneOfS.FieldInfos, *f)
			if len(oneOfS.FieldInfos) > 0 {
				ss = append(ss, *oneOfS)
			}
		} else {
			s.FieldInfos = append(s.FieldInfos, *f)
		}
	}

	typeNames = append(typeNames, CamelCase(dp.GetName()))

	for _, decl := range dp.GetOneofDecl() {
		declS := HandleOneOfDescriptorProto(decl, typeNames)
		if declS == nil {
			continue
		}

		if decl.GetOptions() == nil {
			continue
		}

		oneOfTags, ok := proto.GetExtension(decl.GetOptions(), retagpb.E_OneofTags).([]*retagpb.Tag)
		if !ok || len(oneOfTags) == 0 {
			continue
		}

		info := FieldInfo{
			FieldNameInProto: decl.GetName(),
			FieldNameInGo:    CamelCase(decl.GetName()),
			FieldTag: lo.Map(oneOfTags, func(_ *retagpb.Tag, i int) reflect.StructTag {
				tag := reflect.StructTag{}
				tag.SetName(oneOfTags[i].Name, oneOfTags[i].Value)
				return tag
			}),
		}
		s.FieldInfos = append(s.FieldInfos, info)

		ss = append(ss, *declS)
	}

	if len(s.FieldInfos) > 0 {
		ss = append(ss, s)
	}

	for _, nest := range dp.GetNestedType() {
		ss = append(ss, WalkDescriptorProto(g, nest, typeNames)...)
	}
	return ss
}

func HandleOneOfDescriptorProto(dp *descriptorpb.OneofDescriptorProto, typeNames []string) *StructInfo {
	if dp == nil {
		return nil
	}

	s := StructInfo{}
	s.StructNameInProto = dp.GetName()
	s.StructNameInGo = "is" + CamelCaseSlice(append(typeNames, CamelCase(dp.GetName())))
	return &s
}

func HandleFieldDescriptorProto(field *descriptorpb.FieldDescriptorProto) *FieldInfo {
	if field.GetOptions() == nil {
		return nil
	}

	tags, ok := proto.GetExtension(field.GetOptions(), retagpb.E_Tags).([]*retagpb.Tag)
	if !ok || len(tags) == 0 {
		return nil
	}

	info := &FieldInfo{
		FieldNameInProto: field.GetName(),
		FieldNameInGo:    CamelCase(field.GetName()),
		FieldTag: lo.Map(tags, func(v *retagpb.Tag, _ int) reflect.StructTag {
			tag := reflect.StructTag{}
			tag.SetName(v.Name, v.Value)
			return tag
		}),
	}

	return info
}

func Rewrite(g *protogen.Plugin) {
	var protoFiles []FileInfo

	for _, protoFile := range g.Request.GetProtoFile() {
		if !lo.Contains(g.Request.GetFileToGenerate(), protoFile.GetName()) {
			continue
		}

		var f FileInfo
		f.FileName = protoFile.GetName()

		for _, messageType := range protoFile.GetMessageType() {
			f.StructInfos = append(f.StructInfos, WalkDescriptorProto(g, messageType, nil)...)
		}

		if len(f.StructInfos) > 0 {
			protoFiles = append(protoFiles, f)
		}
	}

	// FIXME: always generate *.pb.go, to replace protoc-go, avoid "Tried to write the same file twice"
	//if len(protoFiles) == 0 {
	//	return
	//}
	//
	//// g.Response() will generate files, so skip this step
	//if len(g.Response().GetFile()) == 0 {
	//	return
	//}

	rewriter := NewGenerator(protoFiles, g)
	for _, f := range g.Response().GetFile() {
		rewriter.ParseGoContent(f)
	}
	rewriter.Generate()
}
