// Package app contains internal services used in other packages
package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
)

type IpfsService struct {
	ctx  context.Context
	ipfs icore.CoreAPI
}

func NewIpfsService() (context.CancelFunc, *IpfsService, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ipfs, err := spawnDefault(ctx)
	if err != nil {
		cancel()
		return nil, nil, errors.Wrap(err, "failed to spawn default node")
	}
	return cancel, &IpfsService{ctx: ctx, ipfs: ipfs}, nil
}

func (s *IpfsService) AddFile(filepath string) (string, error) {
	fnode, err := getUnixfsNode(filepath)
	if err != nil {
		return "", errors.Wrap(err, "failed to get unixfs node for file")
	}
	cid, err := s.ipfs.Unixfs().Add(s.ctx, fnode)
	if err != nil {
		return "", errors.Wrap(err, "failed to add files using unixfs")
	}
	return cid.String(), nil
}

func (s *IpfsService) GetFile(cidStr string) error {
	cid := icorepath.New(cidStr)
	rootNode, err := s.ipfs.Unixfs().Get(s.ctx, cid)
	if err != nil {
		return errors.Wrapf(err, "failed to get file %s from unixfs\n", cidStr)
	}
	outputPath := path.Join(".", cidStr)
	err = files.WriteTo(rootNode, outputPath)
	if err != nil {
		return errors.Wrap(err, "failed to write file to filesystem")
	}
	return nil
}

func (s *IpfsService) Connect(addrStr string) error {
	maddr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		return errors.Wrap(err, "failed to create multiaddr from address string")
	}
	peerAddr, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return errors.Wrap(err, "failed to create addinfo from multiaddr")
	}
	err = s.ipfs.Swarm().Connect(s.ctx, *peerAddr)
	if err != nil {
		return errors.Wrap(err, "failed to connect to peer")
	}
	return nil
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return errors.Wrap(err, "failed to load plugins")
	}

	if err := plugins.Initialize(); err != nil {
		return errors.Wrap(err, "failed to init plugins")
	}

	if err := plugins.Inject(); err != nil {
		return errors.Wrap(err, "failed to inject plugins")
	}

	return nil
}

func createNode(ctx context.Context, repoPath string) (icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open repo")
	}

	// Construct the node

	nodeOptions := &core.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start new node")
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}

func writeKey(key, dst string) error {
	if err := ioutil.WriteFile(dst, []byte(key), 0644); err != nil {
		return errors.Wrap(err, "failed to write to file")
	}
	return nil
}

func setupRepo(repoPath string) error {
	fmt.Printf("setting up new repo at %s\n", repoPath)

	if err := setupPlugins(repoPath); err != nil {
		return errors.Wrap(err, "failed to setup plugins")
	}

	if fsrepo.IsInitialized(repoPath) {
		fmt.Println("repo is already initialized. skipping...")
		return nil
	}

	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}
	cfg.SetBootstrapPeers(nil)

	if err = fsrepo.Init(repoPath, cfg); err != nil {
		return errors.Wrap(err, "failed to init repo")
	}

	// err = copy(path.Join(".", "swarm.key"), path.Join(repoPath, "swarm.key"))
	err = writeKey(swKey, path.Join(repoPath, "swarm.key"))
	if err != nil {
		return errors.Wrap(err, "failed to copy swarm.key file")
	}

	return nil
}

func getUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func spawnDefault(ctx context.Context) (icore.CoreAPI, error) {
	path, err := config.PathRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default path")
	}

	err = setupRepo(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to setup repo")
	}

	return createNode(ctx, path)
}
