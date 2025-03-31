package interfaces

import "net/http"

type IWsUpgrader interface {
	Upgrade(w http.ResponseWriter, r *http.Request)
}
