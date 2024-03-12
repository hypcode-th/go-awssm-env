# go-awssm-env

Inspired by [Berglas](https://github.com/GoogleCloudPlatform/berglas).

This tool will resolve the environment referenced to AWS SecretManager to the secret value.

## AWS Permissions

* `secretsmanager:ListSecrets`
* `secretsmanager:GetSecretValue`

## Usage

The reference must follow this format `{ReferencePrefix}{AWS SecretName}/{AWS SecretKey}`

The default `{ReferencePrefix}` is `awssm://`. To change this prefix use `option.WithReferencePreifx("abc://")`

### Manual resolve

```
    reference := "awssm://my-aws-secret-name/MY_SECRET"
    
    client := awssm.NewClient()
    resolved, ok := client.Resolve(reference)
    
    if !ok {
        fmt.Println("not resolved")
    } else {
        fmt.Printf("resolved to '%s'\n", resolved)
    }    
```

### Auto resolve

To automatic resolve the environment variables, import `auto` package to the main file.


```
// main.go

import (
    _ "github.com/hypcode/go-awssm-env/awssm/auto"
)

func main() {
   ...
}

```