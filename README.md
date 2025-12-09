E-Learning Platform APIThis repository contains the backend API for a robust e-learning platform, built using Go (Golang), the Gin web framework, and GORM for ORM, connecting to a MySQL database.The platform includes a specialized Lecture Upload Service with direct integration to the YouTube Data API v3 for handling large media files securely and efficiently.🚀 Key FeaturesRole-Based Access Control (RBAC): Separate roles for Admin, Teacher, and Student (implicit).Departmental Authorization: Teachers are linked to specific departments and can only create/manage courses within their assigned department.Course Management: CRUD operations for courses, linked to specific departments and teachers.Secure YouTube Integration: Asynchronous lecture video uploads directly to YouTube, tracking status and metadata in the database.RESTful API: Clean, versioned, and documented endpoints.🛠️ Technology StackComponentTechnologyRoleBackend LanguageGo (Golang)Core application logic and performance.Web FrameworkGinRouting and middleware handling.Database ORMGORMDatabase abstraction and migrations.DatabaseMySQLPersistent data storage.External ServiceYouTube Data API v3Secure video hosting and streaming.📦 Project StructureThe project follows a standard Go project layout with separation of concerns:/
├── controllers/       # Handles incoming HTTP requests and responses.
├── models/            # Database structures (GORM structs: Course, Lecture, User, TeacherProfile).
├── service/           # External API interactions (e.g., YouTubeService).
├── database/          # Database connection and migration logic.
├── middleware/        # JWT authentication and authorization checks.
├── main.go            # Application entry point and router setup.
└── .env               # Configuration for database and API keys.
⚙️ Setup and Installation1. PrerequisitesYou must have the following installed:Go (1.18+)MySQL ServerGit2. Get the CodeBashgit clone https://github.com/your-username/your-project.git
cd your-project
3. Configure Environment VariablesCreate a file named .env in the project root. This file must contain your database credentials and, critically, your YouTube OAuth credentials.Code snippet# --- Database Configuration ---
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_root_password
DB_NAME=elearning_db

# --- JWT Configuration ---
JWT_SECRET=a_very_secure_secret_key_for_jwt

# --- YouTube API Configuration (CRITICAL FOR UPLOADS) ---
# Client ID and Secret obtained from Google Cloud Console (Desktop App type)
YOUTUBE_CLIENT_ID="<your_client_id>"
YOUTUBE_CLIENT_SECRET="<your_client_secret>"
# Long-lived token generated via the OAuth flow (needed for server-side uploads)
YOUTUBE_REFRESH_TOKEN="<your_long_lived_refresh_token>" 
4. Run the ApplicationBash# Install dependencies
go mod tidy

# Run the application (This will also run GORM AutoMigrate for all models)
go run main.go
🔐 API Endpoints (Core Routes)The API base path is generally /api.FunctionalityMethodEndpointAuthorizationAuthPOST/api/auth/loginPublicCoursesPOST/api/coursesTeacher/AdminLecturesPOST/api/courses/:courseId/lectures/uploadCourse Owner/AdminLecturesGET/api/lectures/:idAuthenticatedUser ManagementGET/api/users/:idAdmin/Self🎥 Lecture Upload FlowA Teacher sends a POST request to /api/courses/:courseId/lectures/upload with the video file (multipart/form-data).The server verifies the teacher owns the course.A new Lecture record is created with Status: "uploading".The server returns an immediate 202 Accepted response to prevent client timeouts.A Goroutine handles the secure, asynchronous upload to YouTube using the stored YOUTUBE_REFRESH_TOKEN.Once the YouTube upload is complete, the Lecture record is updated with the YouTubeVideoID and YouTubeURL, and Status is changed to "ready".🤝 ContributingThis project is currently under active development. Contributions, suggestions, and feedback are highly encouraged!Fork the repository.Create a new feature branch (git checkout -b feature/AmazingFeature).Commit your changes (git commit -m 'Add some AmazingFeature').Push to the branch (git push origin feature/AmazingFeature).Open a Pull Request.
