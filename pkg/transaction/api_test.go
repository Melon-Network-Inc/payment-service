package transaction

import "testing"

func TestRegisterHandlers(t *testing.T) {
	expected := "Transaction Received!"
	actual := RegisterHandlers()
	if actual != expected {
		t.Errorf("expected %q but got %q", expected, actual)
	}
}
