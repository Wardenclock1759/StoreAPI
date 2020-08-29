package storage

type Storage interface {
	User() UserRepository
	Role() RoleRepository
	Game() GameRepository
	Key() KeyRepository
	Payment() PaymentRepository
}
