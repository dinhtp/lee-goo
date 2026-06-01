# Technical Requirement Document

## Golang Modular Monorepo Boilerplate

## 1\. Objective

Build a reusable Golang backend boilerplate that supports systems of different sizes, from small services to large enterprise platforms.

The boilerplate shall use a monorepo structure and a modular architecture where each module can be installed, enabled, disabled, upgraded, or uninstalled independently.

The modular approach shall be similar in concept to Magento 2 modules, where modules declare their own dependencies, expose service contracts, provide extension points, and extend the behavior of other modules without tightly coupling to their internal implementation.

The boilerplate shall support modules such as:

- User
- Authentication
- Authorization
- Module management
- Notification
- Audit log
- Payment
- Order
- Inventory
- Report

The boilerplate must be suitable for long-term enterprise usage and must enforce clear architectural boundaries between modules.

## 2\. Core Architecture Decision

The codebase shall use a **compiled modular architecture**, not runtime dynamic shared-object loading.

The recommended foundation is:

- **Go workspace / monorepo** for repository-level organization.
- **Custom module manager** for module discovery, dependency validation, lifecycle execution, installation, upgrade, disable, and uninstall.
- **Uber Fx** for dependency injection, runtime composition, lifecycle hooks, and module bootstrapping.
- **golang-migrate or Atlas** for module-level database migration management.
- **Interface-based contracts** for inter-module communication.
- **Event bus and extension-point registry** for loose coupling between modules.
- **Source-code-based module installation and uninstallation**, where install/uninstall changes the codebase and module registry instead of loading third-party binaries at runtime.

The system shall avoid runtime plugin loading as the default module mechanism.

The system shall avoid:

- Go native plugin package.
- HashiCorp go-plugin.
- Runtime dynamic module loading.
- Third-party binary module loading.
- Reflection-heavy service discovery.
- Hidden module registration through blank imports.

The boilerplate shall prefer explicit module registration, static compilation, and predictable build-time composition.

## 3\. Module Design Principle

Each module shall be treated as an independent mini-application domain.

Each module must own its own:

- Domain model
- Use case contracts
- Port interfaces
- Service implementations
- Repository implementations
- HTTP handlers
- Route registration
- Middleware, if needed
- Migration files
- Configuration
- Dependency declarations
- Event handlers
- Tests
- Public contracts
- Optional extension points

A module must not directly access another module's internal packages.

Cross-module interaction must happen only through:

- Public contracts
- Exported service interfaces
- Events
- Declared extension points
- Dependency-injected interfaces
- Stable request/response contracts

A module must be replaceable as long as it satisfies the public contracts required by other modules.

## 4\. Recommended Module Structure

Each business module shall follow a consistent internal structure based on DDD-lite and hexagonal architecture.

The boilerplate shall use a **module-first structure**, where each module owns its own domain, handler, repository, service, router, middleware, migration, configuration, and tests.

The root application shall contain:

/  
├── cmd/  
│ ├── api/  
│ │ └── main.go  
│ ├── worker/  
│ │ └── main.go  
│ └── module/  
│ └── main.go  
├── modules/  
│ ├── user/  
│ ├── authentication/  
│ ├── authorization/  
│ └── module/  
├── system/     
│ ├── config/  
│ ├── database/  
│ ├── eventbus/  
│ ├── http/  
│ ├── grpc/  
│ ├── logger/  
│ ├── security/  
│ └── fx/  
├── pkg/  
├── go.work  
├── go.mod  
└── README.md

The cmd/module command shall be the executable entry point for module management actions.

Example:

go run ./cmd/module list  
go run ./cmd/module install user  
go run ./cmd/module uninstall notification

The module management feature itself shall be implemented as a normal module located at:

/modules/module_manager

This module must follow the same architecture, naming convention, dependency flow, testing rules, and constructor pattern as all other modules.

### 4.1 Standard Module Directory Structure

Each module shall use the following structure:

/modules/{module}/  
├── module.yaml  
├── go.mod  
├── contracts/  
│ ├── service.go  
│ ├── event.go  
│ ├── port.go  
│ └── error.go  
├── internal/  
│ ├── domain/{domain}/  
│ │ ├── domain.go  
│ │ ├── port.go  
│ │ └── use_case.go  
│ ├── handler/{domain}/  
│ │ ├── handler.go  
│ │ ├── contract.go  
│ │ ├── router.go  
│ │ └── handler_test.go  
│ ├── middleware/  
│ │ └── default_http_middleware.go  
│ ├── repository/{domain}/  
│ │ └── \*\_repository.go  
│ ├── router/  
│ │ ├── router.go  
│ │ └── register.go  
│ └── service/{domain}/  
│ ├── service.go  
│ └── service_test.go  
├── migrations/  
│ ├── 000001_create_table.up.sql  
│ └── 000001_create_table.down.sql  
├── config/  
│ └── config.go  
├── fx/  
│ └── module.go  
└── tests/  
└── module_test.go

