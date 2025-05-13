# Nostos Profiles Service

The **Nostos Profiles Service** manages user profiles and social relationships such as followers and following. It enables profile creation, updates, user discovery, and social connectivity.

---

## 🚀 Features

* User profile creation and updates
* Personal data storage (name, avatar, bio, etc.)
* Follow and unfollow functionality
* User search by name or email
* JWT-based authentication and access control

---

## 📌 Endpoints

### 🔹 Profile Management

* **Create Profile**
  `POST /api/profiles`
  Creates a new user profile.

* **Update My Profile**
  `POST /api/profiles/update`
  Updates the profile of the authenticated user.

* **Update Profile by ID**
  `POST /api/profiles/updateProfileByID`
  Updates a profile using the target user's ID.

* **Delete Profile**
  `POST /api/profiles/delete`
  Deletes the current user's profile.

* **Get Profile by User ID**
  `GET /api/profiles/user/:userID`
  Retrieves a user profile by their user ID.

* **Get Profile by Username**
  `GET /api/profiles/username/:username`
  Retrieves a user profile by username.

* **Search Profiles**
  `POST /api/profiles/search`
  Searches for users by name or email.

* **Get User Avatar**
  `GET /api/profiles/userAvatar/:userID`
  Retrieves the avatar of a specific user.

### 🔹 Followers and Following

* **Follow a User**
  `POST /api/follow/:followedID`
  Follows the user with the specified ID.

* **Unfollow a User**
  `POST /api/unfollow/:followedID`
  Unfollows the user with the specified ID.

* **Get Followers**
  `GET /api/:profileID/followers`
  Lists users who follow the given profile ID.

* **Get Following**
  `GET /api/:profileID/following`
  Lists users that the profile ID is following.

---

## ⚙️ Installation and Configuration

### Prerequisites

* Go installed
* PostgreSQL instance
* Docker and Docker Compose (for local development)
* Auth service with JWT support

### Installation

```bash
git clone https://github.com/nostos-globe/NostosProfiles.git
cd NostosProfiles
go mod download
```

### Configuration

Ensure the following environment variables or Vault secrets are set:

* `DATABASE_URL`
* `JWT_SECRET`

Vault can be accessed using a token, AppRole, or Kubernetes auth depending on your setup.

---

## ▶️ Running the Application

```bash
go run cmd/main.go
```

---

## 🧱 Technologies Used

* **Language**: Go
* **Database**: PostgreSQL
* **Authentication**: JWT
* **Orchestration**: Docker

---

## 🏗️ Project Structure

```
NostosProfiles/
├── cmd/                 # Application entry point
│   └── main.go
├── internal/
│   ├── api/             # HTTP route handlers
│   ├── db/              # Database logic
│   ├── models/          # Data structures
│   └── service/         # Business rules
├── pkg/
│   └── config/          # Configuration management
├── Dockerfile
├── go.mod
└── README.md
```
