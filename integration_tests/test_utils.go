package integration

import (
	"bufio"
	"os"
	"strings"
)

// getTestBucket resolves the real bucket name from environment or parent .env.
func getTestBucket() string {
	if b := os.Getenv("TARGET_BUCKET"); b != "" {
		return b
	}

	// Fallback: Parse parent .env manually to avoid dependencies
	file, err := os.Open("../.env")
	if err != nil {
		return "deckardy_private" // Safety fallback inside test
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TARGET_BUCKET=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.Trim(parts[1], " \"'")
				return val
			}
		}
	}
	return "deckardy_private"
}

// getTestPrefix resolves the GCS_PREFIX from environment or parent .env.
func getTestPrefix() string {
	if p := os.Getenv("GCS_PREFIX"); p != "" {
		return p
	}

	file, err := os.Open("../.env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "GCS_PREFIX=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.Trim(parts[1], " \"'")
				return val
			}
		}
	}
	return ""
}

// getAWSAccessKey resolves the AWS_ACCESS_KEY_ID from environment or parent .env.
func getAWSAccessKey() string {
	if k := os.Getenv("AWS_ACCESS_KEY_ID"); k != "" {
		return k
	}

	file, err := os.Open("../.env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "AWS_ACCESS_KEY_ID=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.Trim(parts[1], " \"'")
				return val
			}
		}
	}
	return ""
}

// getAWSSecretKey resolves the AWS_SECRET_ACCESS_KEY from environment or parent .env.
func getAWSSecretKey() string {
	if k := os.Getenv("AWS_SECRET_ACCESS_KEY"); k != "" {
		return k
	}

	file, err := os.Open("../.env")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "AWS_SECRET_ACCESS_KEY=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				val := strings.Trim(parts[1], " \"'")
				return val
			}
		}
	}
	return ""
}
