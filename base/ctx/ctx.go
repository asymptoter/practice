package ctx

import (
	"context"

	"go.uber.org/zap"
)

type M map[string]any

type CTX struct {
	context.Context
	Logger *zap.SugaredLogger
}

func Background() CTX {
	return CTX{
		Context: context.Background(),
		Logger:  zap.NewExample().Sugar(),
	}
}

func With(parent CTX, fields ...any) CTX {
	return CTX{
		Context: parent.Context,
		Logger:  parent.Logger.With(fields),
	}
}

func (c CTX) With(fields ...any) CTX {
	return CTX{
		Context: c.Context,
		Logger:  c.Logger.With(fields),
	}
}

func (c CTX) Debug(args ...any) {
	c.Logger.Debug(args)
}

func (c CTX) Info(args ...any) {
	c.Logger.Info(args)
}

func (c CTX) Warn(args ...any) {
	c.Logger.Warn(args)
}

func (c CTX) Error(args ...any) {
	c.Logger.Error(args)
}

func (c CTX) Fatal(args ...any) {
	c.Logger.Fatal(args)
}
