### Config file about ETCD

Config sample:

```
{
    ......

    "embedEtcd":true,
    "etcd":{
        "name": "testA",
        "advertise":"http://127.0.0.1:2379",
        "advertisePeer":"http://127.0.0.1:2380",
        "clusterToken":"test-cassini-cluster",
        "cluster":"testA=http://127.0.0.1:2380"
    },
    "useEtcd":true,
    "lock":"etcd://127.0.0.1:2379"
    
    ......
}
```

The related configuration of distributed locks is divided into two parts:
lock client and embedded etcd server.

Clients are used to invoke distributed locks to implement code logic.

Embedded server is designed to simplify the complexity of a standalone deployment ETCD cluster.

### Lock client

useEtcd: 

> Whether to use etcd lock or not,  
> Default: false  
> \- Do not use etcd distributed lock

lock:

> Config the lock client,  
> Default: ""  
> \- Do not use etcd lock.  
> Example: "etcd://192.168.1.100:2379,192.168.1.101:2379,192.168.1.102:2379"  
> \- Access to the ETCD cluster through the specified three addresses and ports.

lockTTL:

> Timeout for lock,  
> Default: 5  
> \- When a session is lost to the etcd, the lock is automatically unlocked after 5 seconds.

### Embedded etcd server

name:

> Human-readable name for this member.  
> Default: 'default'  
> ETCD server config: --name 'default'

advertise:

> List of this member's client URLs to advertise to the public.
> This config is use to provide ETCD service.  
> The client URLs advertised should be accessible to machines  
>         that talk to etcd cluster. etcd client libraries parse  
>         these URLs to connect to the cluster.  
> Default: ""  
> \- advertise must be set.  
> Example: "http://192.168.1.100:2379"  
> \- Embedded etcd server will publish the address to the client.  
> ETCD server config: --advertise-client-urls http://192.168.1.100:2379

listen:

> List of URLs to listen on for client traffic.
> This config is use to provide ETCD service.  
> Default: Using advertise's setting.  
> Example: "http://0.0.0.0:2379"  
> \- Embedded etcd server will listen on the address.  
> ETCD server config: --listen-client-urls 'http://0.0.0.0:2379'

advertisePeer:

> List of this member's peer URLs to advertise to the rest of the cluster.
> This config is use to build ETCD cluster network.  
> Default: ""  
> \- advertisePeer must be set.  
> Example: "http://192.168.1.100:2380"  
> \- Embedded etcd server will publish the address to the other peers.  
> ETCD server config: --initial-advertise-peer-urls http://localhost:2380
	
listenPeer：

> List of URLs to listen on for peer traffic.  
> Default: Using advertisePeer's setting  
> Example: "http://192.168.1.100:2380"  
> \- Embedded etcd server will listen on the address for the other peer's connections.  
> ETCD server config: --listen-peer-urls 'http://192.168.1.100:2380'

clusterToken:

> Initial cluster token for the etcd cluster during bootstrap.  
> Specifying this can protect you from unintended cross-cluster  
>         interaction when running multiple clusters.  
> Default: "etcd-cluster"  
> \- Peers in the same cluster must be configured with the same cluster token.  
> ETCD server config: --initial-cluster-token etcd-cluster

cluster:

> Initial cluster configuration for bootstrapping.  
> Default: ""  
> \- cluster must be set.  
> Example: "devA=http://192.168.1.100:2380,devB=http://192.168.1.101:2380,devC=http://192.168.1.102:2380"  
> \- "devA" is setting in name.  
> ETCD server config: --initial-cluster 'default=http://192.168.1.100:2380'
