CREATE TABLE MedicalRecords (
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    age INT,
    sex VARCHAR(10),
    blood_type VARCHAR(3),
    height_cm INT,
    weight_kg INT,
    bmi FLOAT,
    blood_pressure VARCHAR(15),
    heart_rate INT,
    respiratory_rate INT,
    temperature_c FLOAT,
    blood_glucose INT,
    cholesterol INT,
    has_diabetes INT,
    has_heart_disease INT,
    has_asthma INT,
    has_kidney_disease INT,
    has_liver_disease INT,
    has_cancer INT
);
INSERT INTO MedicalRecords VALUES ('David', 'Davis', 64, 'Male', 'A+', 171, 56, 19.2, '123/72', 94, 17, 37.1, 147, 226, 0, 0, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Davis', 34, 'Female', 'O+', 198, 119, 30.4, '111/84', 74, 15, 36.7, 118, 207, 1, 1, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Davis', 24, 'Female', 'A+', 154, 120, 50.6, '134/89', 83, 19, 36.0, 155, 229, 1, 0, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Martinez', 47, 'Male', 'O-', 167, 52, 18.6, '129/81', 79, 12, 36.2, 90, 190, 0, 1, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Rodriguez', 77, 'Male', 'O+', 184, 56, 16.5, '102/74', 86, 17, 37.0, 136, 235, 1, 0, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Garcia', 50, 'Female', 'A+', 173, 78, 26.1, '121/83', 63, 17, 36.2, 85, 191, 1, 1, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('John', 'Smith', 69, 'Male', 'B-', 151, 69, 30.3, '130/74', 92, 16, 36.3, 140, 191, 0, 0, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Martinez', 27, 'Male', 'O+', 178, 51, 16.1, '109/68', 80, 18, 36.7, 141, 155, 0, 1, 0, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Davis', 39, 'Female', 'AB+', 192, 62, 16.8, '131/74', 80, 16, 37.4, 106, 198, 0, 0, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Johnson', 89, 'Female', 'B-', 179, 82, 25.6, '107/88', 86, 18, 37.4, 156, 209, 1, 1, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Martinez', 44, 'Female', 'A-', 158, 97, 38.9, '135/65', 70, 19, 36.6, 82, 194, 0, 0, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Miller', 27, 'Female', 'A+', 181, 61, 18.6, '136/64', 87, 17, 37.2, 136, 184, 1, 0, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Williams', 52, 'Male', 'B+', 179, 85, 26.5, '106/69', 75, 17, 37.7, 122, 192, 1, 0, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Martinez', 42, 'Female', 'A-', 178, 50, 15.8, '119/62', 69, 20, 37.3, 135, 169, 0, 0, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Rodriguez', 62, 'Male', 'O+', 161, 51, 19.7, '119/73', 91, 19, 36.4, 138, 215, 0, 0, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Rodriguez', 81, 'Female', 'AB-', 191, 116, 31.8, '137/64', 92, 12, 36.5, 114, 235, 0, 1, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Miller', 72, 'Female', 'O+', 199, 63, 15.9, '101/82', 84, 17, 37.4, 145, 224, 0, 0, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Davis', 40, 'Female', 'A+', 171, 89, 30.4, '109/70', 72, 14, 37.9, 122, 150, 0, 0, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Laura', 'Johnson', 22, 'Male', 'B-', 196, 79, 20.6, '130/70', 63, 15, 36.2, 157, 168, 0, 0, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Johnson', 89, 'Male', 'A+', 192, 98, 26.6, '101/60', 61, 12, 36.3, 148, 205, 1, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Laura', 'Smith', 18, 'Male', 'O+', 187, 61, 17.4, '136/60', 60, 18, 37.3, 83, 188, 0, 0, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Laura', 'Johnson', 60, 'Female', 'AB+', 165, 90, 33.1, '115/85', 92, 17, 36.8, 149, 167, 1, 1, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Garcia', 40, 'Male', 'A+', 151, 106, 46.5, '104/67', 85, 15, 37.5, 126, 248, 1, 1, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Robert', 'Garcia', 57, 'Male', 'O+', 154, 62, 26.1, '136/67', 84, 20, 36.8, 93, 154, 1, 1, 1, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Laura', 'Garcia', 28, 'Female', 'O+', 155, 65, 27.1, '137/64', 93, 20, 36.7, 72, 175, 0, 1, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Garcia', 77, 'Female', 'AB+', 200, 93, 23.2, '139/75', 84, 19, 37.4, 96, 203, 0, 0, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Miller', 79, 'Female', 'B-', 194, 71, 18.9, '139/82', 93, 15, 37.2, 101, 209, 0, 1, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Garcia', 42, 'Male', 'O-', 194, 106, 28.2, '114/73', 61, 19, 36.1, 160, 187, 1, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('James', 'Jones', 55, 'Male', 'AB-', 183, 71, 21.2, '100/66', 87, 14, 36.8, 127, 150, 1, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Martinez', 72, 'Male', 'O-', 194, 85, 22.6, '105/75', 67, 18, 38.0, 112, 217, 0, 1, 0, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Williams', 21, 'Male', 'A-', 180, 93, 28.7, '109/71', 73, 19, 36.4, 88, 243, 0, 0, 0, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Davis', 82, 'Female', 'AB+', 157, 118, 47.9, '115/62', 62, 18, 36.5, 93, 157, 1, 1, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Johnson', 25, 'Male', 'AB+', 171, 88, 30.1, '114/81', 88, 17, 37.9, 145, 235, 1, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Brown', 71, 'Male', 'B+', 158, 115, 46.1, '137/69', 73, 20, 36.1, 87, 184, 1, 1, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Rodriguez', 55, 'Female', 'O-', 193, 86, 23.1, '121/60', 76, 15, 36.1, 144, 154, 1, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Williams', 54, 'Female', 'O-', 199, 113, 28.5, '118/60', 63, 20, 37.4, 157, 225, 0, 0, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('James', 'Miller', 46, 'Female', 'AB-', 199, 96, 24.2, '113/78', 69, 20, 36.8, 83, 152, 1, 1, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Miller', 69, 'Female', 'O-', 186, 106, 30.6, '136/74', 86, 12, 37.0, 149, 232, 1, 0, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Rodriguez', 74, 'Male', 'B-', 165, 84, 30.9, '126/69', 86, 17, 37.2, 87, 238, 0, 0, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Johnson', 62, 'Male', 'O-', 176, 105, 33.9, '137/88', 97, 12, 37.7, 149, 197, 1, 0, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Miller', 55, 'Female', 'O+', 173, 102, 34.1, '113/60', 70, 16, 36.6, 120, 238, 1, 1, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('John', 'Rodriguez', 73, 'Female', 'AB+', 189, 58, 16.2, '111/79', 86, 17, 37.4, 152, 171, 1, 1, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Garcia', 44, 'Female', 'AB-', 175, 108, 35.3, '129/69', 98, 18, 36.2, 160, 172, 1, 0, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Jones', 40, 'Female', 'O-', 160, 89, 34.8, '118/72', 95, 15, 37.2, 157, 212, 0, 0, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Brown', 31, 'Male', 'A+', 161, 109, 42.1, '126/84', 74, 14, 36.4, 137, 193, 0, 0, 0, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Johnson', 36, 'Female', 'B-', 187, 112, 32.0, '115/61', 85, 12, 37.5, 122, 163, 1, 1, 1, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Garcia', 76, 'Male', 'O+', 179, 86, 26.8, '115/84', 63, 12, 37.2, 158, 171, 0, 0, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Jones', 72, 'Female', 'AB-', 190, 71, 19.7, '127/87', 61, 13, 37.8, 152, 166, 1, 1, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Rodriguez', 20, 'Male', 'AB+', 188, 61, 17.3, '136/83', 60, 20, 37.6, 107, 194, 0, 1, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Robert', 'Smith', 29, 'Female', 'AB-', 195, 82, 21.6, '131/70', 82, 18, 36.4, 73, 199, 0, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Johnson', 77, 'Female', 'AB-', 197, 52, 13.4, '105/81', 83, 12, 36.8, 144, 179, 0, 1, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Williams', 49, 'Female', 'O+', 182, 72, 21.7, '108/67', 95, 15, 37.5, 100, 183, 0, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Williams', 85, 'Female', 'AB-', 188, 51, 14.4, '128/73', 80, 19, 37.5, 100, 239, 1, 0, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Emily', 'Davis', 72, 'Male', 'O+', 167, 108, 38.7, '102/61', 96, 16, 37.3, 142, 155, 1, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Davis', 68, 'Male', 'O+', 185, 86, 25.1, '103/75', 65, 20, 36.9, 102, 239, 1, 0, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Jones', 85, 'Female', 'A-', 200, 82, 20.5, '104/82', 77, 20, 37.1, 126, 232, 0, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('David', 'Martinez', 34, 'Female', 'O+', 181, 74, 22.6, '117/71', 68, 16, 37.2, 120, 238, 0, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Williams', 53, 'Female', 'B-', 168, 92, 32.6, '102/90', 78, 19, 36.5, 146, 218, 0, 1, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Rodriguez', 42, 'Female', 'O+', 151, 112, 49.1, '119/81', 68, 18, 37.1, 144, 242, 0, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Martinez', 78, 'Male', 'AB-', 193, 89, 23.9, '101/81', 98, 16, 36.4, 150, 155, 1, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Smith', 77, 'Female', 'B+', 185, 107, 31.3, '136/76', 64, 17, 37.8, 113, 245, 1, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('James', 'Williams', 73, 'Female', 'O-', 194, 87, 23.1, '106/84', 73, 12, 36.6, 152, 156, 0, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Miller', 62, 'Female', 'O-', 154, 89, 37.5, '113/75', 95, 17, 37.7, 160, 211, 0, 1, 1, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Rodriguez', 25, 'Female', 'AB+', 171, 53, 18.1, '116/76', 88, 20, 36.6, 93, 235, 0, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Rodriguez', 46, 'Female', 'O+', 195, 104, 27.4, '132/81', 82, 12, 37.1, 145, 207, 1, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Williams', 84, 'Male', 'B+', 173, 80, 26.7, '102/74', 91, 17, 38.0, 156, 197, 0, 1, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Brown', 40, 'Male', 'B+', 167, 52, 18.6, '121/61', 83, 17, 36.2, 90, 232, 0, 0, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Miller', 49, 'Female', 'A+', 186, 51, 14.7, '126/80', 85, 19, 37.4, 85, 166, 0, 1, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Jones', 39, 'Male', 'AB+', 183, 118, 35.2, '139/67', 90, 13, 36.1, 137, 154, 0, 1, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Emily', 'Miller', 22, 'Female', 'A-', 183, 91, 27.2, '136/75', 67, 18, 37.7, 119, 200, 0, 1, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Rodriguez', 69, 'Male', 'AB+', 196, 75, 19.5, '104/82', 99, 19, 37.9, 81, 203, 0, 1, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Williams', 80, 'Female', 'B-', 163, 79, 29.7, '117/76', 85, 17, 37.7, 160, 167, 0, 0, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Smith', 72, 'Female', 'B+', 187, 88, 25.2, '137/71', 66, 14, 37.4, 145, 175, 0, 1, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Smith', 52, 'Male', 'AB+', 154, 82, 34.6, '117/65', 82, 16, 37.2, 107, 190, 1, 0, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('David', 'Rodriguez', 73, 'Female', 'O-', 161, 96, 37.0, '132/85', 81, 17, 37.3, 118, 165, 0, 1, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Brown', 88, 'Male', 'AB-', 166, 97, 35.2, '103/88', 72, 13, 37.3, 100, 159, 1, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Miller', 36, 'Female', 'A+', 174, 72, 23.8, '105/70', 71, 15, 37.5, 145, 156, 0, 0, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Brown', 42, 'Male', 'B+', 196, 64, 16.7, '132/63', 79, 14, 38.0, 88, 155, 1, 0, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Rodriguez', 71, 'Female', 'B+', 151, 88, 38.6, '115/84', 100, 20, 37.2, 129, 245, 0, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Brown', 51, 'Female', 'AB-', 165, 76, 27.9, '109/78', 78, 18, 36.5, 108, 234, 0, 1, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Davis', 31, 'Male', 'B-', 186, 90, 26.0, '101/72', 99, 15, 37.5, 136, 230, 0, 1, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Brown', 52, 'Male', 'A+', 186, 61, 17.6, '123/72', 65, 17, 36.7, 99, 228, 0, 0, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Martinez', 53, 'Female', 'A+', 169, 53, 18.6, '126/82', 63, 13, 37.8, 118, 183, 0, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Miller', 79, 'Female', 'B+', 174, 105, 34.7, '133/60', 79, 20, 36.7, 139, 250, 1, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('John', 'Brown', 55, 'Female', 'B-', 173, 59, 19.7, '117/79', 84, 17, 36.2, 116, 248, 1, 1, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Brown', 46, 'Female', 'O-', 191, 50, 13.7, '116/79', 89, 19, 36.9, 74, 168, 0, 1, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Jones', 39, 'Female', 'AB+', 195, 109, 28.7, '107/82', 81, 19, 37.2, 160, 191, 1, 0, 1, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Williams', 37, 'Female', 'AB-', 161, 86, 33.2, '103/73', 69, 16, 36.8, 108, 190, 0, 0, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('John', 'Miller', 19, 'Male', 'AB+', 196, 87, 22.6, '112/85', 79, 19, 36.0, 122, 234, 0, 0, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Johnson', 39, 'Female', 'AB-', 188, 69, 19.5, '118/77', 86, 18, 37.0, 135, 221, 1, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Jones', 69, 'Female', 'B-', 181, 112, 34.2, '119/76', 71, 13, 38.0, 82, 196, 0, 1, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Martinez', 46, 'Male', 'A+', 165, 69, 25.3, '100/66', 86, 15, 37.6, 139, 232, 0, 0, 0, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Miller', 78, 'Female', 'O-', 182, 58, 17.5, '112/64', 92, 12, 37.9, 81, 162, 0, 1, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Smith', 51, 'Female', 'O+', 156, 79, 32.5, '122/73', 81, 12, 36.5, 159, 202, 0, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Davis', 41, 'Male', 'A+', 185, 77, 22.5, '122/77', 99, 15, 37.3, 148, 232, 0, 0, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Davis', 84, 'Male', 'O-', 167, 78, 28.0, '121/85', 87, 17, 37.8, 90, 202, 1, 1, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('David', 'Smith', 26, 'Male', 'B+', 191, 83, 22.8, '135/89', 97, 20, 37.8, 79, 230, 0, 1, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Johnson', 45, 'Female', 'A+', 172, 59, 19.9, '135/90', 65, 16, 36.4, 148, 227, 1, 1, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Martinez', 46, 'Male', 'O-', 180, 90, 27.8, '129/65', 90, 17, 37.6, 159, 200, 0, 1, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Emily', 'Davis', 80, 'Female', 'A-', 199, 50, 12.6, '139/71', 97, 14, 37.3, 139, 170, 1, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Miller', 22, 'Female', 'A+', 154, 96, 40.5, '101/77', 75, 15, 37.6, 153, 212, 1, 1, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Garcia', 77, 'Male', 'A+', 160, 90, 35.2, '100/63', 75, 16, 36.3, 88, 182, 0, 1, 0, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Davis', 30, 'Female', 'B-', 151, 54, 23.7, '140/75', 67, 20, 36.3, 128, 192, 1, 0, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Davis', 66, 'Male', 'O+', 160, 82, 32.0, '118/69', 90, 20, 37.3, 80, 213, 0, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Williams', 30, 'Female', 'B-', 191, 57, 15.6, '129/74', 66, 14, 37.3, 103, 166, 0, 1, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Jones', 83, 'Female', 'AB-', 173, 54, 18.0, '124/60', 93, 15, 36.3, 143, 189, 1, 1, 0, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Emily', 'Williams', 76, 'Female', 'O+', 168, 97, 34.4, '139/86', 85, 12, 37.5, 160, 201, 1, 1, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Smith', 63, 'Male', 'B+', 176, 105, 33.9, '116/76', 76, 19, 36.0, 81, 176, 0, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Brown', 82, 'Female', 'B-', 186, 58, 16.8, '124/72', 93, 19, 37.0, 123, 159, 1, 0, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Miller', 18, 'Female', 'AB-', 179, 108, 33.7, '132/85', 68, 16, 37.9, 122, 185, 1, 1, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Williams', 34, 'Female', 'B-', 159, 85, 33.6, '137/70', 86, 20, 37.2, 102, 157, 1, 0, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Robert', 'Garcia', 56, 'Male', 'A+', 162, 81, 30.9, '135/82', 95, 19, 36.8, 73, 226, 1, 0, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Martinez', 30, 'Female', 'B-', 178, 112, 35.3, '111/79', 96, 18, 37.3, 85, 161, 1, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Brown', 34, 'Female', 'O-', 180, 75, 23.1, '113/79', 78, 12, 36.4, 75, 193, 0, 1, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Martinez', 45, 'Female', 'B-', 183, 93, 27.8, '114/72', 97, 12, 37.8, 136, 229, 0, 1, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Martinez', 60, 'Female', 'B-', 192, 58, 15.7, '127/68', 93, 19, 36.3, 75, 187, 0, 1, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Brown', 39, 'Female', 'O+', 198, 62, 15.8, '113/73', 78, 15, 36.3, 104, 243, 0, 0, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Robert', 'Johnson', 89, 'Male', 'O-', 177, 108, 34.5, '120/62', 75, 14, 37.3, 112, 182, 1, 1, 0, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('James', 'Rodriguez', 25, 'Male', 'AB+', 153, 60, 25.6, '135/60', 100, 12, 36.4, 114, 189, 1, 1, 0, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Davis', 24, 'Male', 'AB-', 163, 84, 31.6, '133/68', 100, 19, 37.4, 118, 197, 1, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Garcia', 20, 'Female', 'B+', 153, 67, 28.6, '117/84', 66, 12, 37.4, 159, 188, 1, 0, 0, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Rodriguez', 49, 'Female', 'B+', 200, 99, 24.8, '112/60', 92, 15, 36.9, 107, 228, 1, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Martinez', 55, 'Female', 'B-', 180, 76, 23.5, '140/69', 93, 15, 37.9, 106, 162, 0, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Miller', 35, 'Male', 'B+', 185, 92, 26.9, '104/84', 81, 16, 37.1, 149, 244, 0, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('John', 'Davis', 62, 'Female', 'AB+', 191, 71, 19.5, '131/84', 73, 20, 36.2, 131, 196, 1, 0, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Robert', 'Davis', 50, 'Female', 'AB+', 157, 101, 41.0, '102/69', 71, 14, 36.4, 116, 206, 1, 1, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Jane', 'Brown', 67, 'Female', 'AB-', 200, 87, 21.8, '129/70', 83, 17, 37.2, 77, 240, 0, 0, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('David', 'Smith', 84, 'Female', 'O+', 198, 84, 21.4, '114/64', 98, 12, 37.0, 79, 150, 1, 0, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Williams', 82, 'Male', 'A+', 171, 85, 29.1, '105/63', 60, 20, 37.8, 112, 161, 1, 1, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Rodriguez', 38, 'Female', 'B+', 160, 102, 39.8, '136/74', 96, 14, 36.7, 96, 197, 1, 0, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Johnson', 44, 'Male', 'AB-', 195, 52, 13.7, '115/63', 94, 16, 36.5, 86, 197, 1, 1, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Mike', 'Martinez', 82, 'Female', 'O+', 196, 105, 27.3, '138/76', 93, 15, 37.2, 97, 200, 1, 1, 1, 1, 1, 0);
INSERT INTO MedicalRecords VALUES ('Emily', 'Johnson', 66, 'Male', 'B-', 196, 52, 13.5, '126/87', 87, 15, 37.1, 147, 216, 0, 1, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('David', 'Martinez', 80, 'Female', 'B-', 164, 62, 23.1, '106/66', 94, 13, 36.5, 151, 195, 0, 1, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Martinez', 56, 'Male', 'A-', 165, 60, 22.0, '115/63', 73, 14, 37.9, 105, 200, 1, 0, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Smith', 74, 'Male', 'B+', 156, 77, 31.6, '136/90', 91, 18, 36.6, 137, 219, 1, 0, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('Laura', 'Rodriguez', 49, 'Male', 'B+', 189, 109, 30.5, '120/76', 67, 14, 37.9, 160, 233, 0, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Garcia', 21, 'Female', 'B+', 182, 89, 26.9, '132/73', 86, 18, 36.3, 157, 222, 0, 1, 0, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Mike', 'Williams', 82, 'Male', 'A-', 186, 118, 34.1, '138/75', 70, 18, 36.4, 124, 165, 1, 1, 1, 0, 1, 1);
INSERT INTO MedicalRecords VALUES ('John', 'Jones', 89, 'Male', 'A-', 154, 51, 21.5, '138/73', 68, 17, 36.3, 112, 199, 0, 0, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('James', 'Smith', 40, 'Female', 'AB+', 181, 114, 34.8, '130/70', 97, 12, 36.3, 151, 161, 0, 1, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Johnson', 67, 'Male', 'B-', 193, 79, 21.2, '120/65', 91, 15, 36.4, 90, 193, 1, 0, 0, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Linda', 'Davis', 34, 'Female', 'A-', 151, 90, 39.5, '125/82', 88, 16, 37.3, 109, 175, 0, 0, 1, 0, 0, 1);
INSERT INTO MedicalRecords VALUES ('Linda', 'Jones', 38, 'Male', 'A+', 168, 78, 27.6, '112/81', 90, 12, 36.8, 133, 228, 1, 0, 1, 1, 0, 1);
INSERT INTO MedicalRecords VALUES ('Robert', 'Martinez', 90, 'Male', 'B-', 151, 104, 45.6, '102/86', 78, 18, 36.1, 83, 157, 0, 1, 1, 1, 1, 1);
INSERT INTO MedicalRecords VALUES ('Sarah', 'Martinez', 73, 'Male', 'AB+', 199, 113, 28.5, '121/74', 95, 13, 36.1, 109, 241, 0, 0, 1, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Garcia', 35, 'Female', 'A-', 152, 87, 37.7, '100/76', 98, 15, 37.9, 91, 205, 0, 0, 1, 0, 0, 0);
INSERT INTO MedicalRecords VALUES ('Jane', 'Johnson', 76, 'Female', 'AB-', 193, 99, 26.6, '116/89', 69, 15, 37.0, 112, 194, 1, 0, 0, 0, 1, 0);
INSERT INTO MedicalRecords VALUES ('Laura', 'Garcia', 86, 'Female', 'O+', 162, 51, 19.4, '104/88', 86, 12, 37.9, 112, 235, 1, 1, 1, 1, 0, 0);
INSERT INTO MedicalRecords VALUES ('Robert', 'Brown', 54, 'Male', 'AB+', 168, 56, 19.8, '120/61', 99, 14, 36.0, 138, 218, 1, 1, 0, 0, 1, 1);
SELECT blood_type AS bt, AVG(has_diabetes), sex AS male_or_female FROM MedicalRecords GROUP BY blood_type, sex;
SELECT blood_type, COUNT(has_diabetes), sex AS male_or_female FROM MedicalRecords GROUP BY blood_type, sex;
