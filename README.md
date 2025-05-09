# Profiles Service (Profiles and Followers)

## Description
The user service manages profiles and relationships between users, such as followers and following. It also allows searching and optimizing profile access through Redis caching.

## Features
- Creation and updating of user profiles.
- Storage of personal data (name, photo, bio, etc.).
- Management of followers and following (follow/unfollow).
- Search users by name or email.
- Profile caching in Redis for optimization.

## Technologies Used
- **Language**: Go
- **Database**: PostgreSQL
- **Cache**: Redis
- **Orchestration**: Docker

## Endpoints
| Method | Route                           | Description |
|--------|--------------------------------|-------------|
| POST   | /api/profiles                  | Creates a new profile |
| POST   | /api/profiles/update           | Updates current user profile |
| POST   | /api/profiles/updateProfileByID | Updates profile by ID |
| POST   | /api/profiles/delete           | Deletes a profile |
| GET    | /api/profiles/user/:userID     | Gets profile by user ID |
| GET    | /api/profiles/username/:username| Gets profile by username |
| POST   | /api/profiles/search           | Searches profiles |
| GET    | /api/profiles/userAvatar/:userID| Gets user's avatar |
| POST   | /api/follow/:followedID        | Follows a user |
| POST   | /api/unfollow/:followedID      | Unfollows a user |
| GET    | /api/:profileID/followers      | Lists user's followers |
| GET    | /api/:profileID/following      | Lists users being followed |

## Security
- Authentication through JWT.
- User search limits to prevent abuse.
- Redis caching to optimize profile loading.
