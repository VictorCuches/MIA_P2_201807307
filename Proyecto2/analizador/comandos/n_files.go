package comandos

type NFiles struct {
	Path string
}

func New_NFiles(path string) *NFiles {
	return &NFiles{path}
}
