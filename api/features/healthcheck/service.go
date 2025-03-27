package healthcheck

type Service interface {
	Check() error
}

type ServiceImpl struct{}

func NewService() *ServiceImpl {
	return &ServiceImpl{}
}

func (s *ServiceImpl) Check() error {
	return nil
}
