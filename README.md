# Com-In Time Service

## Overview

Com-In Time Service is a comprehensive attendance and timesheet management system designed to streamline employee time tracking and management. The service provides an API for recording employee attendance, managing timesheets, and generating reports.

## Features

- **Attendance Management**
  - Record employee check-in and check-out times
  - QR code-based attendance scanning
  - Daily attendance reports and statistics

- **Timesheet Management**
  - Create and manage time entries
  - Automatic calculation of working hours and overtime
  - Approval workflow for timesheets

- **Reporting**
  - Daily, weekly, and monthly attendance reports
  - Department-based attendance summaries
  - Working hours and overtime statistics

## Tech Stack

- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **API**: RESTful API
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)

### Running with Docker

1. Clone the repository:
   ```
   git clone https://github.com/your-org/comin-time-service.git
   cd comin-time-service
   ```

2. Build and start the application:
   ```
   docker-compose up -d
   ```

3. The service will be available at `http://localhost:8080`

### Local Development

1. Set up the PostgreSQL database:
   ```
   docker-compose up -d postgres
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application:
   ```
   go run cmd/api/main.go
   ```

## API Documentation

### Attendance Endpoints

- `POST /api/attendance/mark`: Mark employee attendance
- `GET /api/attendance/records`: Get attendance records with filters
- `GET /api/attendance/stats`: Get attendance statistics
- `GET /api/attendance/employee/:id/:date`: Get employee attendance for a specific date
- `GET /api/attendance/qr-scanner`: Get QR scanner page

### Timesheet Endpoints

- `POST /api/timesheet/entries`: Create a new timesheet entry
- `PUT /api/timesheet/entries/:id`: Update a timesheet entry
- `GET /api/timesheet/entries`: Get timesheet entries with filters
- `GET /api/timesheet/entries/:id`: Get a timesheet entry by ID
- `POST /api/timesheet/entries/:id/approve`: Approve a timesheet entry
- `POST /api/timesheet/entries/:id/reject`: Reject a timesheet entry
- `GET /api/timesheet/stats`: Get timesheet statistics

### QR Code Endpoints

- `GET /api/qr/employee/:id`: Generate QR code for an employee

## Project Structure

```
├── cmd/               # Application entrypoints
│   └── api/           # API server
├── config/            # Configuration handling
├── internal/          # Private application code
│   ├── api/           # API handlers
│   ├── model/         # Data models
│   ├── repository/    # Database repositories
│   └── service/       # Business logic services
├── migrations/        # Database migrations
├── pkg/               # Public libraries
│   └── database/      # Database utilities
└── utils/             # Utility functions
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.