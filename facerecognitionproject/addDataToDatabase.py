import firebase_admin
from firebase_admin import credentials
from firebase_admin import db

cred = credentials.Certificate("serviceAccountKey.json")
firebase_admin.initialize_app(cred, {
    'databaseURL': 'https://facerecognition-f0579-default-rtdb.asia-southeast1.firebasedatabase.app/'
})

ref = db.reference('Students')

data = {
    "777777":
        {"name": "Nurzhan Danenov",
         "date_of_birth": "2003-10-19",
         "gender": "male",
         "total_attendance": 7,
         "starting_year": 2020,
         "last_attendance_time": "2024-02-16 16:04:30"},
    "324567":
        {"name": "Elon Musk",
         "date_of_birth": "1971-06-28",
         "gender": "male",
         "total_attendance": 5,
         "starting_year": 2018,
         "last_attendance_time": "2024-02-16 16:04:30"},
    "980745":
        {"name": "Megan Fox",
         "date_of_birth": "1986-05-16",
         "gender": "female",
         "total_attendance": 9,
         "starting_year": 2022,
         "last_attendance_time": "2024-02-16 16:04:30"}
}

for key, value in data.items():
    ref.child(key).set(value)