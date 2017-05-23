// Copyright © 2017 Meyer Zinn <meyerzinn@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-kit/kit/log"
	"github.com/ory/hydra/sdk"
	"github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"github.com/studiously/classsvc/ddl"
	"github.com/studiously/classsvc/service"
)

var addr string

// hostCmd represents the host command
var hostCmd = &cobra.Command{
	Use:   "host",
	Short: "Start the service.",
	Long: `Starts the service on all transports and connects to a database backend.

This command exposes several environmental variables for controls. You can set environments using "export KEY=VALUE" (Linux/macOS) or "set KEY=VALUE" (Windows). On Linux, you can also set environments by prepending key value pairs: "KEY=VALUE KEY2=VALUE2 classsvc".

All possible controls are listed below. The host process additionally exposes a few flags, which are listed below the controls section.

CORE CONTROLS
=============
- DATABASE_URL: A URL to a persistent backend. Classsvc supports PostgreSQL currently.

HYDRA CONTROLS
==============
A Hydra server is required to perform token introspection and thus authorization. Most endpoints (excepting health and unauthenticated ones) will fail without a valid Hydra server.

- HYDRA_CLIENT_ID: ID for Hydra client.
- HYDRA_CLIENT_SECRET: Secret for Hydra client.
- HYDRA_CLUSTER_URL: URL of Hydra cluster.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(logrus.StandardLogger().Out)
			logger = log.With(logger, "ts", log.DefaultTimestampUTC)
			logger = log.With(logger, "caller", log.DefaultCaller)
		}
		var s service.Service
		{
			// Set up database
			var driver = os.Getenv("DATABASE_DRIVER")
			var config = os.Getenv("DATABASE_CONFIG")

			db, err := sql.Open(driver, config)
			if err != nil {
				logrus.Fatalln("database connection failed", err)
			}
			if err := pingDatabase(db); err != nil {
				logrus.WithError(err).Fatalln("database ping attempts failed")
			}
			if err := setupDatabase(driver, db); err != nil {
				logrus.WithError(err).Fatalln("migration failed")
			}
			s = service.NewPostgres(db)
		}
		// Set up Hydra
		tlsVerify, err := strconv.ParseBool(os.Getenv("HYDRA_TLS_VERIFY"))
		if err != nil {
			tlsVerify = false
		}
		sdk.Connect()
		client, err := sdk.Connect(
			sdk.ClientID(os.Getenv("HYDRA_CLIENT_ID")),
			sdk.ClientSecret(os.Getenv("HYDRA_CLIENT_SECRET")),
			sdk.ClusterURL(os.Getenv("HYDRA_CLUSTER_URL")),
			sdk.SkipTLSVerify(tlsVerify),
		)
		if err != nil {
			logrus.WithError(err).Fatal("could not connect to hydra")
		}
		var h = service.MakeHTTPHandler(s, logger, client)
		errs := make(chan error)
		go func() {
			c := make(chan os.Signal, 100)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			errs <- fmt.Errorf("%s", <-c)
		}()

		go func(address string) {
			logger.Log("transport", "HTTP", "addr", addr)
			errs <- http.ListenAndServe(address, h)
		}(addr)

		logger.Log("exit", <-errs)
	},
}

func init() {
	RootCmd.AddCommand(hostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	hostCmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "HTTP bind address")
}

func setupDatabase(driver string, db *sql.DB) error {
	var migrations = &migrate.AssetMigrationSource{
		Asset:    ddl.Asset,
		AssetDir: ddl.AssetDir,
		Dir:      driver,
	}
	_, err := migrate.Exec(db, driver, migrations, migrate.Up)
	return err
}

func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 30; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		logrus.Infof("database ping failed. retry in 1s")
		time.Sleep(time.Second)
	}
	return
}
