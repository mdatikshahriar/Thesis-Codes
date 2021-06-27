const express = require('express');
const app = express();
const cors = require('cors');
const dotenv = require('dotenv').config();
const routes = require('./home.js');

//Middleware
app.use(cors());
app.use(express.urlencoded({extended: true}));
app.use(express.json());

app.use('/', routes);

const PORT = process.env.PORT || 3050;

app.listen(PORT, console.log(`Server started on port ${PORT} in the link http://localhost:${PORT}/`));