**Замеры скорости вставки** с большим количеством параметров для подходов описанных в статье https://klotzandrew.com/blog/postgres-passing-65535-parameter-limit

RegularInsert - функция с обычной вставкой, используя неоптимальное формирование строки

RegularInsertSB - обычная вставка, используя strings.Builder (оптимизация, предложенная в комментах)

UnnestInsert - вставка, используя unnest

|                 | 50_000 | 134_000 | 300_000 |
|-----------------|--------|---------|---------|
| RegularInsert   | 5s     | 34s*    | -       |
| RegularInsertSB | 89ms   | 116ms*  | 249ms*  |
| UnnestInsert    | 24ms   | 70ms    | 161ms   |

В RegularInsert и RegularInsertSB на значениях больше 65535 происходит ошибка: _extended protocol limited to 65535 parameters_

Замеры на localhost docker postgres

cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz

###  Profiling

`go test . -bench=. -benchmem -cpuprofile=cpu.out -memprofile=mem.out`

| 50_000 parameters        |    ns/op |     B/op | allocs/op |
|--------------------------|---------:|---------:|----------:|
| BenchmarkRegularInsert   |        - |        - |         - |
| BenchmarkRegularInsertSB | 35444475 | 19124752 |    149040 |
| BenchmarkUnnestInsert    | 23984159 |  5148908 |     99579 |

### RegularInsertSB

#### CPU

`go test . -bench=. -benchmem -cpuprofile=cpuRegularInsertSB.out`

`go tool pprof cpuRegularInsertSB.out`

```
      flat  flat%   sum%        cum   cum%
      60ms  7.50%  7.50%      150ms 18.75%  runtime.mallocgc
      60ms  7.50% 15.00%       60ms  7.50%  runtime.memclrNoHeapPointers
      40ms  5.00% 20.00%       80ms 10.00%  runtime.greyobject
      30ms  3.75% 23.75%      220ms 27.50%  github.com/jackc/pgx/v5/pgtype.(*Map).Encode
      30ms  3.75% 27.50%       60ms  7.50%  runtime.findObject
      30ms  3.75% 31.25%       70ms  8.75%  runtime.mapaccess1
      30ms  3.75% 35.00%       30ms  3.75%  runtime.memhash32
      30ms  3.75% 38.75%       30ms  3.75%  runtime.memmove
      20ms  2.50% 41.25%      100ms 12.50%  github.com/jackc/pgx/v5/pgtype.(*Map).PlanEncode
      20ms  2.50% 43.75%       20ms  2.50%  runtime.(*itabTableType).find
      20ms  2.50% 46.25%       20ms  2.50%  runtime.(*mspan).base (inline)
      20ms  2.50% 48.75%       20ms  2.50%  runtime.(*mspan).divideByElemSize (inline)
      20ms  2.50% 51.25%       20ms  2.50%  runtime.futex
      20ms  2.50% 53.75%       20ms  2.50%  runtime.heapBits.nextFast (inline)
      20ms  2.50% 56.25%       20ms  2.50%  runtime.interhash
      20ms  2.50% 58.75%       20ms  2.50%  runtime.madvise
      20ms  2.50% 61.25%       60ms  7.50%  runtime.mapaccess2_fast32
      20ms  2.50% 63.75%       20ms  2.50%  runtime.procyield
      20ms  2.50% 66.25%      180ms 22.50%  runtime.scanobject
      10ms  1.25% 67.50%       10ms  1.25%  crypto/sha256.block
      10ms  1.25% 68.75%       10ms  1.25%  encoding/binary.bigEndian.PutUint32 (inline)
      10ms  1.25% 70.00%      390ms 48.75%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).Build
      10ms  1.25% 71.25%      370ms 46.25%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).appendParam
      10ms  1.25% 72.50%      240ms 30.00%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).encodeExtendedParamValue
      10ms  1.25% 73.75%       10ms  1.25%  github.com/jackc/pgx/v5/internal/anynil.Is
      10ms  1.25% 75.00%       10ms  1.25%  github.com/jackc/pgx/v5/internal/anynil.NormalizeSlice (inline)
      10ms  1.25% 76.25%       10ms  1.25%  memeqbody
      10ms  1.25% 77.50%       10ms  1.25%  runtime.(*gcBits).bytep (inline)
      10ms  1.25% 78.75%       20ms  2.50%  runtime.(*mspan).markBitsForIndex (inline)
      10ms  1.25% 80.00%       10ms  1.25%  runtime.(*unwinder).resolveInternal
      10ms  1.25% 81.25%       10ms  1.25%  runtime.add (inline)
      10ms  1.25% 82.50%       70ms  8.75%  runtime.convT64
      10ms  1.25% 83.75%       40ms  5.00%  runtime.deductAssistCredit
      10ms  1.25% 85.00%       10ms  1.25%  runtime.evacuated (inline)
      10ms  1.25% 86.25%       40ms  5.00%  runtime.findRunnable
      10ms  1.25% 87.50%       10ms  1.25%  runtime.heapBits.next
      10ms  1.25% 88.75%       10ms  1.25%  runtime.ifaceeq
      10ms  1.25% 90.00%       10ms  1.25%  runtime.mapaccess1_fast32
      10ms  1.25% 91.25%       10ms  1.25%  runtime.pageIndexOf (inline)
      10ms  1.25% 92.50%       10ms  1.25%  runtime.spanOf (inline)
      10ms  1.25% 93.75%       10ms  1.25%  runtime.stealWork
      10ms  1.25% 95.00%       10ms  1.25%  runtime/internal/atomic.(*Uint32).Add
      10ms  1.25% 96.25%       40ms  5.00%  strconv.FormatInt
      10ms  1.25% 97.50%       30ms  3.75%  strconv.formatBits
      10ms  1.25% 98.75%       20ms  2.50%  strings.(*Builder).WriteString (inline)
      10ms  1.25%   100%      130ms 16.25%  unnest.valuesToRowsStringBuilder

```

