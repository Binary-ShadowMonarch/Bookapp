package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var ErrNotFound = errors.New("object not found")

type Namespace string

const (
	NamespaceUser   Namespace = "user"
	NamespaceGlobal Namespace = "global"
)

type Scope struct {
	Namespace Namespace
	UserID    int
}

func UserScope(userID int) Scope {
	return Scope{Namespace: NamespaceUser, UserID: userID}
}

func GlobalScope() Scope {
	return Scope{Namespace: NamespaceGlobal}
}

type ObjectInfo struct {
	Key          string
	Size         int64
	ContentType  string
	LastModified time.Time
}

type ObjectReader struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64
}

type Provider interface {
	Put(ctx context.Context, scope Scope, objectName string, reader io.Reader, size int64, contentType string) error
	Get(ctx context.Context, scope Scope, objectName string) (*ObjectReader, error)
	List(ctx context.Context, scope Scope) ([]ObjectInfo, error)
	Stat(ctx context.Context, scope Scope, objectName string) (*ObjectInfo, error)
	Delete(ctx context.Context, scope Scope, objectName string) error
	EnsureScope(ctx context.Context, scope Scope) error
}
