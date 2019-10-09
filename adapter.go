package dynamodbadapter

import (
	"errors"
	"regexp"

	"github.com/google/uuid"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/guregu/dynamo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type (
	// Adapter structs holds dynamoDB config and service	Adapter struct {
	Adapter struct {
		Config         *aws.Config
		Service        *dynamodb.DynamoDB
		DB             *dynamo.DB
		DataSourceName string
	}

	CasbinRule struct {
		ID    string `dynamo:"ID"`
		PType string `dynamo:"PType"`
		V0    string `dynamo:"V0"`
		V1    string `dynamo:"V1"`
		V2    string `dynamo:"V2"`
		V3    string `dynamo:"V3"`
		V4    string `dynamo:"V4"`
		V5    string `dynamo:"V5"`
	}
)

// NewAdapter is the constructor for adapter
func NewAdapter(config *aws.Config, ds string) *Adapter {
	a := &Adapter{}
	a.Config = config
	a.DataSourceName = ds
	a.Service = dynamodb.New(session.New(config), a.Config)
	a.DB = dynamo.New(session.New(), a.Config)
	return a
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

func (a *Adapter) LoadPolicy(model model.Model) {
	p, err := a.getAllItems()
	if err != nil {
		panic(err)
	}

	for _, v := range p {
		loadPolicyLine(v, model)
	}

	return
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	id := uuid.New().String()
	line := CasbinRule{
		ID: id,
	}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

func (a *Adapter) SavePolicy(model model.Model) error {
	a.DeleteTable()
	a.CreateTable()

	var lines []CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	_, err := a.saveItems(lines)
	a.LoadPolicy(model)
	return err
}

func (a *Adapter) saveItems(rules []CasbinRule) (int, error) {
	items := make([]interface{}, len(rules))

	for i := 0; i < len(rules); i++ {
		items[i] = rules[i]
	}

	return a.DB.Table(a.DataSourceName).Batch().Write().Put(items...).Run()
}

func (a *Adapter) getAllItems() ([]CasbinRule, error) {
	var rule []CasbinRule
	err := a.DB.Table(a.DataSourceName).Scan().All(&rule)
	if err != nil {
		return nil, err
	}
	return rule, nil
}

// CreateTable has response for create new table for store
func (a *Adapter) CreateTable() (*dynamodb.CreateTableOutput, error) {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(a.DataSourceName),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}

	out, err := a.Service.CreateTable(params)

	if err != nil {
		matched, err := regexp.MatchString("ResourceInUseException: Cannot create preexisting table", err.Error())
		if err != nil {
			return nil, err
		}

		if !matched {
			return nil, err
		}
	}

	return out, nil
}

// DeleteTable should delete a table
func (a *Adapter) DeleteTable() error {
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(a.DataSourceName),
	}
	_, err := a.Service.DeleteTable(params)
	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
