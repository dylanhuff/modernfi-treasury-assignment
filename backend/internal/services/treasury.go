package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"modernfi-treasury-app/internal/models"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	treasuryURLTemplate  = "https://home.treasury.gov/resource-center/data-chart-center/interest-rates/pages/xml?data=daily_treasury_yield_curve&field_tdr_date_value=%d"
	httpTimeout          = 10 * time.Second
	httpTimeoutMultiYear = 30 * time.Second // Longer timeout for multi-year requests
	cacheDuration        = 1 * time.Hour
	iso8601DateLength    = 10 // Length of "YYYY-MM-DD"
)

// historicalCacheEntry stores cached historical yield data with a timestamp
type historicalCacheEntry struct {
	data      *models.HistoricalYieldData
	timestamp time.Time
}

// TreasuryService handles fetching and caching of treasury yield data
type TreasuryService struct {
	cacheData      *models.YieldData
	cacheTimestamp time.Time
	cacheDuration  time.Duration
	mu             sync.RWMutex
	httpClient     *http.Client

	historicalCache map[string]*historicalCacheEntry
	historicalMu    sync.RWMutex
}

var historicalPeriods = []string{"1W", "1M", "3M", "6M", "1Y", "5Y", "10Y", "30Y"}

func NewTreasuryService() *TreasuryService {
	return &TreasuryService{
		cacheDuration: cacheDuration,
		httpClient: &http.Client{
			Timeout: httpTimeout,
		},
		historicalCache: make(map[string]*historicalCacheEntry),
	}
}

// calculateDateRange returns start and end dates for the given period
func calculateDateRange(period string) (startDate, endDate time.Time, err error) {
	endDate = time.Now()

	switch period {
	case "1W":
		startDate = endDate.AddDate(0, 0, -7)
	case "1M":
		startDate = endDate.AddDate(0, -1, 0)
	case "3M":
		startDate = endDate.AddDate(0, -3, 0)
	case "6M":
		startDate = endDate.AddDate(0, -6, 0)
	case "1Y":
		startDate = endDate.AddDate(-1, 0, 0)
	case "5Y":
		startDate = endDate.AddDate(-5, 0, 0)
	case "10Y":
		startDate = endDate.AddDate(-10, 0, 0)
	case "30Y":
		startDate = endDate.AddDate(-30, 0, 0)
	default:
		return time.Time{}, time.Time{}, fmt.Errorf("invalid period: %s (must be 1W, 1M, 3M, 6M, 1Y, 5Y, 10Y, or 30Y)", period)
	}

	return startDate, endDate, nil
}

func (s *TreasuryService) fetchFromAPI() (*models.TreasuryFeed, error) {
	url := fmt.Sprintf(treasuryURLTemplate, time.Now().Year())
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch treasury data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("treasury API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed models.TreasuryFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	if len(feed.Entries) == 0 {
		return nil, fmt.Errorf("no entries found in treasury feed")
	}

	return &feed, nil
}

// fetchFromAPIForYears fetches and combines data from multiple years in parallel
func (s *TreasuryService) fetchFromAPIForYears(startYear, endYear int) (*models.TreasuryFeed, error) {
	client := &http.Client{
		Timeout: httpTimeoutMultiYear,
	}

	yearCount := endYear - startYear + 1

	type yearResult struct {
		year    int
		entries []models.Entry
		err     error
	}
	results := make(chan yearResult, yearCount)

	for year := startYear; year <= endYear; year++ {
		go func(y int) {
			url := fmt.Sprintf(treasuryURLTemplate, y)
			resp, err := client.Get(url)
			if err != nil {
				results <- yearResult{year: y, err: fmt.Errorf("failed to fetch treasury data for year %d: %w", y, err)}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- yearResult{year: y, err: fmt.Errorf("treasury API returned status %d for year %d", resp.StatusCode, y)}
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				results <- yearResult{year: y, err: fmt.Errorf("failed to read response body for year %d: %w", y, err)}
				return
			}

			var feed models.TreasuryFeed
			if err := xml.Unmarshal(body, &feed); err != nil {
				results <- yearResult{year: y, err: fmt.Errorf("failed to parse XML for year %d: %w", y, err)}
				return
			}

			results <- yearResult{year: y, entries: feed.Entries, err: nil}
		}(year)
	}

	yearData := make(map[int][]models.Entry)
	var errors []error

	for i := 0; i < yearCount; i++ {
		result := <-results
		if result.err != nil {
			errors = append(errors, result.err)
		} else {
			yearData[result.year] = result.entries
		}
	}

	if len(errors) > 0 {
		return nil, errors[0]
	}

	var combinedFeed models.TreasuryFeed
	for year := startYear; year <= endYear; year++ {
		if entries, ok := yearData[year]; ok {
			combinedFeed.Entries = append(combinedFeed.Entries, entries...)
		}
	}

	if len(combinedFeed.Entries) == 0 {
		return nil, fmt.Errorf("no entries found in treasury feed for years %d-%d", startYear, endYear)
	}

	return &combinedFeed, nil
}

// convertToYieldData transforms the most recent XML entry into YieldData format
func (s *TreasuryService) convertToYieldData(feed *models.TreasuryFeed) (*models.YieldData, error) {
	if len(feed.Entries) == 0 {
		return nil, fmt.Errorf("no entries to convert")
	}

	entry := feed.Entries[len(feed.Entries)-1]

	date := entry.Date
	if len(date) > iso8601DateLength {
		date = date[:iso8601DateLength]
	}

	yields := []models.YieldPoint{
		{Term: "1M", Rate: entry.BC1Month},
		{Term: "3M", Rate: entry.BC3Month},
		{Term: "6M", Rate: entry.BC6Month},
		{Term: "1Y", Rate: entry.BC1Year},
		{Term: "2Y", Rate: entry.BC2Year},
		{Term: "5Y", Rate: entry.BC5Year},
		{Term: "10Y", Rate: entry.BC10Year},
		{Term: "30Y", Rate: entry.BC30Year},
	}

	return &models.YieldData{
		Date:   date,
		Yields: yields,
	}, nil
}

