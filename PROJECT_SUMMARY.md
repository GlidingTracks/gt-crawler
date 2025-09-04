# Project Summary

## Current Features
- Crawls gliding track links using a headless Chrome instance driven by [chromedp](https://github.com/chromedp/chromedp).
- Currently supports scraping daily flight records from **XContest** via the `XContestChrome` implementation.
- Stores previously crawled dates using `jConfig-go` to avoid reprocessing the same day.
- Generates authenticated upload requests by retrieving Firebase tokens and posting IGC files to a backend service.
- Provides helper utilities for date calculations, file downloading, and multipart form creation.
- Includes basic unit tests for date utilities and upload helpers.

## Feature Blueprint
- **Multi-site crawling:** Add new `sites.ChromeSite` implementations for other gliding websites and allow runtime configuration of which sites to crawl.
- **Improved scheduling:** Persist crawler state and schedule recurring runs, enabling automatic daily or hourly execution.
- **Robust error handling and retries:** Detect transient network or rendering errors and retry failed crawls or uploads.
- **Concurrent processing:** Streamline link extraction and file uploads to leverage Go's concurrency for better throughput.
- **Extensive testing:** Expand unit and integration test coverage, especially around crawler and uploader behaviour.
- **Configurability:** Centralize configuration (e.g., disallowed domains, time windows) and support environment variable overrides for container deployments.

