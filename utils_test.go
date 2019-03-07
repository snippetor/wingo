package wingo

import "testing"

func BenchmarkBytes2String(b *testing.B) {
	bytes := []byte("testOK 拉斯卡的积分类似地方 历史课lskdjflskjdflskjdflsdsldfjsldf ")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = Bytes2String(bytes)
	}
	b.StopTimer()
}

func BenchmarkString2Bytes(b *testing.B) {
	str := "testOK 拉斯卡的积分类似地方 历史课lskdjflskjdflskjdflsdsldfjsldf "
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = String2Bytes(str)
	}
	b.StopTimer()
}

func BenchmarkBytes2String1(b *testing.B) {
	bytes := []byte("testOK 拉斯卡的积分类似地方 历史课lskdjflskjdflskjdflsdsldfjsldf ")
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = string(bytes)
	}
	b.StopTimer()
}

func BenchmarkString2Bytes1(b *testing.B) {
	str := "testOK 拉斯卡的积分类似地方 历史课lskdjflskjdflskjdflsdsldfjsldf "
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = []byte(str)
	}
	b.StopTimer()
}