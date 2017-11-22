package dynamodbadapter

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/util"
)

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	myRes := e.GetPolicy()

	if !util.Array2DEquals(res, myRes) {
		t.Error("Policy: ", myRes, ", supposed to be ", res)
	}
}

func TestAdapter(t *testing.T) {
	// load from CSV file
	e := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	dbEndpoint := "http://localhost:14045"
	region := "ap-southeast-1"

	a := NewAdapter(
		&aws.Config{Endpoint: aws.String(dbEndpoint), Region: aws.String(region)},
		"casbin-rules",
	)

	a.DeleteTable()

	m := e.GetModel()

	a.SavePolicy(m)

	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// Load the policy from DB.
	a.LoadPolicy(e.GetModel())
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	a = NewAdapter(
		&aws.Config{Endpoint: aws.String(dbEndpoint), Region: aws.String(region)},
		"casbin-rules",
	)
	e = casbin.NewEnforcer("examples/rbac_model.conf", a)
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}
