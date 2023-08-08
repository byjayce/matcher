# Matcher

Gale-Shapley algorithm implementation in Go.

## Installation

```bash
go get github.com/byjayce/matcher
```

## Usage

```go
package person

import (
    "fmt"

    "github.com/byjayce/matcher"
)


// Man is matcher.Evaluator, which evaluates opponent.
type Man struct {
    id string
    // ...
}

func (m *Man) ID() string {
    return m.id
}

func (m *Map) Evaluate(w *Woman) int {
    // calculate the score
    return calc(w)
}
```

```go
package main

import (
	"example.com/person"
	"github.com/byjayce/matcher"
)

func main() {
	var (
		women []*person.Woman = person.InitWomen()
		men   []*person.Man   = person.InitMen()
	)

	res := matcher.Match(men, women)
	// ... use matched result.
}
```

## Note

The order of proposals affects the outcome. If a proposer with the same rating is selected when he proposes first.
