# Go Task Manager API

This project provides a simple REST API for managing tasks. It is built using Go, MySQL, and the Gorilla Mux router.

## Getting Started

To get started, you will need Go installed on your local machine. MySQL will also need to be set up and running. Clone the repository to your local machine.

```bash
git clone https://github.com/ahmdalgendi/todo-list.git
```

## Running the Application

Navigate to the cloned repository and run the following command:

```bash
go run main.go
```

The application will start running at `http://localhost:8080`.

## API Endpoints

- `GET /tasks`: Fetch all tasks.
- `POST /tasks`: Create a new task.
- `GET /tasks/{id}`: Fetch a specific task by ID.
- `PUT /tasks/{id}`: Update a specific task by ID.
- `DELETE /tasks/{id}`: Delete a specific task by ID.

## TODO

Here are some improvements to be made:

- [ ]  Add more detailed error handling for better debugging.
- [ ]  Set maximum connections, connection lifetimes, etc. for the database.
- [ ]  Hide sensitive details such as the database connection string using environment variables or a secure vault.
- [ ]  Add separate routers for each resource.
- [ ]  Return a status code of `201 Created` when a task is successfully created.
- [ ]  Add validation for incoming requests using a package like `go-playground/validator`.
- [ ]  Add additional struct tags for SQL DB to control how Go struct fields are translated to SQL fields.
- [ ]  Separate the database operations into another package.
- [ ]  Add more comments to explain what each function does and what the expected inputs and outputs are.
- [ ] Add tests using Go's built-in testing package to ensure the app works as expected.
- [ ] Add middleware for logging, authentication, CORS, etc.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
