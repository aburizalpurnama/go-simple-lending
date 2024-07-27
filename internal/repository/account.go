package repository

var _ Account = new(accountImpl)

type (
	Account interface {
		Create()
		GetById()
	}

	accountImpl struct{}
)

func NewAccount() *accountImpl {
	return &accountImpl{}
}

func (a *accountImpl) Create() {}

func (a *accountImpl) GetById() {}