{module} is the module name, such as:

user  
authentication  
authorization  
module  
notification  
audit_log

{domain} is a short lowercase noun, such as:

system  
user  
order  
role  
permission  
session  
module

A module may contain one or more domains.

### 4.2 Internal Directory Structure

Each module shall use the following internal directory layout:

internal/  
├── domain/{domain}/ # DDD-lite contracts used inside the module  
│ ├── domain.go # Domain model structs  
│ ├── port.go # Repository/service port interfaces  
│ └── use_case.go # UseCase interface for business operations  
├── handler/{domain}/ # HTTP transport layer  
│ ├── handler.go # Echo handler methods  
│ ├── contract.go # JSON request/response structs  
│ ├── router.go # Route registration, implements HandlerRouter  
│ └── handler_test.go # Handler unit tests  
├── middleware/ # Echo middleware provider functions  
│ └── default_http_middleware.go  
├── repository/{domain}/ # Port implementations and data access  
│ └── \*\_repository.go  
├── router/ # Shared routing infrastructure inside the module  
│ ├── router.go # HandlerRouter interface definition  
│ └── register.go # Register() mounts routes and delegates to domain routers  
└── service/{domain}/ # UseCase implementations and business logic  
├── service.go  
└── service_test.go

### 4.3 Public Contract Directory

Each module shall expose cross-module contracts through a dedicated contracts/ directory.

contracts/  
├── service.go # Public service interfaces exposed to other modules  
├── event.go # Public event payload definitions  
├── port.go # Public extension interfaces, if required  
└── error.go # Public module error definitions, if required

Other modules may import:

modules/user/contracts

Other modules must not import:

modules/user/internal/domain/user  
modules/user/internal/service/user  
modules/user/internal/repository/user  
modules/user/internal/handler/user

This rule keeps module internals private while still allowing controlled extension and dependency between modules.

### 4.4 Layer Responsibilities

| Layer            | Package path                                  | Responsibility                                                                |
| ---------------- | --------------------------------------------- | ----------------------------------------------------------------------------- |
| Public contracts | modules/{module}/contracts                    | Interfaces, events, public types, and errors that other modules may import    |
| Domain           | modules/{module}/internal/domain/{domain}     | Internal domain models, port interfaces, and UseCase interface                |
| Repository       | modules/{module}/internal/repository/{domain} | Implement domain port interfaces and access databases or external systems     |
| Service          | modules/{module}/internal/service/{domain}    | Implement domain UseCase and orchestrate port calls                           |
| Handler          | modules/{module}/internal/handler/{domain}    | HTTP only: parse request, call UseCase, map response and status code          |
| Router           | modules/{module}/internal/router              | Wire Echo route groups and delegate to domain routers                         |
| Middleware       | modules/{module}/internal/middleware          | Configure module-level Echo middleware                                        |
| Migration        | modules/{module}/migrations                   | Own module database migration files                                           |
| Config           | modules/{module}/config                       | Own module configuration schema and defaults                                  |
| Fx               | modules/{module}/fx                           | Provide Fx module options for dependency injection and lifecycle registration |
| Tests            | modules/{module}/tests                        | Integration, contract, and module boot tests                                  |

### 4.5 Dependency Flow

The dependency flow inside each module must be strict and must never be reversed.

domain ← repository implements domain port interfaces  
domain ← service implements domain.UseCase  
domain ← handler depends on domain.UseCase  
router ← handler handler/router.go implements router.HandlerRouter  
fx ← repository, service, handler, router

The handler must not depend directly on service or repository implementation.

The service must not know HTTP, Echo, SQL, Redis, Kafka, or external infrastructure details.

The repository must implement domain ports and isolate persistence or external system access.

The Fx package is allowed to import implementation packages for wiring only.

### 4.6 Import Rules

| Package                      | May import                                                          |
| ---------------------------- | ------------------------------------------------------------------- |
| contracts                    | stdlib only, shared primitive packages if approved                  |
| internal/domain/{domain}     | stdlib, context only                                                |
| internal/repository/{domain} | internal/domain/{domain}, approved pkg/\*, infrastructure clients   |
| internal/service/{domain}    | internal/domain/{domain}, module contracts only if needed           |
| internal/handler/{domain}    | internal/domain/{domain}, internal/router                           |
| internal/middleware          | config, server, logger, security packages                           |
| internal/router              | stdlib, Echo only                                                   |
| fx                           | module internal implementation packages, platform Fx helpers        |
| Other modules                | only modules/{module}/contracts, never another module's internal/\* |

### 4.7 Coding Style

- **Interface-first:** define contracts in domain/ or contracts/ before writing implementations.
- **Dependency injection via constructors:** no init(), no package-level singletons in internal/.
- **Thin service layer:** service only orchestrates port calls; no SQL and no HTTP logic.
- **Thin handler layer:** handler only translates HTTP to domain and domain to HTTP.
- **No blank imports** except explicitly approved infrastructure registration cases.
- **No dot imports.**
- **Receiver names:** single lowercase letter or short abbreviation matching the type, such as h for Handler, r for repository/router, and s for service.
- **No direct cross-module internal imports.**
- **No hidden module registration through package side effects.**

