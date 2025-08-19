### Hope Backend (Go + gRPC + MySQL)

Backend for a simple ride sharing workflow built with Go, gRPC, and MySQL. It supports Google Sign‑In (via ID token) to issue backend JWTs, and provides APIs for users, rides (offers/requests), matching, chat, reviews, and location.

### What this is ?
This service is a focused, production-style gRPC backend for coordinating rides between people in the same network (e.g., a campus or company domain). Users authenticate with Google, the backend verifies the ID token with Google, then issues its own JWT for subsequent calls. Once authenticated, users can:
- Create ride offers (drivers) or ride requests (riders)
- Get matched in two ways: riders can request to join an offer, or drivers can accept a rider’s request (which creates a match)
- Chat only when a match is accepted/completed (enforced by business rules)
- Share last-known location using geohash prefix queries for nearby lookups
- Review each other after a ride

The codebase is organized with clear layering: gRPC handlers -> services (business rules) -> repositories (GORM) -> MySQL. Cross-cutting concerns (auth) are handled with a gRPC interceptor, and dependencies are wired with Google Wire for clean construction and testability.

### End-to-end user journey
1) Login: Client gets a Google ID token and calls `AuthService/Login`. The server verifies the token with Google, checks the `aud` against `GOOGLE_CLIENT_ID`, ensures the email is verified and belongs to an `ALLOWED_DOMAINS` domain, creates the user if needed, and returns a backend JWT.
2) Profile & presence: Authenticated users can fetch/update their profile and upsert their current location. Location is stored along with a geohash for prefix-based proximity queries.
3) Supply and demand:
   - Drivers post ride offers with route, time, seats, and optional fare.
   - Riders post ride requests with route, time, and seats.
4) Matching:
   - Riders can `RequestToJoin` a driver’s offer (driver later `AcceptRequest` or `RejectRequest`).
   - Drivers can `AcceptRideRequest` directly on a rider’s request, which creates an offer+match and marks the request as matched.
5) Communication: After a match is accepted/completed, participants can use `ChatService/SendMessage` on that ride ID. The service enforces that only the matched rider/driver can send.
6) Wrap-up: Driver (or system) marks the match `completed`, and users can submit reviews. Listing RPCs exist to retrieve a user’s data, nearby offers/requests, messages, matches, and reviews.

### Architecture at a glance
- gRPC server with reflection enabled. Middleware intercepts all non-public RPCs and enforces JWT auth. The interceptor injects `user_id` and `email` into the request context for handlers.
- Handlers translate protobufs and call services. Services enforce business rules like match eligibility, message permissions, and status transitions. Repositories perform GORM queries on MySQL. Auto-migrations run on startup.
- Dependency Injection via Google Wire assembles handlers, services, and repositories from a single provider set for a clean, testable composition.

### Tech stack
- Go (gRPC, Protobuf)
- GORM (MySQL)
- Google Wire (DI)
- JWT (HS256)
- godotenv

### Project layout
- `main.go`: gRPC server bootstrap, auth interceptor, service registration
- `api/`: gRPC handlers (one per service)
- `service/`: business logic
- `repository/`: data access with GORM
- `db/`: GORM models and hooks
- `config/`: environment config and DB initialization
- `di/`: dependency injection via Wire (`wire.go`, generated `wire_gen.go`)
- `proto/v1/`: protobuf definitions and generated code

### Requirements
- Go 1.24+
- MySQL instance accessible from the app

### Configuration (.env)
The server loads environment variables from `.env` (via `godotenv`).

```env
# gRPC
GRPC_PORT=8080

# DB
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=hope

# Auth
JWT_SECRET=your-long-random-secret
GOOGLE_CLIENT_ID=your-google-oauth-client-id
ALLOWED_DOMAINS=example.com,another.com
```

Notes:
- Database DSN used: `user:password@tcp(host:port)/db?parseTime=True&loc=Local`.
- Auto-migrations run on startup for all models in `db/`.

### Run locally
```bash
go run ./...
```

The server listens on `:${GRPC_PORT}` (default `:8080`). gRPC reflection is enabled.

### Authentication
Only `proto.v1.AuthService/Login` is public. All other RPCs require a Bearer token in the metadata header:

```
authorization: Bearer <JWT>
```

Login flow:
1) Client obtains a Google ID token.
2) `Login` verifies token via Google, checks `aud` against `GOOGLE_CLIENT_ID`, ensures `email_verified`, enforces `ALLOWED_DOMAINS`.
3) A user is created if not present; then a backend JWT (HS256, 24h) is returned.

Example (login):
```bash
grpcurl -plaintext \
  -d '{"id_token":"<google-id-token>"}' \
  localhost:8080 proto.v1.AuthService/Login
```

Example (authenticated call):
```bash
grpcurl -plaintext \
  -H "authorization: Bearer <jwt>" \
  -d '{}' \
  localhost:8080 proto.v1.UserService/GetMe
```

### Services and RPCs

All services are under the `proto.v1` package.

- AuthService
  - `Login(LoginRequest) -> LoginResponse` (public)

- UserService
  - `GetMe(GetMeRequest) -> GetMeResponse` (auth)
  - `GetUser(GetUserRequest) -> GetUserResponse` (auth)
  - `UpdateMe(UpdateMeRequest) -> UpdateMeResponse` (auth)
  - `ListUsers(ListUsersRequest) -> ListUsersResponse` (auth)

- RideService
  - Offers
    - `CreateOffer(CreateOfferRequest) -> CreateOfferResponse` (auth)
    - `GetOffer(GetOfferRequest) -> GetOfferResponse` (auth)
    - `UpdateOffer(UpdateOfferRequest) -> UpdateOfferResponse` (auth)
    - `DeleteOffer(DeleteOfferRequest) -> DeleteOfferResponse` (auth)
    - `ListNearbyOffers(ListNearbyOffersRequest) -> ListNearbyOffersResponse` (auth)
    - `ListMyOffers(ListMyOffersRequest) -> ListMyOffersResponse` (auth)
  - Requests
    - `CreateRequest(CreateRequestRequest) -> CreateRequestResponse` (auth)
    - `GetRequest(GetRequestRequest) -> GetRequestResponse` (auth)
    - `UpdateRequestStatus(UpdateRequestStatusRequest) -> UpdateRequestStatusResponse` (auth)
    - `DeleteRequest(DeleteRequestRequest) -> DeleteRequestResponse` (auth)
    - `ListNearbyRequests(ListNearbyRequestsRequest) -> ListNearbyRequestsResponse` (auth)
    - `ListMyRequests(ListMyRequestsRequest) -> ListMyRequestsResponse` (auth)

- MatchService
  - `RequestToJoin(RequestToJoinRequest) -> RequestToJoinResponse` (auth)
  - `AcceptRideRequest(AcceptRideRequestRequest) -> AcceptRideRequestResponse` (auth)
  - `AcceptRequest(AcceptRequestRequest) -> AcceptRequestResponse` (auth)
  - `RejectRequest(RejectRequestRequest) -> RejectRequestResponse` (auth)
  - `CompleteMatch(CompleteMatchRequest) -> CompleteMatchResponse` (auth)
  - `GetMatch(GetMatchRequest) -> GetMatchResponse` (auth)
  - `ListMatchesByRide(ListMatchesByRideRequest) -> ListMatchesByRideResponse` (auth)
  - `ListMatchesByRider(ListMatchesByRiderRequest) -> ListMatchesByRiderResponse` (auth)
  - `ListMyMatches(ListMyMatchesRequest) -> ListMyMatchesResponse` (auth)

- ChatService
  - `SendMessage(SendMessageRequest) -> SendMessageResponse` (auth)
  - `ListMessagesByRide(ListMessagesByRideRequest) -> ListMessagesByRideResponse` (auth)
  - `ListMessagesBySender(ListMessagesBySenderRequest) -> ListMessagesBySenderResponse` (auth)
  - `ListChatsForUser(ListChatsForUserRequest) -> ListChatsForUserResponse` (auth)

- LocationService
  - `UpsertLocation(UpsertLocationRequest) -> UpsertLocationResponse` (auth)
  - `GetLocationByUser(GetLocationByUserRequest) -> GetLocationByUserResponse` (auth)
  - `ListNearby(ListNearbyRequest) -> ListNearbyResponse` (auth)
  - `DeleteMyLocation(DeleteMyLocationRequest) -> DeleteMyLocationResponse` (auth)

