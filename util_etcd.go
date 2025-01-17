package go_atomos

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// etcdKeepalive
// 保持活跃
// 1. 注册服务
// 2. 保持活跃
func etcdKeepalive(cli *clientv3.Client, key, value string, ttl int64) (*clientv3.LeaseGrantResponse, <-chan *clientv3.LeaseKeepAliveResponse, *Error) {
	// grant a lease.
	lease, er := cli.Grant(context.Background(), ttl)
	if er != nil {
		return nil, nil, NewErrorf(ErrCosmosEtcdKeepaliveFailed, "etcd: Failed to grant lease. err=(%s)", er).AddStack(nil)
	}

	// register the service with etcd.
	_, er = cli.Put(context.Background(), key, value, clientv3.WithLease(lease.ID))
	if er != nil {
		return nil, nil, NewErrorf(ErrCosmosEtcdKeepaliveFailed, "etcd: Failed to register service. err=(%s)", er).AddStack(nil)
	}

	// keepalive the service so that the lease does not expire.
	keepAliveCh, er := cli.KeepAlive(context.Background(), lease.ID)
	if er != nil {
		return nil, nil, NewErrorf(ErrCosmosEtcdKeepaliveFailed, "etcd: Failed to keepalive. err=(%s)", er).AddStack(nil)
	}

	return lease, keepAliveCh, nil
}

// etcdPut
// Put a string into etcd.
func etcdPut(cli *clientv3.Client, key string, value []byte) *Error {
	// Put the string into etcd.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, er := cli.Put(ctx, key, string(value))
	cancel()
	if er != nil {
		return NewErrorf(ErrCosmosEtcdPutFailed, "etcd: Put failed. key=(%s),err=(%v)", key, er).AddStack(nil)
	}
	return nil
}

// etcdGet
// Get a string from etcd.
func etcdGet(cli *clientv3.Client, key string) ([]byte, int64, *Error) {
	// Get the string back from etcd.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, er := cli.Get(ctx, key)
	cancel()
	if er != nil {
		return nil, 0, NewErrorf(ErrCosmosEtcdGetFailed, "etcd: Get failed. key=(%s),err=(%v)", key, er).AddStack(nil)
	}
	for _, kv := range resp.Kvs {
		if string(kv.Key) == key {
			return kv.Value, kv.Version, nil
		}
	}
	return nil, 0, nil
}

// etcdDelete
// Delete a string from etcd.
func etcdDelete(cli *clientv3.Client, key string) *Error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, er := cli.Delete(ctx, key)
	cancel()
	if er != nil {
		return NewErrorf(ErrCosmosEtcdDeleteFailed, "etcd: Delete failed. key=(%s),err=(%v)", key, er).AddStack(nil)
	}
	return nil
}
