package service

import (
	"time"
	"uas-prestasi/app/repository"

	gocache "github.com/patrickmn/go-cache"
)

type PermissionService struct {
	Repo  *repository.PermissionRepository
	cache *gocache.Cache
}

func NewPermissionService(r *repository.PermissionRepository) *PermissionService {
    //cache
	c := gocache.New(5*time.Minute, 10*time.Minute)
	return &PermissionService{Repo: r, cache: c}
}

func (s *PermissionService) HasPermission(roleID, requiredPerm string) (bool, error) {
	// cek cache
	if x, found := s.cache.Get(roleID); found {
		perms := x.([]string)
		for _, p := range perms {
			if p == requiredPerm {
				return true, nil
			}
		}
		return false, nil
	}

	perms, err := s.Repo.GetPermissionsByRole(roleID)
	if err != nil {
		return false, err
	}

	s.cache.Set(roleID, perms, gocache.DefaultExpiration)

	for _, p := range perms {
		if p == requiredPerm {
			return true, nil
		}
	}
	return false, nil
}

// Optional: invalidate cache when role-permission berubah (optional aja ini teh)
// func (s *PermissionService) InvalidateRole(roleID string) {
//    s.cache.Delete(roleID)
// }