// sampleDataPoints reduces data density for long periods (30Y: monthly, 10Y/5Y: weekly)
func sampleDataPoints(dataPoints []map[string]interface{}, period string) []map[string]interface{} {
	if period == "1W" || period == "1M" || period == "3M" || period == "6M" || period == "1Y" {
		return dataPoints
	}

	if len(dataPoints) == 0 {
		return dataPoints
	}

	var samplingInterval int
	switch period {
	case "30Y":
		samplingInterval = 30
	case "10Y", "5Y":
		samplingInterval = 7
	default:
		return dataPoints
	}

	intervalMap := make(map[string]map[string]interface{})

	for _, point := range dataPoints {
		dateStr, ok := point["date"].(string)
		if !ok {
			continue
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		var intervalKey string
		if samplingInterval == 30 {
			intervalKey = date.Format("2006-01")
		} else {
			year, week := date.ISOWeek()
			intervalKey = fmt.Sprintf("%d-W%02d", year, week)
		}

		if existing, exists := intervalMap[intervalKey]; exists {
			existingDate, _ := existing["date"].(string)
			if dateStr > existingDate {
				intervalMap[intervalKey] = point
			}
		} else {
			intervalMap[intervalKey] = point
		}
	}

	sampledPoints := make([]map[string]interface{}, 0, len(intervalMap))
	for _, point := range intervalMap {
		sampledPoints = append(sampledPoints, point)
	}

	sort.Slice(sampledPoints, func(i, j int) bool {
		dateI, _ := sampledPoints[i]["date"].(string)
		dateJ, _ := sampledPoints[j]["date"].(string)
		return dateI < dateJ
	})

	return sampledPoints
}

// convertToHistoricalData builds time-series dataset from feed entries
func (s *TreasuryService) convertToHistoricalData(
	feed *models.TreasuryFeed,
	startDate, endDate time.Time,
	period string,
) (*models.HistoricalYieldData, error) {
	var dataPoints []map[string]interface{}

	for _, entry := range feed.Entries {
		dateStr := entry.Date
		if len(dateStr) > iso8601DateLength {
			dateStr = dateStr[:iso8601DateLength]
		}

		entryDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		if entryDate.Before(startDate) || entryDate.After(endDate) {
			continue
		}

		point := map[string]interface{}{
			"date": dateStr,
			"10Y":  entry.BC10Year,
			"5Y":   entry.BC5Year,
			"2Y":   entry.BC2Year,
		}
		dataPoints = append(dataPoints, point)
	}

	sampledPoints := sampleDataPoints(dataPoints, period)

	return &models.HistoricalYieldData{
		Period:    period,
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		Terms:     []string{"10Y", "5Y", "2Y"},
		Data:      sampledPoints,
	}, nil
}

// GetHistoricalYields fetches historical yield data with permanent caching
func (s *TreasuryService) GetHistoricalYields(period string) (*models.HistoricalYieldData, error) {
	s.historicalMu.RLock()
	if cached, exists := s.historicalCache[period]; exists {
		data := cached.data
		s.historicalMu.RUnlock()
		return data, nil
	}
	s.historicalMu.RUnlock()

	s.historicalMu.Lock()
	defer s.historicalMu.Unlock()

	if cached, exists := s.historicalCache[period]; exists {
		return cached.data, nil
	}

	fmt.Printf("Fetching historical yields for period %s (cache miss)\n", period)

	startDate, endDate, err := calculateDateRange(period)
	if err != nil {
		return nil, err
	}

	var feed *models.TreasuryFeed
	startYear := startDate.Year()
	endYear := endDate.Year()

	if startYear == endYear {
		feed, err = s.fetchFromAPI()
	} else {
		feed, err = s.fetchFromAPIForYears(startYear, endYear)
	}

	if err != nil {
		return nil, err
	}

	data, err := s.convertToHistoricalData(feed, startDate, endDate, period)
	if err != nil {
		return nil, err
	}

	s.historicalCache[period] = &historicalCacheEntry{
		data:      data,
		timestamp: time.Now(),
	}

	return data, nil
}

// GetLatestYields returns latest yields with 1-hour caching
func (s *TreasuryService) GetLatestYields() (*models.YieldData, error) {
	s.mu.RLock()
	if s.cacheData != nil && time.Since(s.cacheTimestamp) < s.cacheDuration {
		data := s.cacheData
		s.mu.RUnlock()
		return data, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cacheData != nil && time.Since(s.cacheTimestamp) < s.cacheDuration {
		return s.cacheData, nil
	}

	feed, err := s.fetchFromAPI()
	if err != nil {
		return nil, err
	}

	data, err := s.convertToYieldData(feed)
	if err != nil {
		return nil, err
	}

	s.cacheData = data
	s.cacheTimestamp = time.Now()

	return data, nil
}

// WarmCache pre-fetches all historical data in background on startup
func (s *TreasuryService) WarmCache() {
	log.Println("Starting historical yield cache warming for all periods...")

	for _, period := range historicalPeriods {
		go func(p string) {
			log.Printf("Warming cache for period: %s", p)
			start := time.Now()

			if _, err := s.GetHistoricalYields(p); err != nil {
				log.Printf("ERROR: Failed to warm cache for period %s: %v", p, err)
			} else {
				log.Printf("Cache warmed successfully for period %s in %v", p, time.Since(start))
			}
		}(period)
	}
}
