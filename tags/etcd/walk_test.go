package etcd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/dubbikins/envy/v2/tag"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
)



func TestWalk(t *testing.T) {
	var client *clientv3.Client
	var err error
	var _ListenClientUrls = "localhost:42379"
	var ListenClientUrl *url.URL
	var _ListenPeerURL = "localhost:42380"
	var ListenPeerURL *url.URL
	if ListenClientUrl, err = url.Parse(fmt.Sprintf("http://%s", _ListenClientUrls)); err != nil {
		t.Fatal(err)
	}
	if ListenPeerURL, err = url.Parse(fmt.Sprintf("http://%s", _ListenPeerURL)); err != nil {
		t.Fatal(err)
	}
	cfg := embed.NewConfig()
	cfg.ListenPeerUrls = []url.URL{*ListenPeerURL}
	cfg.ListenClientUrls= []url.URL{*ListenClientUrl}
	cfg.Dir = t.TempDir()
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	<-e.Server.ReadyNotify()

	type EtdcExample struct {
		Field string `etcd:"/path/test;endpoints='localhost:42379'"`
		SkippedField string `etcd:"-"`
	}
	key :=  "/path/test"
	have := EtdcExample{}
	want := EtdcExample{
		Field: "foobar",
	}
	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{_ListenClientUrls},
		DialTimeout: 5 * time.Second,
	}); err != nil {
		return
	}
	var watcher = clientv3.NewWatcher(client)
	var watcherUpdates = watcher.Watch(t.Context(), key)
	
	defer client.Close()
	if _, err = client.Put(t.Context(), key, want.Field); err != nil {
		t.Fatal(err)
	}
	<-watcherUpdates
	if err := tag.WalkContext(t.Context(), tag.Chain(tag.UnmarshalText, WalkFn), &have); err != nil {
		t.Fatal(err)
	}
	if have != want {
		t.Fatalf("Expected %v, got %v", want, have)
	}
}



func TestWatchWalk(t *testing.T) {
	
	
	var client *clientv3.Client
	var err error
	var _ListenClientUrls = "localhost:42379"
	var ListenClientUrl *url.URL
	var _ListenPeerURL = "localhost:42380"
	var ListenPeerURL *url.URL
	if ListenClientUrl, err = url.Parse(fmt.Sprintf("http://%s", _ListenClientUrls)); err != nil {
		t.Fatal(err)
	}
	if ListenPeerURL, err = url.Parse(fmt.Sprintf("http://%s", _ListenPeerURL)); err != nil {
		t.Fatal(err)
	}
	type EtdcExample struct {
		Field string `etcd:"/test;endpoints='localhost:42379'"`
		SkippedField string `etcd:"-"`
	}
	cfg := embed.NewConfig()
	cfg.ListenPeerUrls = []url.URL{*ListenPeerURL}
	cfg.ListenClientUrls= []url.URL{*ListenClientUrl}
	cfg.Dir = t.TempDir()
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	<-e.Server.ReadyNotify()
	key :=  "/test"
	have := EtdcExample{}
	want := EtdcExample{
		Field: "foobar",
	}
	
	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{_ListenClientUrls},
		DialTimeout: 5 * time.Second,
	}); err != nil {
		return
	}
	defer client.Close()
	var watcher = clientv3.NewWatcher(client)
	var watcherUpdates = watcher.Watch(t.Context(), key)
	var callback = make(chan clientv3.WatchResponse, 2) //watcher.Watch(t.Context(), key)
	
	if _, err = client.Put(t.Context(), key, want.Field); err != nil {
		t.Fatal(err)
	}
	<-watcherUpdates
	if err := tag.WalkContext(context.WithValue(t.Context(), _updatesCtxKey, callback), tag.Chain(func (n *tag.Node) ( error) {return nil}, WatchWalkFn), &have); err != nil {
		t.Fatal(err)
	}
	if have != want {
		t.Fatalf("Expected %v, got %v", want, have)
	}

	want.Field = "foobaz"
	if _, err = client.Put(t.Context(), key, want.Field); err != nil {
		t.Fatal(err)
	}
	<-watcherUpdates
	<-callback
	if have != want {
		t.Fatalf("Expected %v, got %v", want, have)
	}
}


