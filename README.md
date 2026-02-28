# 🎓 Digital Evaluation System

> **A Next-Generation, Secure, and Automated Platform for Academic Evaluations.**

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Backend-Go-00ADD8.svg?logo=go&logoColor=white)
![Python](https://img.shields.io/badge/Validator-Python-3776AB.svg?logo=python&logoColor=white)
![React](https://img.shields.io/badge/Frontend-React-61DAFB.svg?logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-336791.svg?logo=postgresql&logoColor=white)
![IPFS](https://img.shields.io/badge/Storage-IPFS-65C2CB.svg?logo=ipfs&logoColor=white)
![jq](https://img.shields.io/badge/Tools-jq-3776AB.svg?logo=jq&logoColor=white)
![openssl](https://img.shields.io/badge/Tools-openssl-3776AB.svg?logo=openssl&logoColor=white)
![make](https://img.shields.io/badge/Tools-make-3776AB.svg?logo=make&logoColor=white)
![MSYS2](https://img.shields.io/badge/Tools-MSYS2-3776AB.svg?logo=msys2&logoColor=white)

---

## 📖 Overview

The **Digital Evaluation System** is a robust, distributed platform designed to modernize the academic evaluation process. It streamlines the workflow from answer script assignment to result declaration, ensuring transparency, security, and efficiency.

By leveraging **Go** for high-performance orchestration, **Python (FastAPI + GenAI)** for intelligent automated validation, and a modern **React** frontend, this system facilitates a seamless experience for students, evaluators, and academic authorities.

---

## ✨ Key Features

- **🔐 Role-Based Access Control (RBAC):** converting granular permissions for:
    - **Students:** View results and submit scripts.
    - **Evaluators:** Request subjects, evaluate scripts, and submit marks.
    - **Examiners/Authorities:** Oversee the process, approve requests, and release results.
    - **Admins:** Manage users and system configurations.

- **🤖 Intelligent Validation:**
    - Automated script processing using **Python, OpenCV, and PyMuPDF**.
    - AI-powered insights and validation using **Google GenAI**.

- **⚡ High-Performance Core:**
    - Backend services built with **Go** for speed and concurrency.
    - Secure **JWT Authentication** and **mTLS** communication between nodes.

- **📊 Efficient Workflow Management:**
    - **Evaluation Requests:** Evaluators can request specific courses/semesters.
    - **Script Assignment:** Automated or manual assignment of scripts to evaluators.
    - **Audit Logging:** Comprehensive tracking of all critical system actions.

- **🎨 Modern User Interface:**
    - built with **React, Vite, and Tailwind CSS** for a responsive and intuitive user experience.

---

## 🏗️ System Architecture and Structure

The system follows a microservices-inspired architecture:

1.  **Core Service (`go-node`)**:
    - Written in **Go**.
    - Handles API requests, user authentication (JWT), and business logic.
    - Manages the **PostgreSQL** database.
    - Orchestrates communication with the Validator service.

2.  **Validator Service (`python-validator`)**:
    - Written in **Python (FastAPI)**.
    - Performs heavy-lifting tasks: PDF processing, Image recognition (OpenCV), and AI analysis (GenAI).

3.  **Frontend (`digital-eval-ui`)**:
    - **React** application served via Vite.
    - Interacts with the Core Service via REST APIs.

4.  **Infrastructure (`infra`)**:
    - **PostgreSQL** for relational data (Users, Results, Audit Logs).
    - **IPFS** (InterPlanetary File System) integration for decentralized storage of artifacts (optional/planned).

You can view the structure of the project here [STRUCTURE](STRUCTURE.md).

---

## 🛠️ Technology Stack

| Component        | technologies                                         |
| :--------------- | :--------------------------------------------------- |
| **Backend Core** | Go, Gorilla Mux, JWT, Logrus                         |
| **AI Validator** | Python 3.12+, FastAPI, OpenCV, PyMuPDF, Google GenAI |
| **Frontend**     | React 18, TypeScript, Tailwind CSS, Vite             |
| **Database**     | PostgreSQL                                           |
| **Tools**        | Make, Bash, IPFS, jq                                 |

---

## 🚀 Getting Started

### Prerequisites

Ensure you have the following installed:

- [Go](https://go.dev/dl/) (v1.25+)
- [Python](https://www.python.org/downloads/) (v3.12+)
- [Node.js](https://nodejs.org/) (v20+) & npm
- [PostgreSQL](https://www.postgresql.org/download/) (v18)
- [IPFS](https://github.com/ipfs/kubo/releases) (Latest)
- [jq](https://stedolan.github.io/jq/download/) (Latest)
- [openssl](https://openssl-library.org/source/) (Latest)
- [make](https://www.msys2.org/) (Latest)

### 📥 Installation & Setup

1.  **Clone the Repository**

    ```bash
    git clone https://github.com/vishwanathdvgmm/digital-eval-system.git
    cd digital-eval-system
    ```

2.  **Run Local Setup Script**
    This script checks for required tools like `psql`, `ipfs`, `openssl` and `jq`.

    ```bash
    chmod +x ./tools/local_setup.sh
    ./tools/local_setup.sh
    ```

    **Note:** after installing add it to Path.

3.  **Generate Certificates**
    Generate TLS certificates and JWT keys for secure communication.

    (You need to install make tool for this. Better run this in MSYS2 MSYS terminal. Open MSYS2 MSYS terminals seprately.)

    Link to download MSYS2
    - [MSYS2](https://www.msys2.org/)

    ```msys2
    cd digital-eval-system/
    make keys
    ```

    It will create the necessary keys.

4.  **Backend Setup**
    - **Config file and .env file setup:**
        - Open `services/go-node/configs`
        - In that set all the necessary fields.
        - Don't change these `ports`, `host`, `enabled`, `url`, `issuer` field.
        - For postgresql create a database and put it in this field `dsn` with the password and user name same.
        - You should create a new user and password for that user.
        - Open .env and make necessary chnages.
    - **Creating and applying migrations for postgresql:**
        - Open terminal and run this:

        ```bash
        psql -U admin -d digital_eval
        ```

        - Enter the password when prompted and when logged in run the below commands.

        ```sql
        CREATE EXTENSION IF NOT EXISTS pgcrypto;
        \q
        ```

        - In the terminal run this.

        ```bash
        cd infra/migrations/postgres/
        psql -U admin -d digital_eval -f run_all.sql
        ```

        - Enter the password when prompted.
        - Now the migrations are applied.

        - Now insert the users:

        ```bash
        psql -U admin -d digital_eval
        ```

        - Enter the password and run this.

        ```sql
        INSERT INTO users(user_id, email, role, password_hash)
        VALUES
        ('admin_1', 'admin@example.com', 'admin', crypt('admin123', gen_salt('bf'))),
        ('authority_1', 'authority@example.com', 'authority', crypt('auth123', gen_salt('bf'))),
        ('examiner_1', 'examiner@example.com', 'examiner', crypt('exam123', gen_salt('bf'))),
        ('evaluator_1', 'evaluator@example.com', 'evaluator', crypt('eval123', gen_salt('bf'))),
        ('student_1', 'student@example.com', 'student', crypt('stud123', gen_salt('bf')));
        \q
        ```

    - **Go Node:**
        ```bash
        cd services/go-node
        go mod download
        go mod tidy
        ```
    - **Python Validator:**
        ```bash
        cd services/python-validator
        python -m venv .venv
        ```
        Restart the terminal:
        ```bash
        pip install -r requirements.txt
        ```

### ▶️ Running the Application

_Note: The project uses a Makefile for convenience._

1.  **Start Services (Draft/Phase 0)**
    Currently, the `make start` command is a placeholder. You may need to start services individually during development:

    Open one terminal run this and close it:

    ```bash
    ipfs init
    ```

    **Run Go Service: (Terminal 1)**

    ```bash
    cd services/go-node
    go clean -cache # if needed (optional)
    cd ../..
    bash scripts/build.sh
    cd services/go-node
    ./node.exe -config configs/config.yaml # This runs both forntend and backend.
    ```

2.  In GUI login with the admin credentials:
    - Email: admin@example.com
    - Password: admin123

    - After that start the services `ipfs`, `python-validator` and `python-extractor`.
    - Then logout.
    - Now follow the below tutorial video.

## 🎥 Tutorial Video

## 🎥 Tutorial Video

[![Watch the tutorial](assets/thumbnail.png)](https://private-user-images.githubusercontent.com/133391353/556414211-bc0584d2-37d1-45f9-96a2-9e4a1782103c.mp4?jwt=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NzIyOTU3NTIsIm5iZiI6MTc3MjI5NTQ1MiwicGF0aCI6Ii8xMzMzOTEzNTMvNTU2NDE0MjExLWJjMDU4NGQyLTM3ZDEtNDVmOS05NmEyLTllNGExNzgyMTAzYy5tcDQ_WC1BbXotQWxnb3JpdGhtPUFXUzQtSE1BQy1TSEEyNTYmWC1BbXotQ3JlZGVudGlhbD1BS0lBVkNPRFlMU0E1M1BRSzRaQSUyRjIwMjYwMjI4JTJGdXMtZWFzdC0xJTJGczMlMkZhd3M0X3JlcXVlc3QmWC1BbXotRGF0ZT0yMDI2MDIyOFQxNjE3MzJaJlgtQW16LUV4cGlyZXM9MzAwJlgtQW16LVNpZ25hdHVyZT04YzQ2OGUwNDFiNDI3ZjlmMGJlOTUxOTBiM2M3ZDI4NDM4NzZiNTRlOTVmMjRiOTc3MWM3NDAwMWU4NjRkZTBmJlgtQW16LVNpZ25lZEhlYWRlcnM9aG9zdCJ9.eDMeYkvoJk2nq7uyiFD1w1EK5V9jpXY7k3GsULymlHU)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

### CONTRIBUTIONS BY:

[VISHWANATH M M](https://github.com/vishwanathdvgmm)<br>
[SWAYAM R VERNEKAR](https://github.com/ver1619)<br>
[SAMPATKUMAR H ANGADI](https://github.com/Raone2005)<br>
[PRAJWAL M](https://github.com/Prajwalpraju17)<br>
