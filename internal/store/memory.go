package store

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain"
)

var ErrRouteUnavailable = errors.New("route capacity is unavailable")

type RouteStore interface {
	Load(context.Context) (domain.Plan, error)
	Reserve(context.Context, string, string) (func(bool), error)
	RouteLoads(context.Context) ([]domain.RouteLoad, error)
}

type MemoryStore struct {
	mu       sync.Mutex
	plan     domain.Plan
	used     map[string]int
	pending  map[string]int
	capacity map[string]int
}

func NewMemory(plan domain.Plan) *MemoryStore {
	capacity := make(map[string]int, len(plan.Routes))
	for _, route := range plan.Routes {
		capacity[route.ID] = route.Capacity
	}
	return &MemoryStore{
		plan:     domain.ClonePlan(plan),
		used:     make(map[string]int, len(plan.Routes)),
		pending:  make(map[string]int, len(plan.Routes)),
		capacity: capacity,
	}
}

func (s *MemoryStore) Load(ctx context.Context) (domain.Plan, error) {
	if err := ctx.Err(); err != nil {
		return domain.Plan{}, err
	}
	plan := domain.ClonePlan(s.plan)
	if err := ctx.Err(); err != nil {
		return domain.Plan{}, err
	}
	return plan, nil
}

func (s *MemoryStore) Reserve(ctx context.Context, routeID, shipmentID string) (func(bool), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s.used[routeID]+s.pending[routeID] >= s.capacity[routeID] {
		return nil, ErrRouteUnavailable
	}
	s.pending[routeID]++
	if err := ctx.Err(); err != nil {
		s.pending[routeID]--
		return nil, err
	}

	var once sync.Once
	return func(commit bool) {
		once.Do(func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			s.pending[routeID]--
			if commit {
				s.used[routeID]++
			}
		})
	}, nil
}

func (s *MemoryStore) RouteLoads(ctx context.Context) ([]domain.RouteLoad, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	loads := make([]domain.RouteLoad, 0, len(s.capacity))
	for routeID, capacity := range s.capacity {
		loads = append(loads, domain.RouteLoad{
			RouteID:  routeID,
			Used:     s.used[routeID],
			Capacity: capacity,
		})
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].RouteID < loads[j].RouteID })
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return loads, nil
}
