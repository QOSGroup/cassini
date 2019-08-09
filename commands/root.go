package commands

import (
	"context"
	"os"
	"strings"

	"github.com/QOSGroup/cassini/common"
	"github.com/QOSGroup/cassini/config"
	"github.com/QOSGroup/cassini/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (

	// DefaultEventSubscribe events 默认订阅条件
	DefaultEventSubscribe string = "tm.event='Tx' AND qcp.to='qos'"
)

// Runner 通过配置数据执行方法，返回运行过程中出现的错误，如果返回空则代表运行成功。
type Runner func() (context.CancelFunc, error)

// NewRootCommand 创建 root/默认 命令
//
// 实现默认功能，显示帮助信息，预处理配置初始化，日志配置初始化。
func NewRootCommand(versioner Runner) *cobra.Command {
	root := &cobra.Command{
		Use:   CmdCassini,
		Short: ShortDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if viper.GetBool(CmdVersion) {
				versioner()
				return
			}
			cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			// binding flags

			if err = viper.BindPFlags(cmd.Flags()); err != nil {
				log.Error("bind flags error: ", err)
				return err
			}
			if strings.EqualFold(cmd.Use, CmdCassini) ||
				strings.EqualFold(cmd.Use, CmdVersion) ||
				strings.HasPrefix(cmd.Use, CmdHelp) ||
				strings.HasPrefix(cmd.Use, CmdTx) {
				// doesn't need init log and config
				return nil
			}

			// init & binding config

			err = initConfig()
			if err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					log.Warn(err.Error())
					// create & write default config

				} else {
					log.Error("Load config error: ", err.Error())
				}
				return
			}

			// init logger

			initLogger()

			return
		},
	}

	root.Flags().BoolP(CmdVersion, "v", false, "Show version info")

	return root
}

func initConfig() error {
	// init config

	log.Debug("home: ", viper.GetString(FlagHome))
	viper.Set(FlagHome, viper.GetString(FlagHome))

	// // Sets name for the config file.
	// // Does not include extension.
	// viper.SetConfigName("config")
	// // Adds a path for Viper to search for the config file in.
	// viper.AddConfigPath(filepath.Join(homeDir, "config"))
	// // Can be called multiple times to define multiple search paths.
	// // viper.AddConfigPath(homeDir)

	log.Debug("Init config: ", viper.GetString(FlagConfig))
	viper.SetConfigFile(viper.GetString(FlagConfig))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func initLogger() {
	log.Debug("Init log: ", viper.GetString(FlagLog))
	logger, err := log.LoadLogger(viper.GetString(FlagLog))
	if err != nil {
		log.Warn("Used the default logger because error: ", err)
	} else {
		log.Replace(logger)
	}
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run()
	if err != nil {
		log.Error("Run command error: ", err.Error())
		return err
	}
	if isKeepRunning {
		common.KeepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}

func reconfigMock(node string) (mock *config.MockConfig) {
	conf := config.GetConfig()
	if len(conf.Mocks) < 1 {
		mock = &config.MockConfig{
			RPC: &config.RPCConfig{
				NodeAddress: node}}
		conf.Mocks = []*config.MockConfig{mock}
	}
	if mock == nil {
		conf.Mocks = conf.Mocks[:1]
		mock = conf.Mocks[0]
		mock.RPC.NodeAddress = node
	}
	return
}