### 4.8 File Naming Convention

All file names must use snake_case.

| File                  | Content                                              |
| --------------------- | ---------------------------------------------------- |
| domain.go             | Domain model structs                                 |
| port.go               | Port interfaces                                      |
| use_case.go           | UseCase interface                                    |
| handler.go            | Handler struct and methods                           |
| contract.go           | Request/response JSON structs                        |
| router.go             | Route registration                                   |
| handler_test.go       | Handler tests                                        |
| service.go            | Service implementation                               |
| service_test.go       | Service tests                                        |
| {name}\_repository.go | Repository implementation, e.g. health_repository.go |
| module.go             | Fx module provider                                   |
| config.go             | Module configuration schema                          |

### 4.9 Package Naming Convention

Package name shall be the last path segment, lowercase, without underscores.

Examples:

modules/user/internal/domain/system -> package system  
modules/user/internal/router -> package router  
modules/user/internal/middleware -> package middleware  
modules/user/contracts -> package contracts  
modules/user/fx -> package fx

### 4.10 Type, Interface, and Struct Naming Convention

| Concept              | Rule                              | Example                       |
| -------------------- | --------------------------------- | ----------------------------- |
| Domain model struct  | PascalCase noun                   | Health, Version, User         |
| Port interface       | PascalCase + Port suffix          | HealthPort, VersionPort       |
| UseCase interface    | Fixed name UseCase per domain     | UseCase                       |
| Service struct       | Unexported lowercase domain name  | service                       |
| Repository struct    | Unexported lowercase + Repository | healthRepository              |
| Handler struct       | Exported Handler                  | Handler                       |
| Router struct        | Unexported router                 | router                        |
| JSON response struct | PascalCase + Response             | HealthResponse, ErrorResponse |
| JSON request struct  | PascalCase + Request              | CreateUserRequest             |
| Fx module function   | Exported Module                   | Module()                      |

### 4.11 Constructor Pattern

Constructors must use dependency injection and must not rely on init() or package-level singletons.

// Always return the interface type, not the concrete struct.  
func NewService(healthPort domain.HealthPort, versionPort domain.VersionPort) domain.UseCase {  
return &service{  
healthPort: healthPort,  
versionPort: versionPort,  
}  
}  
<br/>func NewHealthRepository(connection database.Connection) domain.HealthPort {  
return &healthRepository{connection: connection}  
}  
<br/>func NewHandler(useCase domain.UseCase) \*Handler {  
// Handler is the exception: returns \*Handler because caller needs to pass it to NewRouter.  
return &Handler{useCase: useCase}  
}  
<br/>func NewRouter(handler \*Handler) internalRouter.HandlerRouter {  
return &router{handler: handler}  
}

Constructor rules:

- Prefix always uses New + exported type name.
- Service and repository constructors return interfaces.
- Handler constructor may return \*Handler.
- Struct field names in the returned literal must match parameter names exactly.
- Constructors must not create hidden dependencies internally.
- Constructors must not read environment variables directly.
- Constructors must not open database connections directly unless the type is explicitly an infrastructure provider.

### 4.12 Compile-Time Interface Assertion

Each implementation must include a compile-time interface assertion immediately after the struct declaration and before the first method.

type healthRepository struct {  
connection database.Connection  
}  
<br/>var \_ domain.HealthPort = (\*healthRepository)(nil)

Every implementation file must include one assertion per interface it implements.

Required files include:

service.go  
\*\_repository.go  
router.go

### 4.13 Import Alias Convention

Import aliases must be used when package names are ambiguous or collide.

import (  
domain "github.com/company/project/modules/user/internal/domain/system"  
internalRouter "github.com/company/project/modules/user/internal/router"  
systemHandler "github.com/company/project/modules/user/internal/handler/system"  
systemRepository "github.com/company/project/modules/user/internal/repository/system"  
systemService "github.com/company/project/modules/user/internal/service/system"  
)

Rules:

- Alias format should be {layer}{Domain} where practical.
- Use shortened aliases only when unambiguous.
- Prefer full qualified aliases for domain packages to avoid shadowing built-ins.
- Never alias standard library packages.

### 4.14 Example: User Module

/modules/user  
├── module.yaml  
├── go.mod  
├── contracts  
│ ├── service.go  
│ ├── event.go  
│ └── error.go  
├── internal  
│ ├── domain/user  
│ │ ├── domain.go  
│ │ ├── port.go  
│ │ └── use_case.go  
│ ├── handler/user  
│ │ ├── handler.go  
│ │ ├── contract.go  
│ │ ├── router.go  
│ │ └── handler_test.go  
│ ├── repository/user  
│ │ └── user_repository.go  
│ ├── router  
│ │ ├── router.go  
│ │ └── register.go  
│ └── service/user  
│ ├── service.go  
│ └── service_test.go  
├── migrations  
│ ├── 000001_create_users.up.sql  
│ └── 000001_create_users.down.sql  
├── config  
│ └── config.go  
├── fx  
│ └── module.go  
└── tests  
└── user_module_test.go

### 4.15 Example: Authentication Depends on User

The authentication module must not import the internal implementation of the user module.

Allowed:

import userContracts "github.com/company/project/modules/user/contracts"

Forbidden:

import "github.com/company/project/modules/user/internal/service/user"  
import "github.com/company/project/modules/user/internal/repository/user"  
import "github.com/company/project/modules/user/internal/domain/user"

Authentication shall depend on the public user contract, not the user implementation.

type service struct {  
userService userContracts.UserService  
}

This ensures that authentication can extend the user module while preserving module independence and replaceability.

### 4.16 Example: Module Management Module

The module management feature shall be implemented as a module.

/modules/module_manager  
├── module.yaml  
├── go.mod  
├── contracts  
│ ├── service.go  
│ ├── event.go  
│ └── error.go  
├── internal  
│ ├── domain/module  
│ │ ├── domain.go  
│ │ ├── port.go  
│ │ └── use_case.go  
│ ├── handler/module  
│ │ ├── handler.go  
│ │ ├── contract.go  
│ │ ├── router.go  
│ │ └── handler_test.go  
│ ├── repository/module  
│ │ └── module_repository.go  
│ ├── router  
│ │ ├── router.go  
│ │ └── register.go  
│ └── service/module  
│ ├── service.go  
│ └── service_test.go  
├── migrations  
│ ├── 000001_create_modules.up.sql  
│ └── 000001_create_modules.down.sql  
├── config  
│ └── config.go  
├── fx  
│ └── module.go  
└── tests  
└── module_management_test.go

The cmd/module executable shall call the use cases exposed by this module.

The module management module is responsible for:

- Discovering modules.
- Reading module manifests.
- Validating dependencies.
- Installing modules.
- Enabling modules.
- Disabling modules.
- Upgrading modules.
- Uninstalling modules.
- Removing module source code during uninstall.
- Updating the module registry table.
- Updating workspace and generated registration files.
- Running module migrations.

## 5\. Module Manifest

Each module must include a module.yaml file.

Example:

name: user  
version: 1.0.0  
description: Core user management module  
status: stable  
<br/>dependencies:  
required: \[\]  
optional: \[\]  
<br/>provides:  
services:  
\- UserService  
\- UserRepository  
events:  
\- user.created  
\- user.updated  
<br/>extension*points:  
\- user.profile.validator  
\- user.after_created  
<br/>migrations:  
path: ./migrations  
transactional: true  
<br/>config:  
prefix: USER*

Example for authentication module:

name: authentication  
version: 1.0.0  
description: Authentication module  
<br/>dependencies:  
required:  
\- user  
optional: \[\]  
<br/>extends:  
\- user  
<br/>provides:  
services:  
\- AuthService  
\- TokenService  
<br/>migrations:  
path: ./migrations

Example for authorization module:

name: authorization  
version: 1.0.0  
description: Authorization and permission module  
<br/>dependencies:  
required:  
\- user  
optional: \[\]  
<br/>provides:  
services:  
\- RoleService  
\- PermissionService  
\- PolicyService  
<br/>migrations:  
path: ./migrations

Example for module management module:

name: module  
version: 1.0.0  
description: Module management module  
<br/>dependencies:  
required: \[\]  
optional: \[\]  
<br/>provides:  
services:  
\- ModuleService  
\- ModuleRepository  
<br/>migrations:  
path: ./migrations

## 6\. Module Dependency Rules

The module manager must support two dependency types.

### 6.1 Required Dependency

A required dependency means the dependent module cannot be installed, enabled, or booted unless the required module is already installed and enabled.

Example:

authorization -> user  
authentication -> user

Authorization must require the user module because authorization cannot exist without a user identity.

Authentication must require the user module because authentication needs to identify a user account before issuing credentials or tokens.

### 6.2 Optional Dependency

An optional dependency means a module can enhance its behavior when another module exists, but it can still work without it.

Example:

notification -> user  
audit_log -> authentication

Optional dependencies shall be resolved at runtime through service discovery, event subscribers, or extension-point registration.

### 6.3 Dependency Validation Rules

The module manager must validate:

- Required dependencies exist.
- Required dependencies are installed.
- Required dependencies are enabled before the dependent module is enabled.
- Dependency versions are compatible.
- Circular dependencies do not exist.
- Disabled modules are not required by enabled modules.
- Uninstalled modules are not required by installed modules.
- A module cannot depend on another module's internal package.

## 7\. Inheritance and Extension Model

The architecture shall not use classical inheritance between modules.

Go favors composition and interfaces.

