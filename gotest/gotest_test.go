package gotest

import "testing"

func Test_one(t *testing.T) {
	t.Log("测试通过")
}
func Test_two(t *testing.T) {
	t.Log("测试也通过")
}
func Benchmark_one(b *testing.B) {
	for i := 0; i < 1000; i++ {
		f1()
	}
}