### Deep dive: how I implemented each RPC and why

Below is how I designed and implemented each RPC end‑to‑end. I describe the handler (gRPC edge), service (business rules), and repository (DB), and why I made those choices.

#### AuthService
- Login
  - What: Exchange a Google ID token for my backend JWT and create a user record if needed.
  - How: `api/auth_handlers.go` validates payload and calls `service/auth_service.go`:
    - I verify the Google token with Google (`verifyGoogleIDToken` via `https://oauth2.googleapis.com/tokeninfo`).
    - I check `aud == GOOGLE_CLIENT_ID`, ensure `email_verified`, and enforce `ALLOWED_DOMAINS`.
    - I upsert the user through `UserRepository` (create if not found), then issue HS256 JWT with claims `sub`, `email`, `name`.
  - Why: Delegating identity to Google reduces auth surface area. Domain allowlist keeps the product scoped (e.g., campus/company). JWT keeps the server stateless.

#### UserService
- GetMe
  - What: Return the authenticated user.
  - How: I read `user_id` from context (set by `middleware.AuthInterceptor`) and fetch via `UserService.GetUserByID` → `UserRepository.FindByID`.
  - Why: Context‑injected identity avoids passing user IDs over the wire for self‑queries.
- GetUser
  - What: Fetch any user by ID.
  - How: Simple pass‑through to service/repo with not‑found handling.
  - Why: Needed to render counterpart profiles.
- UpdateMe
  - What: Update my profile fields (name, photo_url, geohash).
  - How: Load current user, mutate only provided fields, save via `UserRepository.Update`.
  - Why: Partial updates prevent unintended overwrites and keep the RPC small.
- ListUsers
  - What: Batch fetch a small set of users by IDs.
  - How: Iterate IDs and call `GetUserByID` for each.
  - Why: The list is expected to be small (chat headers, match cards); this keeps logic simple.

#### RideService — Offers
- CreateOffer
  - What: Drivers post an offer (route, time, seats, fare).
  - How: Handler validates required fields and auth; service trims input, enforces future time and positive seats, generates ID, verifies driver exists via `UserRepository`, and persists via `RideOfferRepository.Create`.
  - Why: Validations protect integrity; verifying driver existence catches dangling refs.
- GetOffer
  - How/Why: Lookup by ID via repo; returns `NotFound` if missing. Straightforward read path.
- UpdateOffer
  - What: Partial update (fare/seats/status) by the owner.
  - How: I load current offer, authorize that caller is the driver (in handler), apply only provided fields, and `Save`.
  - Why: Owner‑only updates and partial mutation keep state consistent.
- DeleteOffer
  - How/Why: Owner check in handler, then `Delete` by ID. Prevents unauthorized deletions.
- ListNearbyOffers
  - What: Query offers by `from_geo` geohash prefix.
  - How: Repo uses `LIKE geohash_prefix%` with optional `LIMIT`, ordered by time ASC.
  - Why: Prefix queries are a simple, fast approximation for proximity without a geo index.
- ListMyOffers
  - What: Caller’s offers.
  - How: Read `callerID` from context; repo filters by `driver_id` with optional limit.
  - Why: Common dashboard view for drivers.

#### RideService — Requests
- CreateRequest
  - What: Riders post a request (route, time, seats).
  - How: Service validates, trims, ensures future time and positive seats, sets default `status=active`, generates ID, verifies user exists, then persists via `RideRequestRepository.Create`.
  - Why: Symmetric to offers; keeps model consistent.
- GetRequest
  - How/Why: Lookup by ID; errors map to `NotFound` at the handler.
- UpdateRequestStatus
  - What: Change status of my request.
  - How: Handler authorizes ownership; service trims inputs and updates status via repo.
  - Why: Only request owner should transition their request.
- DeleteRequest
  - How/Why: Owner check in handler; repo delete by ID.
- ListNearbyRequests
  - How/Why: Same geohash prefix approach as offers, ordered by time ASC.
- ListMyRequests
  - How/Why: Caller’s requests filtered by `user_id` with optional limit.

