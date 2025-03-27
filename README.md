# AWS RDS AND DYNAMO DEMO

## 👋 Welcome!

This project is a demonstration of how to build a full-stack web application that interacts with **two different types of databases**: a relational database (MySQL managed by AWS RDS) and a NoSQL database (Amazon DynamoDB). It showcases how a single frontend can seamlessly communicate with both to manage different sets of data.

This demo features:

* **A Go-based backend for RDS (MySQL):** Manages Pokémon data with full CRUD (Create, Read, Update, Delete) operations.
* **A Go-based backend for DynamoDB:** Manages Plant data, also with complete CRUD functionality.
* **A React-based frontend:** Provides a user interface to interact with both the Pokémon and Plant data through the respective backends.

This project is intended to illustrate the process of integrating different database technologies within a single application, rather than comparing their specific performance or use cases.

## 🛠️ Technologies Used

* **Go:** For building the backend APIs.
* **React:** For creating the interactive user interface.
* **MySQL:** Relational database managed by AWS RDS.
* **Amazon DynamoDB:** NoSQL database service.
* **Gorilla Mux:** For routing in the Go backends.
* **`go-sql-driver/mysql`:** MySQL driver for Go.
* **AWS SDK for Go v2:** For interacting with DynamoDB.
* **Vite:** For the React frontend development.
* **Bun:** Package manager for the frontend.

## ⚙️ Setup Instructions

Follow these steps to get the `aws-rds-and-dynamo-demo` project up and running on your local machine:

### 1. Clone the Repository

```bash
git clone https://github.com/vivalchemy/aws-rds-and-dynamo-demo.git
cd aws-rds-and-dynamo-demo
```

### 2. Configure Environment Variables

You need to set up `.env` files for both the `rds` and `dynamo-db` backends with your database credentials.

#### RDS Configuration (`rds/.env`)

1.  Navigate to the `rds` directory:
    ```bash
    cd rds
    ```
2.  Copy the example environment file:
    ```bash
    cp .env.example .env
    ```
3.  Open the `.env` file and replace the placeholder values with your actual AWS RDS MySQL database credentials:

    ```env
    MYSQL_USER=your_rds_username
    MYSQL_PASSWORD=your_rds_password
    MYSQL_HOST=your_rds_endpoint
    MYSQL_PORT=your_rds_port
    MYSQL_DATABASE=your_rds_database_name
    ```

#### DynamoDB Configuration (`dynamo-db/.env`)

1.  Navigate to the `dynamo-db` directory:
    ```bash
    cd ../dynamo-db
    ```
2.  Copy the example environment file:
    ```bash
    cp .env.example .env
    ```
3.  Open the `.env` file. You might need to configure your AWS credentials and region here, depending on your setup. Ensure your AWS environment is configured to allow access to DynamoDB in the specified region (currently `ap-south-1` in the code).

    ```env
    AWS_REGION=your_aws_region # e.g., ap-south-1
    # AWS_ACCESS_KEY_ID=your_access_key # Optional if AWS CLI or IAM roles are configured
    # AWS_SECRET_ACCESS_KEY=your_secret_key # Optional if AWS CLI or IAM roles are configured
    PORT=8081 # Optional: You can change the default port for the DynamoDB backend
    ```

### 3. Configure Frontend API URL

The frontend needs to know where to reach your backend servers.

1.  Navigate to the `frontend` directory:
    ```bash
    cd ../frontend
    ```
2.  Open the `frontend/src/lib/config.ts` file.
3.  Update the `apiURL` variable to the base URL where your Go backends will be hosted. **Make sure this URL is accessible from your frontend development environment.**

    ```typescript
    // frontend/src/lib/config.ts
    export const apiURL = "http://your-backend-host"; // Replace with your backend URL
    ```

    **Note:** If you are running the backends locally, this might be something like `http://localhost:8080` (for RDS) and `http://localhost:8081` (for DynamoDB), and you might need to adjust your `config.ts` to handle both. Alternatively, you can host both backends under a common base URL and configure your routing accordingly.

### 4. Run the Backends(tested by deploying on ec2)

Open two separate terminal windows to run the RDS and DynamoDB backends concurrently.

#### RDS Backend

1.  Navigate to the `rds` directory:
    ```bash
    cd rds
    ```
2.  Build and run the Go application:
    ```bash
    go build -o rds-server main.go
    ./rds-server
    ```

    You should see a message indicating that the server is running on port 8080 and connected to the RDS MySQL database.

#### DynamoDB Backend

1.  Navigate to the `dynamo-db` directory:
    ```bash
    cd ../dynamo-db
    ```
2.  Build and run the Go application:
    ```bash
    go build -o dynamo-server main.go
    ./dynamo-server
    ```

    You should see a message indicating that the server is running on port 8081 (or the port you configured in `.env`) and that the Plants table is being created or already exists.

### 5. Run the Frontend(tested by deploying on beanstalk using docker)

1.  Navigate to the `frontend` directory:
    ```bash
    cd ../frontend
    ```
2.  Install the frontend dependencies:
    ```bash
    bun install # Or npm install or yarn install
    ```
3.  Start the frontend development server:
    ```bash
    bun run dev # Or npm run dev or yarn dev
    ```

    The frontend application should now be accessible in your browser, usually at `http://localhost:5173`.

## 🚀 Using the Application

Once the frontend is running, you should be able to:

* View a list of Pokémon fetched from the RDS MySQL database.
* View details of individual Pokémon.
* Create, update, and delete Pokémon.

* View a list of Plants fetched from the DynamoDB database.
* View details of individual Plants.
* Create, update, and delete Plants.

The navigation within the frontend should allow you to switch between the Pokémon and Plant management sections.

## 💡 Contributing

Feel free to contribute to this project by opening issues or submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
