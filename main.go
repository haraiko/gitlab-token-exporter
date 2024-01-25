package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	listenAddress    = ":9960"
	metricsEndpoint  = "/metrics"
	metricNamePrefix = "gitlab_access_tokens"
)

var (
	gitlabAPIRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: metricNamePrefix + "_request_duration_seconds",
			Help: "Duration of GitLab API requests",
		},
		[]string{"endpoint"},
	)

	gitlabAPIErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricNamePrefix + "_api_errors_total",
			Help: "Total number of errors when making GitLab API requests",
		},
		[]string{"endpoint"},
	)

	gitlabTokenCreationDate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricNamePrefix + "_creation_date_seconds",
			Help: "Creation date of GitLab token in seconds since epoch",
		},
		[]string{"token_id"},
	)

	gitlabTokenExpiryDate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: metricNamePrefix + "_expiry_date_seconds",
			Help: "Expiry date of GitLab token in seconds since epoch",
		},
		[]string{"token_id"},
	)
)

func init() {
	prometheus.MustRegister(gitlabAPIRequestDuration)
	prometheus.MustRegister(gitlabAPIErrors)
        prometheus.MustRegister(gitlabTokenCreationDate)
	prometheus.MustRegister(gitlabTokenExpiryDate)
}

func main() {
	err := loadEnvFile("token.env")
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	http.Handle(metricsEndpoint, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h1>Hello, Prometheus!</h1>")
	})

	go func() {
		fmt.Printf("Listening on %s\n", listenAddress)
		http.ListenAndServe(listenAddress, nil)
	}()

	for {
		err := fetchGitLabAPIData()
		if err != nil {
			fmt.Printf("Error fetching GitLab API data: %v\n", err)
			gitlabAPIErrors.WithLabelValues(metricsEndpoint).Inc()
		}
		time.Sleep(time.Minute * 5) // Fetch data every 5 minutes
	}
}


func fetchGitLabAPIData() error {
	gitlabAPIURL := os.Getenv("GITLAB_API_URL")
	privateToken := os.Getenv("PRIVATE_TOKEN")
	userID := os.Getenv("USER_ID")

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Seconds()
		gitlabAPIRequestDuration.WithLabelValues(metricsEndpoint).Observe(duration)
	}()

	client := &http.Client{}
	req, err := http.NewRequest("GET", gitlabAPIURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("PRIVATE-TOKEN", privateToken)
	req.URL.Query().Add("user_id", userID)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}


	for _, token := range result {
		name, ok := token["name"].(string)
		if !ok {
			fmt.Println("Error getting token name")
			continue
		}

		creationDateString, ok := token["created_at"].(string)
		if !ok {
			fmt.Println("Error getting creation date string")
			continue
		}
		creationDate, err := time.Parse(time.RFC3339, creationDateString)
		if err == nil {
			gitlabTokenCreationDate.WithLabelValues(name).Set(float64(creationDate.UnixNano()) / float64(time.Second))
		} else {
			fmt.Printf("Error parsing creation date: %v\n", err)
		}

		expiryDateString, ok := token["expires_at"].(string)
		if !ok {
			fmt.Println("Error getting expiry date string")
			continue
		}
		// Try parsing the date with multiple layouts
		expiryDate, err := time.Parse("2006-01-02", expiryDateString)
		if err != nil {
			expiryDate, err = time.Parse(time.RFC3339, expiryDateString)
		}
		if err == nil {
			gitlabTokenExpiryDate.WithLabelValues(name).Set(float64(expiryDate.UnixNano()) / float64(time.Second))
		} else {
			fmt.Printf("Error parsing expiry date: %v\n", err)
		}
	}

	fmt.Printf("GitLab API response: %+v\n", result)
	return nil
}

func loadEnvFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				os.Setenv(key, value)
			}
		}
	}

	return nil
}