#### Memory

`go test . -bench=. -benchmem -memprofile=memRegularInsertSB.out`

`go tool pprof memRegularInsertSB.out`

```
Showing nodes accounting for 638.11MB, 99.20% of 643.24MB total
Dropped 26 nodes (cum <= 3.22MB)
      flat  flat%   sum%        cum   cum%
  249.21MB 38.74% 38.74%   291.31MB 45.29%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).appendParam
  155.90MB 24.24% 62.98%   239.66MB 37.26%  unnest.valuesToRowsStringBuilder
   71.26MB 11.08% 74.06%    71.26MB 11.08%  strings.(*Builder).WriteString (inline)
   59.76MB  9.29% 83.35%    80.14MB 12.46%  github.com/jackc/pgx/v5/pgproto3.(*Bind).Encode
   26.60MB  4.14% 87.48%    26.60MB  4.14%  github.com/jackc/pgx/v5/internal/pgio.AppendUint32 (inline)
   20.38MB  3.17% 90.65%    20.38MB  3.17%  github.com/jackc/pgx/v5/internal/pgio.AppendUint16 (inline)
   15.50MB  2.41% 93.06%    42.10MB  6.55%  github.com/jackc/pgx/v5/pgtype.(*wrapIntEncodePlan).Encode
   14.93MB  2.32% 95.38%    27.01MB  4.20%  fmt.Sprintf
   12.50MB  1.94% 97.33%    12.50MB  1.94%  strconv.formatBits
   12.08MB  1.88% 99.20%    12.08MB  1.88%  fmt.(*buffer).writeString (inline)
         0     0% 99.20%    12.08MB  1.88%  fmt.(*fmt).fmtS
         0     0% 99.20%    12.08MB  1.88%  fmt.(*fmt).padString
         0     0% 99.20%    12.08MB  1.88%  fmt.(*pp).doPrintf
         0     0% 99.20%    12.08MB  1.88%  fmt.(*pp).fmtString
         0     0% 99.20%    12.08MB  1.88%  fmt.(*pp).printArg
         0     0% 99.20%   375.53MB 58.38%  github.com/jackc/pgx/v5.(*Conn).Exec
         0     0% 99.20%   375.53MB 58.38%  github.com/jackc/pgx/v5.(*Conn).exec
         0     0% 99.20%   371.45MB 57.75%  github.com/jackc/pgx/v5.(*Conn).execPrepared
         0     0% 99.20%   291.31MB 45.29%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).Build
         0     0% 99.20%    42.10MB  6.55%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).encodeExtendedParamValue
         0     0% 99.20%    20.38MB  3.17%  github.com/jackc/pgx/v5/internal/pgio.AppendInt16 (inline)
         0     0% 99.20%    26.60MB  4.14%  github.com/jackc/pgx/v5/internal/pgio.AppendInt32 (inline)
         0     0% 99.20%    80.14MB 12.46%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecPrepared
         0     0% 99.20%    80.14MB 12.46%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).SendBind
         0     0% 99.20%    42.10MB  6.55%  github.com/jackc/pgx/v5/pgtype.(*Map).Encode
         0     0% 99.20%    26.60MB  4.14%  github.com/jackc/pgx/v5/pgtype.encodePlanInt4CodecBinaryInt64Valuer.Encode
         0     0% 99.20%    12.50MB  1.94%  strconv.FormatInt
         0     0% 99.20%    12.50MB  1.94%  strconv.Itoa (inline)
         0     0% 99.20%   624.73MB 97.12%  testing.(*B).launch
         0     0% 99.20%    18.01MB  2.80%  testing.(*B).run1.func1
         0     0% 99.20%   642.74MB 99.92%  testing.(*B).runN
         0     0% 99.20%   642.74MB 99.92%  unnest.BenchmarkRegularInsertSB
         0     0% 99.20%   642.20MB 99.84%  unnest.RegularInsertSB

```

