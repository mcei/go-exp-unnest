// go test . -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out
// go tool pprof -noinlines cpu.out
// go tool pprof -http :3000 cpu.out

package main

import (
	"context"
	"testing"
)

var num = 25_000

//usersNum = 67_000
//usersNum = 150_000

var values = makeTestUsers(num)

// Capture the time it takes to execute RegularInsert.
func BenchmarkRegularInsert(b *testing.B) {
	conn := setup()
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < b.N; i++ {
		// обычная вставка, используя неоптимальное формирование строки
		RegularInsert(conn, values)
	}
}

// Capture the time it takes to execute RegularInsertSB.
func BenchmarkRegularInsertSB(b *testing.B) {
	conn := setup()
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < b.N; i++ {
		// обычная вставка, используя strings.Builder
		RegularInsertSB(conn, values)
	}
}

// Capture the time it takes to execute UnnestInsert.
func BenchmarkUnnestInsert(b *testing.B) {
	conn := setup()
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < b.N; i++ {
		// вставка, используя unnest
		UnnestInsert(conn, values)
	}
}
