package repositories

import (
	"cli-project/internal/config/queries"
	database "cli-project/internal/db"
	"cli-project/internal/domain/interfaces"
	"cli-project/internal/domain/models"
	"cli-project/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
	"time"
)

type userRepo struct{}

func NewUserRepo() interfaces.UserRepository {
	return &userRepo{}
}

func (r *userRepo) getDBConnection() (*sql.DB, error) {
	return database.ClientGetter()
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.StandardUser) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseInsert, map[string]string{
		"table":   "Users",
		"columns": "id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned",
		"values":  "$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11",
	})

	_, err = db.ExecContext(ctx, query,
		user.StandardUser.ID,
		strings.ToLower(user.StandardUser.Username),
		user.StandardUser.Password,
		user.StandardUser.Name,
		strings.ToLower(user.StandardUser.Email),
		user.StandardUser.Role,
		user.LastSeen,
		user.StandardUser.Organisation,
		user.StandardUser.Country,
		user.LeetcodeID,
		user.StandardUser.IsBanned,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUserCreationFailed, err)
	}

	return nil
}

// DeleteUser deletes a user by the given ID
func (r *userRepo) DeleteUser(ctx context.Context, userID string) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	deleteQuery := queries.QueryBuilder(queries.BaseDelete, map[string]string{
		"table":      "users",
		"conditions": "id = $1",
	})

	result, err := db.ExecContext(ctx, deleteQuery, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDeletingUserFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDeletingUserFailed, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: user ID %s", errs.ErrUserNotFound, userID)
	}

	return nil
}

func (r *userRepo) UpdateUserProgress(ctx context.Context, userID uuid.UUID, newSlugs []string) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var existingSlugs []string
	fetchQuery := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "title_slugs",
		"table":      "users_progress",
		"conditions": "user_id = $1",
	})

	err = tx.QueryRowContext(ctx, fetchQuery, userID).Scan(pq.Array(&existingSlugs))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertQuery := queries.QueryBuilder(queries.BaseInsert, map[string]string{
				"table":   "users_progress",
				"columns": "user_id, title_slugs",
				"values":  "$1, $2",
			})
			_, err = tx.ExecContext(ctx, insertQuery, userID, pq.Array(newSlugs))
			if err != nil {
				return fmt.Errorf("%w: %v", errs.ErrUpdatingUserProgressFailed, err)
			}
			return tx.Commit()
		}
		return fmt.Errorf("%w: %v", errs.ErrUpdatingUserProgressFailed, err)
	}

	existingSlugSet := make(map[string]struct{}, len(existingSlugs))
	for _, slug := range existingSlugs {
		existingSlugSet[slug] = struct{}{}
	}

	var slugsToAdd []string
	for _, slug := range newSlugs {
		if _, exists := existingSlugSet[slug]; !exists {
			slugsToAdd = append(slugsToAdd, slug)
		}
	}

	if len(slugsToAdd) > 0 {
		updateQuery := queries.QueryBuilder(queries.BaseUpdate, map[string]string{
			"table":       "users_progress",
			"assignments": "title_slugs = array(SELECT DISTINCT unnest(title_slugs) || unnest($1::text[]))",
			"conditions":  "user_id = $2",
		})
		_, err = tx.ExecContext(ctx, updateQuery, pq.Array(slugsToAdd), userID)
		if err != nil {
			return fmt.Errorf("%w: %v", errs.ErrUpdatingUserProgressFailed, err)
		}
	}

	return tx.Commit()
}

func (r *userRepo) FetchAllUsers(ctx context.Context) ([]models.StandardUser, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelect, map[string]string{
		"columns": "id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned",
		"table":   "Users",
	})

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var users []models.StandardUser

	for rows.Next() {
		var user models.StandardUser
		err := rows.Scan(
			&user.StandardUser.ID,
			&user.StandardUser.Username,
			&user.StandardUser.Password,
			&user.StandardUser.Name,
			&user.StandardUser.Email,
			&user.StandardUser.Role,
			&user.LastSeen,
			&user.StandardUser.Organisation,
			&user.StandardUser.Country,
			&user.LeetcodeID,
			&user.StandardUser.IsBanned,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
	}

	return users, nil
}

func (r *userRepo) FetchUserByID(ctx context.Context, userID string) (*models.StandardUser, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned",
		"table":      "Users",
		"conditions": "id = $1",
	})

	row := db.QueryRowContext(ctx, query, userID)

	var user models.StandardUser
	err = row.Scan(
		&user.StandardUser.ID,
		&user.StandardUser.Username,
		&user.StandardUser.Password,
		&user.StandardUser.Name,
		&user.StandardUser.Email,
		&user.StandardUser.Role,
		&user.LastSeen,
		&user.StandardUser.Organisation,
		&user.StandardUser.Country,
		&user.LeetcodeID,
		&user.StandardUser.IsBanned,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrFetchingUserFailed, err)
	}

	return &user, nil
}

