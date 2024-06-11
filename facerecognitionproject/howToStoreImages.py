import cv2
import os

import firebase_admin
from firebase_admin import credentials
from firebase_admin import db
from firebase_admin import storage

cred = credentials.Certificate("serviceAccountKey.json")
firebase_admin.initialize_app(cred, {
    'databaseURL': 'https://facerecognition-f0579-default-rtdb.asia-southeast1.firebasedatabase.app/',
    'storageBucket': 'facerecognition-f0579.appspot.com'
})
bucket = storage.bucket()
blob = bucket.blob("980745/980745.jpg")
blob.upload_from_filename("Images/980745.jpg")