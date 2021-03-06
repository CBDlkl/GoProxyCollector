package main

import (
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/scheduler"
	"github.com/AceDarkkinght/GoProxyCollector/server"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/verifier"

	"github.com/cihub/seelog"
)

func main() {
	// Load log.
	scheduler.SetLogger("logConfig.xml")
	defer seelog.Flush()

	// Load database.
	database, err := storage.NewBoltDbStorage("proxy.db", "IpList")
	if err != nil {
		seelog.Critical(err)
		panic(err)
	}

	// Sync data.
	database.SyncKeys()
	seelog.Infof("database initialize finish.")
	defer database.Close()

	// Start server
	go server.NewServer(database)

	// Sync DB every 5min.
	syncTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range syncTicker.C {
			verifier.VerifyAndDelete(database)
			database.SyncKeys()
			seelog.Debug("verify and sync database.")
		}
	}()

	configs := collector.NewCollectorConfig("collectorConfig.xml")
	scheduler.Run(configs, database)
}
