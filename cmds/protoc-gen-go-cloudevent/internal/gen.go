package internal

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/pubgo/funk/assert"
	"github.com/pubgo/funk/component/cloudevent"
	"github.com/pubgo/funk/errors/errcheck"
	cloudeventpb "github.com/pubgo/funk/proto/cloudevent"
	"github.com/pubgo/funk/stack"
)

var cloudeventPkg = reflect.TypeOf(cloudevent.Client{}).PkgPath()
var jobTypesPkg = reflect.TypeOf(cloudeventpb.PushEventOptions{}).PkgPath()
var resultTypesPkg = stack.CallerWithFunc(errcheck.Check).Pkg
var ctxPkg = stack.CallerWithFunc(context.WithTimeout).Pkg
var assertPkt = stack.CallerWithFunc(assert.Assert).Pkg
var protojsonPkt = stack.CallerWithFunc(protojson.Marshal).Pkg

type eventInfo struct {
	srv            *protogen.Service
	mth            *protogen.Method
	job            *cloudeventpb.CloudEventServiceOptions
	subject        *cloudeventpb.CloudEventMethodOptions
	jobName        string
	jobSubjectName string
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".cloudevent.pb.go"
	genFile := jen.NewFile(string(file.GoPackageName))
	genFile.HeaderComment("Code generated by protoc-gen-go-cloudevent. DO NOT EDIT.")
	genFile.HeaderComment("versions:")
	genFile.HeaderComment(fmt.Sprintf("  - protoc-gen-go-cloudevent %s", Version))
	genFile.HeaderComment(fmt.Sprintf("  - protoc               %s", protocVersion(gen)))
	if file.Proto.GetOptions().GetDeprecated() {
		genFile.HeaderComment(fmt.Sprintf("%s is a deprecated file.", file.Desc.Path()))
	} else {
		genFile.HeaderComment(fmt.Sprintf("source: %s", file.Desc.Path()))
	}

	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.Skip()

	if len(file.Services) == 0 {
		return g
	}

	var events = make(map[string]map[string]*eventInfo)
	for _, srv := range file.Services {
		job, ok := proto.GetExtension(srv.Desc.Options(), cloudeventpb.E_Job).(*cloudeventpb.CloudEventServiceOptions)
		if !ok || job == nil {
			continue
		}

		jobName := strings.TrimSpace(job.Name)

		if events[jobName] == nil {
			events[jobName] = map[string]*eventInfo{}
		}

		for _, m := range srv.Methods {
			jobSubject, ok := proto.GetExtension(m.Desc.Options(), cloudeventpb.E_Subject).(*cloudeventpb.CloudEventMethodOptions)
			if !ok || jobSubject == nil {
				continue
			}

			jobSubjectName := strings.TrimSpace(jobSubject.Name)
			if events[jobName][jobSubjectName] != nil {
				gen.Error(fmt.Errorf("job:%s subject:%s exists", jobName, jobSubjectName))
				return g
			}

			events[jobName][jobSubjectName] = &eventInfo{
				srv:            srv,
				mth:            m,
				job:            job,
				subject:        jobSubject,
				jobName:        jobName,
				jobSubjectName: jobSubjectName,
			}
		}
	}

	if len(events) == 0 {
		return g
	}

	jobNames := lo.Keys(events)
	sort.Strings(jobNames)

	for _, jobName := range jobNames {
		subjects := events[jobName]
		if len(subjects) == 0 {
			continue
		}
		g.Unskip()

		srvInfo := getSrv(subjects)
		jobKeyPrefix := strings.ReplaceAll(srvInfo.GoName, "InnerService", "")
		jobKeyPrefix = strings.ReplaceAll(jobKeyPrefix, "Inner", "")
		jobKeyName := fmt.Sprintf("%sCloudEventKey", jobKeyPrefix)
		genFile.Const().
			Id(jobKeyName).
			Op("=").
			Lit(jobName)

		subjectNames := lo.Keys(subjects)
		sort.Strings(subjectNames)

		for _, subName := range subjectNames {
			info := subjects[subName]
			var keyName = fmt.Sprintf("%sCloudEventKey", info.mth.GoName)
			genFile.Commentf("%s /%s/%s", keyName, info.srv.Desc.FullName(), info.mth.GoName)
			genFile.Commentf(strings.TrimSpace(info.mth.Comments.Leading.String()))
			genFile.Const().
				Id(keyName).
				Op("=").
				Lit(subName)
		}

		//for _, subName := range subjectNames {
		//	//info := subjects[subName]
		//	//var dd = string(assert.Must1(protojson.Marshal(info.subject)))
		//	//var keyName = fmt.Sprintf("%sCloudEventKey", info.mth.GoName)
		//	//genFile.Var().Id("_").Op("=").
		//	//	Qual(cloudeventPkg, "RegisterSubject").
		//	//	Call(
		//	//		jen.Id(keyName),
		//	//		jen.Lit(fmt.Sprintf("/%s/%s", info.srv.Desc.FullName(), info.mth.GoName)),
		//	//		jen.Func().Params().Params(jen.Op("*").Qual(jobTypesPkg, "CloudEventSubject")).
		//	//			BlockFunc(func(group *jen.Group) {
		//	//				group.Var().Id("data=[]byte").Call(jen.Lit(fmt.Sprintf(`%s`, dd)))
		//	//				group.Var().Id("p").Qual(jobTypesPkg, "CloudEventSubject")
		//	//				group.Qual(assertPkt, "Must").Call(jen.Qual(protojsonPkt, "Unmarshal").Call(jen.Id("data, &p")))
		//	//				group.Return().Id("&p")
		//	//			}).Call(),
		//	//	).Line()
		//}

		subjectValues := lo.Values(subjects)
		var cloudEventName = fmt.Sprintf("%sCloudEvent", subjectValues[0].srv.GoName)
		genFile.Type().Id(cloudEventName).StructFunc(func(group *jen.Group) {
			for _, ss := range subjectValues {
				group.Id("On"+ss.mth.GoName).Func().Params(
					jen.Id("ctx").Qual(ctxPkg, "Context"),
					jen.Id("req").Op("*").Add(getPkg(file, ss.mth.Input.GoIdent)),
				).Params(jen.Error())
			}
		})

		genFile.Func().
			Id(fmt.Sprintf("Register%s", cloudEventName)).
			Params(
				jen.Id("jobCli").Op("*").Qual(cloudeventPkg, "Client"),
				jen.Id("event").Id(cloudEventName),
				jen.Id("opts").Op("...").Op("*").Qual(jobTypesPkg, "RegisterJobOptions"),
			).BlockFunc(func(group *jen.Group) {
			for _, ss := range subjectValues {
				var keyName = fmt.Sprintf("%sCloudEventKey", ss.mth.GoName)
				group.If(jen.Id("event").Dot("On" + ss.mth.GoName)).Op("!=").Nil().BlockFunc(func(group *jen.Group) {
					group.Qual(cloudeventPkg, "RegisterJobHandler").Call(
						jen.Id("jobCli"),
						jen.Id(jobKeyName),
						jen.Id(keyName),
						jen.Id("event").Dot("On"+ss.mth.GoName),
						//jen.Qual(cloudeventPkg, "WrapHandler").Call(jen.Id("event").Dot(ss.mth.GoName)),
						jen.Id("opts").Op("..."),
					)
				}).Line()
			}
		})

		var publisher = fmt.Sprintf("%sPublisher", cloudEventName)
		genFile.Type().Id(publisher).StructFunc(func(group *jen.Group) {
			group.Id("Client").Op("*").Qual(cloudeventPkg, "Client")
		})
		for _, ss := range subjectValues {
			var mthName = fmt.Sprintf("Push%sEvent", ss.mth.GoName)
			var keyName = fmt.Sprintf("%sCloudEventKey", ss.mth.GoName)
			genFile.Func().
				Params(jen.Id(fmt.Sprintf("a %s", publisher))).
				Id(mthName).
				Params(
					jen.Id("ctx").Qual("context", "Context"),
					jen.Id("req").Op("*").Id(ss.mth.Input.GoIdent.GoName),
					jen.Id("opts").Op("...").Op("*").Qual(jobTypesPkg, "PushEventOptions"),
				).
				Params(jen.Op("*").Qual(cloudeventPkg, "PubAckInfo"), jen.Error()).
				Block(jen.Return().Id("a.Client").Dot("Publish").Call(
					jen.Id("ctx"),
					jen.Id(keyName),
					jen.Id("req"),
					jen.Id("opts").Op("..."),
				))
		}
	}

	g.P(genFile.GoString())
	return g
}

func handlerPushEventName(name string, prefix string) string {
	if strings.HasPrefix(name, prefix) {
		return name
	}
	return fmt.Sprintf("%s%s", prefix, name)
}

func getSrv(data map[string]*eventInfo) *protogen.Service {
	for _, srv := range data {
		return srv.srv
	}
	return nil
}

func getPkg(file *protogen.File, goIdent protogen.GoIdent) *jen.Statement {
	var pkgName = ""
	if file.GoImportPath != goIdent.GoImportPath {
		pkgName = string(goIdent.GoImportPath)
	}

	return jen.Qual(pkgName, goIdent.GoName)
}
