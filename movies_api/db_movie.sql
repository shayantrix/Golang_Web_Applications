DROP TABLE IF EXISTS movie;
DROP TABLE IF EXISTS director;

-- First, create the 'director' table
CREATE TABLE director (
    ID_D int NOT NULL,
    FirstName varchar(100),
    LastName varchar(100),
    PRIMARY KEY (ID_D)
);

-- Then, create the 'movie' table referencing 'director'
CREATE TABLE movie (
    ID int NOT NULL,
    Isbn int,
    Title varchar(200),
    ID_D int NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (ID_D) REFERENCES director(ID_D) ON DELETE CASCADE
);

-- Insert directors first (to satisfy foreign key constraints)
INSERT INTO director(ID_D, FirstName, LastName) VALUES (1, "Edvard", "Holand");
INSERT INTO director(ID_D, FirstName, LastName) VALUES (2, "Shayan", "Amir Shahkarami");

-- Now, insert movies referencing existing directors
INSERT INTO movie(ID, Isbn, Title, ID_D) VALUES (1, 1235748, "Omni man: The return", 1);
INSERT INTO movie(ID, Isbn, Title, ID_D) VALUES (2, 4581564, "Invincible", 2);
