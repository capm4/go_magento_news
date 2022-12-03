package model

type SenderInterface interface {
	Send(msg string) (bool, error)
	IsRun() (bool, error)
}
