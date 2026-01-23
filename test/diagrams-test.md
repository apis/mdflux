# Diagram Test Document

A comprehensive test document for various diagram syntaxes in markdown.

---

## 1. Mermaid - Flowchart

```mermaid
flowchart TD
    A[Start] --> B{Is it working?}
    B -->|Yes| C[Great!]
    B -->|No| D[Debug]
    D --> B
    C --> E[End]
```

---

## 2. Mermaid - Sequence Diagram

```mermaid
sequenceDiagram
    participant U as User
    participant F as Frontend
    participant A as API
    participant D as Database

    U->>F: Click Login
    F->>A: POST /auth
    A->>D: Query user
    D-->>A: User data
    A-->>F: JWT Token
    F-->>U: Redirect to Dashboard
```

---

## 3. Mermaid - Class Diagram

```mermaid
classDiagram
    class Animal {
        +String name
        +int age
        +makeSound()
    }
    class Dog {
        +String breed
        +bark()
    }
    class Cat {
        +String color
        +meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
```

---

## 4. Mermaid - State Diagram

```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Processing: Start
    Processing --> Success: Complete
    Processing --> Error: Fail
    Error --> Idle: Reset
    Success --> [*]
```

---

## 5. Mermaid - Entity Relationship

```mermaid
erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE_ITEM : contains
    PRODUCT ||--o{ LINE_ITEM : includes
    CUSTOMER {
        int id PK
        string name
        string email
    }
    ORDER {
        int id PK
        int customer_id FK
        date created_at
    }
    PRODUCT {
        int id PK
        string name
        decimal price
    }
```

---

## 6. Mermaid - Gantt Chart

```mermaid
gantt
    title Project Timeline
    dateFormat YYYY-MM-DD
    section Planning
        Requirements    :a1, 2024-01-01, 7d
        Design          :a2, after a1, 5d
    section Development
        Backend         :b1, after a2, 14d
        Frontend        :b2, after a2, 14d
    section Testing
        QA Testing      :c1, after b1, 7d
        UAT             :c2, after c1, 5d
```

---

## 7. Mermaid - Pie Chart

```mermaid
pie showData
    title Browser Market Share
    "Chrome" : 65
    "Safari" : 19
    "Firefox" : 8
    "Edge" : 5
    "Other" : 3
```

---

## 8. Mermaid - Git Graph

```mermaid
gitGraph
    commit id: "Initial"
    branch develop
    checkout develop
    commit id: "Feature A"
    commit id: "Feature B"
    checkout main
    merge develop id: "Release v1.0"
    branch hotfix
    commit id: "Bugfix"
    checkout main
    merge hotfix id: "v1.0.1"
```

---

## 9. Mermaid - Mind Map

```mermaid
mindmap
  root((Project))
    Frontend
      React
      TypeScript
      Tailwind
    Backend
      Node.js
      PostgreSQL
      Redis
    DevOps
      Docker
      Kubernetes
      CI/CD
```

---

## 10. Mermaid - Timeline

```mermaid
timeline
    title Company History
    2020 : Founded
         : First Product
    2021 : Series A
         : 50 Employees
    2022 : International Expansion
         : 200 Employees
    2023 : IPO
```

---

## 11. Mermaid - Quadrant Chart

```mermaid
quadrantChart
    title Feature Priority Matrix
    x-axis Low Effort --> High Effort
    y-axis Low Impact --> High Impact
    quadrant-1 Do First
    quadrant-2 Plan
    quadrant-3 Delegate
    quadrant-4 Eliminate
    Feature A: [0.8, 0.9]
    Feature B: [0.3, 0.8]
    Feature C: [0.7, 0.3]
    Feature D: [0.2, 0.2]
```

---

## 12. Mermaid - User Journey

```mermaid
journey
    title User Onboarding Experience
    section Sign Up
      Visit website: 5: User
      Fill form: 3: User
      Verify email: 4: User
    section First Use
      Complete tutorial: 4: User
      Create first project: 5: User
    section Engagement
      Invite team: 3: User
      Upgrade plan: 4: User
```

---

## 13. D2 - Architecture Diagram

```d2
direction: right

client: Client {shape: person}
cdn: CDN {shape: cloud}
lb: Load Balancer {shape: hexagon}

backend: Backend {
  api1: API Server 1
  api2: API Server 2
}

data: Data Layer {
  db: PostgreSQL {shape: cylinder}
  cache: Redis {shape: cylinder}
}

client -> cdn -> lb
lb -> backend.api1
lb -> backend.api2
backend.api1 -> data.db
backend.api2 -> data.db
backend.api1 -> data.cache
backend.api2 -> data.cache
```

---

## 14. D2 - Sequence Diagram

```d2
shape: sequence_diagram

user: User
app: Mobile App
api: REST API
db: Database

user -> app: Open App
app -> api: GET /profile
api -> db: SELECT * FROM users
db -> api: User record
api -> app: JSON response
app -> user: Display profile
```

---

## 15. D2 - Flowchart with Styles

```d2
start: Start {
  shape: oval
  style.fill: "#e8f5e9"
}

process: Process Data {
  style.fill: "#e3f2fd"
}

decision: Valid? {
  shape: diamond
  style.fill: "#fff3e0"
}

success: Success {
  style.fill: "#c8e6c9"
}

error: Error {
  style.fill: "#ffcdd2"
}

end: End {
  shape: oval
  style.fill: "#f3e5f5"
}

start -> process -> decision
decision -> success: yes
decision -> error: no
error -> process: retry
success -> end
```

---

## 16. PlantUML - Component Diagram

```plantuml
@startuml
package "Frontend" {
  [Web App]
  [Mobile App]
}

package "Backend" {
  [API Gateway]
  [Auth Service]
  [User Service]
  [Order Service]
}

database "Data" {
  [PostgreSQL]
  [Redis]
}

[Web App] --> [API Gateway]
[Mobile App] --> [API Gateway]
[API Gateway] --> [Auth Service]
[API Gateway] --> [User Service]
[API Gateway] --> [Order Service]
[User Service] --> [PostgreSQL]
[Order Service] --> [PostgreSQL]
[Auth Service] --> [Redis]
@enduml
```

---

## 17. PlantUML - Activity Diagram

```plantuml
@startuml
start
:Receive Order;
if (In Stock?) then (yes)
  :Process Payment;
  if (Payment OK?) then (yes)
    :Ship Order;
    :Send Confirmation;
  else (no)
    :Cancel Order;
    :Notify Customer;
  endif
else (no)
  :Backorder Item;
  :Notify Customer;
endif
stop
@enduml
```

---

## 18. PlantUML - Use Case Diagram

```plantuml
@startuml
left to right direction
actor Customer
actor Admin

rectangle "E-Commerce System" {
  Customer --> (Browse Products)
  Customer --> (Add to Cart)
  Customer --> (Checkout)
  Customer --> (View Orders)
  
  Admin --> (Manage Products)
  Admin --> (Process Orders)
  Admin --> (View Reports)
}
@enduml
```

---

*End of diagram test document*
