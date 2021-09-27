package sts

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
)

var (
	testAccount = "test-account"
)

func (m *StsMock) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, opts ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	args := m.Called(ctx, input)
	output := args.Get(0)
	err := args.Error(1)

	if output != nil {
		return output.(*sts.GetCallerIdentityOutput), err
	}
	return nil, err
}

func TestClient_GetAccount(t *testing.T) {
	client := NewMockClient()
	getCallerIdentityOutput := &sts.GetCallerIdentityOutput{Account: &testAccount}
	client.sts.(*StsMock).On("GetCallerIdentity", context.Background(), &sts.GetCallerIdentityInput{}).Return(getCallerIdentityOutput, nil)
	account, err := client.GetAccount()
	assert.NoError(t, err)
	assert.Equal(t, testAccount, account)
	client.sts.(*StsMock).AssertExpectations(t)
}

func TestClient_GetAccount_WithError(t *testing.T) {
	client := NewMockClient()
	getCallerIdentityError := fmt.Errorf("some error")
	client.sts.(*StsMock).On("GetCallerIdentity", context.Background(), &sts.GetCallerIdentityInput{}).Return(nil, getCallerIdentityError)
	_, err := client.GetAccount()
	assert.Error(t, err, getCallerIdentityError)
}

func TestClient_GetAccount_WithNullAccount(t *testing.T) {
	client := NewMockClient()
	getCallerIdentityOutput := &sts.GetCallerIdentityOutput{}
	client.sts.(*StsMock).On("GetCallerIdentity", context.Background(), &sts.GetCallerIdentityInput{}).Return(getCallerIdentityOutput, nil)
	_, err := client.GetAccount()
	assert.Error(t, err, fmt.Errorf("unable to determine account ID"))
	client.sts.(*StsMock).AssertExpectations(t)
}
