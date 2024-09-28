# Metrics Service

This project is a simple metrics service that allows for the storage, retrieval, and updating of various types of metrics, such as gauges and counters. Built using Go and the Chi routing package, this service can be integrated into larger applications for monitoring and metrics management.

## Features

- **Store Metrics**: Supports gauge and counter metrics.
- **Retrieve Metrics**: Fetch the current value of a specific metric by type and name.
- **List Metrics**: Serve an HTML page listing all known metrics and their values.
- **Concurrency Support**: Uses memory storage with proper synchronization for concurrent access.

## Installation

To get started, clone this repository:

```bash
git clone <repository-url>
cd metrics-service
```

## Prerequisites

Make sure you have Go installed on your machine. You can download it from the official Go website.

## Usage

Run the Service: Use the following command to start the service:

```bash
go run cmd/agent/main.go
```
Update Metrics: Send a POST request to update a metric. For example:

```bash
curl -X POST http://localhost:8080/update/gauge/testGauge/10.5
```
Retrieve a Metric Value: Send a GET request to fetch the value of a specific metric:

```bash
curl http://localhost:8080/value/gauge/testGauge
```

List All Metrics: Access the root URL to see a list of all metrics in HTML format:

```bash
curl http://localhost:8080/
```

## Contributing

Contributions are welcome! If you'd like to contribute to this project, please fork the repository and create a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
Feel free to modify the content to better fit your project's specifics or to add any additional sections that you think might be relevant!