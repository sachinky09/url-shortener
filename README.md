# 🔗 LinkShrink – URL Shortener

A modern URL shortener built with **Go**, **React**, and **Supabase**.  
Shorten long links, track click counts, and share clean URLs instantly.

---

## 🚀 Features
- Shorten long URLs into clean short codes
- Track number of clicks per shortened link
- Modern frontend with **React + Tailwind (white & purple theme)**
- Backend powered by **Go (net/http)**
- **Supabase** as database & storage
- Copy-to-clipboard functionality
- Production-ready CORS and API setup

---

## 📂 Project Structure
```
url-shortener/
├── backend/        # Go server
│   ├── main.go     # Backend entrypoint
│   ├── go.mod      # Go module file
│   ├── go.sum
│   └── .env        # Supabase keys
├── frontend/       # React app
│   ├── src/
│   │   ├── components/
│   │   │   └── Shortener.jsx
│   │   └── App.css
│   └── package.json
└── README.md
```

---

## 🛠️ Setup Instructions

### 1. Clone the repo
```bash
git clone https://github.com/yourusername/url-shortener.git
cd url-shortener
```

### 2. Backend (Go)
```bash
cd backend
go mod tidy
```

Create a `.env` file in `backend/`:
```env
SUPABASE_URL=your-supabase-url
SUPABASE_KEY=your-supabase-key
```

Run the server:
```bash
go run main.go
```
Backend will run on: **http://localhost:8080**

---

### 3. Database (Supabase)
Create the following tables in Supabase SQL editor:
```sql
create table urls (
  id serial primary key,
  code text unique,
  long_url text not null,
  clicks int default 0
);
```

---

### 4. Frontend (React)
```bash
cd frontend
npm install
npm start
```
Frontend will run on: **http://localhost:3000**

---

## ⚡ Usage
1. Enter a long URL in the input field.  
2. Get a short link instantly.  
3. Copy and share.  
4. When someone clicks the short link, they’ll be redirected and the click will be tracked in Supabase.
