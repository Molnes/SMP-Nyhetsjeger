package user_roles

type Role int

const (
	User              Role = 0
	QuizAdmin         Role = 1
	OrganizationAdmin Role = 2

	ROLE_CONTEXT_KEY = "user-role" // The key used to store the user role in the context
)

const (
	userString              = "user"
	quizAdminString         = "quiz_admin"
	organizationAdminString = "organization_admin"
)

func (r Role) String() string {
	var roleString string
	switch r {
	case User:
		roleString = userString
	case QuizAdmin:
		roleString = quizAdminString
	case OrganizationAdmin:
		roleString = organizationAdminString
	}
	return roleString
}

func RoleFromString(role string) Role {
	var r Role
	switch role {
	case userString:
		r = User
	case quizAdminString:
		r = QuizAdmin
	case organizationAdminString:
		r = OrganizationAdmin
	}
	return r
}

func (r Role) IsAdministrator() bool {
	return r == QuizAdmin || r == OrganizationAdmin
}