package anyhow

import (
	"fmt"
	"log"
)

// GenericErrorHandler demonstrates polymorphic error handling using Catchable interface
func GenericErrorHandler(catchable Catchable, context string) (stdErr error, anyhowErr Error) {
	// Use Catch method to extract standard error
	if catchable.Catch(&stdErr) {
		log.Printf("[%s] Standard error caught: %v", context, stdErr)
	}

	// Use CatchErr method to extract anyhow Error
	if catchable.CatchErr(&anyhowErr) {
		log.Printf("[%s] Anyhow error caught: %v", context, anyhowErr)
	}

	return stdErr, anyhowErr
}

// ProcessMultipleCatchables demonstrates handling multiple Catchable types
func ProcessMultipleCatchables(catchables []Catchable) (stdErrors []error, anyhowErrors []Error) {
	for i, catchable := range catchables {
		var stdErr error
		var anyhowErr Error

		// Both methods can be called on any Catchable
		stdCaught := catchable.Catch(&stdErr)
		anyhowCaught := catchable.CatchErr(&anyhowErr)

		if stdCaught {
			stdErrors = append(stdErrors, fmt.Errorf("item %d: %w", i, stdErr))
		}

		if anyhowCaught {
			anyhowErrors = append(anyhowErrors, ErrorFromFormat("item %d: %v", i, anyhowErr))
		}
	}

	return stdErrors, anyhowErrors
}

// ErrorCollectionWithCatchable demonstrates collecting errors using Catchable interface
func ErrorCollectionWithCatchable(catchables ...Catchable) (totalErrors int, firstError error) {
	for _, catchable := range catchables {
		var err error
		if catchable.Catch(&err) {
			totalErrors++
			if firstError == nil {
				firstError = err
			}
		}
	}

	return totalErrors, firstError
}

// DemoCatchableInterface demonstrates various Catchable interface usages
func DemoCatchableInterface() {
	fmt.Printf("=== Catchable Interface 演示 ===\n\n")

	// Create different types that implement Catchable
	successResult := Ok("success")
	failResult := Fail[string](fmt.Errorf("result error"))
	okError := ErrorOf(nil)
	failError := ErrorOf(fmt.Errorf("error type error"))

	catchables := []Catchable{successResult, failResult, okError, failError}

	// Demo: Generic error handling
	fmt.Printf("通用错误处理:\n")
	for i, catchable := range catchables {
		stdErr, anyhowErr := GenericErrorHandler(catchable, fmt.Sprintf("item_%d", i))
		if stdErr != nil {
			fmt.Printf("   - 标准错误: %v\n", stdErr)
		}
		if anyhowErr.IsError() {
			fmt.Printf("   - Anyhow错误: %v\n", anyhowErr)
		}
		if stdErr == nil && !anyhowErr.IsError() {
			fmt.Printf("   - 无错误\n")
		}
	}

	fmt.Printf("\n✅ Catchable 接口演示完成\n")
}
