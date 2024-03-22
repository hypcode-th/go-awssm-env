package auto

import (
	"context"
	"fmt"
	"github.com/hypcode-th/go-awssm-env/awssm"
	"os"
	"strings"
)

func init() {
	client := awssm.NewClient()
	ctx := context.Background()

	for _, e := range os.Environ() {
		p := strings.SplitN(e, "=", 2)
		if len(p) < 2 {
			continue
		}
		k, v := p[0], p[1]
		if !client.IsReference(v) {
			continue
		}

		resolved, err := client.Resolve(ctx, v)
		if err != nil {
			fmt.Printf("failed to resolve a reference '%s'. %s", v, err.Error())
			continue
		}

		if err := os.Setenv(k, resolved); err != nil {
			fmt.Printf("failed to set environment '%s'. %s", k, err.Error())
			continue
		}

	}
}