Therefore, module "inheritance" shall be implemented through:

- Required dependencies
- Public interfaces
- Service contracts
- Decorators
- Event subscribers
- Extension points
- Middleware chains
- Policy hooks

Example:

Authentication extends User by depending on the user service contract:

type UserService interface {  
FindByEmail(ctx context.Context, email string) (\*User, error)  
FindByID(ctx context.Context, id string) (\*User, error)  
}

Authorization extends User by requiring a user identity:

type IdentityProvider interface {  
CurrentUser(ctx context.Context) (\*UserIdentity, error)  
}

The module manager must not implement extension through source-code inheritance.

The system shall use interface composition and explicit dependency registration instead.

## 8\. Module Runtime Contract

Each module must expose a module provider.

type Module interface {  
Name() string  
Version() string  
Dependencies() \[\]Dependency  
FxOptions() \[\]fx.Option  
Migrations() MigrationSource  
}

Each module shall expose its Fx options through its fx package.

Example:

func Module() fx.Option {  
return fx.Module(  
"user",  
fx.Provide(  
NewUserRepository,  
NewUserService,  
NewUserHandler,  
),  
fx.Invoke(RegisterUserRoutes),  
)  
}

The root application shall compose enabled modules:

app := fx.New(  
infrafx.Options(),  
moduleManager.EnabledFxOptions()...,  
)

The root application must not manually instantiate module internals.

The root application must only ask the module manager for enabled module options.

## 9\. Module Manager Requirements

Module management shall be implemented as a dedicated module named:

module

The module shall be located at:

/modules/module_manager

The module management module shall follow the same architecture as every other module.

It must have:

- contracts
- internal/domain/module
- internal/service/module
- internal/repository/module
- internal/handler/module
- internal/router
- migrations
- config
- fx
- tests

The module manager shall be responsible for:

- Discovering modules.
- Reading module manifests.
- Validating module dependencies.
- Detecting circular dependencies.
- Sorting modules by dependency order.
- Installing modules.
- Running module migrations.
- Enabling modules.
- Disabling modules.
- Upgrading modules.
- Uninstalling modules.
- Removing uninstalled module source code from the codebase.
- Registering module services.
- Registering extension points.
- Maintaining module state in the database.
- Updating generated module registration files.
- Updating go.work after install or uninstall.
- Updating root dependency references when needed.
- Validating that the application can still compile after module operations.

The module management module must not expose admin HTTP APIs by default.

Module management actions shall be executed through the local CLI entry point:

cmd/module

## 10\. Module Registry Table

The platform shall maintain a module registry table named:

modules

Example:

CREATE TABLE modules (  
name VARCHAR(100) PRIMARY KEY,  
version VARCHAR(50) NOT NULL,  
status VARCHAR(30) NOT NULL,  
path VARCHAR(255) NOT NULL,  
checksum VARCHAR(255) NULL,  
installed_at TIMESTAMP NULL,  
enabled_at TIMESTAMP NULL,  
disabled_at TIMESTAMP NULL,  
upgraded_at TIMESTAMP NULL,  
uninstalled_at TIMESTAMP NULL,  
removed_from_codebase_at TIMESTAMP NULL,  
created_at TIMESTAMP NOT NULL,  
updated_at TIMESTAMP NOT NULL  
);

Supported statuses:

discovered  
installed  
enabled  
disabled  
upgrading  
failed  
uninstalled  
removed

Status meaning:

| Status      | Meaning                                                    |
| ----------- | ---------------------------------------------------------- |
| discovered  | Module exists in the codebase but is not installed         |
| installed   | Module is registered and migrations have been applied      |
| enabled     | Module is installed and active during application boot     |
| disabled    | Module is installed but not active during application boot |
| upgrading   | Module is currently being upgraded                         |
| failed      | Module lifecycle operation failed                          |
| uninstalled | Module was uninstalled logically                           |
| removed     | Module was removed from the codebase                       |

The module manager must update this table during lifecycle operations.

## 11\. Installation Flow

Module installation shall follow this flow:

1\. Discover module source code under /modules/{module}  
2\. Read module.yaml  
3\. Validate required dependencies  
4\. Validate optional dependencies  
5\. Validate version compatibility  
6\. Detect circular dependencies  
7\. Sort modules by dependency graph  
8\. Run pre-install hook  
9\. Run module migrations  
10\. Insert or update modules table  
11\. Generate or update module registration file  
12\. Update go.work if needed  
13\. Register services and extension points  
14\. Run post-install hook  
15\. Mark module as installed  
16\. Run compile check

Example CLI:

go run ./cmd/module install user  
go run ./cmd/module install authentication  
go run ./cmd/module install authorization

If a required dependency is missing, installation must fail.

Example:

go run ./cmd/module install authorization

Expected result:

Error: authorization requires user module to be installed and enabled

Installation must be a source-code-aware operation.

The module must exist in the codebase before it can be installed.

Example:

/modules/user  
/modules/authentication  
/modules/authorization

