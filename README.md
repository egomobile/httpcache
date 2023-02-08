# httpcache

An easy to use httpcache module that can cache a key value pair.

# Example usage

```go
package main

import (
    "github.com/egomobile/httpcache"
)

func main() {
	httpcache.Server()
}
```

# API

```bash
curl --location --request PUT 'http://localhost/datasets' \
--header 'Content-Type: application/json' \
--data-raw '{
    "key": "foo3",
    "value": "bar3"
}'
```

```bash
curl --location --request GET 'http://localhost/datasets'
```