# bug25519
Possible bug in Go x/crypto/curve25519. Go issue [#30095](https://golang.org/issue/30095).

## Problem

The [testvectors.json](testdata/testvectors.json) file contains test vectors generated from the curve25519 implementation in Boring SSL. These test vectors pass on `amd64`:

```
$ go version
go version go1.11 darwin/amd64
$ go test -v -run Current
=== RUN   TestTestVectorsCurrent
--- PASS: TestTestVectorsCurrent (0.00s)
    testvectors_test.go:47: failed 0 of 32
PASS
ok  	github.com/mmcloughlin/bug25519	0.006s
```

However the pure Go version does not (induced with the `appengine` tag)

```
$ go test -v -run Current -tags appengine | head
=== RUN   TestTestVectorsCurrent
--- FAIL: TestTestVectorsCurrent (0.01s)
    testvectors_test.go:39:     in = 668fb9f76ad971c81ac900071a1560bce2ca00cac7e67af99348913761434014
    testvectors_test.go:40:   base = db5f32b7f841e7a1a00968effded12735fc47a3eb13b579aacadeae80939a7dd
    testvectors_test.go:41:    got = 78202e24db99e237f2a14f9ec61b051814ec8fd23a5e8e68add48d66fd09fc12
    testvectors_test.go:42: expect = 090d85e599ea8e2beeb61304d37be10ec5c905f9927d32f42a9a0afb3e0b4074
    testvectors_test.go:39:     in = 203161c3159a876a2beaec29d2427fb0c7c30d382cd013d27cc3d393db0daf6f
    testvectors_test.go:40:   base = 6ab95d1abe68c09b005c3db9042cc91ac849f7e94a2a4a9b893678970b7b95bf
    testvectors_test.go:41:    got = 2ec45ca394a3febc6d63b8995ae63b38c7ba909bafed2a039dd54973f2b5be73
    testvectors_test.go:42: expect = 11edaedc95ff78f563a1c8f15591c071dea092b4d7ecaac8e0387b5a160c4e5d
```

## Fix

Comparison against the `ref10` implementation suggests the following fix.

```
diff --git a/curve25519/curve25519.go b/curve25519/curve25519.go
index cb8fbc5..75f24ba 100644
--- a/curve25519/curve25519.go
+++ b/curve25519/curve25519.go
@@ -86,7 +86,7 @@ func feFromBytes(dst *fieldElement, src *[32]byte) {
        h6 := load3(src[20:]) << 7
        h7 := load3(src[23:]) << 5
        h8 := load3(src[26:]) << 4
-       h9 := load3(src[29:]) << 2
+       h9 := (load3(src[29:]) & 0x7fffff) << 2

        var carry [10]int64
        carry[9] = (h9 + 1<<24) >> 25
```

This is implemented in the [`fixed`](fixed) package.

```
$ go test -v -run Fixed -tags appengine
=== RUN   TestTestVectorsFixed
--- PASS: TestTestVectorsFixed (0.01s)
    testvectors_test.go:47: failed 0 of 32
PASS
ok  	github.com/mmcloughlin/bug25519	0.010s
```
