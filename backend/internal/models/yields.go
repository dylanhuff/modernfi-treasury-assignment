package models

import "encoding/xml"

// YieldPoint represents a single term and its corresponding yield rate
type YieldPoint struct {
	Term string  `json:"term"` // e.g., "1M", "3M", "6M"
	Rate float64 `json:"rate"` // e.g., 4.45
}

// YieldData represents the complete yield data for a specific date
type YieldData struct {
	Date   string       `json:"date"`   // ISO 8601 date
	Yields []YieldPoint `json:"yields"` // Array of yield points
}

// TreasuryFeed represents the XML feed structure from Treasury.gov
type TreasuryFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

// Entry represents a single entry in the Treasury XML feed
type Entry struct {
	Date     string  `xml:"content>properties>NEW_DATE"`
	BC1Month float64 `xml:"content>properties>BC_1MONTH"`
	BC3Month float64 `xml:"content>properties>BC_3MONTH"`
	BC6Month float64 `xml:"content>properties>BC_6MONTH"`
	BC1Year  float64 `xml:"content>properties>BC_1YEAR"`
	BC2Year  float64 `xml:"content>properties>BC_2YEAR"`
	BC5Year  float64 `xml:"content>properties>BC_5YEAR"`
	BC10Year float64 `xml:"content>properties>BC_10YEAR"`
	BC30Year float64 `xml:"content>properties>BC_30YEAR"`
}

// HistoricalYieldData represents time-series yield data for a specific period
// The data is formatted for direct consumption by Tremor LineChart component
// Data array contains flattened objects: {date: "2025-01-02", "10Y": 4.25, "5Y": 4.10, "2Y": 4.05}
type HistoricalYieldData struct {
	Period    string                   `json:"period"`    // "1M", "3M", "6M", or "1Y"
	StartDate string                   `json:"startDate"` // YYYY-MM-DD format
	EndDate   string                   `json:"endDate"`   // YYYY-MM-DD format
	Terms     []string                 `json:"terms"`     // e.g., ["10Y", "5Y", "2Y"]
	Data      []map[string]interface{} `json:"data"`      // Flattened for Tremor chart compatibility
}
