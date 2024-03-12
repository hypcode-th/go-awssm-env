package auto

import (
	"context"
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

		if resolved, ok := client.Resolve(ctx, v); ok {
			if err := os.Setenv(k, resolved); err != nil {
				continue
			}
		}
	}
}
