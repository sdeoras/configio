package configfile

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/sdeoras/configio/simpleconfig"
	"github.com/sirupsen/logrus"
)

func TestNewWriter(t *testing.T) {
	writer, err := NewWriter(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()

	if err := writer.Marshal(new(simpleconfig.Config).Rand()); err != nil {
		t.Fatal(err)
	}
}

func TestNewReader(t *testing.T) {
	config := new(simpleconfig.Config).Rand()
	config2 := new(simpleconfig.Config).Init(config.Key())

	reader, err := NewReader(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()

	writer, err := NewWriter(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer writer.Close()

	if err := writer.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if err := reader.Unmarshal(config2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, config2) {
		t.Fatal("read config not same as the one set previously")
	}
}

func TestNewReadWriter(t *testing.T) {
	config := new(simpleconfig.Config).Rand()
	config2 := new(simpleconfig.Config).Init(config.Key())

	readWriter, err := NewReadWriter(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer readWriter.Close()

	if err := readWriter.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if err := readWriter.Unmarshal(config2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, config2) {
		t.Fatal("read config not same as the one set previously")
	}
}

func TestNewWatcher(t *testing.T) {
	a, b, c := callbackFunc("a"), callbackFunc("b"), callbackFunc("c")
	config := new(simpleconfig.Config).Rand()

	manager, err := NewManager(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	// register watch functions
	ca, cb, cc := manager.Watch("a", nil, a), manager.Watch("b", nil, b), manager.Watch("c", nil, c)

	// trigger config change
	if err := manager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	// make sure watch notifications are received
	select {
	case <-ca:
	case <-time.After(time.Second):
		t.Fatal("did not receive watch in 1 second")
	}

	select {
	case <-cb:
	case <-time.After(time.Second):
		t.Fatal("did not receive watch in 1 second")
	}

	select {
	case <-cc:
	case <-time.After(time.Second):
		t.Fatal("did not receive watch in 1 second")
	}

	time.Sleep(time.Second)
}

func callbackFunc(name string) func(ctx context.Context, data interface{}, err error) <-chan error {
	return func(ctx context.Context, data interface{}, err error) <-chan error {
		status := make(chan error)
		go func() {
			select {
			case <-ctx.Done():
				logrus.WithField("callback", name).Info("received context status")
			case status <- nil:
				logrus.WithField("callback", name).Info("sending status")
			}
		}()
		return status
	}
}
