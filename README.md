# **🎟 Event Ticket Booking Website**

A full-stack **Event Ticket Booking Website** built with **Go (Kratos Framework)** for the backend and **React + Tailwind CSS** for the frontend.
It allows users to **signup, signin, view events, and book tickets**, with data stored in **PostgreSQL** and managed using **Docker** containers.

---

## **🚀 Features**
 
* 🔐 **User Authentication**

  * Signup
  * Sign-in / Login
  * JWT or session-based access (depending on backend)
* 🎫 **Browse available events**
* 📦 **Docker containerized backend**
* 🧰 **RESTful API developed using Go + Kratos**
* 🗄 **PostgreSQL database**
* 🖥 **Frontend built with React**
* 🎨 **Tailwind CSS UI + Responsive design**
* ⚡ **Vite dev server for fast frontend development**

---

## **🛠 Tech Stack**

### **Frontend**

* React (Vite)
* Tailwind CSS
* React Router DOM
* Fetch API for backend requests

### **Backend**

* Go (Kratos Framework)
* REST APIs
* Docker & Docker Compose

### **Database**

* PostgreSQL

---

## **📁 Project Structure**

```
miniproject/
 ┣ backend/ (Go + Kratos service)
 ┣ frontend/
 ┃ ┣ src/
 ┃ ┃ ┣ pages/
 ┃ ┃ ┃ ┣ SignUp.jsx
 ┃ ┃ ┃ ┣ SignIn.jsx
 ┃ ┃ ┃ ┗ Events.jsx
 ┃ ┃ ┣ App.jsx
 ┃ ┃ ┣ main.jsx
 ┃ ┣ index.css
 ┃ ┣ tailwind.config.js
 ┃ ┗ vite.config.js
 ┗ README.md
```

---

## **▶️ Running the Project**

### **Backend Setup**

```bash
cd backend
docker compose up --build
```

Backend will run at:

```
http://localhost:8000
```

---

### **Frontend Setup**

```bash
cd frontend
npm install
npm run dev
```

Frontend will run at:

```
http://localhost:5173
```

---

## **🌐 API Endpoints**

| Method | Endpoint       | Description                |
| ------ | -------------- | -------------------------- |
| `POST` | `/users`       | Create a new user (Signup) |
| `POST` | `/users/login` | Authenticate user (Signin) |
| `GET`  | `/events`      | Fetch events list          |

---

## **📸 UI Preview**

(Add screenshots here)

---

## **💡 Future Enhancements**

* Booking confirmation and ticket generation
* Admin dashboard to manage events
* Online mock payment integration
* Email OTP verification

---

## **👨‍💻 Author**

**Mano**
Full-Stack Developer (Go / React)

🔗 LinkedIn: *(Add your link here)*
📧 Email: *(Add email here)*

---

## **⭐ Show Support**

If you like this project, please give a **⭐ star** in GitHub — it motivates further development!



