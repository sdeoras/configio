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
	if err := NewWriter(context.Background()).Marshal(new(simpleconfig.Config).Rand()); err != nil {
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
	if err := NewReader(context.Background()).Unmarshal(config2); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, config2) {
		t.Fatal("read config not same as the one set previously")
	}
}

func TestNewReadWriter(t *testing.T) {
	config := new(simpleconfig.Config).Rand()
	readWriter := NewReadWriter(context.Background())
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
	manager := NewManager(context.Background())
	ca, cb, cc := manager.Watch("a", nil, a), manager.Watch("b", nil, b), manager.Watch("c", nil, c)

	if err := manager.Marshal(config); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, <-ca) {
		t.Fatal("set and received config not matching")
	}

	if !reflect.DeepEqual(config, <-cb) {
		t.Fatal("set and received config not matching")
	}

	if !reflect.DeepEqual(config, <-cc) {
		t.Fatal("set and received config not matching")
	}

	time.Sleep(time.Second)
}

func callbackFunc(name string) func(ctx context.Context, data interface{}, err error) <-chan error {
	return func(ctx context.Context, data interface{}, err error) <-chan error {
		done := make(chan error)
		go func() {
			select {
			case <-ctx.Done():
				logrus.Info("context done ", name)
			case done <- nil:
				logrus.Info("executed ", name)
			}
		}()
		return done
	}
}
