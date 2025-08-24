### **Tech Stack Guideline: Narrative Agent Application**

**Last Updated:** 24 August 2025

#### **1. Guiding Philosophy**

This stack is optimized for a single developer to rapidly build a high-performance, cross-platform, local-first desktop application. Every choice prioritizes **speed of feature development**, reliability, and a clean separation of concerns, minimizing time spent on configuration, styling, and boilerplate.

---

#### **2. Core Technologies**

| Layer | Technology | Rationale |
| :--- | :--- | :--- |
| **Backend Logic** | **Go** | **Performance & Concurrency.** Go's goroutines are perfectly suited for the multi-agent architecture. It compiles to a fast, single binary, ensuring a responsive core application. |
| **Desktop Framework** | **Wails** | **The Go-to-UI Bridge.** Wails is a lightweight framework designed to wrap a Go backend with a web frontend. It's more memory-efficient than alternatives like Electron and provides seamless bindings between Go and JavaScript. |
| **Frontend Framework** | **React** | **Maturity & Developer Familiarity.** A robust, mature library with a vast ecosystem. Your familiarity with React eliminates the learning curve, allowing you to build the UI immediately. |
| **UI Components & Styling** | **shadcn/ui & Tailwind CSS** | **Maximum Development Velocity.** This is the key to not focusing on styling. `shadcn/ui` provides beautiful, accessible, copy-and-pasteable components. Tailwind CSS enables rapid styling directly in the markup. This combination drastically reduces the time required to build a polished UI. |
| **Local Database** | **SQLite** | **Zero-Configuration & Reliability.** The ideal choice for a local-first application. The database is a single, portable file with no server process to manage. It's fast, robust, and battle-tested. |
| **Database Interaction** | **sqlc** | **Type-Safe SQL.** `sqlc` generates fully type-safe, idiomatic Go code from your raw SQL queries. This gives you the safety of an ORM with the simplicity and performance of writing pure SQL, preventing common bugs and speeding up data layer development. |
| **Context Management** | **Context Manager (Go)** | **Narrative-aware memory.** Builds task-scoped context bundles: token budgeting, prompt assembly (beats, characters, prior scenes), model selection. |
| **Local Vector DB** | **sqlite-vec** | **RAG without new infra.** Embeddings stored in SQLite using the sqlite-vec extension enable fast local similarity search per project. |

---

#### **3. Development Workflow Summary**

1.  **Backend:** Define core logic and agent behaviour in **Go**. Expose functions that can be called from the frontend.
2.  **Database:** Write `CREATE TABLE` and query `.sql` files. Run **sqlc** to generate type-safe Go methods for database access.
3.  **Frontend:** Use the **Wails** CLI to manage the project and run the development server.
4.  **UI Construction:** Build React components. When a UI element is needed (e.g., a button, dialog, or data table), add it from **shadcn/ui** and style it instantly with **Tailwind CSS**.
5.  **Integration:** Call Go functions directly from your React components via the Wails bridge to interact with the backend and the SQLite database.
6.  **Context/RAG:** Chunk and embed relevant documents (scenes, notes); store vectors with sqlite-vec; query via a Retriever to assemble task-specific context.
7.  **Model Selection:** Use a simple policy (task, complexity, budget) to choose between local Ollama and user-provided API keys.

This stack provides a cohesive and powerful environment, allowing you to focus your energy where it matters most: building a great product.