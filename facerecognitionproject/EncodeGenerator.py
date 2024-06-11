import cv2
import face_recognition
import pickle
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

#importing the mode images into a list
pathList = os.listdir("Images")
imgList = []
studentIds = []
for path in pathList:
    imgList.append(cv2.imread(os.path.join('Images', path)))
    studentIds.append(os.path.splitext(path)[0])

    fileName = f'Images/{path}'
    bucket = storage.bucket()
    blob = bucket.blob(fileName)
    blob.upload_from_filename(fileName)

# print(len(imgList))
# print(studentIds)

def findEncoding(imagesList):
    encodings = []
    for img in imagesList:
        img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
        encode = face_recognition.face_encodings(img)[0]
        encodings.append(encode)
    return encodings

print("Encoding started ...")
encodingsKnown = findEncoding(imgList)
encodingsKnownWithIds = [encodingsKnown, studentIds]
print("Encoding completed")

file = open("encodeFile.p", "wb")
pickle.dump(encodingsKnownWithIds, file)
file.close()
print("File saved")