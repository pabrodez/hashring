Consistent hashing implemented in Go.

```go
import (
    "github.com/pabrodez/hashring"
)

hashring := hashring.NewHashRing([]string{"node1", "node2", "node3", 100})

```