package sqlite

import (
	"context"

	"github.com/0xdod/fileserve"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type fileService struct {
	db *DB
}

func NewFileService(db *DB) *fileService {
	return &fileService{
		db: db,
	}
}

func (fs *fileService) CreateFile(ctx context.Context, file *fileserve.File) error {
	file.ID = uuid.NewString()
	file.CreatedAt = fs.db.Now()
	file.UpdatedAt = fs.db.Now()

	_, err := fs.db.db.NamedExecContext(ctx, `
		INSERT INTO files (id, name, url, created_at, updated_at)
		VALUES (:id, :name, :url, :created_at, :updated_at)
	`, file)

	return err
}

func (fs *fileService) GetFiles(ctx context.Context, param fileserve.GetFilesParam) ([]*fileserve.File, error) {
	query, args, err := sqlx.Named(`
		SELECT * FROM files
		WHERE (:name IS NULL OR name = :name)
		LIMIT :limit OFFSET :offset
	`, param)

	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)

	if err != nil {
		return nil, err
	}

	query = fs.db.db.Rebind(query)

	var files []*fileserve.File

	err = fs.db.db.SelectContext(ctx, &files, query, args...)

	return files, err
}

func (fs *fileService) GetFile(ctx context.Context, id string) (*fileserve.File, error) {
	file := &fileserve.File{}

	err := fs.db.db.GetContext(ctx, file, `
		SELECT * FROM files
		WHERE id = ?
	`, id)

	return file, err
}

func (fs *fileService) UpdateFile(ctx context.Context, id string, param fileserve.UpdateFileParam) error {
	query, args, err := sqlx.Named(`
		UPDATE files
		SET name = :name
		WHERE id = :id
	`, struct {
		ID   string
		Name *string
	}{
		ID:   id,
		Name: param.Name,
	})

	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)

	if err != nil {
		return err
	}

	query = fs.db.db.Rebind(query)

	_, err = fs.db.db.ExecContext(ctx, query, args...)

	return err
}

func (fs *fileService) DeleteFile(ctx context.Context, id string) error {
	_, err := fs.db.db.ExecContext(ctx, `
		DELETE FROM files
		WHERE id = ?
	`, id)

	return err
}
