package service

import (
	"context"
	"errors"
	"github.com/ggfedortsov/tx-runner/internal/model"
)

type Runner interface {
	Run(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type UserStorage interface {
	CreateUser(ctx context.Context, u model.User) error
	GatAll(ctx context.Context) ([]model.User, error)
}

type Service struct {
	UserStorage UserStorage
	Runner      Runner
}

func (s *Service) MethodOk(ctx context.Context, u model.User) ([]model.User, error) {
	var res []model.User
	err := s.Runner.Run(ctx, func(ctx context.Context) error {
		if err := s.UserStorage.CreateUser(ctx, u); err != nil {
			return err
		}

		all, err := s.UserStorage.GatAll(ctx)
		if err != nil {
			return err
		}
		res = all

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) MethodError(ctx context.Context, u model.User) ([]model.User, error) {
	var res []model.User
	err := s.Runner.Run(ctx, func(ctx context.Context) error {
		if err := s.UserStorage.CreateUser(ctx, u); err != nil {
			return err
		}

		all, err := s.UserStorage.GatAll(ctx)
		if err != nil {
			return err
		}
		res = all

		if len(all) >= 2 {
			return errors.New("app error")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
