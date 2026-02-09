# Lawbook - Virtual Moot Court Platform

A modern, role-based authentication system for a virtual moot court platform built with Go, following Alex Edwards' "Let's Go" best practices.

## 🎯 Project Overview

Lawbook is an interactive legal education platform that connects:
- **Students**: Practice moot courts and build portfolios
- **Lawyers**: Sharpen skills and showcase expertise to recruiters  
- **Recruiters**: Evaluate candidates based on AI-powered performance metrics

## ✨ Features

### Authentication & Authorization
- ✅ Role-based authentication (Student, Lawyer, Recruiter)
- ✅ Secure session management with SCS
- ✅ CSRF protection
- ✅ Password hashing with bcrypt
- ✅ Email validation
- ✅ Account activation/deactivation

### User Roles
- **Students**: Access moot court practice, track progress
- **Lawyers**: Practice sessions, performance analytics, recruiter visibility
- **Recruiters**: View lawyer evaluations, search candidates

### Virtual Moot Court (Coming Soon)
- AI-powered judge and opponents
- Multiple session types (solo, dual player, trio)
- Performance evaluation and scoring
- Various case types (Constitutional, Criminal, Civil, etc.)

## 📁 Project Structure

```
lawbook/
├── cmd/
│   └── web/
│       ├── main.go           # Application entry point
│       ├── handlers.go       # HTTP handlers
│       ├── routes.go         # Route definitions
│       ├── middleware.go     # Middleware (auth, role checks)
│       ├── helpers.go        # Helper functions
│       ├── templates.go      # Template management
│       └── context.go        # Context keys
├── internal/
│   ├── models/
│   │   ├── users.go          # User model
│   │   ├── sessions.go       # Session model
│   │   ├── models.go         # Model wrapper
│   │   └── errors.go         # Error definitions
│   └── validator/
│       └── validator.go      # Form validation
├── ui/
│   ├── html/
│   │   ├── base.tmpl.html    # Base template
│   │   ├── pages/            # Page templates
│   │   └── partials/         # Partial templates
│   └── static/
│       ├── css/
│       │   └── main.css      # Styles
│       └── js/
│           └── main.js       # JavaScript
├── migrations/
│   ├── 001_initial.up.sql    # Database schema
│   └── 001_initial.down.sql  # Migration rollback
├── go.mod
├── Makefile
└── README.md
```

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- MySQL 8.0 or higher

### Installation

1. **Clone or extract the project**
   ```bash
   cd lawbook
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   ```bash
   # Create database and tables
   mysql -u root -p < migrations/001_initial.up.sql
   ```

4. **Configure environment**
   ```bash
   # Set your database connection string
   export LAWBOOK_DB_DSN="root:yourpassword@tcp(localhost:3306)/lawbookauth?parseTime=true"
   ```

5. **Run the application**
   ```bash
   make run
   # OR
   go run ./cmd/web
   ```

6. **Access the application**
   ```
   http://localhost:4000
   ```

## 🔧 Configuration

### Database Connection
Set the `LAWBOOK_DB_DSN` environment variable:
```bash
export LAWBOOK_DB_DSN="username:password@tcp(host:port)/database?parseTime=true"
```

### Server Port
Default port is `:4000`. Change with the `-addr` flag:
```bash
go run ./cmd/web -addr=":8080"
```

## 📝 Available Make Commands

```bash
make help          # Show available commands
make run           # Run the application
make build         # Build the binary
make test          # Run tests
make clean         # Clean build artifacts
make migrate-up    # Run database migrations
make migrate-down  # Rollback migrations
make deps          # Download dependencies
make fmt           # Format code
make vet           # Run go vet
make dev           # Run in development mode
```

## 🗃️ Database Schema

### Core Tables
- **users**: User accounts with role-based access
- **sessions**: Session management
- **student_profiles**: Student-specific data
- **lawyer_profiles**: Lawyer-specific data
- **recruiter_profiles**: Recruiter-specific data

### Moot Court Tables
- **moot_sessions**: Virtual court sessions
- **session_participants**: Session participants
- **performance_evaluations**: AI-generated evaluations

## 🔐 Security Features

- **Password Security**: bcrypt hashing (cost 12)
- **Session Security**: Secure, HTTP-only cookies with 12-hour expiry
- **CSRF Protection**: Token-based CSRF prevention
- **SQL Injection**: Prepared statements throughout
- **XSS Protection**: Template auto-escaping
- **Secure Headers**: CSP, X-Frame-Options, etc.

## 🎨 User Interface

### Pages
- **Public**: Home, About, Login, Signup
- **Student**: Dashboard, Moot Court Setup
- **Lawyer**: Dashboard, Moot Court Setup, Performance Analytics
- **Recruiter**: Dashboard, Candidate Search
- **Shared**: Account Settings

### Responsive Design
Mobile-friendly design with CSS Grid and Flexbox

## 🚧 Roadmap

### Phase 1: Authentication ✅ (Complete)
- [x] User registration with role selection
- [x] Login/logout functionality
- [x] Role-based access control
- [x] Session management
- [x] Account pages

### Phase 2: Moot Court (In Progress)
- [ ] AI integration (OpenAI/Anthropic)
- [ ] Real-time session interface
- [ ] Audio/video support
- [ ] Performance evaluation system
- [ ] Session recording and playback

### Phase 3: Advanced Features
- [ ] Profile management
- [ ] Performance analytics dashboard
- [ ] Recruiter search and filtering
- [ ] Email notifications
- [ ] Advanced case scenarios
- [ ] Collaborative features

## 📚 Technologies Used

- **Backend**: Go 1.21
- **Database**: MySQL 8.0
- **Session**: alexedwards/scs
- **Router**: julienschmidt/httprouter
- **Middleware**: justinas/alice
- **CSRF**: justinas/nosurf
- **Templates**: html/template
- **Forms**: go-playground/form

## 🤝 Contributing

This is a learning project adapted from Alex Edwards' "Let's Go" book. The authentication system follows professional Go patterns and best practices.

## 📖 Learning Resources

- [Let's Go by Alex Edwards](https://lets-go.alexedwards.net/)
- [Let's Go Further](https://lets-go-further.alexedwards.net/)
- [Go Documentation](https://go.dev/doc/)

## 🔗 Original Inspiration

Based on the twitbox project structure, adapted for role-based legal education platform.

## 📄 License

This project is for educational purposes.

## 🙋 Support

For issues or questions related to the authentication system:
1. Check the README
2. Review the code comments
3. Consult Alex Edwards' "Let's Go" book

## 🎯 Next Steps

1. **Set up your database** using the provided migration
2. **Run the application** with `make run`
3. **Create an account** by visiting `/user/signup`
4. **Explore the dashboards** based on your selected role
5. **Start building** the moot court functionality!

---

Built with ❤️ for legal education in India