The module manager must not download, load, or execute third-party module binaries.

## 12\. Enable and Disable Flow

Enabling a module shall validate:

- Module exists in the codebase.
- Module is installed.
- Required dependencies are installed.
- Required dependencies are enabled.
- No dependency conflict exists.
- Module configuration is valid.
- Module can be included in the compiled application.
- Module Fx options can be resolved.

Example:

go run ./cmd/module enable authentication

Disabling a module shall validate:

- No enabled module depends on it.
- Or the dependent modules are also disabled.
- The operation does not break system integrity.
- The generated module registration file is updated.
- The application can still compile after the module is disabled.

Example:

go run ./cmd/module disable user

Expected result:

Error: cannot disable user because authentication and authorization depend on it

Disable must not delete source code.

Disable only changes module activation state.

## 13\. Uninstall Flow

Uninstall shall be treated as a destructive local codebase operation.

Unlike disable, uninstall must remove the module from the codebase.

Uninstall shall include both:

1\. Logical uninstall from module registry  
2\. Physical removal from source code

Default behavior should be:

disable first, uninstall later

The uninstall flow shall include:

1\. Check whether the module exists in /modules/{module}  
2\. Check dependent modules  
3\. Stop if another installed or enabled module depends on it  
4\. Require explicit confirmation  
5\. Run pre-uninstall hook  
6\. Backup or archive module-owned data if configured  
7\. Run uninstall migration only if allowed  
8\. Mark module as uninstalled in modules table  
9\. Remove module from generated registration file  
10\. Remove module from go.work  
11\. Remove module dependency references if applicable  
12\. Remove module source directory from /modules/{module}  
13\. Run go mod tidy where applicable  
14\. Run compile check  
15\. Mark module as removed in modules table  
16\. Run post-uninstall hook

Example:

go run ./cmd/module uninstall notification

Expected result:

Module notification has been disabled, uninstalled, removed from registration, removed from go.work, and deleted from /modules/notification.

The uninstall command must not be exposed as a default runtime HTTP API.

The uninstall command is intended for local development, controlled codebase maintenance, CI automation, or build-time module composition.

Production systems should avoid running physical source-code removal commands against a deployed runtime environment.

### 13.1 Uninstall Safety Rules

The module manager must prevent uninstall when:

- Another installed module requires the target module.
- Another enabled module requires the target module.
- The module is a protected core module.
- The module is the module management module itself.
- The module has unapplied down migrations and uninstall requires rollback.
- Source code removal would break compilation.
- The module directory is outside the allowed /modules path.
- The module manifest checksum does not match the expected checksum.

### 13.2 Protected Modules

The following modules should be protected by default:

module  
user  
authentication  
authorization

Protected modules may only be uninstalled with a force flag in non-production environments.

Example:

go run ./cmd/module uninstall user --force --env=local

## 14\. Migration Requirements

Each module shall own its own migration directory.

Example:

/modules/user/migrations  
000001_create_users.up.sql  
000001_create_users.down.sql  
<br/>/modules/authentication/migrations  
000001_create_auth_sessions.up.sql  
000001_create_auth_sessions.down.sql  
<br/>/modules/module_manager/migrations  
000001_create_modules.up.sql  
000001_create_modules.down.sql

Migration execution order shall follow module dependency order.

Example:

module migrations  
user migrations  
authentication migrations  
authorization migrations

Migration requirements:

- Migrations must be versioned.
- Migrations must be idempotent where possible.
- Migration status must be trackable per module.
- Failed migrations must stop installation.
- Rollback strategy must be explicitly defined.
- Production rollback must be treated carefully and should prefer roll-forward fixes.
- Migration files must remain inside the owning module.
- One module must not modify another module's tables unless explicitly approved through a shared contract.

Migration tooling options:

Simple SQL migration execution:  
\- golang-migrate  
<br/>Advanced migration planning, linting, and CI checks:  
\- Atlas

The boilerplate shall support golang-migrate as the default migration execution engine.

Atlas may be used in CI/CD for migration linting, safety checks, schema planning, and migration review.

## 15\. Extension Point Requirements

The platform shall support extension points.

Example:

type ExtensionPoint\[T any\] interface {  
Register(name string, priority int, handler T)  
Resolve() \[\]T  
}

Use cases:

- Add validation to user profile creation.
- Add login risk checks.
- Add authorization policy rules.
- Add audit logging.
- Add notification hooks.
- Add module lifecycle hooks.

Example:

type UserCreatedHook interface {  
Handle(ctx context.Context, user User) error  
}

Modules can register handlers:

fx.Invoke(func(registry ExtensionRegistry) {  
registry.Register("user.after_created", auditLogHandler, 100)  
})

Extension-point rules:

- Extension points must be explicitly declared.
- Extension handlers must have priority ordering.
- Extension handlers must be replaceable.
- Extension handlers must not depend on another module's internal implementation.
- Extension handlers must be testable in isolation.
- Critical business logic must not rely on hidden side effects.

