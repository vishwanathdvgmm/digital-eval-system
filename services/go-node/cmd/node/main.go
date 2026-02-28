package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"digital-eval-system/services/go-node/internal/admin"
	"digital-eval-system/services/go-node/internal/api"
	"digital-eval-system/services/go-node/internal/auth"
	"digital-eval-system/services/go-node/internal/authority"
	"digital-eval-system/services/go-node/internal/chain"
	"digital-eval-system/services/go-node/internal/core"
	"digital-eval-system/services/go-node/internal/db"
	"digital-eval-system/services/go-node/internal/evaluator"
	"digital-eval-system/services/go-node/internal/examiner"
	"digital-eval-system/services/go-node/internal/logger"
	"digital-eval-system/services/go-node/internal/pybridge"
	"digital-eval-system/services/go-node/internal/rootdir"
	"digital-eval-system/services/go-node/internal/storage"
	"digital-eval-system/services/go-node/internal/student"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
		TLS  struct {
			Enabled  bool   `yaml:"enabled"`
			CertPath string `yaml:"cert_path"`
			KeyPath  string `yaml:"key_path"`
		} `yaml:"tls"`
	} `yaml:"server"`
	Storage struct {
		BoltDBPath        string `yaml:"boltdb_path"`
		BoltDBTimeoutSecs int    `yaml:"boltdb_timeout_seconds"`
	} `yaml:"storage"`
	Block struct {
		SignerID      string `yaml:"signer_id"`
		SignatureAlgo string `yaml:"signature_algo"`
	} `yaml:"block"`
	PythonExtractor struct {
		URL string `yaml:"url"`
	} `yaml:"python_extractor"`
	PythonValidator struct {
		URL string `yaml:"url"`
	} `yaml:"python_validator"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
	Postgres struct {
		DSN                    string `yaml:"dsn"`
		MaxOpenConns           int    `yaml:"max_open_conns"`
		MaxIdleConns           int    `yaml:"max_idle_conns"`
		ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_seconds"`
	} `yaml:"postgres"`
	Auth struct {
		PrivKeyPath       string `yaml:"priv_key_path"`
		PubKeyPath        string `yaml:"pub_key_path"`
		Issuer            string `yaml:"issuer"`
		AccessTTLSeconds  int    `yaml:"access_ttl_seconds"`
		RefreshTTLSeconds int    `yaml:"refresh_ttl_seconds"`
	} `yaml:"auth"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// resolveConfigPaths converts any relative paths in the config to absolute
// paths anchored at the auto-detected project root.
func resolveConfigPaths(cfg *Config) {
	resolve := func(p string) string {
		if p == "" || filepath.IsAbs(p) {
			return p
		}
		return rootdir.Resolve(p)
	}
	cfg.Server.TLS.CertPath = resolve(cfg.Server.TLS.CertPath)
	cfg.Server.TLS.KeyPath = resolve(cfg.Server.TLS.KeyPath)
	cfg.Storage.BoltDBPath = resolve(cfg.Storage.BoltDBPath)
	cfg.Auth.PrivKeyPath = resolve(cfg.Auth.PrivKeyPath)
	cfg.Auth.PubKeyPath = resolve(cfg.Auth.PubKeyPath)
}

func setupLogger(levelStr string) {
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logger.PrettyFormatter{})
}

func main() {
	cfgPath := flag.String("config", "services/go-node/configs/config.yaml", "path to config yaml")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(2)
	}
	resolveConfigPaths(cfg)
	logrus.Infof("project root detected: %s", rootdir.Root())

	setupLogger(cfg.Logging.Level)
	logrus.Infof("starting go-node with signer=%s", cfg.Block.SignerID)

	// initialize BoltDB storage
	store, err := storage.NewBoltDB(cfg.Storage.BoltDBPath, time.Duration(cfg.Storage.BoltDBTimeoutSecs)*time.Second)
	if err != nil {
		logrus.Fatalf("failed to open boltdb: %v", err)
	}
	defer store.Close()

	// create core service registry
	registry := core.NewServiceRegistry()
	logrus.Info("service registry initialized")

	// python extractor client
	pyURL := cfg.PythonExtractor.URL
	if pyURL == "" {
		pyURL = "http://127.0.0.1:8081"
	}
	pyClient := pybridge.NewClient(pyURL, 300*time.Second)
	registry.Register("pybridge_extractor", pyClient)
	logrus.Infof("registered python extractor client at %s", pyURL)

	pyValidatorURL := cfg.PythonValidator.URL
	if pyValidatorURL == "" {
		pyValidatorURL = "http://127.0.0.1:8082"
	}
	pyValidatorClient := pybridge.NewClient(pyValidatorURL, 60*time.Second)
	registry.Register("pybridge_validator", pyValidatorClient)
	logrus.Infof("registered python validator client at %s", pyValidatorURL)

	// examiner upload service
	examSvc, err := examiner.NewService(chain.NewChain(store), store, pyClient, cfg.Block.SignerID, "")
	if err != nil {
		logrus.Fatalf("failed to create examiner upload service: %v", err)
	}
	registry.Register("examiner_upload_service", examSvc)
	logrus.Info("examiner upload service registered")

	// -----------------------------------------
	// Phase 5 – PostgreSQL Setup
	// -----------------------------------------
	pgCfg := db.PostgresConfig{
		DSN:             cfg.Postgres.DSN,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: time.Duration(cfg.Postgres.ConnMaxLifetimeSeconds) * time.Second,
	}

	pgDB, err := db.NewPostgres(pgCfg)
	if err != nil {
		logrus.Fatalf("failed to connect postgres: %v", err)
	}
	registry.Register("postgres", pgDB)
	logrus.Info("postgres connected")

	// load RSA keys into jwt manager
	jwtMgr, err := auth.NewManagerFromFiles(
		cfg.Auth.PrivKeyPath,
		cfg.Auth.PubKeyPath,
		cfg.Auth.AccessTTLSeconds/60,     // minutes
		cfg.Auth.RefreshTTLSeconds/86400, // days
		cfg.Auth.Issuer,
		"",
	)
	if err != nil {
		logrus.Fatalf("failed creating jwt manager: %v", err)
	}
	authStore := auth.NewPostgresStore(pgDB.DB)
	authSvc := auth.NewService(authStore, jwtMgr)
	registry.Register("auth_service", authSvc)
	logrus.Info("auth service registered")

	// -----------------------------------------
	// Phase 5 – Authority Service
	// -----------------------------------------
	authoritySvc := authority.NewService(pgDB, store)
	registry.Register("authority_service", authoritySvc)
	logrus.Info("authority service registered")

	// -----------------------------------------
	// Phase 5 – Evaluator Service
	// -----------------------------------------

	evSvc := evaluator.NewService(pgDB, store, pyValidatorClient, chain.NewChain(store))
	registry.Register("evaluator_service", evSvc)
	logrus.Info("evaluator service registered")

	submitSvc := evaluator.NewSubmitService(pgDB, store, pyValidatorClient, chain.NewChain(store))
	registry.Register("evaluator_submit_service", submitSvc)
	logrus.Info("evaluator submit service registered")

	// Evaluator Upload Service (New)
	evalUploadSvc := evaluator.NewUploadService(pyClient) // uses pyClient (extractor) for IPFS upload
	registry.Register("evaluator_upload_service", evalUploadSvc)
	logrus.Info("evaluator upload service registered")

	// Release service
	releaseSvc := authority.NewReleaseService(pgDB, chain.NewChain(store))
	registry.Register("authority_release_service", releaseSvc)
	logrus.Info("authority release service registered")

	// student service
	studentSvc := student.NewService(pgDB)
	registry.Register("student_service", studentSvc)
	logrus.Info("student service registered")

	// -----------------------------------------
	// Phase 6 – Admin Service (Process Orchestration)
	// -----------------------------------------
	adminSvc := admin.NewService()
	// Register processes (Paths are relative to where the binary is run, usually root of repo)
	// IPFS
	adminSvc.RegisterProcess("IPFS", ".", "ipfs", "daemon")
	// Python Extractor
	adminSvc.RegisterProcess("Python Extractor", ".", "python", "../python-validator/app/main.py")
	// Python Validator
	adminSvc.RegisterProcess("Python Validator", ".", "python", "../python-validator/app/validate_evaluation.py")

	registry.Register("admin_service", adminSvc)
	logrus.Info("admin service registered")

	// API handler with registry injected
	handler := api.NewHandlerWithRegistry(store, cfg.Block.SignerID, registry)
	router := handler.WithRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// start server
	go func() {
		if cfg.Server.TLS.Enabled {
			logrus.Infof("listening (tls) on %s", srv.Addr)
			if err := srv.ListenAndServeTLS(cfg.Server.TLS.CertPath, cfg.Server.TLS.KeyPath); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("server failed: %v", err)
			}
		} else {
			logrus.Infof("listening on %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Fatalf("server failed: %v", err)
			}
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logrus.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("server shutdown failed: %v", err)
	}
	logrus.Info("server exited properly")
}
