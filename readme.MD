# Event Planning BackEnd Api for Klender App

Event Planning App adalah sebuah aplikasi yang digunakan untuk membuat dan management sebuah event.

Membuat API menggunakan bahasa Go dengan menerapkan layered architecture untuk memudahkan pemeliharaan dan peningkatan Fitur

## Acknowledgements

- [Wireframe Example](https://whimsical.com/event-planning-app-RnJLuh3osvawc733NVaQzE)
- [Open API (Referensi untuk Graphql)](https://app.swaggerhub.com/apis-docs/iswanulumam/EventPlanningApp/1.0.0)
- [Gqlgen](https://gqlgen.com)
- [ERD Eventapp API](https://app.diagrams.net/#G1Tgzu7nNEl_V5AEURr1VyVNsI7bc7NeAL)
- [Task Docs](https://docs.google.com/document/d/1SubVXXZ-1KNGTj6s9EyO7rXUd3Pe8mPOMjuhxPy2REQ/edit#)

## Tech Stack

- Go Language
- Labstack Echo
- Mysql
- Docker

## Features

- Bcrypt
- JWT
- GraphQL
- CRUD User
- CRUD Event
- Join/Going Event
- Post and Read Comment
<!--

## Queries for _client_ 🏄🏻‍♀️

### User 🧛🏾‍♂️

#### Register

- _Mutation_

  ```graphql
  mutation registerNewUser($input: UserRegister!) {
  	register(input: $input) {
  		result
  	}
  }
  ```

  ```jsx
    // Query variable example

    {
      "input": {
        "username": "kholid",
        "email": "fikr@gmail.com",
        "password": "123456"
      }
    }
  ```

#### Login

- _Mutation_

  ```graphql
  mutation loginUser($input: UserLogin!) {
  	login(input: $input) {
  		result
  	}
  }
  ```

  ```jsx
    // Query variable example

    {
      "input": {
        "email": "fikr@gmail.com",
        "password": "123456"
      }
    }
  ```

#### Logout

- _Mutation_
  ```graphql
  mutation logoutUser {
  	logout {
  		result
  	}
  }
  ```

### Event 🎸

#### New Event

- _Mutation_

  ```graphql
  mutation createEvent($event: EventInput!) {
  	createEvent(event: $event) {
  		id
  		title
  		description
  		timestamp
  	}
  }
  ```

  ```jsx
    // Query variable example

    {
      "event": {
        "title": "tekken 7",
        "description": "Abis main, pergi ke pasar buat nyari makan sama temen g",
        "timestamp": "1589814419531"
      }
    }
  ```

#### Get All Events

- _Query_
  ```graphql
  query allEvents {
  	getEvents {
  		id
  		email
  		title
  		description
  		timestamp
  	}
  }
  ```

#### Get Event by ID

- _Query_

  ```graphql
  query eventById($id: ID!) {
  	getEventsById(id: $id) {
  		id
  		title
  		description
  		timestamp
  	}
  }
  ```

  ```jsx
    // Query variable example

    {
      "id": "5ec2a500e1a9271cba37af03"
    }
  ```

#### Edit Event by ID

- _Mutation_

  ```graphql
  mutation editEvent($id: ID!, $event: EventInput!) {
  	editEvent(id: $id, event: $event) {
  		id
  		title
  		description
  		timestamp
  	}
  }
  ```

  ```jsx
    // Query variable example

    {
      "id": "5ec2a500e1a9271cba37af03",
      "event": {
        "title": "trine 4",
        "description": "Ya gitu dehhh",
        "timestamp": "1589814733208"
      }
    }
  ```

#### Delete Event

- _Mutation_

  ```graphql
  mutation deleteEvent($id: ID!) {
  	deleteEvent(id: $id) {
  		result
  	}
  }
  ```

  ````jsx
    // Query variable example

    {
      "id": "5ec2a5c7e1a9271cba37af04"
    }
  ``` -->
  ````

## Authors

- [@bagusbpg](https://www.github.com/bagusbpg)
- [@zakariyahya](https://www.github.com/zakariyahya)
- [@najibjodiansyah](https://www.github.com/najibjodiansyah)

## Environment Variables

To run this project, you will need to set up the following environment variables to your ./config/config.go file

```bash
defaultConfig.Port
defaultConfig.Database.Driver
defaultConfig.Database.Name
defaultConfig.Database.Address
defaultConfig.Database.Port
defaultConfig.Database.Username
defaultConfig.Database.Password
```

## Run Locally

Clone the project

```bash
  git clone https://github.com/najibjodiansyah/eventapp.git
```

Go to the project directory

```bash
  cd eventapp
```

Install dependencies

```bash
  go mod tidy
```

Start the server

```bash
  go run ./app/main.go
```

## Running Tests

To run tests, run the following command

```bash
  go test ./...
```

## Used By

This project is used by the following companies:

- Group 5 Projet Event Planning App : Klender

## Lessons Learned

What did you learn while building this project? What challenges did you face and how did you overcome them?

## Feedback

If you have any feedback, please reach out to us at najibjodiansyah11@gmail.com

## Support

For support, email najibjodiansyah11@gmail.com or join our Discord channel.
