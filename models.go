type UserContext struct {
	UserID             string
	Roles              []string
	AllowedDepartments []string
}

type Document struct {
	DocID        string
	Content      string
	Ownership    string
	SecurityTags []string
}

type PermissionProvider interface {
	GetAllowedTags(user UserContext) ([]string, error)
}

type MockPermissionProvider struct {
	Rules map[string][]string
}

func (m *MockPermissionProvider) GetAllowedTags(user UserContext) ([]string, error) {
	return m.Rules[user.Role], nil
}