#### MatchService
- RequestToJoin
  - What: Rider asks to join an existing offer.
  - How: Handler sets `rider_id` from context; service validates, loads offer by `ride_id`, copies `driver_id` from the offer, sets `status=requested`, stamps `created_at`, and creates the match.
  - Why: Centralizes driver identity on the server, avoids spoofing.
- AcceptRideRequest
  - What: Driver accepts a rider’s request (creates an offer+match and marks the request matched).
  - How: Service loads the request, checks `active`, prevents self‑accept, synthesizes a new offer for the driver (status `matched`), creates a match with `accepted` status, and updates the original request to `matched`.
  - Why: Supports the inverse flow (driver initiates) while preserving invariants atomically at the service layer.
- AcceptRequest / RejectRequest
  - What: Driver decision on a `requested` match.
  - How: Service enforces caller is the `driver_id`, checks current `status=requested`, then moves to `accepted` or `rejected`.
  - Why: Prevents riders from self‑approving; keeps a clear state machine.
- CompleteMatch
  - What: Mark a match completed.
  - How: Service updates status to `completed`.
  - Why: Minimal flow completion; permissioning is intentionally simple here.
- GetMatch / ListMatchesByRide / ListMatchesByRider / ListMyMatches
  - How/Why: Standard reads. `ListMyMatches` uses caller identity for convenience.

#### ChatService
- SendMessage
  - What: Send a message scoped to a ride.
  - How: Service generates ID and timestamp, then authorizes the sender by loading matches for the ride and ensuring the sender is either the rider or driver on a match with `accepted` or `completed` status; writes via `ChatMessageRepository.Create`.
  - Why: Enforces that only matched participants can chat; prevents arbitrary ride spam.
- ListMessagesByRide / ListMessagesBySender / ListChatsForUser
  - How: Repos filter by ride, sender, or user with `before` and `limit`, ordered by timestamp DESC.
  - Why: Pagination‑ready access patterns without complex indices.

#### LocationService
- UpsertLocation
  - What: Save my last known location and geohash.
  - How: Handler pulls `user_id` from context and builds the location; service validates lat/lon bounds, sets `updated_at`, and calls repo upsert (MySQL `ON CONFLICT`‑style via GORM clauses).
  - Why: Idempotent writes let clients update frequently without worrying about record existence.
- GetLocationByUser
  - How/Why: Simple lookup by `user_id`, with a clear `NotFound` mapping.
- ListNearby
  - How/Why: Geohash prefix search with optional limit, ordered by `updated_at` DESC to show freshest first.
- DeleteMyLocation
  - How/Why: Remove my row; useful for privacy or sign‑out flows.

### Example RPC payloads

- RideService/CreateOffer
```bash
grpcurl -plaintext \
  -H "authorization: Bearer <jwt>" \
  -d '{
    "from_geo":"u4pruy",
    "to_geo":"u4pruv",
    "fare": 10.5,
    "time": {"seconds": 1893456000},
    "seats": 3
  }' \
  localhost:8080 proto.v1.RideService/CreateOffer
```

- LocationService/UpsertLocation
```bash
grpcurl -plaintext \
  -H "authorization: Bearer <jwt>" \
  -d '{"latitude":28.61, "longitude":77.20, "geohash":"tsq4"}' \
  localhost:8080 proto.v1.LocationService/UpsertLocation
```

### Data models (GORM)
- `User`: id, name, email (unique), photo_url, geohash, last_seen; has one `UserLocation`
- `RideOffer`: id, driver_id, from_geo, to_geo, fare, time, seats, status
- `RideRequest`: id, user_id, from_geo, to_geo, time, seats, status
- `Match`: id, rider_id, driver_id, ride_id, status, created_at
- `ChatMessage`: id, ride_id, sender_id, content, timestamp
- `Review`: id, ride_id, from_user_id, to_user_id, score, comment, created_at
- `UserLocation`: user_id, latitude, longitude, geohash, updated_at

Auto-migrations run on startup for all the above.

### Dependency Injection
Constructed via `di/wire.go` and generated `di/wire_gen.go`. If you change providers or constructor signatures, regenerate with Wire.


