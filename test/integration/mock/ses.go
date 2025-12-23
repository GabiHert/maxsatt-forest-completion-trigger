package mock

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/google/uuid"
)

type SesClient struct {
	SendTemplatedEmailError    error
	SendTemplatedEmailFunc     func(ctx context.Context, params *ses.SendTemplatedEmailInput, optFns ...func(*ses.Options)) (*ses.SendTemplatedEmailOutput, error)
	SendTemplatedEmailResponse *ses.SendTemplatedEmailOutput
	SendTemplatedEmailCalls    []SendTemplatedEmailCall
	mu                         sync.Mutex
}

type SendTemplatedEmailCall struct {
	Ctx    context.Context
	Params *ses.SendTemplatedEmailInput
	OptFns []func(*ses.Options)
}

func NewSesClient() *SesClient {
	return &SesClient{
		SendTemplatedEmailCalls: make([]SendTemplatedEmailCall, 0),
	}
}

func (m *SesClient) SendTemplatedEmail(ctx context.Context, params *ses.SendTemplatedEmailInput, optFns ...func(*ses.Options)) (*ses.SendTemplatedEmailOutput, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.SendTemplatedEmailCalls = append(m.SendTemplatedEmailCalls, SendTemplatedEmailCall{
		Ctx:    ctx,
		Params: params,
		OptFns: optFns,
	})

	if m.SendTemplatedEmailFunc != nil {
		return m.SendTemplatedEmailFunc(ctx, params, optFns...)
	}

	if m.SendTemplatedEmailError != nil {
		return nil, m.SendTemplatedEmailError
	}

	if m.SendTemplatedEmailResponse != nil {
		return m.SendTemplatedEmailResponse, nil
	}

	messageId := uuid.New().String()
	return &ses.SendTemplatedEmailOutput{
		MessageId: &messageId,
	}, nil
}

func (m *SesClient) GetSendTemplatedEmailCalls() []SendTemplatedEmailCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.SendTemplatedEmailCalls
}

func (m *SesClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SendTemplatedEmailCalls = make([]SendTemplatedEmailCall, 0)
}
