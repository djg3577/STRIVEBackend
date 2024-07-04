package user

type UserService struct {
    repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *User) error {
    return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByID(id string) (*User, error) {
    return s.repo.GetUserByID(id)
}
