# D2 Diagrams Test Document

A comprehensive test document for D2 diagram rendering in markdown.

---

## 1. Basic Connections

```d2
server -> database: queries
database -> cache: sync
cache -> server: response
```

---

## 2. Shape Types

```d2
rect: Rectangle
oval: Oval {shape: oval}
cylinder: Database {shape: cylinder}
queue: Queue {shape: queue}
cloud: Cloud {shape: cloud}
person: User {shape: person}
diamond: Decision {shape: diamond}
```

---

## 3. Nested Containers

```d2
aws: AWS {
  vpc: VPC {
    web: Web Tier {
      lb: Load Balancer
    }
    app: App Tier {
      api: API Server
    }
    data: Data Tier {
      db: PostgreSQL {shape: cylinder}
    }
  }
}

aws.vpc.web.lb -> aws.vpc.app.api
aws.vpc.app.api -> aws.vpc.data.db
```

---

## 4. Styled Elements

```d2
primary: Primary Button {
  style: {
    fill: "#1976d2"
    font-color: white
    border-radius: 8
  }
}

secondary: Secondary Button {
  style: {
    fill: "#f5f5f5"
    stroke: "#9e9e9e"
  }
}

primary -> secondary: click {
  style: {
    stroke: "#ff5722"
    stroke-dash: 4
  }
}
```

---

## 5. Connection Types

```d2
a: Node A
b: Node B
c: Node C
d: Node D

a -> b: forward arrow
b <- c: backward arrow
c <-> d: bidirectional
a -- d: no arrow
```

---

## 6. Sequence Diagram

```d2
shape: sequence_diagram

user: User
browser: Browser
server: Server
db: Database

user -> browser: Enter credentials
browser -> server: POST /login
server -> db: Validate user
db -> server: User found
server -> browser: Set session cookie
browser -> user: Show dashboard
```

---

## 7. SQL Tables

```d2
users: users {
  shape: sql_table
  id: int {constraint: primary_key}
  name: varchar(100)
  email: varchar(255) {constraint: unique}
}

orders: orders {
  shape: sql_table
  id: int {constraint: primary_key}
  user_id: int {constraint: foreign_key}
  total: decimal
}

users.id <-> orders.user_id
```

---

## 8. Grid Layout

```d2
features: Features {
  grid-columns: 2
  
  a: Fast
  b: Secure
  c: Reliable
  d: Scalable
}
```

---

## 9. Microservices Architecture

```d2
direction: right

client: Client {shape: person}
gateway: API Gateway {shape: hexagon}

services: Services {
  auth: Auth Service
  users: User Service
  orders: Order Service
  payments: Payment Service
}

data: Data Stores {
  userdb: User DB {shape: cylinder}
  orderdb: Order DB {shape: cylinder}
  redis: Redis {shape: cylinder}
}

client -> gateway
gateway -> services.auth
gateway -> services.users
gateway -> services.orders
gateway -> services.payments

services.users -> data.userdb
services.orders -> data.orderdb
services.auth -> data.redis
```

---

## 10. Simple Flowchart

```d2
start: Start {shape: oval}
input: Get Input
check: Valid? {shape: diamond}
process: Process Data
error: Show Error
output: Display Result
end: End {shape: oval}

start -> input
input -> check
check -> process: yes
check -> error: no
error -> input
process -> output
output -> end
```

---

## 11. CI/CD Pipeline

```d2
direction: right

code: Code {
  icon: https://icons.terrastruct.com/dev%2Fgit.svg
}
build: Build
test: Test
deploy: Deploy
prod: Production {
  style.fill: "#c8e6c9"
}

code -> build -> test -> deploy -> prod
```

---

## 12. Class Diagram

```d2
Animal: Animal {
  shape: class
  +name: string
  +age: int
  +speak(): void
}

Dog: Dog {
  shape: class
  +breed: string
  +bark(): void
}

Cat: Cat {
  shape: class
  +color: string
  +meow(): void
}

Animal <-- Dog: extends
Animal <-- Cat: extends
```

---

*End of D2 diagrams test document*
