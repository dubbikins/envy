package etcd

import (
	"log/slog"
	"strings"
	"time"

	"github.com/dubbikins/envy/v2/tag"
	"github.com/dubbikins/envy/v2/tag/text"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// type etcdOptions struct {
// 	Endpoints []string `default:"localhost:2379,localhost:22379,localhost:32379"`
// 	DialTimeout time.Duration `default:"5s"`
// }


func WalkFn( next tag.WalkFn) tag.WalkFn {
	return func(node *tag.Node) (err error) {
		var client *clientv3.Client
		//If this isn't a struct field element, then we should stop walking the tag
		if node.Field() == nil { //&& node.Value().Kind() != reflect.Pointer
			return 
		}
		if err = text.Parse("etcd", node, LexEctdTag); err != nil || node.Skipped(){
			return 
		}
		var endpoints []string 
		if _endpoints, set := node.Option("endpoints"); !set {
			endpoints= []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		}else {
			endpoints = strings.Split(_endpoints, ",")
		}
		slog.Info("Connection to ETCD", "endpoints", endpoints)
		if client, err = clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}); err != nil {
			return
		}
		defer client.Close()
			
		for _, key := range node.TagValues() {
			var resp *clientv3.GetResponse 
			if resp, err = client.Get(node.Context(), key); err != nil {
				return
			}
			for _,kv := range resp.Kvs {
				node.Write(kv.Value)
			}
		}
		return next(node)
	}
}
type updatesCtxKey struct {}
var _updatesCtxKey = updatesCtxKey{}

func WatchWalkFn( next tag.WalkFn) tag.WalkFn {
	return func( node *tag.Node) (err error) {
		var client *clientv3.Client
		var postUpdates chan clientv3.WatchResponse
		var sendPostUpdates bool
		if postUpdates, sendPostUpdates = node.Context().Value(_updatesCtxKey).(chan clientv3.WatchResponse); sendPostUpdates {
			slog.Info("sending updates to response channel")
		}
		//If this isn't a struct field element, then we should stop walking the tag
		if node.Field() == nil { //&& node.Value().Kind() != reflect.Pointer
			return 
		}
		if err = text.Parse("etcd", node, LexEctdTag); err != nil || node.Skipped(){
			return 
		}
		var endpoints []string 
		if _endpoints, set := node.Option("endpoints"); !set {
			endpoints= []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		}else {
			endpoints = strings.Split(_endpoints, ",")
		}
		slog.Info("Connection to ETCD", "endpoints", endpoints)
		if client, err = clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		}); err != nil {
			return
		}
	
		for _, key := range node.TagValues() {
			var resp *clientv3.GetResponse 
			if resp, err = client.Get(node.Context(), key); err != nil {
				return
			}
			var watcher =  clientv3.NewWatcher(client)		
			var key_updates = watcher.Watch(node.Context(), key)

			go func() {
				var update clientv3.WatchResponse
				// var ok bool
				defer client.Close()
				defer watcher.Close()
				for {
					select {
					case <-node.Context().Done():
						return
					case update = <- key_updates:
						
							if resp, err = client.Get(node.Context(), key); err != nil {
								return
							}
							for _,kv := range resp.Kvs {
								slog.Info("ETCD Key Change", "update", update)
								if  err = node.UnmarshalText(kv.Value); err != nil {
									return
								}
							}
							if sendPostUpdates {
								slog.Info("Sending ETCD Key Change Callback", "update", update)

								postUpdates <- update
								slog.Info("ETCD Key Change Callback sent", "update", update)

							}
					}
				}
			}()
			for _,kv := range resp.Kvs {
				if  err = node.UnmarshalText(kv.Value); err != nil {
					return
				}
			}
			

		}
		return next(node)
	}
}