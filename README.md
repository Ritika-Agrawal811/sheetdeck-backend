# 📚 SheetDeck Backend Server

**sheetdeck-backend** is a production-oriented backend server written in **Go** for [Sheetdeck](https://sheetdeck.vercel.app) — a website built by me that delivers concise, structured cheat sheets for modern web development.

It showcases real-world backend engineering practices including API design, modular architecture, and maintainable code organization.

## 🎯 Why This Project?

The goals of this project include :

- Supporting the long-term growth of **Sheetdeck** as a content-driven web platform.
- Enabling structured storage and easy addition of new cheat sheets
- Tracking usage and analytics to understand user behavior
- Providing a clean foundation for adding future features without major refactors


## 🧠 Key Highlights

- Designed and implemented a complete **cheat sheets API** with support for single and bulk uploads, updates, and efficient data retrieval
- Built flexible fetch endpoints with filtering by category, subcategory, view count, and download count
- Implemented **analytics APIs** to track user interactions including page views, clicks, and downloads
- Exposed aggregated analytics insights such as usage by country, device type, operating system, and browser
- Added **configuration and system metrics APIs** to monitor database usage and storage consumption
- Wrote **unit tests** for core services and reusable packages to ensure correctness and maintainability
- Implemented middleware for **rate limiting** and **origin validation** to improve security and protect backend resources
- Used environment-based configuration for local and production setups

## 🛠 Tech Stack

- **Language:** Go (Golang)
- **Web Framework:** Gin (HTTP router and middleware)
- **API Style:** RESTful APIs
- **Backend Architecture:** Modular, service-oriented design
- **Database:** Supabase (PostgreSQL) for cheat sheets, analytics, and configuration data
- **Storage:** Supabase Storage for cheat sheet assets
- **Testing:** Go `testing` package for unit tests
- **Middleware:** Custom middleware for rate limiting and origin validation
- **External Services:** ipinfo.io for IP-based country resolution in analytics workflows
- **Tooling:** Go modules and standard Go tooling

## 🏗 Architecture Overview

The backend follows a layered architecture with a clear separation of concerns to keep the system maintainable and extensible.

- **HTTP Layer:** Incoming requests are handled by Gin, with middleware applied for rate limiting and origin validation.
- **API / Handler Layer:** Route handlers are responsible for request validation and response formatting.
- **Service Layer:** Core business logic lives in services, including cheat sheet management, analytics processing, and configuration handling.
- **Data Access Layer:** Services interact with Supabase (PostgreSQL) for persistent storage and Supabase Storage for file assets.
- **External Integrations:** IP-based location data is resolved using ipinfo.io as part of analytics processing.
- **Cross-Cutting Concerns:** Logging, error handling, and configuration are shared across layers.

## 📈 Future Improvements

- Add **Redis caching** for read-heavy cheat sheet APIs
- Introduce **cron jobs** for analytics data backups and retention management
- Improve analytics aggregation and system observability


## 📄 License

The **codebase** of this project is open source and licensed under the [MIT License](./LICENSE).


## 🚫 Contributions

This is a personal project. I'm NOT accepting pull requests, issues, or external contributions.

<p align="center">• • •</p>

## ☕️ Support Me

Give a ⭐ if you like this project, and share it with your friends. Your support means a lot and helps me create more useful resources!

<p align="left">
  <a href="https://buymeacoffee.com/ritikaagrawal08"><img alt="Buy me a coffee" title="Buy me a coffee" src="https://img.shields.io/badge/-Buy%20me%20a%20coffee-yellow?style=for-the-badge&logo=buymeacoffee&logoColor=white"/></a>
</p>

**Thanks so much! Happy Coding!** :sparkles: