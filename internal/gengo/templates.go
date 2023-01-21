package gengo

type GoDropper interface {
	WriteSrc(string) error
	WriteSharedSrc(string) error
}
