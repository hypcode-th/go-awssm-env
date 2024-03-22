package awssm

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
	"github.com/hypcode-th/go-awssm-env/awssm/internal"
	"github.com/hypcode-th/go-awssm-env/awssm/option"
	"strings"
	"sync"
)

type Client interface {
	// IsReference return true when the given value starts with ReferencePrefix.
	IsReference(value string) bool

	// Resolve check value IsReference and resolve to the AWS secret value.
	//
	// returns an empty string with ok = false when failed to resolve the value.
	Resolve(ctx context.Context, value string) (resolved string, err error)
}

func NewClient(opts ...option.ClientOption) Client {
	settings := &internal.Settings{
		ReferencePrefix: ReferencePrefix, // default prefix
	}
	for _, opt := range opts {
		opt.Apply(settings)
	}

	smCache, _ := secretcache.New()
	return &client{
		settings:                   settings,
		secretNameToSecretKeyValue: &sync.Map{},
		smCache:                    smCache,
	}
}

type secretKeyValue map[string]string

type client struct {
	settings *internal.Settings

	secretNameToSecretKeyValue *sync.Map

	smCache *secretcache.Cache
}

func (c *client) IsReference(value string) bool {
	return strings.HasPrefix(value, c.settings.ReferencePrefix)
}

func (c *client) Resolve(ctx context.Context, value string) (resolved string, err error) {
	if !c.IsReference(value) {
		return "", errors.New("value is not a reference")
	}

	ref := c.parseReference(value)
	if ref == nil {
		return "", errors.New("failed to parse a reference")
	}

	v, _ := c.secretNameToSecretKeyValue.Load(ref.SecretName)
	secretKv, ok := v.(*secretKeyValue)
	if !ok {
		secrets, err := c.smCache.Client.ListSecretsWithContext(ctx, &secretsmanager.ListSecretsInput{
			Filters: []*secretsmanager.Filter{
				{
					Key: aws.String("name"),
					Values: []*string{
						aws.String(ref.SecretName),
					},
				},
			},
			MaxResults: aws.Int64(1),
		})
		if err != nil {
			return "", err
		}
		if len(secrets.SecretList) == 0 {
			return "", errors.New("secret not found")
		}

		secretString, err := c.smCache.GetSecretStringWithContext(ctx, *secrets.SecretList[0].ARN)
		if err != nil {
			return "", err
		}

		secretKv = &secretKeyValue{}
		if err := json.Unmarshal([]byte(secretString), secretKv); err != nil {
			return "", err
		}
		c.secretNameToSecretKeyValue.Store(ref.SecretName, secretKv)
	}

	if secretKv != nil {
		kv := *secretKv
		secretValue, ok := kv[ref.SecretKey]
		if !ok {
			return "", errors.New("secret key not found")
		}
		return secretValue, nil
	}

	return "", errors.New("secret not found")
}

// parseReference parses `awssm://secretName/secretKey` to Reference.
// returns `nil` when the value is malformed.
func (c *client) parseReference(value string) *Reference {
	if !c.IsReference(value) {
		return nil
	}

	p := strings.SplitN(strings.TrimPrefix(value, c.settings.ReferencePrefix), "/", 2)
	if len(p) != 2 {
		return nil
	}
	return &Reference{
		SecretName: p[0],
		SecretKey:  p[1],
	}
}