### UnnestInsert

#### CPU

`go test . -bench=. -benchmem -cpuprofile=cpuUnnestInsert.out`

`go tool pprof cpuUnnestInsert.out`

```
      flat  flat%   sum%        cum   cum%
      60ms 13.95% 13.95%      120ms 27.91%  runtime.mallocgc
      40ms  9.30% 23.26%      280ms 65.12%  github.com/jackc/pgx/v5/pgtype.(*encodePlanArrayCodecBinary).Encode
      30ms  6.98% 30.23%       30ms  6.98%  github.com/jackc/pgx/v5/pgtype.encodePlanInt4CodecBinaryInt64Valuer.Encode
      30ms  6.98% 37.21%       30ms  6.98%  runtime.ifaceeq
      30ms  6.98% 44.19%       30ms  6.98%  runtime.memmove
      30ms  6.98% 51.16%       30ms  6.98%  runtime.nextFreeFast (inline)
      30ms  6.98% 58.14%       30ms  6.98%  runtime/internal/syscall.Syscall6
      20ms  4.65% 62.79%       50ms 11.63%  reflect.packEface
      10ms  2.33% 65.12%       70ms 16.28%  reflect.Value.Interface (inline)
      10ms  2.33% 67.44%       60ms 13.95%  reflect.valueInterface
      10ms  2.33% 69.77%       10ms  2.33%  runtime.(*mspan).markBitsForIndex (inline)
      10ms  2.33% 72.09%       10ms  2.33%  runtime.acquirem (inline)
      10ms  2.33% 74.42%      100ms 23.26%  runtime.convT64
      10ms  2.33% 76.74%       10ms  2.33%  runtime.findObject
      10ms  2.33% 79.07%       60ms 13.95%  runtime.gcBgMarkWorker.func2
      10ms  2.33% 81.40%       20ms  4.65%  runtime.greyobject
      10ms  2.33% 83.72%       10ms  2.33%  runtime.madvise
      10ms  2.33% 86.05%       10ms  2.33%  runtime.memclrNoHeapPointers
      10ms  2.33% 88.37%       10ms  2.33%  runtime.mget
      10ms  2.33% 90.70%       10ms  2.33%  runtime.publicationBarrier
      10ms  2.33% 93.02%       10ms  2.33%  runtime.releasem (inline)
      10ms  2.33% 95.35%       40ms  9.30%  runtime.scanobject
      10ms  2.33% 97.67%       10ms  2.33%  runtime.shrinkstack
      10ms  2.33%   100%       10ms  2.33%  runtime.typedmemmove

```

