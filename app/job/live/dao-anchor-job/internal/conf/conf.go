package conf

import (
	"errors"
	"flag"

	"go-common/library/database/bfs"
	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/rpc/warden"

	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log            *log.Config
	BM             *bm.ServerConfig
	Verify         *verify.Config
	Tracer         *trace.Config
	Ecode          *ecode.Config
	LiveRpc        map[string]*liverpc.ClientConfig
	GrpcCli        *warden.ClientConfig
	BfsCli         *bfs.Config
	HttpCli        *bm.ClientConfig
	CoverControl   *CoverControl
	Minute3Control *Minute3Control
	MinuteControl  *MinuteControl
}
type CoverControl struct {
	CoverCron string
	PieceSize int
}
type Minute3Control struct {
	Minute3Cron string
	PieceSize   int
}
type MinuteControl struct {
	MinuteCron string
	PieceSize  int
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
