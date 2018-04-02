package configio

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
)

func TestNewConfigWriter(t *testing.T) {
	if err := NewConfigFileWriter(context.Background()).Set(new(Config).Rand()); err != nil {
		t.Fatal(err)
	}
}

func TestNewConfigReader(t *testing.T) {
	config := new(Config).Rand()
	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(configFile, b, 0666); err != nil {
		t.Fatal(err)
	}

	config2, err := NewConfigFileReader(context.Background()).Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, config2) {
		t.Fatal("read config not same as the one set previously")
	}
}

func TestNewConfigReadWriter(t *testing.T) {
	config := new(Config).Rand()
	readWriter := NewConfigFileReadWriter(context.Background())
	if err := readWriter.Set(config); err != nil {
		t.Fatal(err)
	}

	config2, err := readWriter.Get()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(config, config2) {
		t.Fatal("read config not same as the one set previously")
	}
}

func TestNewConfigWatcher(t *testing.T) {
	a, b, c := callbackFunc("a"), callbackFunc("b"), callbackFunc("c")
	config := new(Config).Rand()
	manager := NewConfigFileManager(context.Background())
	ca, cb, cc := manager.Watch("a", nil, a), manager.Watch("b", nil, b), manager.Watch("c", nil, c)

	if err := manager.Set(config); err != nil {
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
