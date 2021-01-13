package nodeprovider

// NodeOpts are options for intializing the node the cli uses
type NodeOpts struct {
	// DataDir sets the path of dir used to store node data
	DataDir string
	// ListenAddress sets the HOST:PORT on which the node will listen
	ListenAddress    string
	ContainerRuntime string
	DisableNode      bool
}
