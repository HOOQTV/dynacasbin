DynamoDB Adapter
====

DynamoDB Adapter is the [DynamoDB](https://aws.amazon.com/dynamodb/) adapter for [Casbin](https://github.com/casbin/casbin). With this library, Casbin can load policy from DynamoDB or save policy to it.

## Installation

    go get github.com/hooqtv/dynacasbin

## Simple Example

```go
package main

import (
	"github.com/casbin/casbin"
	"github.com/hooqtv/dynacasbin"
	"github.com/aws/aws-sdk-go/aws"
)

func main() {
	// Initialize a DynamoDB adapter and use it in a Casbin enforcer:
	config := &aws.Config{} // Your AWS configuration
	ds := "casbin-rules"
	a := dynacasbin.NewAdapter(config, ds) // Your aws configuration and data source.
	e := casbin.NewEnforcer("examples/rbac_model.conf", a)

	// Load the policy from DB.
	e.LoadPolicy()

	// Check the permission.
	e.Enforce("alice", "data1", "read")

	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)

	// Save the policy back to DB.
	e.SavePolicy()
}
```

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
