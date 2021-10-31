-- CreateTablePost
CREATE TABLE post (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title TEXT NOT NULL,
    content LONGTEXT,
    published BOOLEAN NOT NULL DEFAULT false,
    authorId INT NOT NULL
);

-- CreateTableUser
CREATE TABLE user (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    pass VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(100)
);

-- CreateIndex
CREATE UNIQUE INDEX email_unique ON user (email);