//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/internal/entity"
	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/internal/event"
	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/internal/infra/database"
	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/internal/infra/web"
	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/internal/usecase"
	"github.com/felipedias-dev/fullcycle-go-expert-clean-architecture/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEventDependency = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEventDependency,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEventDependency,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
