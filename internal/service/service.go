package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ggfedortsov/tx-runner/internal/model"
)

type UserStorage interface {
	CreateUser(ctx context.Context, u model.User) error
	GatAll(ctx context.Context) ([]model.User, error)
	RunInTx(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type Service struct {
	UserStorage UserStorage
}

func (s *Service) DoubleRunner(ctx context.Context) error {
	err := s.UserStorage.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.UserStorage.CreateUser(ctx, model.User{"user1", 10}); err != nil {
			return err
		}

		s.UserStorage.RunInTx(ctx, func(ctx context.Context) error {
			all, _ := s.UserStorage.GatAll(ctx)
			fmt.Printf("all: %v \n", all)

			if err := s.UserStorage.CreateUser(ctx, model.User{"user2", 10}); err != nil {
				return err
			}

			return errors.New("app error")
		})

		return nil
	})
	return err
}

func (s *Service) MethodOk(ctx context.Context, u model.User) ([]model.User, error) {
	var res []model.User
	err := s.UserStorage.RunInTx(ctx, func(ctx context.Context) error {
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
	err := s.UserStorage.RunInTx(ctx, func(ctx context.Context) error {
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
