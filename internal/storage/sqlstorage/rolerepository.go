package sqlstorage

import (
	"github.com/Wardenclock1759/StoreAPI/internal/model"
	"github.com/google/uuid"
)

type RoleRepository struct {
	storage *Storage
}

func (r *RoleRepository) GrantRole(role *model.UserRole) error {
	r.storage.db.QueryRow(
		"INSERT INTO \"user_roles\" (user_id, role) VALUES ($1, $2)",
		role.ID,
		role.Role,
	)

	return nil
}

func (r *RoleRepository) RevokeRole(role *model.UserRole) error {
	r.storage.db.QueryRow(
		"DELETE FROM \"user_roles\" WHERE user_id = $1 AND role = $2",
		role.ID,
		role.Role,
	)

	return nil
}

func (r *RoleRepository) GetRolesByID(u uuid.UUID) (*model.UserRole, error) {
	roles := &model.UserRole{}
	str := u.String()
	err := r.storage.db.QueryRow(
		"SELECT user_id, role FROM \"user_roles\" WHERE user_id = $1",
		str,
	).Scan(&roles.ID, &roles.Role)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