func (r *userRepo) FetchUserByUsername(ctx context.Context, username string) (*models.StandardUser, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "id, username, password, name, email, role, last_seen, organisation, country, leetcode_id, is_banned",
		"table":      "Users",
		"conditions": "username = $1",
	})

	row := db.QueryRowContext(ctx, query, username)

	var user models.StandardUser
	err = row.Scan(
		&user.StandardUser.ID,
		&user.StandardUser.Username,
		&user.StandardUser.Password,
		&user.StandardUser.Name,
		&user.StandardUser.Email,
		&user.StandardUser.Role,
		&user.LastSeen,
		&user.StandardUser.Organisation,
		&user.StandardUser.Country,
		&user.LeetcodeID,
		&user.StandardUser.IsBanned,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrFetchingUserFailed, err)
	}

	return &user, nil
}

func (r *userRepo) FetchUserProgress(ctx context.Context, userID string) (*[]string, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "title_slugs",
		"table":      "users_progress",
		"conditions": "user_id = $1",
	})

	row := db.QueryRowContext(ctx, query, userID)

	var titleSlugs pq.StringArray
	err = row.Scan(&titleSlugs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrFetchingUserFailed, err)
	}

	titleSlugList := []string(titleSlugs)
	return &titleSlugList, nil
}

func (r *userRepo) UpdateUserDetails(ctx context.Context, user *models.StandardUser) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	if user.StandardUser.ID == "" {
		return fmt.Errorf("%w: user ID is required", errs.ErrUserNotFound)
	}

	query := queries.QueryBuilder(queries.BaseUpdate, map[string]string{
		"table":       "Users",
		"assignments": "username = $1, email = $2, password = $3, name = $4, organisation = $5, country = $6, leetcode_id = $7, last_seen = $8",
		"conditions":  "id = $9",
	})

	_, err = db.ExecContext(
		ctx,
		query,
		user.StandardUser.Username,
		user.StandardUser.Email,
		user.StandardUser.Password,
		user.StandardUser.Name,
		user.StandardUser.Organisation,
		user.StandardUser.Country,
		user.LeetcodeID,
		user.LastSeen,
		user.StandardUser.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUpdatingUserDetailsFailed, err)
	}

	return nil
}

func (r *userRepo) BanUser(ctx context.Context, userID string) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	if userID == "" {
		return fmt.Errorf("%w: user ID is required", errs.ErrUserNotFound)
	}

	query := queries.QueryBuilder(queries.BaseUpdate, map[string]string{
		"table":       "Users",
		"assignments": "is_banned = TRUE",
		"conditions":  "id = $1 and role = 'user'",
	})

	result, err := db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrBanningUserFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: user with ID %s not found", errs.ErrUserNotFound, userID)
	}

	return nil
}

func (r *userRepo) UnbanUser(ctx context.Context, userID string) error {
	db, err := r.getDBConnection()
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	if userID == "" {
		return fmt.Errorf("%w: user ID is required", errs.ErrUserNotFound)
	}

	query := queries.QueryBuilder(queries.BaseUpdate, map[string]string{
		"table":       "Users",
		"assignments": "is_banned = FALSE",
		"conditions":  "id = $1 and role = 'user'",
	})

	result, err := db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUnbanningUserFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: user with ID %s not found", errs.ErrUserNotFound, userID)
	}

	return nil
}

func (r *userRepo) CountActiveUsersInLast24Hours(ctx context.Context) (int, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	now := time.Now().UTC()
	twentyFourHoursAgo := now.Add(-24 * time.Hour)

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "COUNT(*)",
		"table":      "Users",
		"conditions": "last_seen >= $1",
	})

	row := db.QueryRowContext(ctx, query, twentyFourHoursAgo)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
	}

	return count, nil
}

func (r *userRepo) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "COUNT(*)",
		"table":      "Users",
		"conditions": "email = $1",
	})

	row := db.QueryRowContext(ctx, query, email)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
	}

	return count == 0, nil
}

func (r *userRepo) IsUsernameUnique(ctx context.Context, username string) (bool, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "COUNT(*)",
		"table":      "Users",
		"conditions": "username = $1",
	})

	row := db.QueryRowContext(ctx, query, username)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
	}

	return count == 0, nil
}

func (r *userRepo) IsLeetcodeIDUnique(ctx context.Context, LeetcodeID string) (bool, error) {
	db, err := r.getDBConnection()
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrDatabaseConnection, err)
	}

	query := queries.QueryBuilder(queries.BaseSelectWhere, map[string]string{
		"columns":    "COUNT(*)",
		"table":      "Users",
		"conditions": "leetcode_id = $1",
	})

	row := db.QueryRowContext(ctx, query, LeetcodeID)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrFetchingUsersFailed, err)
	}

	return count == 0, nil
}
