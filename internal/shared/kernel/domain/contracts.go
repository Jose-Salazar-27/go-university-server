package domain

type Storage interface {
	Save()
	BuildObjectURL(folderName, id, filename string) (url string)
}
