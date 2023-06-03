package _interface

type DtoInterface[T any] interface {
	GetMany(filter map[string]interface{}, page uint, pageSize uint) ([]T, error)
	GetSingle(id string) (T, error)
	Create(data T) (T, error)
	Update(id string, data T) (T, error)
	Delete(id string) (T, error)
}
