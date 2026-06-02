package storage

import (
	"fmt"
	"strings"
)

func ValidateObjectName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return fmt.Errorf("object name is required")
	}
	if trimmed != name {
		return fmt.Errorf("object name must not include leading or trailing spaces")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("object name must not include path separators")
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("object name must not include traversal sequences")
	}
	return nil
}

func ValidateScope(scope Scope) error {
	switch scope.Namespace {
	case NamespaceUser:
		if scope.UserID <= 0 {
			return fmt.Errorf("user scope requires a valid user id")
		}
	case NamespaceGlobal:
		return nil
	default:
		return fmt.Errorf("unknown namespace: %q", scope.Namespace)
	}
	return nil
}