## 16\. Event-Driven Communication

Modules should communicate through events when synchronous dependency is not required.

Example events:

user.created  
user.updated  
auth.login_succeeded  
auth.login_failed  
authorization.role_assigned  
module.installed  
module.enabled  
module.disabled  
module.uninstalled  
module.removed

Rules:

- Events must be versioned.
- Events must have stable payload contracts.
- Event handlers must be idempotent.
- Event failure handling must be defined.
- Critical domain operations should not depend only on best-effort events.
- Events must be declared in the module's public contracts when consumed by other modules.

Example:

type UserCreatedEvent struct {  
UserID string  
Email string  
}

Event naming convention:

{module}.{action}

Examples:

user.created  
authentication.login_succeeded  
authorization.role_assigned  
module.installed

## 17\. Public Contract Rule

Each module shall expose only its public contract.

Recommended package structure:

/modules/user/contracts  
service.go  
identity.go  
event.go  
error.go

Other modules may import:

modules/user/contracts

Other modules must not import:

modules/user/internal/domain/user  
modules/user/internal/service/user  
modules/user/internal/repository/user  
modules/user/internal/handler/user

Public contracts may include:

- Service interfaces
- Event payloads
- Shared value objects
- Error definitions
- Extension-point interfaces

Public contracts must not include:

- Repository implementation
- Service implementation
- HTTP handler implementation
- Database-specific structs
- ORM models
- SQL queries
- Infrastructure clients

## 18\. Configuration Requirements

Each module shall define its own configuration namespace.

Example:

user:  
allow_registration: true  
<br/>authentication:  
jwt_secret: \${AUTH_JWT_SECRET}  
access_token_ttl: 15m  
<br/>authorization:  
default_role: user  
<br/>module:  
protected_modules:  
\- module  
\- user  
\- authentication  
\- authorization  
allow_source_removal: false

The platform must validate configuration at startup before enabling a module.

Configuration requirements:

- Each module owns its own configuration schema.
- Configuration must be validated before application boot.
- Configuration must not be read directly inside constructors.
- Configuration must be injected into services through typed config structs.
- Sensitive configuration must be loaded through approved secret providers.
- Module configuration keys must be namespaced by module name.

Example:

