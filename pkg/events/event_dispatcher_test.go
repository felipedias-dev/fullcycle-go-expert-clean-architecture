package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) SetPayload(payload interface{}) {
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) HandleEvent(event EventInterface, wg *sync.WaitGroup) error {
	wg.Done()
	return nil
}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.event = TestEvent{Name: "test", Payload: "test"}
	suite.event2 = TestEvent{Name: "test2", Payload: "test2"}
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.eventDispatcher = NewEventDispatcher()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	assert.Equal(suite.T(), suite.eventDispatcher.handlers[suite.event.GetName()][0], &suite.handler)
	assert.Equal(suite.T(), suite.eventDispatcher.handlers[suite.event.GetName()][1], &suite.handler2)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_Duplicate() {
	err := suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler)
	suite.Equal(ErrorHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	err := suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RegisterHandler(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.ClearHandlers()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_HasHandler() {
	err := suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.HasHandler(suite.event.GetName(), &suite.handler))
	assert.True(suite.T(), suite.eventDispatcher.HasHandler(suite.event.GetName(), &suite.handler2))
	assert.False(suite.T(), suite.eventDispatcher.HasHandler(suite.event.GetName(), &suite.handler3))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) HandleEvent(event EventInterface, wg *sync.WaitGroup) error {
	m.Called(event)
	wg.Done()
	return nil
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eh := &MockHandler{}
	eh.On("HandleEvent", &suite.event)
	err := suite.eventDispatcher.RegisterHandler(suite.event.GetName(), eh)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	eh2 := &MockHandler{}
	eh2.On("HandleEvent", &suite.event)
	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), eh2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	eh3 := &MockHandler{}
	eh3.On("HandleEvent", &suite.event)
	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), eh3)
	suite.Nil(err)
	suite.Equal(3, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	eh4 := &MockHandler{}
	eh4.On("HandleEvent", &suite.event)
	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), eh4)
	suite.Nil(err)
	suite.Equal(4, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.DispatchEvent(&suite.event)
	suite.Nil(err)
	eh.AssertExpectations(suite.T())
	eh.AssertNumberOfCalls(suite.T(), "HandleEvent", 1)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	err := suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RegisterHandler(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.RemoveHandler(&suite.event, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
	assert.Equal(suite.T(), suite.eventDispatcher.handlers[suite.event.GetName()][0], &suite.handler2)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
