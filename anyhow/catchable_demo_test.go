package anyhow

import (
	"fmt"
	"testing"
)

// TestCatchableDemo demonstrates the Catchable interface usage
func TestCatchableDemo(t *testing.T) {
	fmt.Println("\n=== Catchable Interface 使用演示 ===")

	// 创建不同类型的 Catchable 实例
	successResult := Ok("成功的结果")
	failResult := Fail[string](fmt.Errorf("Result 错误"))
	okError := ErrorOf(nil)
	failError := ErrorOf(fmt.Errorf("Error 类型错误"))

	// 将它们放在一个 Catchable 切片中
	catchables := []Catchable{successResult, failResult, okError, failError}
	names := []string{"成功的Result", "失败的Result", "成功的Error", "失败的Error"}

	fmt.Println("\n1. 使用 Catch 方法提取标准错误:")
	for i, catchable := range catchables {
		var stdErr error
		caught := catchable.Catch(&stdErr)

		if caught {
			fmt.Printf("   %s: 捕获到错误 - %v\n", names[i], stdErr)
		} else {
			fmt.Printf("   %s: 无错误\n", names[i])
		}
	}

	fmt.Println("\n2. 使用 CatchErr 方法提取 anyhow.Error:")
	for i, catchable := range catchables {
		var anyhowErr Error
		caught := catchable.CatchErr(&anyhowErr)

		if caught {
			fmt.Printf("   %s: 捕获到 anyhow.Error - %v\n", names[i], anyhowErr)
		} else {
			fmt.Printf("   %s: 无错误\n", names[i])
		}
	}

	fmt.Println("\n3. 多态错误处理函数:")
	handleAnyError := func(c Catchable, name string) {
		var stdErr error
		var anyhowErr Error

		stdCaught := c.Catch(&stdErr)
		anyhowCaught := c.CatchErr(&anyhowErr)

		if stdCaught || anyhowCaught {
			fmt.Printf("   %s: 检测到错误\n", name)
			if stdCaught {
				fmt.Printf("     - 标准错误: %v\n", stdErr)
			}
			if anyhowCaught {
				fmt.Printf("     - Anyhow错误: %v\n", anyhowErr)
			}
		} else {
			fmt.Printf("   %s: 一切正常\n", name)
		}
	}

	for i, catchable := range catchables {
		handleAnyError(catchable, names[i])
	}

	fmt.Println("\n✅ Catchable 接口演示完成")

	// 确保测试通过
	t.Log("Catchable interface demo completed successfully")
}
