package service

import (
	"encoding/json"
	"github.com/bopoh24/ma_1/notifier/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTplToString(t *testing.T) {

	message := "{\"offer\":{\"id\":123,\"service_id\":20,\"service_name\":\"Кардиолог\",\"company_id\":97,\"company_name\":\"\",\"datetime\":\"2024-02-20T21:00:00Z\",\"description\":\"\",\"price\":12.5,\"status\":\"failed\"},\"type\":\"booking_failed\",\"fail_reason\":\"insufficient funds\",\"company_contacts\":{\"name\":\"Lemke, Crooks and Turcotte\",\"email\":\"\",\"phone\":\"669-524-0414\",\"address\":\"983 Dejah Forks\"},\"customer_contacts\":{\"email\":\"mikel_gaylord@gmail.com\",\"first_name\":\"Frances\",\"last_name\":\"Kerluke\"},\"company_manager_contacts\":[{\"email\":\"odessa.damore45@hotmail.com\"}]}"

	notification := model.BookingNotification{}
	notification.Status = "failed"
	err := json.Unmarshal([]byte(message), &notification)
	assert.NoError(t, err)

	srv := Service{}
	tpl := "company"
	got, err := srv.tplToString(tpl, notification)
	assert.NoError(t, err)
	assert.NotEmpty(t, got)
	t.Log(got)
}
