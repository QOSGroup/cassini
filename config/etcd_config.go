package config

// EtcdConfig Embed-etcd config
type EtcdConfig struct {

	// Name - human-readable name for this member.
	//
	// --name 'default'
	Name string `yaml:"name,omitempty"`

	// Listen - list of URLs to listen on for client traffic.
	//
	// --listen-client-urls 'http://localhost:2379'
	Listen string `yaml:"listen,omitempty"`

	// Advertise list of this member's client URLs to advertise to the public.
	// The client URLs advertised should be accessible to machines that talk to etcd cluster. etcd client libraries parse these URLs to connect to the cluster.
	//
	// --advertise-client-urls http://localhost:2379
	Advertise string `yaml:"advertise,omitempty"`

	// ListenPeer list of URLs to listen on for peer traffic.
	//
	// --listen-peer-urls http://localhost:2380
	ListenPeer string `yaml:"listenPeer,omitempty"`

	// AdvertisePeer list of this member's peer URLs to advertise to the rest of the cluster.
	//
	// --initial-advertise-peer-urls http://localhost:2380
	AdvertisePeer string `yaml:"advertisePeer,omitempty"`

	// ClusterToken initial cluster token for the etcd cluster during bootstrap.
	// Specifying this can protect you from unintended cross-cluster interaction when running multiple clusters.
	//
	// --initial-cluster-token etcd-cluster
	ClusterToken string `yaml:"clusterToken,omitempty"`

	// Cluster initial cluster configuration for bootstrapping.
	//
	// --initial-cluster 'default=http://localhost:2380'
	Cluster string `yaml:"cluster,omitempty"`
}
