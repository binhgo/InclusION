package model

type FcmMessage struct {

	Content string

	// if Username != nil -> send to user's devices
	Username string

	// if device token != nil -> send to 1 device
	DeviceToken string
}
