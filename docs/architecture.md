# Project: Aegis RAG (Permission-Aware Vector Orchestrator)

## 1. Executive Summary & Vision
Aegis RAG is a high-performance, Go-based middleware orchestrator designed to solve the "Enterprise Gap" in Retrieval-Augmented Generation (RAG) systems. It enforces strict Data Governance and Role-Based Access Control (RBAC) *before* vector retrieval occurs.

The ultimate vision is to evolve this project into a comprehensive **Knowledge Governance Platform** that not only filters data by user permission but also intelligently governs context based on logical truth, data staleness, and negation.

## 2. The Core Problem (The "Why")
Current RAG architectures are built for accuracy, not security. They blindly pass user queries to Vector Databases.
* **The Post-Validation Flaw:** Retrieving all nearest neighbors and then hiding unauthorized documents causes massive data leakage and destroys the relevance of the LLM prompt (e.g., if top 5 results are blocked, the LLM gets nothing).
* **The HNSW / RLS Bottleneck:** Relying entirely on the database's Row-Level Security (RLS) forces the HNSW (Hierarchical Navigable Small World) graph to filter nodes dynamically. If a query hits a densely restricted graph, latency spikes exponentially because the algorithm has to "hop" over forbidden nodes.

## 3. The Architectural Solution (The "How")
Aegis implements the **Filter-Decorator Pattern (Identity-Carrying Queries)**. 

Rather than treating the Vector DB as an autonomous security boundary, Aegis acts as a Policy Decision Point (PDP).
1.  **Intercept:** It intercepts the natural language query and the user's identity token.
2.  **Resolve:** It looks up the user's claims (Roles, Departments, Clearances).
3.  **Decorate:** It translates those claims into a highly optimized database-specific Metadata Filter (e.g., `WHERE department IN ('engineering', 'public')`).
4.  **Execute:** It binds this strict filter to the vector search, narrowing the global search space so the HNSW algorithm only traverses nodes the user is explicitly allowed to see.

## 4. Technical Stack & Engineering Principles
* **Language:** Go (Golang). Chosen for high-throughput concurrency, strict typing, and efficient middleware routing.
* **Architecture Pattern:** Interface-driven, decoupled modules. The Engine must never know *how* permissions are stored or *which* database is used.
* **Core Principle 1:** Interfaces over implementations. Everything (Storage, Identity, Logging) must be swappable.
* **Core Principle 2:** Pre-filtering over post-filtering. Unauthorized data must never enter the application's RAM.

## 5. Implementation Roadmap (Step-by-Step)

### Phase 1: The Policy Core (In-Memory Prototype)
*Goal: Prove the Filter-Decorator logic works using pure Go, with no external dependencies.*
* **Milestone 1.1:** Define core domain structs (`UserContext`, `Document`).
* **Milestone 1.2:** Create the `PermissionProvider` interface and a mock implementation.
* **Milestone 1.3:** Build the `QueryPlanner` function that takes a User and outputs an allowed `Tags` array.
* **Milestone 1.4:** Create a CLI loop demonstrating a user with "Manager" clearance retrieving different documents than a "Guest".

### Phase 2: Vector Storage Integration
*Goal: Move from Go slices to a real mathematical vector space.*
* **Milestone 2.1:** Abstract the storage layer by creating a `VectorStore` interface.
* **Milestone 2.2:** Spin up `pgvector` (PostgreSQL) or `Qdrant` via Docker.
* **Milestone 2.3:** Write the translation layer that converts Aegis's internal `QueryPlan` into the specific JSONB/Metadata query language of the database.
* **Milestone 2.4:** Benchmark pre-filtering latency vs. full table scans.

### Phase 3: Identity & Staleness (The Real World)
*Goal: Handle dynamic permissions and hierarchical access.*
* **Milestone 3.1:** Implement a JWT Parser to extract `UserContext` securely from HTTP headers instead of mock strings.
* **Milestone 3.2:** Build the `SyncEngine`. An event-listener that accepts webhooks (e.g., "Folder A permission changed") and updates the metadata tags on all corresponding vector chunks in the database.
* **Milestone 3.3:** Implement recursive group unrolling (e.g., User is in "Engineering", "Engineering" has access to "Project X", therefore User has access to "Project X").

### Phase 4: Observability & Auditability
*Goal: Make the system enterprise-ready for CISOs and compliance.*
* **Milestone 4.1:** Implement structured logging (e.g., `slog` in Go).
* **Milestone 4.2:** Create an "Audit Trail" interceptor that logs: Timestamp, UserID, Query, Applied Filter, and Number of Chunks Retrieved/Blocked.

### Phase 5: Knowledge Governance Platform (Future Scope)
* **Milestone 5.1:** Integrate an LLM "Judge" or deterministic logic to evaluate document staleness (Temporal Filtering).
* **Milestone 5.2:** Implement contradiction resolution (if Chunk A and Chunk B conflict, apply priority rules before sending to the generative LLM).

## 6. Directory Structure (Starting Point)
```text
aegis-rag/
├── docs/
│   └── architecture.md       # This living document
├── main.go                   # CLI entry point and test harness
├── models.go                 # Core domain structs (User, Document)
├── engine.go                 # The Filter-Decorator logic
└── go.mod
