package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Settings contains all the configuration for the proxy
type Settings struct {
	Port                string
	GCPProjectID        string
	TargetBucket        string
	StorageBaseURL      string // For testing or custom setups
	GCSPrefix           string // Subfolder prefix for testing or namespacing
	DryRun              bool   // DryRun mode disables real GCS API hits
	DebugLogging        bool   // DebugLogging enables verbose output
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	ProxyAccessKey      string // For re-signing
	ProxySecretKey      string // For re-signing
	JSONKey             string // Path to GCS Service Account JSON key

	// Request data logging (SOH-delimited CSV file via ymlog)
	ReqLogEnabled   bool   // REQUEST_LOG_ENABLED,      default true
	ReqLogPath      string // REQUEST_LOG_PATH,          default "./logs/req_%Y%M%D.csv"
	ReqLogMaxSizeMB int    // REQUEST_LOG_MAX_SIZE_MB,   default 512
	ReqLogMaxBackup int    // REQUEST_LOG_MAX_BACKUP,    default 5
	ReqLogChanBuf   int    // REQUEST_LOG_CHAN_BUF,       default 10240
	ReqLogKeepDays  int    // REQUEST_LOG_KEEP_DAYS,     default 7
}

var Config *Settings

// LoadConfig initialize the settings from a .env file or environment variables
func LoadConfig() {
	// Load from .env if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables directly.")
	}

	dryRunStr := getEnv("DRY_RUN", "true") // Default to true if not set (Safe for laptop testing)
	dryRun := dryRunStr == "true"
	debugLogging := getEnv("DEBUG_LOGGING", "false") == "true"

	maxIdleConns, _ := strconv.Atoi(getEnv("MAX_IDLE_CONNS", "1000"))
	maxIdleConnsPerHost, _ := strconv.Atoi(getEnv("MAX_IDLE_CONNS_PER_HOST", "1000"))

	reqLogEnabled := getEnv("REQUEST_LOG_ENABLED", "true") == "true"
	reqLogMaxSizeMB := 512
	if v := getEnv("REQUEST_LOG_MAX_SIZE_MB", "512"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			reqLogMaxSizeMB = n
		}
	}
	reqLogMaxBackup := 5
	if v := getEnv("REQUEST_LOG_MAX_BACKUP", "5"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			reqLogMaxBackup = n
		}
	}
	reqLogChanBuf := 10240
	if v := getEnv("REQUEST_LOG_CHAN_BUF", "10240"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			reqLogChanBuf = n
		}
	}
	reqLogKeepDays := 7
	if v := getEnv("REQUEST_LOG_KEEP_DAYS", "7"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			reqLogKeepDays = n
		}
	}

	Config = &Settings{
		Port:                getEnv("PORT", "8080"),
		GCPProjectID:        getEnv("GCP_PROJECT_ID", ""),
		TargetBucket:        getEnv("TARGET_BUCKET", ""),
		StorageBaseURL:      getEnv("STORAGE_BASE_URL", "https://storage.googleapis.com"),
		GCSPrefix:           getEnv("GCS_PREFIX", ""),
		DryRun:              dryRun,
		DebugLogging:        debugLogging,
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		ProxyAccessKey:      getEnv("PROXY_AWS_ACCESS_KEY_ID", getEnv("AWS_ACCESS_KEY_ID", "")),
		ProxySecretKey:      getEnv("PROXY_AWS_SECRET_ACCESS_KEY", getEnv("AWS_SECRET_ACCESS_KEY", "")),
		JSONKey:             getEnv("JSON_KEY", ""),
		ReqLogEnabled:       reqLogEnabled,
		ReqLogPath:          getEnv("REQUEST_LOG_PATH", "./logs/req_%Y%M%D.csv"),
		ReqLogMaxSizeMB:     reqLogMaxSizeMB,
		ReqLogMaxBackup:     reqLogMaxBackup,
		ReqLogChanBuf:       reqLogChanBuf,
		ReqLogKeepDays:      reqLogKeepDays,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
