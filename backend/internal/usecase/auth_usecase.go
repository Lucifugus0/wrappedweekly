package usecase

import (
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"wrappedweekly/backend/internal/domain"
	"wrappedweekly/backend/pkg/apperror"
)

type AuthUsecase struct {
	users domain.UserRepository
	jwt   *JWTManager
}

func NewAuthUsecase(users domain.UserRepository, jwt *JWTManager) *AuthUsecase {
	return &AuthUsecase{users: users, jwt: jwt}
}

type RegisterInput struct {
	Email    string
	Password string
	Name     string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string
	User  *domain.User
}

func (u *AuthUsecase) Register(in RegisterInput) (*domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	name := strings.TrimSpace(in.Name)

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, apperror.BadRequest("format email tidak valid")
	}
	if len(in.Password) < 8 {
		return nil, apperror.BadRequest("password minimal 8 karakter")
	}
	if name == "" {
		return nil, apperror.BadRequest("nama wajib diisi")
	}

	existing, err := u.users.FindByEmail(email)
	if err != nil {
		return nil, apperror.Internal("gagal memeriksa email")
	}
	if existing != nil {
		return nil, apperror.Conflict("email sudah terdaftar")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("gagal memproses password")
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
	}
	if err := u.users.Create(user); err != nil {
		return nil, apperror.Internal("gagal membuat user")
	}
	return user, nil
}

func (u *AuthUsecase) Login(in LoginInput) (*AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	user, err := u.users.FindByEmail(email)
	if err != nil {
		return nil, apperror.Internal("gagal memeriksa email")
	}
	if user == nil {
		return nil, apperror.Unauthorized("email atau password salah")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, apperror.Unauthorized("email atau password salah")
	}

	token, err := u.jwt.Generate(user.ID)
	if err != nil {
		return nil, apperror.Internal("gagal membuat token")
	}

	return &AuthResult{Token: token, User: user}, nil
}

func (u *AuthUsecase) Me(userID string) (*domain.User, error) {
	user, err := u.users.FindByID(userID)
	if err != nil {
		return nil, apperror.Internal("gagal mengambil user")
	}
	if user == nil {
		return nil, apperror.NotFound("user tidak ditemukan")
	}
	return user, nil
}
