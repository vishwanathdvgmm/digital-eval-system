# File Tree: digital-eval-system

**Root Path:** `g:\digital-eval-system`

```
├── 📁 assets
│   └── 🎬 tutorial.mp4
├── 📁 data
│   ├── 📁 boltdb
│   │   └── 📄 blocks.db
│   └── 📁 uploads
├── 📁 digital-eval-ui
│   ├── 📁 src
│   │   ├── 📁 api
│   │   │   ├── 📄 admin.ts
│   │   │   ├── 📄 auth.ts
│   │   │   ├── 📄 authority.ts
│   │   │   ├── 📄 evaluator.ts
│   │   │   ├── 📄 examiner.ts
│   │   │   ├── 📄 http.ts
│   │   │   └── 📄 student.ts
│   │   ├── 📁 auth
│   │   │   └── 📄 AuthProvider.tsx
│   │   ├── 📁 components
│   │   │   ├── 📁 authority
│   │   │   │   └── 📄 RequestRow.tsx
│   │   │   ├── 📁 evaluator
│   │   │   │   └── 📄 ScriptCard.tsx
│   │   │   ├── 📁 examiner
│   │   │   │   └── 📄 ScriptUploadItem.tsx
│   │   │   ├── 📁 navigation
│   │   │   │   ├── 📄 AuthorityNav.tsx
│   │   │   │   ├── 📄 EvaluatorNav.tsx
│   │   │   │   ├── 📄 ExaminerNav.tsx
│   │   │   │   ├── 📄 RoleMenu.tsx
│   │   │   │   ├── 📄 StudentNav.tsx
│   │   │   │   └── 📄 TopbarUser.tsx
│   │   │   ├── 📁 student
│   │   │   │   └── 📄 ResultCard.tsx
│   │   │   ├── 📄 Button.tsx
│   │   │   ├── 📄 Card.tsx
│   │   │   ├── 📄 FileUploader.tsx
│   │   │   ├── 📄 Input.tsx
│   │   │   ├── 📄 Loader.tsx
│   │   │   ├── 📄 Navbar.tsx
│   │   │   ├── 📄 ScriptUploadItem.tsx
│   │   │   ├── 📄 Select.tsx
│   │   │   ├── 📄 Sidebar.tsx
│   │   │   └── 📄 Table.tsx
│   │   ├── 📁 context
│   │   │   └── 📄 AuthContext.tsx
│   │   ├── 📁 hooks
│   │   │   └── 📄 useAuth.ts
│   │   ├── 📁 layouts
│   │   │   ├── 📄 AdminLayout.tsx
│   │   │   ├── 📄 AuthLayout.tsx
│   │   │   ├── 📄 AuthorityLayout.tsx
│   │   │   ├── 📄 EvaluatorLayout.tsx
│   │   │   ├── 📄 ExaminerLayout.tsx
│   │   │   ├── 📄 GlobalLayout.tsx
│   │   │   └── 📄 StudentLayout.tsx
│   │   ├── 📁 pages
│   │   │   ├── 📁 admin
│   │   │   │   └── 📄 UserManagement.tsx
│   │   │   ├── 📁 auth
│   │   │   │   ├── 📄 Login.tsx
│   │   │   │   └── 📄 Logout.tsx
│   │   │   ├── 📁 authority
│   │   │   │   ├── 📄 ApproveRequest.tsx
│   │   │   │   ├── 📄 PendingRequests.tsx
│   │   │   │   └── 📄 ReleaseResults.tsx
│   │   │   ├── 📁 dashboard
│   │   │   │   ├── 📄 AdminDashboard.tsx
│   │   │   │   ├── 📄 AuthorityDashboard.tsx
│   │   │   │   ├── 📄 EvaluatorDashboard.tsx
│   │   │   │   ├── 📄 ExaminerDashboard.tsx
│   │   │   │   └── 📄 StudentDashboard.tsx
│   │   │   ├── 📁 errors
│   │   │   │   ├── 📄 NotFound.tsx
│   │   │   │   └── 📄 Unauthorized.tsx
│   │   │   ├── 📁 evaluator
│   │   │   │   ├── 📄 AssignedScripts.tsx
│   │   │   │   ├── 📄 EvaluateScript.tsx
│   │   │   │   ├── 📄 EvaluationPage.tsx
│   │   │   │   ├── 📄 RequestEvaluation.tsx
│   │   │   │   └── 📄 RequestHistory.tsx
│   │   │   ├── 📁 examiner
│   │   │   │   ├── 📄 UploadHistory.tsx
│   │   │   │   └── 📄 UploadScripts.tsx
│   │   │   ├── 📁 student
│   │   │   │   ├── 📄 DownloadPDF.tsx
│   │   │   │   ├── 📄 StudentResults.tsx
│   │   │   │   └── 📄 ViewResults.tsx
│   │   │   └── 📄 HealthCheck.tsx
│   │   ├── 📁 routes
│   │   │   ├── 📄 AppRoutes.tsx
│   │   │   ├── 📄 DashboardRoute.tsx
│   │   │   ├── 📄 PrivateRoute.tsx
│   │   │   ├── 📄 ProtectedRoute.tsx
│   │   │   └── 📄 RoleGuard.tsx
│   │   ├── 📁 services
│   │   │   └── 📄 apiClient.ts
│   │   ├── 📁 store
│   │   │   ├── 📄 app.ts
│   │   │   ├── 📄 auth.ts
│   │   │   └── 📄 user.ts
│   │   ├── 📁 styles
│   │   │   ├── 🎨 layouts.css
│   │   │   └── 🎨 tailwind.css
│   │   ├── 📁 types
│   │   │   ├── 📄 auth.ts
│   │   │   ├── 📄 authority.ts
│   │   │   ├── 📄 common.ts
│   │   │   ├── 📄 env.d.ts
│   │   │   ├── 📄 evaluator.ts
│   │   │   ├── 📄 examiner.ts
│   │   │   └── 📄 student.ts
│   │   ├── 📁 utils
│   │   │   ├── 📄 download.ts
│   │   │   ├── 📄 roles.ts
│   │   │   ├── 📄 token.ts
│   │   │   └── 📄 validators.ts
│   │   ├── 📄 App.tsx
│   │   ├── 🎨 index.css
│   │   └── 📄 main.tsx
│   ├── ⚙️ .gitignore
│   ├── 🌐 index.html
│   ├── ⚙️ package-lock.json
│   ├── ⚙️ package.json
│   ├── 📄 postcss.config.js
│   ├── 📄 tailwind.config.js
│   ├── ⚙️ tsconfig.app.json
│   ├── ⚙️ tsconfig.json
│   └── 📄 vite.config.ts
├── 📁 infra
│   └── 📁 migrations
│       └── 📁 postgres
│           ├── 📄 V001__create_users.sql
│           ├── 📄 V002__create_results.sql
│           ├── 📄 V003__audit_logs.sql
│           ├── 📄 V004__authority_evaluator.sql
│           ├── 📄 V005__evaluations_table.sql
│           ├── 📄 V006__results_release.sql
│           └── 📄 run_all.sql
├── 📁 scripts
│   ├── 📄 bash.sh
├── 📁 services
│   ├── 📁 go-node
│   │   ├── 📁 cmd
│   │   │   └── 📁 node
│   │   │       └── 🐹 main.go
│   │   ├── 📁 configs
│   │   │   └── ⚙️ config.yaml
│   │   ├── 📁 internal
│   │   │   ├── 📁 admin
│   │   │   │   └── 🐹 service.go
│   │   │   ├── 📁 api
│   │   │   │   ├── 🐹 handlers_blocks.go
│   │   │   │   ├── 🐹 handlers_chain.go
│   │   │   │   ├── 🐹 handlers_validation.go
│   │   │   │   ├── 🐹 register_admin_routes.go
│   │   │   │   ├── 🐹 register_auth_routes.go
│   │   │   │   ├── 🐹 spa.go
│   │   │   │   ├── 🐹 register_authority_routes.go
│   │   │   │   ├── 🐹 register_evaluator_routes.go
│   │   │   │   ├── 🐹 register_examiner_routes.go
│   │   │   │   ├── 🐹 register_results_routes.go
│   │   │   │   ├── 🐹 register_student_routes.go
│   │   │   │   ├── 🐹 routes.go
│   │   │   │   └── 🐹 server.go
│   │   │   ├── 📁 auth
│   │   │   │   ├── 🐹 cookie_writer.go
│   │   │   │   ├── 🐹 handlers.go
│   │   │   │   ├── 🐹 jwt.go
│   │   │   │   ├── 🐹 jwt_manager.go
│   │   │   │   ├── 🐹 middleware.go
│   │   │   │   ├── 🐹 models.go
│   │   │   │   ├── 🐹 password.go
│   │   │   │   ├── 🐹 queries_users.go
│   │   │   │   ├── 🐹 refresh_handler.go
│   │   │   │   ├── 🐹 refresh_service.go
│   │   │   │   ├── 🐹 roles.go
│   │   │   │   ├── 🐹 service.go
│   │   │   │   └── 🐹 store_postgres.go
│   │   │   ├── 📁 authority
│   │   │   │   ├── 🐹 handlers.go
│   │   │   │   ├── 🐹 models.go
│   │   │   │   ├── 🐹 release_handler.go
│   │   │   │   ├── 🐹 release_service.go
│   │   │   │   └── 🐹 service.go
│   │   │   ├── 📁 block
│   │   │   │   ├── 🐹 block.go
│   │   │   │   ├── 🐹 hash.go
│   │   │   │   ├── 🐹 signature.go
│   │   │   │   └── 🐹 transaction.go
│   │   │   ├── 📁 chain
│   │   │   │   ├── 🐹 chain.go
│   │   │   │   ├── 🐹 iterator.go
│   │   │   │   └── 🐹 validate.go
│   │   │   ├── 📁 core
│   │   │   │   ├── 🐹 context.go
│   │   │   │   ├── 🐹 envelope.go
│   │   │   │   ├── 🐹 errors.go
│   │   │   │   ├── 🐹 middlewares.go
│   │   │   │   └── 🐹 service_registry.go
│   │   │   ├── 📁 db
│   │   │   │   ├── 🐹 postgres.go
│   │   │   │   ├── 🐹 queries_authority.go
│   │   │   │   ├── 🐹 queries_evaluator.go
│   │   │   │   └── 🐹 queries_results.go
│   │   │   ├── 📁 evaluator
│   │   │   │   ├── 🐹 handlers.go
│   │   │   │   ├── 🐹 models.go
│   │   │   │   ├── 🐹 service.go
│   │   │   │   ├── 🐹 submit_handler.go
│   │   │   │   ├── 🐹 submit_service.go
│   │   │   │   ├── 🐹 upload_handler.go
│   │   │   │   └── 🐹 upload_service.go
│   │   │   ├── 📁 examiner
│   │   │   │   ├── 🐹 models.go
│   │   │   │   ├── 🐹 storage_temp.go
│   │   │   │   ├── 🐹 upload_handler.go
│   │   │   │   └── 🐹 upload_service.go
│   │   │   ├── 📁 logger
│   │   │   │   ├── 🐹 formatter.go
│   │   │   ├── 📁 pybridge
│   │   │   │   ├── 🐹 client.go
│   │   │   │   └── 🐹 types.go
│   │   │   ├── 📁 rootdir
│   │   │   │   └── 🐹 rootdir.go
│   │   │   ├── 📁 storage
│   │   │   │   ├── 🐹 boltdb.go
│   │   │   │   └── 🐹 buckets.go
│   │   │   └── 📁 student
│   │   │       ├── 📁 assets
│   │   │       │   └── 📁 fonts
│   │   │       │       ├── 📄 NotoSansKannada-Regular.ttf
│   │   │       │       ├── 📄 Roboto-Bold.ttf
│   │   │       │       └── 📄 Roboto-Regular.ttf
│   │   │       ├── 🐹 handlers.go
│   │   │       ├── 🐹 pdf_generator.go
│   │   │       └── 🐹 service.go
│   │   ├── 📁 ui
│   │   │   └── 🐹 embed.go
│   │   ├── 📄 go.mod
│   │   └── 📄 go.sum
│   └── 📁 python-validator
│       ├── 📁 app
│       │   ├── 🐍 __init__.py
│       │   ├── 🐍 cli.py
│       │   ├── 🐍 extractor.py
│       │   ├── 🐍 genai_client.py
│       │   ├── 🐍 ipfs_client.py
│       │   ├── 🐍 main.py
│       │   ├── 🐍 schema.py
│       │   ├── 🐍 validate_evaluation.py
│       │   └── 🐍 validator.py
│       └── 📄 requirements.txt
├── 📁 tools
│   ├── 📁 keygen
│   │   └── 📄 gen_keys.sh
│   └── 📄 local_setup.sh
├── ⚙️ .gitignore
├── 📄 LICENSE
├── 📄 Makefile
├── 📝 README.md
└── 📝 STRUCTURE.md
```
