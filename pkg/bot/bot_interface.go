package bot

type Bot interface {
	SendMessage(message string) (bool, error)
}
