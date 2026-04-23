package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	now := s.now()

	rec := normalized.Recurrence
	if rec != nil && rec.Type == taskdomain.RecurrenceDaily && rec.StartDate == nil {
		today := taskdomain.Date{Time: now.Truncate(24 * time.Hour)}
		rec.StartDate = &today
	}

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Recurrence:  rec,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	rec := normalized.Recurrence
	if rec != nil && rec.Type == taskdomain.RecurrenceDaily && rec.StartDate == nil {
		today := taskdomain.Date{Time: s.now().Truncate(24 * time.Hour)}
		rec.StartDate = &today
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Recurrence:  rec,
		UpdatedAt:   s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) ListByDate(ctx context.Context, date time.Time) ([]taskdomain.Task, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]taskdomain.Task, 0)
	for _, t := range all {
		if t.Recurrence != nil && t.Recurrence.IsScheduledFor(date) {
			result = append(result, t)
		}
	}

	return result, nil
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if err := validateRecurrence(input.Recurrence); err != nil {
		return CreateInput{}, err
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if err := validateRecurrence(input.Recurrence); err != nil {
		return UpdateInput{}, err
	}

	return input, nil
}

func validateRecurrence(r *taskdomain.Recurrence) error {
	if r == nil {
		return nil
	}

	switch r.Type {
	case taskdomain.RecurrenceDaily:
		if r.EveryNDays <= 0 {
			return fmt.Errorf("%w: recurrence.every_n_days must be >= 1 for type 'daily'", ErrInvalidInput)
		}

	case taskdomain.RecurrenceMonthly:
		if r.DayOfMonth < 1 || r.DayOfMonth > 30 {
			return fmt.Errorf("%w: recurrence.day_of_month must be between 1 and 30", ErrInvalidInput)
		}

	case taskdomain.RecurrenceSpecificDates:
		if len(r.Dates) == 0 {
			return fmt.Errorf("%w: recurrence.dates must not be empty for type 'specific_dates'", ErrInvalidInput)
		}

	case taskdomain.RecurrenceEvenOdd:
		if r.Parity != taskdomain.ParityEven && r.Parity != taskdomain.ParityOdd {
			return fmt.Errorf("%w: recurrence.parity must be 'even' or 'odd'", ErrInvalidInput)
		}

	default:
		return fmt.Errorf("%w: unknown recurrence type %q", ErrInvalidInput, r.Type)
	}

	return nil
}
