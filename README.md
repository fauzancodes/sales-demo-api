# Sales Demo API Documentation

Welcome to the Sales Demo API documentation. This API is built using the [Go](https://go.dev/) programming language with the [Echo](https://echo.labstack.com/) framework, leveraging [GORM](https://gorm.io/) as the ORM for database connections. [PostgreSQL](https://www.postgresql.org/) is used as the main database, ensuring high performance and scalability. The API is designed to efficiently and securely support various data management features, including authentication, product management, sales, customers, and other functionalities to meet business needs.

## Features

- **Secure Authentication:** Manage user access with login, register, email verification, and password recovery endpoints, powered by [JWT](https://jwt.io/) authentication and [bcrypt](https://en.wikipedia.org/wiki/Bcrypt) hashing.
- **Comprehensive User Management:** Endpoints for email verification, password recovery, accessing personal data, updating profiles, and deleting accounts.
- **Product Management:** CRUD endpoints for managing product categories, individual products, stock levels, and image uploads.
- **Real-Time Stock Monitoring:** Track inventory changes with a comprehensive stock change history and automatic stock validation.
- **Seamless Product Image Uploads:** Upload product images to [Cloudinary](https://cloudinary.com/) for secure and optimized cloud storage.
- **Customer Management:** CRUD endpoints for managing customer information.
- **Sales Management:** Handle sales and transactions with CRUD endpoints.
- **Sales Validation:** Validate subtotals, discounts, taxes, and totals for accurate calculations.
- **Invoicing:** Automatically send invoices to customers via email, including downloadable links.
- **Data Import & Export:** Import/export product categories, products, and customer data using Excel or CSV files.
- **Payment Gateway Integration:** Secure payment solutions with [Midtrans](https://midtrans.com/), [IPaymu](https://ipaymu.com/), and [Xendit](https://www.xendit.co/).
- **Dynamic Data Retrieval:** GET endpoints support dynamic pagination, filtering, searching, and sorting.
- **Robust Input Validation:** Comprehensive validation on all POST endpoints.
- **SQL Injection Prevention:** Query parameter sanitization to guard against SQL injection attacks.
- **Secure One-Time API Key:** Protect endpoints with HMAC and SHA-256 generated single-use API keys.
- **Automatic Database Migration:** Automatic table structure and relationship synchronization.
- **Modular and Flexible Architecture:** Support for containerization with [Docker](https://www.docker.com/) and [Kubernetes](https://kubernetes.io/).
- **Comprehensive Documentation:** Clear [Postman](https://www.postman.com/) documentation for simplified testing.

## Additional Notes

### For Frontend Developers
- You can use `SPECIAL_API_KEY (Uh/UB%SKft3CU3e0zJAvBhp3cyo/un2021/zLQf1BKGZZuQ6w5P9VAM6Sj0CcQCm)`, put it directly in the http request header as `X-Api-Key`.
- Alternatively, if you want to try the One-Time API Key feature, the way to create the `X-Api-Key` are:
  1. Generate a random string.
  2. Calculate the HMAC signature between the random string and the `HMAC_KEY (dI62Fk_8wb2uL8CLmSLFkDoAO/tfDeod)` using SHA-256.
  3. The result of the HMAC calculation is combined with the random string with the pattern `random_string:hmac_result`.
  4. Then, encode the pattern with the Base64 algorithm, the endcode result is the `X-Api-key`.

### For Backend Developers
- if you don't want to use the  One-Time API key feature, don't forget to set `SPECIAL_API_KEY` in .env to the request header as `X-Api-Key` or change `ENABLE_API_KEY` in .env to `false` or you will not be able to access all endpoints at all.
- And don't forget to:
  - Crete [Cloudinary](https://cloudinary.com/) account for image uploads.
  - Create [Backblaze](https://www.backblaze.com/) account for file system needs.
  - Set up [Gmail SMTP](https://www.digitalocean.com/community/tutorials/how-to-use-google-s-smtp-server) for email sending.

---

© 2024 Sales Demo API Project. All rights reserved. By [fauzancodes](https://fauzancodes.id/)
