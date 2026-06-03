package color

import "testing"

func TestIsTerminal(t *testing.T) {
	// 在测试环境中，stdout 通常是管道，不是终端
	result := IsTerminal()
	// 测试环境中应该返回 false（因为是管道输出）
	if result {
		t.Log("Running in terminal, IsTerminal() = true")
	} else {
		t.Log("Running in pipe, IsTerminal() = false")
	}
}

func TestSupportsColor(t *testing.T) {
	result := SupportsColor()
	t.Logf("SupportsColor() = %v", result)
}
