# LMS Backend

Backend system for the Learning Management System (LMS) project. This system is designed to provide a scalable and modular backend architecture for managing educational content, user interactions, and administrative functionalities.

---

## Technologies Used

- **Programming Language:** Go  
- **Framework:** Fiber (for REST APIs)  
- **Database:** PostgreSQL  
- **JWT:** Authentication  
- **ORM:** No ORM ORM Gang 
- **Logging:** Zerolog

---

# Getting Started

Follow these steps to set up and run the LMS backend project locally.

## Prerequisites

- Install [Go](https://golang.org/doc/install) version >= 1.18
- PostgreSQL installed and running

## Steps

### 1. Clone the repository

Clone the project repository to your local machine:

```bash
git clone <repository-url>
cd lms-backend
```
### 2. Install dependencies

Install all project dependencies:

```bash
go mod tidy
```

### 3. Set up environment variables

Create a .env file in the project root based on the .env.example file. Update the values to match your local environment:

```bash
JWT_SECRET_KEY=your_secret_key
DB_HOST=localhost
DB_PORT=5432
DB_USER=username
DB_PASSWORD=password
DB_NAME=lms_db
```

### 4. Run the application

Start the application using the following command:

```bash
go run cmd/bin/main.go
```

---

# Developer Guide

## General Rules

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org) for commit messages:

- **`feat:`** For new features.  
  Example: `feat: add user login functionality`
- **`fix:`** For bug fixes.  
  Example: `fix: resolve null pointer exception in login service`
- **`docs:`** For documentation updates.  
  Example: `docs: update README with setup instructions`
- **`style:`** For code style updates (e.g., formatting, missing semicolons).  
  Example: `style: format code in user handler`
- **`refactor:`** For code refactoring without functional changes.  
  Example: `refactor: optimize database query structure`
- **`test:`** For test file updates or additions.  
  Example: `test: add unit tests for user service`

### Changelog

Use [git-chglog](https://github.com/git-chglog/git-chglog) to generate a changelog before merging to the release branch.

```bash
git-chglog > CHANGELOG.md
```
