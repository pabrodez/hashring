# hashring
Consistent hash ring implementation. Useful for distributing cache among nodes.
## How to use

```go
import (
    "github.com/pabrodez/hashring"
)
// we are using a high replication factor as it ensures a lower standard deviation of the distribution of keys among nodes
nodeHashring := hashring.New([]string{"node1", "node2", "node3", 150})

nodeAddr := nodeHashring.GetNodeForKey("cacheKey")

nodeHashring.RemoveNode("node2")

nodeHashring.AddNode("node4")

```

## References:
- https://authzed.com/blog/consistent-hash-load-balancing-grpc
- https://www.metabrew.com/article/libketama-consistent-hashing-algo-memcached-clients
- https://github.com/RJ/ketama/blob/master/libketama/ketama.c#L194
- https://github.com/Doist/hash_ring/blob/master/README.markdown
- https://web.archive.org/web/20150316185345/http://amix.dk/blog/post/19369
- https://web.archive.org/web/20121130192740/http://amix.dk/blog/viewEntry/19367
- https://michaelnielsen.org/blog/consistent-hashing/
- https://web.archive.org/web/20090122110831/http://lexemetech.com/2007/11/consistent-hashing.html
