package service

type service struct {
	db Repo
}

func New(db Repo) Service {
	return &service{db: db}
}

func (s *service) CrateAlias(url string) (string, error) {
	return "", nil
}
func (s *service) GetOrigURL(alias string) (string, error) {
	return "", nil
}
