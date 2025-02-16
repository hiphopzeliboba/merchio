package user

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"merchio/internal/model"
	"testing"
	"time"
)

type mockPgxPool struct {
	mock pgxmock.PgxPoolIface
}

func TestRepo_CreateUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &repo{
		db: &pgxpool.Pool{},
	}
	ctx := context.Background()

	tests := []struct {
		name          string
		username      string
		password      string
		mockSetup     func(pgxmock.PgxPoolIface)
		expectedID    int64
		expectedError bool
	}{
		{
			name:     "Успешное создание пользователя",
			username: "testuser",
			password: "password",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(int64(1))
				m.ExpectQuery("INSERT INTO users").
					WithArgs("testuser", "password").
					WillReturnRows(rows)
			},
			expectedID:    1,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			id, err := repo.CreateUser(ctx, tt.username, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepo_GetUserByUsername(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &repo{
		db: &pgxpool.Pool{},
	}
	ctx := context.Background()

	expectedUser := &model.User{
		ID:        1,
		Username:  "testuser",
		Password:  "hashedpass",
		Coins:     1000,
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name          string
		username      string
		mockSetup     func(pgxmock.PgxPoolIface)
		expectedUser  *model.User
		expectedError bool
	}{
		{
			name:     "Успешное получение пользователя",
			username: "testuser",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "username", "password", "coins", "created_at"}).
					AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Password,
						expectedUser.Coins, expectedUser.CreatedAt)
				m.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expectedUser:  expectedUser,
			expectedError: false,
		},
		{
			name:     "Пользователь не найден",
			username: "nonexistent",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs("nonexistent").
					WillReturnError(pgx.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			user, err := repo.GetUserByUsername(ctx, tt.username)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Username, user.Username)
				assert.Equal(t, tt.expectedUser.Password, user.Password)
				assert.Equal(t, tt.expectedUser.Coins, user.Coins)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepo_IsUserPresent(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	repo := &repo{
		db: &pgxpool.Pool{},
	}
	ctx := context.Background()

	tests := []struct {
		name           string
		username       string
		mockSetup      func(pgxmock.PgxPoolIface)
		expectedExists bool
		expectedError  bool
	}{
		{
			name:     "Пользователь существует",
			username: "existinguser",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
				m.ExpectQuery("SELECT EXISTS").
					WithArgs("existinguser").
					WillReturnRows(rows)
			},
			expectedExists: true,
			expectedError:  false,
		},
		{
			name:     "Пользователь не существует",
			username: "nonexistent",
			mockSetup: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
				m.ExpectQuery("SELECT EXISTS").
					WithArgs("nonexistent").
					WillReturnRows(rows)
			},
			expectedExists: false,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			exists, err := repo.IsUserPresent(ctx, tt.username)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
