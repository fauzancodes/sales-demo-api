package controllers

import (
	"strings"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	if strings.ToLower(config.LoadConfig().Env) == "development" {
		return c.File("assets/html/index.html")
	}

	return c.HTML(200, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Sales Demo API Documentation</title>
				<style>
						body {
								font-family: Arial, sans-serif;
								background-color: #ECECEC;
								margin: 0;
								padding: 0;
						}
						header {
								background-color: #006BFF;
								color: white;
								padding: 20px;
								text-align: center;
						}
						main {
								padding: 40px 20px;
								text-align: center;
						}
						p {
								line-height: 1.6;
								color: #555;
						}
						.feature-list {
								text-align: left;
								display: grid;
								grid-template-columns: 1fr;
								gap: 20px;
								max-width: 100%;
								margin: 20px auto;
								padding: 0;
								list-style-type: disc;
						}
						.feature-list li {
								margin-left: 20px;
								color: #333;
						}
						@media screen and (min-width: 768px) {
								.feature-list {
										grid-template-columns: 1fr 1fr 1fr;
										max-width: 90%;
								}
						}
						.button-group {
								margin: 30px 0;
						}
						a.button {
								display: inline-block;
								margin: 10px;
								padding: 14px 24px;
								font-size: 16px;
								text-decoration: none;
								background: #006BFF;
								color: white;
								border-radius: 5px;
								transition: background 0.3s;
						}
						a.button:hover {
								background: #004bb3;
						}
						footer {
								background-color: #E0E0E0;
								color: #333;
								padding: 10px;
								text-align: center;
						}
				</style>
		</head>
		<body>
				<header>
						<h1>Sales Demo API Documentation</h1>
				</header>
				<main>
						<p>Welcome to the Sales Demo API documentation. This API is built using the <a href="https://go.dev/" target="_blank">Go</a> programming language with the <a href="https://echo.labstack.com/" target="_blank">Echo</a> framework, leveraging <a href="https://gorm.io/" target="_blank">GORM</a> as the ORM for database connections. <a href="https://www.postgresql.org/" target="_blank">PostgreSQL</a> is used as the main database, ensuring high performance and scalability. The API is designed to efficiently and securely support various data management features, including authentication, product management, sales, customers, and other functionalities to meet business needs.</p>
						<p>Below are the features of the API:</p>
						<ul class="feature-list">
								<li><strong>Secure Authentication:</strong> Effortlessly manage user access with robust login and registration endpoints, powered by <a href="https://jwt.io/" target="_blank">JWT</a> authentication for a seamless user experience.</li>
								<li><strong>Password Protection:</strong> Safeguard sensitive information with industry-standard <a href="https://en.wikipedia.org/wiki/Bcrypt" target="_blank">bcrypt</a> hashing, ensuring user passwords are encrypted and secure.</li>
								<li><strong>Comprehensive User Management:</strong> Empower users with endpoints for email verification, password recovery, accessing personal data, updating profiles, and deleting accounts with ease.</li>
								<li><strong>Email Verification Made Easy:</strong> Automatically send verification emails to users and seamlessly check their verification status to enhance security and trust.</li>
								<li><strong>Password Recovery Simplified:</strong> Streamline the password reset process by sending users a secure verification token via email, enabling them to create new passwords effortlessly.</li>
								<li><strong>Product Management at Your Fingertips:</strong> Enjoy complete control with CRUD endpoints for managing product categories, individual products, stock levels, and image uploads.</li>
								<li><strong>Stock Monitoring with History:</strong> Keep track of inventory changes with a comprehensive stock change history, ensuring accurate product management.</li>
								<li><strong>Customer Management:</strong> Easily manage customer information with dedicated CRUD endpoints, allowing you to maintain detailed records.</li>
								<li><strong>Sales Management Made Simple:</strong> Handle sales and detailed transactions with ease through dedicated CRUD endpoints for streamlined operations.</li>
								<li><strong>Sales Validation Built-In:</strong> Ensure accurate calculations with validation for subtotals, discounts, taxes, and totals, so you can focus on what matters most.</li>
								<li><strong>Real-Time Stock Validation:</strong> Stay informed with product stock validation and automatic stock reduction during sales, preventing overselling and enhancing customer satisfaction.</li>
								<li><strong>Invoicing Made Effortless:</strong> Automatically send invoices to customers via email, complete with downloadable links for easy access.</li>
								<li><strong>Seamless Product Image Uploads:</strong> Effortlessly upload product images to <a href="https://cloudinary.com/" target="_blank">Cloudinary</a>, ensuring your images are stored securely in the cloud while optimizing load times for a smooth and lightweight user experience.</li>
								<li><strong>Effortless Data Import:</strong> Seamlessly import product categories, products, and customer data using Excel or CSV files, simplifying data entry and migration for a smoother workflow.</li>
								<li><strong>Trusted Payment Gateway Integration:</strong> Leverage <a href="https://midtrans.com/" target="_blank">Midtrans</a>, Indonesia’s most complete payment gateway, to provide a secure and reliable payment solution, making transactions seamless and hassle-free for your customers.</li>
								<li><strong>Dynamic Data Retrieval:</strong> All GET endpoints feature dynamic pagination, filtering, searching, and sorting options, making data retrieval intuitive and efficient.</li>
								<li><strong>Robust Input Validation:</strong> Ensure data integrity with comprehensive validation on all POST endpoints, preventing errors and enhancing user experience.</li>
								<li><strong>SQL Injection Prevention:</strong> All endpoints is protected with query parameter sanitization, guarding against SQL injection attacks for enhanced security.</li>
								<li><strong>Automatic Database Migration:</strong> Enjoy worry-free database management with automatic migration, ensuring your table structures and relationships stay perfectly in sync.</li>
								<li><strong>Modular Design for Flexibility:</strong> Crafted with modularity in mind, our functions can be easily copied and adapted for use in your own <a href="https://go.dev/" target="_blank">Golang</a> projects, or even translated into other programming languages with ease.</li>
								<li><strong>Comprehensive Documentation:</strong> Navigate through the API with ease using clear, organized <a href="https://www.postman.com/" target="_blank">Postman</a> documentation designed to simplify testing and integration, ensuring a smooth development experience for all users.</li>
						</ul>
						<div class="button-group">
								<a href="/docs/Sale Demo API.postman_collection.json" download="Sale Demo API.postman_collection.json" class="button">Download Postman Collection</a>
								<a href="/docs/Sales Demo API.postman_environment.json" download="Sales Demo API.postman_environment.json" class="button">Download Postman Environment</a>
								<a href="https://github.com/fauzancodes/sales-demo-api" class="button" target="_blank">Go to GitHub Repository</a>
						</div>
				</main>
				<footer>
						<p>&copy; 2024 Sales Demo API Project. All rights reserved. By <a href="https://fauzancodes.id/" target="_blank">fauzancodes</a></p>
				</footer>
		</body>
		</html>
	`)
}
