# RSSHub

![RSSHub](https://img.shields.io/badge/status-competed-green)

## Overview

**RSSHub** is a **CLI RSS feed aggregator** that:

- Fetches and parses RSS feeds
- Stores articles in PostgreSQL
- Aggregates feeds using a scalable worker pool
- Supports dynamic interval and worker configuration

This tool helps users stay informed without visiting multiple websites, making it ideal for journalists, researchers, and tech enthusiasts.

Features
--------

- Command-line interface (CLI)
- Background feed aggregation with a worker pool
- Configurable fetch interval and number of workers
- Add, list, delete RSS feeds
- View latest articles per feed
- Graceful shutdown for running aggregator


Technologies
---------

- **Languages:** Go, 
- **Database:** PostgreSQL  
- **Tools:** Docker, Docker Compose  

Commands
------------
### - Start Background Fetchin

```sh
rsshub fetch
```
- Launches background feed aggregation with a worker pool
- Default interval: 3 minutes
- Default workers: 3

### - Add New RSS Feed
```sh
rsshub add --name "tech-crunch" --url "https://techcrunch.com/feed/"
```

### - Set Fetch Interval
```sh
rsshub set-interval 2m
```
- Changes fetch interval dynamically without stopping the process

### - Set Number of Workers
```sh
rsshub set-workers 5
```
- Dynamically resizes the worker pool

### - List RSS Feeds
```sh
rsshub list --num 5
```
- Shows recently added feeds (default shows all if ```--num``` not set)

### - Delete RSS Feed
```sh
rsshub delete --name "tech-crunch"
```
- Displays latest articles for a specific feed

### - Show Latest Articles
```sh
rsshub articles --feed-name "tech-crunch" --num 5
```

### - CLI Help
```sh
rsshub --help
```

Database Scheme
---------
### Database Schema

**Table: `feeds`**
| Field      | Type          | Description                    |
| ---------- | ------------- | ------------------------------ |
| id         | UUID (PK)     | Unique identifier for the feed |
| created_at | TIMESTAMP     | When the feed was added        |
| updated_at | TIMESTAMP     | Last update timestamp          |
| name       | TEXT (unique) | Human-readable feed name       |
| url        | TEXT          | RSS feed URL                   |

**Table: `articles`**
| Field        | Type      | Description                        |
| ------------ | --------- | ---------------------------------- |
| id           | UUID (PK) | Unique identifier for the article  |
| created_at   | TIMESTAMP | When the article was stored        |
| updated_at   | TIMESTAMP | Last modification timestamp        |
| title        | TEXT      | Article title                      |
| link         | TEXT      | Article URL                        |
| published_at | TIMESTAMP | Original publication date          |
| description  | TEXT      | Short description from RSS feed    |
| feed_id      | UUID      | Foreign key referencing `feeds.id` |


Contribute
----------

- [Yerlan](https://github.com/Yerl1)
- [Adilet](https://github.com/adal4ik)
