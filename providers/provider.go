package providers

type Provider interface {
	TtsOnline(word string, filePath string) error 
}
