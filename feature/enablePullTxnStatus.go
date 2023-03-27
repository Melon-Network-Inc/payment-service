package feature

import (
	"github.com/Melon-Network-Inc/common/pkg/utils"
)

type EnablePullTxnStatusExperiment struct {
	Production  bool
	Staging     bool
	Development bool
	Testing     bool
}

var EnablePullTxnStatus = EnablePullTxnStatusExperiment{
	Production:  false,
	Staging:     true,
	Development: true,
	Testing:     true,
}

func (t EnablePullTxnStatusExperiment) Set(env string, state bool) {
	switch env {
	case "PROD":
		t.Production = state
	case "STAGING":
		t.Staging = state
	case "DEV":
		t.Development = state
	case "TEST":
		t.Testing = state
	}
}

func (t EnablePullTxnStatusExperiment) Get() bool {
	switch utils.GetEnvironment() {
	case "PROD":
		return t.Production
	case "STAGING":
		return t.Staging
	case "DEV":
		return t.Development
	case "TEST":
		return t.Testing
	default:
		return t.Testing
	}
}
