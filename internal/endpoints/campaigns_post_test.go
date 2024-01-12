package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrique1501/email-N/internal/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type serviceMock struct {
	mock.Mock
}

func (r *serviceMock) Create(newCampaing contracts.NewCampaingDTO) (string, error) {
	args := r.Called(newCampaing)

	return args.String(0), args.Error(1)
}

func Test_CampaignsPost_Should_Save_New_Campaign(t *testing.T) {
	assert := assert.New(t)
	body := contracts.NewCampaingDTO{
		Name:    "test",
		Content: "Hi there!",
		Emails:  []string{"test@test.com"},
	}
	service := new(serviceMock)
	service.On("Create", mock.MatchedBy(func(request contracts.NewCampaingDTO) bool {
		if request.Name == body.Name && request.Content == body.Content {
			return true
		} else {
			return false
		}
	})).Return("34x", nil)
	handler := Handler{CampaignService: service}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)
	rr := httptest.NewRecorder()

	_, status, err := handler.CampaignPost(rr, req)

	assert.Equal(http.StatusCreated, status)
	assert.Nil(err)
}

func Test_CampaignsPost_Should_Inform_Error_When_Exists(t *testing.T) {
	assert := assert.New(t)
	body := contracts.NewCampaingDTO{
		Name:    "test",
		Content: "Hi there!",
		Emails:  []string{"test@test.com"},
	}
	service := new(serviceMock)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := Handler{CampaignService: service}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)
	rr := httptest.NewRecorder()

	_, _, err := handler.CampaignPost(rr, req)

	// assert.Equal(id, status)
	assert.NotNil(err)
}