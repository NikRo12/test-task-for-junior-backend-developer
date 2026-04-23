package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type recurrenceDTO struct {
	Type       taskdomain.RecurrenceType `json:"type"`
	EveryNDays int                       `json:"every_n_days,omitempty"`
	StartDate  *taskdomain.Date          `json:"start_date,omitempty"`
	DayOfMonth int                       `json:"day_of_month,omitempty"`
	Dates      []taskdomain.Date         `json:"dates,omitempty"`
	Parity     taskdomain.Parity         `json:"parity,omitempty"`
}

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *recurrenceDTO    `json:"recurrence"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *recurrenceDTO    `json:"recurrence,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.Recurrence != nil {
		dto.Recurrence = domainRecurrenceToDTO(task.Recurrence)
	}

	return dto
}

func domainRecurrenceToDTO(r *taskdomain.Recurrence) *recurrenceDTO {
	return &recurrenceDTO{
		Type:       r.Type,
		EveryNDays: r.EveryNDays,
		StartDate:  r.StartDate,
		DayOfMonth: r.DayOfMonth,
		Dates:      r.Dates,
		Parity:     r.Parity,
	}
}

func dtoRecurrenceToDomain(dto *recurrenceDTO) *taskdomain.Recurrence {
	if dto == nil {
		return nil
	}
	return &taskdomain.Recurrence{
		Type:       dto.Type,
		EveryNDays: dto.EveryNDays,
		StartDate:  dto.StartDate,
		DayOfMonth: dto.DayOfMonth,
		Dates:      dto.Dates,
		Parity:     dto.Parity,
	}
}
