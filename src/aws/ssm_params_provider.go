package aws

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/samber/lo"
	"pht/comments-processor/utils"
	"strings"
)

type SSMParamsProvider struct {
	client     *ssm.Client
	paramsPath string
	params     map[string]string
}

func NewSSMParamsProvider(client *ssm.Client, paramsPath string) *SSMParamsProvider {
	return &SSMParamsProvider{
		client:     client,
		paramsPath: paramsPath,
	}
}

func (p *SSMParamsProvider) PrefetchParams() error {
	output, err := p.client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           &p.paramsPath,
		WithDecryption: utils.Ptr(true),
		Recursive:      utils.Ptr(true),
	})
	if err != nil {
		return err
	}
	if len(output.Parameters) == 0 {
		return errors.New("no params fetched")
	}

	err = verifyParams(output)
	if err != nil {
		return errors.Join(errors.New("error when reading params"), err)
	}

	p.params = lo.Associate(output.Parameters, func(param types.Parameter) (string, string) {
		return strings.Replace(*param.Name, fmt.Sprintf("%s/", p.paramsPath), "", 1), *param.Value
	})

	return nil
}

func (p *SSMParamsProvider) GetParam(key string) string {
	return p.params[key]
}

func (p *SSMParamsProvider) UpdateParam(key string, value string) error {
	if !validateParam(&key, &value) {
		return errors.New(fmt.Sprintf("bad param `%s`", key))
	}

	paramPath := fmt.Sprintf("%s/%s", p.paramsPath, key)
	_, err := p.client.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:      &paramPath,
		Value:     &value,
		Overwrite: utils.Ptr(true),
	})
	if err != nil {
		return err
	}

	p.params[key] = value
	return nil
}

func verifyParams(output *ssm.GetParametersByPathOutput) error {
	badParams := lo.Filter(output.Parameters, func(param types.Parameter, _ int) bool {
		return !validateParam(param.Name, param.Value)
	})

	if len(badParams) == 0 {
		return nil
	}

	return errors.Join(lo.Map(badParams, func(param types.Parameter, _ int) error {
		return errors.New(fmt.Sprintf("bad param `%s`", *param.Name))
	})...)
}

func validateParam(name *string, value any) bool {
	return name != nil && value != nil
}