#### Memory

`go test . -bench=. -benchmem -memprofile=memUnnestInsert.out`

`go tool pprof memUnnestInsert.out`

```
Showing nodes accounting for 359.61MB, 99.63% of 360.95MB total
Dropped 5 nodes (cum <= 1.80MB)
      flat  flat%   sum%        cum   cum%
  131.61MB 36.46% 36.46%   131.61MB 36.46%  github.com/jackc/pgx/v5/internal/pgio.AppendUint32 (inline)
  124.49MB 34.49% 70.95%   359.61MB 99.63%  unnest.UnnestInsert
   43.50MB 12.05% 83.00%    43.50MB 12.05%  github.com/jackc/pgx/v5/pgproto3.(*Bind).Encode
   31.50MB  8.73% 91.73%    31.50MB  8.73%  reflect.packEface
   28.50MB  7.90% 99.63%    86.84MB 24.06%  github.com/jackc/pgx/v5/pgtype.(*wrapIntEncodePlan).Encode
         0     0% 99.63%   235.12MB 65.14%  github.com/jackc/pgx/v5.(*Conn).Exec
         0     0% 99.63%   235.12MB 65.14%  github.com/jackc/pgx/v5.(*Conn).exec
         0     0% 99.63%   235.12MB 65.14%  github.com/jackc/pgx/v5.(*Conn).execPrepared
         0     0% 99.63%   191.62MB 53.09%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).Build
         0     0% 99.63%   191.62MB 53.09%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).appendParam
         0     0% 99.63%   191.62MB 53.09%  github.com/jackc/pgx/v5.(*ExtendedQueryBuilder).encodeExtendedParamValue
         0     0% 99.63%   131.61MB 36.46%  github.com/jackc/pgx/v5/internal/pgio.AppendInt32 (inline)
         0     0% 99.63%    43.50MB 12.05%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecPrepared
         0     0% 99.63%    43.50MB 12.05%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).SendBind
         0     0% 99.63%   191.62MB 53.09%  github.com/jackc/pgx/v5/pgtype.(*Map).Encode
         0     0% 99.63%   191.62MB 53.09%  github.com/jackc/pgx/v5/pgtype.(*encodePlanArrayCodecBinary).Encode
         0     0% 99.63%   191.62MB 53.09%  github.com/jackc/pgx/v5/pgtype.(*wrapSliceEncodeReflectPlan).Encode
         0     0% 99.63%    31.50MB  8.73%  github.com/jackc/pgx/v5/pgtype.anySliceArrayReflect.Index
         0     0% 99.63%    58.34MB 16.16%  github.com/jackc/pgx/v5/pgtype.encodePlanInt4CodecBinaryInt64Valuer.Encode
         0     0% 99.63%    31.50MB  8.73%  reflect.Value.Interface (inline)
         0     0% 99.63%    31.50MB  8.73%  reflect.valueInterface
         0     0% 99.63%   354.29MB 98.15%  testing.(*B).launch
         0     0% 99.63%     5.32MB  1.47%  testing.(*B).run1.func1
         0     0% 99.63%   359.61MB 99.63%  testing.(*B).runN
         0     0% 99.63%   359.61MB 99.63%  unnest.BenchmarkUnnestInsert

```