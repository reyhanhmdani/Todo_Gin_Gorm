CREATE DATABASE IF NOT EXISTS Gin_todo;
CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY 'Pastibisa';
GRANT ALL PRIVILEGES ON Gin_todo.* TO 'root'@'%';
