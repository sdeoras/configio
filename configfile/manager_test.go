package configfile

import (
	"context"
	"encoding/json"
	"io/ioutil"
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

	if err := writer.Marshal(new(simpleconfig.Config).Rand()); err != nil {
		t.Fatal(err)
	}
}

func TestNewReader(t *testing.T) {
	config := new(simpleconfig.Config).Rand()
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(DefaultConfigFile, b, 0666); err != nil {
		t.Fatal(err)
	}

	config2 := new(simpleconfig.Config)
	reader, err := NewReader(context.Background())
	if err != nil {
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

	readWriter, err := NewReadWriter(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if err := readWriter.Marshal(config); err != nil {
		t.Fatal(err)
	}

	config2 := new(simpleconfig.Config)
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
	default:
		<-ca
		<-cb
		<-cc
		time.Sleep(time.Second)
	case <-time.After(time.Second):
		t.Fatal("did not receive watch in 1 second")
	}
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