type Config struct {  
AllowRegistration bool \`mapstructure:"allow_registration"\`  
}

## 19\. CLI Requirements

The boilerplate shall provide a separate module management command.

The command shall be located at:

/cmd/module

The command shall be executable with:

go run ./cmd/module

The command shall delegate business logic to the module management module located at:

/modules/module_manager

The command must not contain business logic directly.

The command may only:

- Parse CLI arguments.
- Load configuration.
- Initialize required dependencies.
- Call module management use cases.
- Print results.
- Return process exit codes.

### 19.1 Required Commands

The module command shall support:

go run ./cmd/module list  
go run ./cmd/module status user  
go run ./cmd/module install user  
go run ./cmd/module enable user  
go run ./cmd/module disable user  
go run ./cmd/module uninstall user  
go run ./cmd/module upgrade user  
go run ./cmd/module migrate user  
go run ./cmd/module migrate all  
go run ./cmd/module graph  
go run ./cmd/module doctor

### 19.2 Source-Code-Aware Commands

Because modules are compiled into the codebase, the module command shall also support source-code-aware operations:

go run ./cmd/module make notification  
go run ./cmd/module remove notification  
go run ./cmd/module sync  
go run ./cmd/module tidy  
go run ./cmd/module compile-check

Command responsibilities:

| Command            | Responsibility                                                                    |
| ------------------ | --------------------------------------------------------------------------------- |
| list               | List discovered, installed, enabled, disabled, and removed modules                |
| status {module}    | Show module status, version, dependencies, and path                               |
| install {module}   | Register module, run migrations, update generated registration                    |
| enable {module}    | Activate module during application boot                                           |
| disable {module}   | Deactivate module during application boot                                         |
| uninstall {module} | Disable, uninstall, remove registration, update workspace, and delete source code |
| upgrade {module}   | Apply version upgrade flow                                                        |
| migrate {module}   | Run migrations for one module                                                     |
| migrate all        | Run migrations in dependency order                                                |
| graph              | Show module dependency graph                                                      |
| doctor             | Validate module manifests, dependencies, migrations, workspace, and compile state |
| make {module}      | Generate a new module skeleton                                                    |
| remove {module}    | Remove module source code only when safe                                          |
| sync               | Rebuild generated module registration from discovered modules                     |
| tidy               | Run workspace/module dependency cleanup                                           |
| compile-check      | Verify the codebase compiles after module changes                                 |

### 19.3 Example CLI Usage

Create a module:

go run ./cmd/module make notification

Install a module:

go run ./cmd/module install notification

Enable a module:

go run ./cmd/module enable notification

Disable a module:

go run ./cmd/module disable notification

Uninstall and remove a module from the codebase:

go run ./cmd/module uninstall notification

Run health checks on the module system:

go run ./cmd/module doctor

## 20\. Testing Requirements

Each module must include:

- Unit tests for domain logic.
- Service tests.
- Repository tests.
- Handler tests.
- Contract tests.
- Migration tests.
- Integration tests with required dependencies.
- Module boot tests.
- Dependency graph tests.

The platform must provide a test harness for loading only selected modules.

Example:

testapp.New(  
WithModules(user.Module(), authentication.Module()),  
)

The module management module must include tests for:

- Module discovery.
- Manifest parsing.
- Dependency validation.
- Circular dependency detection.
- Install flow.
- Enable flow.
- Disable flow.
- Uninstall flow.
- Source code removal.
- go.work update.
- Generated registration update.
- Compile-check validation.
- Migration execution order.
- Protected module prevention.

Example:

func TestCannotUninstallRequiredModule(t \*testing.T) {  
app := testapp.New(  
WithModules(  
module.Module(),  
user.Module(),  
authentication.Module(),  
),  
)  
<br/>err := app.ModuleService().Uninstall(context.Background(), "user")  
<br/>require.Error(t, err)  
}

## 21\. Security Requirements

The module system must enforce:

- No loading of untrusted modules in production.
- No runtime execution of third-party plugin binaries.
- Module checksum verification.
- Dependency validation.
- Administrative permission for local install/uninstall execution.
- Audit logging for all module lifecycle changes.
- Safe migration execution.
- Secret isolation per module.
- No direct database access across module boundaries unless explicitly allowed.
- No direct import of another module's internal packages.
- No source-code removal outside the /modules directory.
- No uninstall of protected modules without force mode.
- No force mode in production.

Uninstall security rules:

- The target path must be resolved and validated before deletion.
- The target path must be inside /modules.
- The target path must match the module manifest.
- The module name must match the directory name.
- The command must reject path traversal attempts.
- The command must reject symbolic-link deletion outside the workspace.
- The command must create a backup or Git diff warning before deletion if configured.
- The command must require explicit confirmation unless --yes is provided in CI.

Example dangerous input that must be rejected:

go run ./cmd/module uninstall ../../system

Expected result:

Error: invalid module path

## 22\. Recommended Library Decision

The boilerplate shall use:

Dependency injection / lifecycle:  
\- Uber Fx  
<br/>Module manager:  
\- Custom internal implementation inside /modules/module_manager  
<br/>Database migration:  
\- golang-migrate for simple SQL migration execution  
\- Atlas optional for schema planning, linting, CI validation, and advanced migration workflow  
<br/>CLI:  
\- Cobra or standard library flag package  
\- The selected CLI library must only handle command parsing  
\- Business logic must remain in /modules/module_manager  
<br/>Avoid as default:  
\- Go native plugin package  
\- HashiCorp go-plugin  
\- Runtime dynamic plugin loading  
\- Third-party binary module loading  
\- Blank-import-based module registration  
\- Reflection-heavy automatic module discovery

Rationale:

- Go modules and workspaces are suitable for source-code-based modular monorepo development.
- Uber Fx is suitable for explicit dependency injection and application lifecycle composition.
- A custom module manager is required because module installation, uninstallation, dependency graph validation, source-code removal, and registry updates are application-specific requirements.
- Runtime plugin systems are not suitable for this boilerplate because the target architecture requires source-controlled, compiled, auditable modules.
- Module uninstall must remove source code from the codebase, which is incompatible with runtime plugin-style module loading.

## 23\. Acceptance Criteria

The boilerplate is accepted when:

- A new module can be created using a generator command.
- Module management is implemented as a normal module under /modules/module_manager.
- The module management command can be executed with go run ./cmd/module.
- A module can be discovered from the /modules directory.
- A module can be installed independently.
- A module can be enabled independently.
- A module can be disabled independently.
- A module can be upgraded independently.
- A module can be uninstalled independently.
- Uninstall removes the module from the module registry.
- Uninstall removes the module from generated registration.
- Uninstall removes the module from go.work.
- Uninstall removes the module source directory from the codebase.
- A module can declare required and optional dependencies.
- The system prevents circular dependencies.
- The system prevents disabling a module required by another enabled module.
- The system prevents uninstalling a module required by another installed or enabled module.
- The system prevents unsafe source-code deletion outside /modules.
- Module migrations run in dependency order.
- The module registry table is named modules.
- Modules can expose public service contracts.
- Modules can register event handlers.
- Modules can register extension-point handlers.
- The root application can boot with a selected list of enabled modules.
- Authentication can depend on User.
- Authorization can depend on User.
- Module management can manage User, Authentication, and Authorization through declared contracts.
- Tests can run per module and per selected module group.
- The module lifecycle is auditable.
- The system does not require Go native plugin package.
- The system does not require HashiCorp go-plugin.
- The system can compile after install, enable, disable, and uninstall operations.
