# Sales Demo API Documentation

Welcome to the Sales Demo API documentation. This API is built using the [Go](https://go.dev/) programming language with the [Echo](https://echo.labstack.com/) framework, leveraging [GORM](https://gorm.io/) as the ORM for database connections. [PostgreSQL](https://www.postgresql.org/) is used as the main database, ensuring high performance and scalability. The API is designed to efficiently and securely support various data management features, including authentication, product management, sales, customers, and other functionalities to meet business needs.

## Features of the API

- **Secure Authentication**: Effortlessly manage user access with login, register, email verification, and password recovery endpoints, powered by [JWT](https://jwt.io/) authentication and industry-standard [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) hashing for a secure user experience.
- **Comprehensive User Management**: Empower users with endpoints for email verification, password recovery, accessing personal data, updating profiles, and deleting accounts with ease.
- **Product Management at Your Fingertips**: Enjoy complete control with CRUD endpoints for managing product categories, individual products, stock levels, and image uploads.
- **Real-Time Stock Monitoring**: Keep track of inventory changes with a comprehensive stock change history and automatic stock validation, ensuring accurate product management and preventing overselling.
- **Seamless Product Image Uploads**: Effortlessly upload product images to [Cloudinary](https://cloudinary.com/), ensuring your images are stored securely in the cloud while optimizing load times for a smooth and lightweight user experience.
- **Customer Management Made Easy**: Easily manage customer information with dedicated CRUD endpoints, allowing you to maintain detailed records.
- **Sales Management Made Simple**: Handle sales and detailed transactions with ease through dedicated CRUD endpoints for streamlined operations.
- **Sales Validation Built-In**: Ensure accurate calculations with validation for subtotals, discounts, taxes, and totals, so you can focus on what matters most.
- **Invoicing Made Effortless**: Automatically send invoices to customers via email, complete with downloadable links for easy access.
- **Effortless Data Import & Export**: Seamlessly import & export product categories, products, and customer data using Excel or CSV files, simplifying data entry and migration for a smoother workflow.
- **Trusted Payment Gateway Integration**: Leverage [Midtrans](https://midtrans.com/), [IPaymu](https://ipaymu.com/), and [Xendit](https://www.xendit.co/), Indonesia’s most complete payment gateways, to provide a secure and reliable payment solution, making transactions seamless and hassle-free for your customers.
- **Dynamic Data Retrieval**: All GET endpoints feature dynamic pagination, filtering, searching, and sorting options, making data retrieval intuitive and efficient.
- **Robust Input Validation**: Ensure data integrity with comprehensive validation on all POST endpoints, preventing errors and enhancing user experience.
- **SQL Injection Prevention**: All endpoints are protected with query parameter sanitization, guarding against SQL injection attacks for enhanced security.
- **Secure One-Time API Key**: Each request to the endpoints is protected with a single-use API key generated using [HMAC](https://en.wikipedia.org/wiki/HMAC) and [SHA-256](https://en.wikipedia.org/wiki/SHA-2) methods, ensuring robust security and integrity for every transaction while safeguarding your data from unauthorized access.
- **Automatic Database Migration**: Enjoy worry-free database management with automatic migration, ensuring your table structures and relationships stay perfectly in sync.
- **Modular And Flexible Architecture Support**: Crafted with modularity in mind, automating deployment, scaling, and management of containerized applications, easily managed through [Docker](https://www.docker.com/) and [Kubernetes](https://kubernetes.io/) for a streamlined deployment and maintenance experience.
- **Comprehensive Documentation**: Navigate through the API with ease using clear, organized [Postman](https://www.postman.com/) documentation designed to simplify testing and ensure a smooth development experience for all developers.

## Downloads and Resources

- [Download Postman Collection](./docs/Sale%20Demo%20API.postman_collection.json)
- [Download Postman Environment](./docs/Sales%20Demo%20API.postman_environment.json)
- [Go to GitHub Repository](https://github.com/fauzancodes/sales-demo-api)

---

**Additional Note:**  
Don't forget to set `SPECIAL_API_KEY` in `.env` to the request header as `X-Api-Key` or change `ENABLE_API_KEY` in `.env` to `false`, or you will not be able to access all endpoints. For those accessing the endpoint not through API documentation but from the frontend application, directly implement the use of the single-use `X-Api-Key`. To create the `X-Api-Key`, generate a random string, calculate the HMAC between the random string and the `HMAC_KEY` in `.env` with the SHA-256 algorithm. Combine the random string with the HMAC result in the pattern `random_string: hmac_result`, encode the pattern with Base64, and the encoded result is the `X-Api-key`.

---

&copy; 2024 Sales Demo API Project. All rights reserved. By [fauzancodes](https://fauzancodes.id/)
