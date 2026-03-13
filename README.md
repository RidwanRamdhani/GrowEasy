# GrowEasy

GrowEasy is an intelligent agricultural assistant backend designed to empower farmers with data-driven insights. By combining real-time weather data, soil analysis, machine learning predictions, and Generative AI (Gemini), GrowEasy provides actionable recommendations for crop selection and farming practices.

## Key Features

*   **Crop Prediction:** utilizing Machine Learning models to predict the most suitable crops based on specific soil and weather conditions.
*   **AI Agronomist:** Powered by **Google Gemini**, providing personalized farming advisory summaries (Fertilizer recommendations, Risk assessment, Planting schedules).
*   **Weather Integration:** Real-time and historical weather data fetches (Temperature, Rainfall, Humidity) via **Open-Meteo**.
*   **Soil Data Analysis:** Integration with soil data services **(SoilGrids)** to analyze soil properties for better decision making.
*   **Interactive Chat:** Context-aware chat interface allowing users to ask follow-up questions about their specific field analysis.
*   **Secure Authentication:** User registration and login protected by JWT authentication.
*   **History Tracking:** Saves analysis history for users to track field conditions over time.

## Tech Stack

*   **Language:** [Go (Golang)](https://go.dev/) v1.25+
*   **Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
*   **Database:** PostgreSQL (using [GORM](https://gorm.io/))
*   **AI Integration:** [Google Gemini API](https://ai.google.dev/)
*   **ML Integration:** Custom ML Service (Prediction) https://github.com/Rayya12/Hackathon-AI
*   **Authentication:** JWT (JSON Web Tokens)
*   **External APIs:**
    *   Open-Meteo (Weather)
    *   SoilGrids (Soil Data)

## Related Projects

*   **Frontend UI:** GrowEasy-UI (https://github.com/RidwanRamdhani/GrowEasy-UI) 
*   **ML Service:** Custom ML Service (Prediction) https://github.com/Rayya12/Hackathon-AI

## Prerequisites

Before running the application, ensure you have the following installed:

*   **Go** (version 1.25 or higher)
*   **PostgreSQL** database instance
*   **Git**

## Configuration

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/groweasy.git
    cd groweasy
    ```

2.  Create a `.env` file in the root directory. You can use a template or add the following variables:

    ```env
    # Server Configuration
    PORT=8080

    # Database Configuration
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=groweasy_db
    DB_PORT=5432

    # JWT Secret
    JWT_SECRET=your_super_secret_key

    # Google Gemini API Key
    GEMINI_API_KEY=your_gemini_api_key
    GIN_MODE=release
    GEMINI_MODEL=gemini-3.1-flash-lite-preview
    # you can change the model to Pro if you have the Gemini Pro/Paid API key
    ```

## Running the Application

  1.  **Install Dependencies:**
     ```bash
     go mod download
     ```

  2.  **Database Setup:**
     Before running the application, create the PostgreSQL database specified in your .env file (default: `groweasy_db`).
     The application automatically handles database migrations using GORM on startup. Ensure your PostgreSQL database is running and accessible.

  3.  **Run the Server:**
     ```bash
     go run main.go
     ```
     The server will start on port `8080`
  

## API Documentation

### Public Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/auth/register` | Register a new user |
| `POST` | `/api/auth/login` | Login and receive JWT token |

### Protected Endpoints (Requires Bearer Token)

All endpoints below require the `Authorization` header: `Bearer <token>`

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/weather` | Fetch weather data for coordinates |
| `POST` | `/api/soil` | Fetch soil data for coordinates |
| `POST` | `/api/predict` | **Core Feature:** Run full analysis (Weather + Soil + ML + AI Summary) |
| `GET` | `/api/history` | Get user's past analysis history |
| `POST` | `/api/chat` | Chat with AI about the latest analysis |
| `GET` | `/api/chat/history` | Get current chat session history (add `?all=true` for all sessions grouped by session) |
| `POST` | `/api/chat/reset` | Reset the current chat session context |

## Project Structure

```
GrowEasy/
├── config/         # Database and app configuration
├── dto/            # Data Transfer Objects (Request/Response structs)
├── handlers/       # HTTP Request handlers (Controllers)
├── middleware/     # Gin middleware (Auth, CORS, etc.)
├── models/         # Database models (GORM structs)
├── services/       # Business logic and external integrations
│   ├── integration/ # Clients for Gemini, Weather, Soil, ML
├── utils/          # Helper functions (Hash, JWT)
├── main.go         # Application entry point & Routing
└── go.mod          # Dependencies
```


## License

This project is licensed under the GPL-3.0 license - see the [LICENSE](LICENSE) file for details.
