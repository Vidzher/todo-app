package domain

import (
	"fmt"
	"time"
	core_errors "todo-api/internal/core/errors"
)

type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

func NewTask(
	id, version, authorID int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorID,
	}
}

func NewTaskUninitialized(
	title string,
	description *string,
	authorId int,
) Task {
	return NewTask(
		UninitializedID,
		UninitializedVersion,
		authorId,
		title,
		description,
		false,
		time.Now(),
		nil,
	)
}

func (t *Task) Validate() error {
	titleLength := len([]rune(t.Title))
	if titleLength < 1 || titleLength > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLength,
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLength := len([]rune(*t.Description))
		if descriptionLength < 1 || descriptionLength > 1_000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if t.AuthorUserID <= 0 {
		return fmt.Errorf(
			"invalid `AuthorUserID`: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if !t.Completed && t.CompletedAt != nil {
		return fmt.Errorf(
			"invalid `CompletedAt`: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompletedAt` cannot be nil if `Completed`==`true`: %d",
				core_errors.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` cannot be before `CreatedAt`: %d",
				core_errors.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` must be nil if `Completed`=`false`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type TaskPatch struct {
	Title       *string
	Description *string
	Completed   *bool
}

func (t *Task) ApplyPatch(patch TaskPatch) error {
	tmp := *t

	if patch.Title != nil {
		tmp.Title = *patch.Title
	}

	if patch.Description != nil {
		tmp.Description = patch.Description
	}

	if patch.Completed != nil {
		tmp.Completed = *patch.Completed

		if *patch.Completed {
			now := time.Now()
			tmp.CompletedAt = &now
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp

	return nil
}
