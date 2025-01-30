# Verify Influencers Backend

This is a backend application for verifying influencer claims, built using Go and the Fiber framework. The app fetches claims from health-related content and processes them based on the source (currently only Twitter is supported).

## Requirements

- Go 1.18 or higher
- PostgreSQL database
- `.env` file for environment variables

## Setup

### 1. Clone the Repository

Clone the repository to your local machine:


git clone https://github.com/AbdulmalikRaji/verify-influencers-backend/
cd verify-influencers-backend

### 2. Install Dependencies
Make sure Go is installed on your machine. Install the necessary Go modules:

bash
```
go mod tidy
```

### 3. Configure the Database
Ensure you have PostgreSQL running. Create a database and configure the connection details in the .env file. Example:

.env
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=verify_influencers
DB_SSL_MODE=disable
```

### 4. Run the Application
Start the Go application:

bash
```
go run cmd/main.go
```
This will start the server on port 3000.

### 5. Accessing the API
Once the server is running, you can send requests to the API.

#### Endpoint: `/api/v1/claims`
- **Method**: `GET`
- **Description**: Retrieves influencer claims based on the query parameters.

#### Query Parameters:
- `username` (string): The username of the influencer (required).
- `source` (int): The source of the claim. Must be `1` for Twitter (currently the only supported source). `2` is reserved for podcasts (not implemented yet).
- `start_date` (string, format: `yyyy-mm-dd`): The start date for filtering claims (required).
- `end_date` (string, format: `yyyy-mm-dd`): The end date for filtering claims (required).

#### Example Request:
```http
GET http://localhost:3000/api/v1/claims?username=johndoe&source=1&start_date=2025-01-01&end_date=2025-01-31
```
#### Example Response:
json
```
[
    {
        "claim_id": 1,
        "username": "johndoe",
        "claim": "This product cures headaches.",
        "source": 1,
        "date": "2025-01-10"
    },
    {
        "claim_id": 2,
        "username": "johndoe",
        "claim": "This product improves sleep.",
        "source": 1,
        "date": "2025-01-15"
    }
]
```
### 6. Database Migrations
Make sure the database schema is set up before running the app. You can use the gorm migration feature or manually set up the database by running the necessary SQL commands to create tables.

### 7. Additional Notes
The source parameter can only be 1 for Twitter at the moment.
The podcasts source (2) is reserved but not yet implemented.
Ensure that the environment variables are set correctly in the .env file for the database connection to work.

### License
This project is licensed under the MIT License.
