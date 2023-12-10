package usecase

type Service interface {
	CrateAlias(url string) (string, error)
	GetOrigURL(alias string) (string, error)
}
