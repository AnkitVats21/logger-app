# Logger App - Project Overview

## 1. Application Purpose & Stack
**Logger App** is a personal time-tracking and productivity dashboard. It allows a user to log their current physical location (home, office, outside) and automatically calculates the time spent at each location, including commute times. 

**Tech Stack:**
- **Backend:** Go (Golang) standard library (`net/http`, `html/template`, `embed`)
- **Database:** SQLite3 (`github.com/mattn/go-sqlite3`)
- **Frontend:** HTML, Vanilla JavaScript, Tailwind CSS (via CDN), Lucide Icons
- **Deployment:** Compiled as a single standalone Go binary with embedded HTML templates (`//go:embed`).

---

## 2. Directory & Module Structure

```text
logger-app/
├── db/                 # SQLite connection and table initialization
│   └── db.go
├── handlers/           # HTTP Request handlers (Controllers)
│   ├── dashboard.go    # Renders the main dashboard HTML
│   ├── event_handler.go# Handles logging new events and rendering the events list
│   ├── middleware.go   # (If applicable) HTTP middlewares
│   └── summary_handler.go # JSON API to fetch aggregated date-range summaries
├── models/             # Data structures
│   ├── event.go        # Event struct (ID, Place, Timestamp)
│   └── summary.go      # DailySummary, RangeSummary, TimeSegment structs
├── scripts/            # Admin/Utility scripts
│   └── seed.go         # Go script to seed the SQLite DB with 30 days of valid historical data
├── service/            # Core Business Logic 
│   ├── event_service.go# Validates and creates events (enforces state machine rules)
│   ├── range_summary.go# Groups events by day and aggregates summaries
│   ├── summary.go      # Calculates chronological TimeSegments & durations between events
│   └── utils.go        # Template rendering (with embed.FS) and string formatters
├── storage/            # Data Access Layer
│   └── event_repo.go   # SQL queries (Insert, GetLatest, Paginated Fetch, Range Fetch)
├── templates/          # Go HTML Templates (Embedded into binary)
│   ├── components/
│   │   └── nav.html    # Shared Navigation bar
│   ├── pages/
│   │   ├── dashboard.html # The rich Dashboard UI (Date Pickers, Timeline, Tables)
│   │   └── events.html    # Raw paginated list view of all events
│   └── layout.html     # Base HTML layout wrapping all pages
├── embed.go            # go:embed directive exporting the templates FS
├── main.go             # Entry point: Server initialization, Routing binding
├── go.mod / go.sum     # Go module definitions
└── events.db           # SQLite database file (created at runtime)
```

---

## 3. Core Business Logic & Rules

### Event State Machine
To ensure time calculations are accurate, the app enforces a geographical state machine. A user cannot teleport directly from `home` to `office`—they must transition through `outside`. 
- **Valid Transitions**: `home <-> outside <-> office`
- **Invalid Transitions**: `home -> home`, `office -> office`, `home -> office`, `office -> home`.
- *Enforced in*: `service/event_service.go` (`LogEvent` function reads `storage.GetLatestEvent()` to validate).

### Time Calculations & Commutes
- The duration spent in a location is calculated as the time delta between event *N* and event *N+1*.
- **Commute Detection**: If an event sequence is `home -> outside -> office` (or vice versa), the intermediate `outside` segment is classified as `commute` if its duration is below a specific threshold (e.g., 90 minutes).
- *Implemented in*: `service/summary.go` (`CalculateSummary` function).

---

## 4. API Endpoints

1. **`GET /`** -> Redirects to `/dashboard`
2. **`GET /dashboard`** -> Returns the dashboard UI HTML.
3. **`GET /events`** -> Returns the paginated raw events UI HTML (`?page=1`).
4. **`GET /log?place={place}`** -> Simple endpoint to log a new event (e.g. `?place=home`).
5. **`GET /summary/range?start={YYYY-MM-DD}&end={YYYY-MM-DD}`** -> Returns JSON telemetry.
    ```json
    {
      "start_date": "2026-03-01",
      "end_date": "2026-03-10",
      "total_home_time": "120h 30m",
      "total_office_time": "40h 0m",
      "total_outside_time": "10h 15m",
      "total_commute_time": "5h 0m",
      "days": [
        {
          "date": "2026-03-10",
          "home_time": "14h 0m",
          "office_time": "8h 0m",
          "outside_time": "1h 0m",
          "commute_time": "1h 0m",
          "segments": [
            { "place": "home", "start": "00:00", "end": "08:00", "duration": "8h 0m" },
            { "place": "commute", "start": "08:00", "end": "08:30", "duration": "30m" },
            { "place": "office", "start": "08:30", "end": "17:00", "duration": "8h 30m" }
          ]
        }
      ]
    }
    ```

---

## 5. Frontend UI

- **Dashboard (`dashboard.html`)**: 
  - **Top Section**: "Daily Insights", an interactive, vertical, color-coded chronological timeline showing exactly how the day was segmented.
  - **Middle Section**: Date Range picker (`start`/`end`) with quick-select buttons (7D, 30D, Month) that re-fetches data dynamically via JavaScript `fetch()`.
  - **Cards**: High-level aggregate duration totals for the selected range.
  - **Bottom Section**: Data Table breaking down aggregate durations day-by-day.
  - *Interactivity*: Clicking a past day row in the table dynamically updates the Top Section Timeline for that specific day.
- **Events Log (`events.html`)**: Standard HTML table displaying paginated sequential logs of timestamps and locations.
