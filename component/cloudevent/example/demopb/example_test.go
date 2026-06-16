package demopb_test

import (
	"fmt"

	"github.com/pubgo/funk/v2/component/cloudevent/example/demopb"
)

func ExampleHelloExecCloudEventSubjectKey() {
	fmt.Println(demopb.HelloExecCloudEventSubjectKey)

	// Output: demo.hello.exec
}
