package ui

var onError = func(error) {}

// SetOnError sets the error forwarder for widget creation.
func SetOnError(onErr func(error)) {
	onError = onErr
}
