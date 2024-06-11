import os
import pickle
import numpy as np
import cv2
import face_recognition
import cvzone
import firebase_admin
from firebase_admin import credentials
from firebase_admin import db
from firebase_admin import storage
from datetime import datetime
from deepface import DeepFace

cred = credentials.Certificate("serviceAccountKey.json")
firebase_admin.initialize_app(cred, {
    'databaseURL': 'https://facerecognition-f0579-default-rtdb.asia-southeast1.firebasedatabase.app/',
    'storageBucket': 'facerecognition-f0579.appspot.com'
})

bucket = storage.bucket()

cap = cv2.VideoCapture(2)
cap.set(3, 640)
cap.set(4, 480)

imgBackground = cv2.imread('Resources/background.png')

#importing the mode images into a list
modePathList = os.listdir("Resources/Modes")
imgModeList = []
for path in modePathList:
    imgModeList.append(cv2.imread(os.path.join('Resources/Modes', path)))
#print(len(imgModeList))

# load the encoding file
print("Loading encode file")
file = open("encodeFile.p", "rb")
encodingsKnownWithIds = pickle.load(file)
file.close()
encodingsKnown, studentIds = encodingsKnownWithIds
print("Encode file loaded")

modeType = 0
counter = 0
id = -1
imgStudent = []

surprise_count = 0
happy_count = 0
make_photo = 1

while True:
    success, img = cap.read()

    imgSmaller = cv2.resize(img, (0, 0), None, 0.25, 0.25)
    imgSmaller = cv2.cvtColor(imgSmaller, cv2.COLOR_BGR2RGB)

    faceCurrentFrame = face_recognition.face_locations(imgSmaller)
    encodeCurrentFrame = face_recognition.face_encodings(imgSmaller, faceCurrentFrame)

    imgBackground[162:162 + 480, 55:55 + 640] = img
    imgBackground[44:44 + 633, 808:808 + 414] = imgModeList[modeType]

    if faceCurrentFrame:
        for encodeFace, faceLocation in zip(encodeCurrentFrame, faceCurrentFrame):
            matches = face_recognition.compare_faces(encodingsKnown, encodeFace, 0.5)
            faceDistance = face_recognition.face_distance(encodingsKnown, encodeFace)
            # print(matches)
            # print(faceDistance)

            matchIndex = np.argmin(faceDistance)
            if matches[matchIndex]:
                # print("Match found")
                # print(studentIds[matchIndex])
                y1, x2, y2, x1 = faceLocation
                y1, x2, y2, x1 = y1 * 4, x2 * 4, y2 * 4, x1 * 4
                face_img = img[y1:y2, x1:x2]
                bbox = 55 + x1, 162 + y1, x2 - x1, y2 - y1
                imgBackground = cvzone.cornerRect(imgBackground, bbox, rt=0)
                id = studentIds[matchIndex]

                if counter == 0:
                    cvzone.putTextRect(imgBackground, 'Loading', (275, 400))
                    cv2.imshow('Face attendance', imgBackground)
                    cv2.waitKey(1)
                    counter = 1
                    modeType = 1
            else:
                 #print("Match not found")
                 cvzone.putTextRect(imgBackground, 'Match not found', (185, 400))
        if counter != 0:
            if counter == 1:
                #get the data
                studentInfo = db.reference(f"Students/{id}").get()
                print(studentInfo)
                #get the image from the storage
                blob = bucket.blob(f"Images/{id}.jpg")
                array = np.frombuffer(blob.download_as_string(), np.uint8)
                imgStudent = cv2.imdecode(array, cv2.COLOR_BGRA2BGR)
                #update data of attendance
                dateTimeObj = datetime.strptime(studentInfo['last_attendance_time'], "%Y-%m-%d %H:%M:%S")
                secondsElapsed = (datetime.now() - dateTimeObj).total_seconds()
                print(secondsElapsed)
                if secondsElapsed > 30:
                    # ref = db.reference(f"Students/{id}")
                    # studentInfo['total_attendance'] += 1
                    # ref.child('total_attendance').set(studentInfo['total_attendance'])
                    # ref.child('last_attendance_time').set(datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
                    make_photo = 1
                else:
                    modeType = 3
                    counter = 0
                    imgBackground[44:44 + 633, 808:808 + 414] = imgModeList[modeType]

            if modeType != 3:
                if 10 < counter < 20:
                    modeType = 2
                imgBackground[44:44 + 633, 808:808 + 414] = imgModeList[modeType]
                if counter <= 10:
                    cv2.putText(imgBackground, str(studentInfo['total_attendance']), (861, 125),
                                cv2.FONT_HERSHEY_COMPLEX, 1, (255, 255, 255), 1)
                    cv2.putText(imgBackground, str(id), (1006, 493),
                                cv2.FONT_HERSHEY_COMPLEX, 0.5, (255, 255, 255), 1)
                    cv2.putText(imgBackground, str(studentInfo['date_of_birth']), (900, 650),
                                cv2.FONT_HERSHEY_COMPLEX, 0.6, (100, 100, 100), 1)
                    cv2.putText(imgBackground, str(studentInfo['gender']), (1009, 625),
                                cv2.FONT_HERSHEY_COMPLEX, 0.6, (100, 100, 100), 1)
                    cv2.putText(imgBackground, str(studentInfo['starting_year']), (1125, 625),
                                cv2.FONT_HERSHEY_COMPLEX, 0.6, (100, 100, 100), 1)
                    (w, h), _ = cv2.getTextSize(studentInfo['name'], cv2.FONT_HERSHEY_COMPLEX, 1, 1)
                    offset = (414 - w) // 2
                    cv2.putText(imgBackground, str(studentInfo['name']), (808 + offset, 445),
                                cv2.FONT_HERSHEY_COMPLEX, 1, (50, 50, 50), 1)
                    imgBackground[175:175 + 216, 909:909 + 216] = imgStudent
                if make_photo == 1:
                    #saving photo if the emotion is surprise or happy
                    #while surprise_count < 6 or happy_count < 10:
                    try:
                        # Analyze emotions if face is detected
                        result = DeepFace.analyze(face_img, actions=['emotion'], enforce_detection=False)
                        current_emotion = result[0]['dominant_emotion']
                        if current_emotion != 'neutral':
                            happy_count += 1
                            filename = current_emotion + datetime.now().strftime("%Y-%m-%d_%H-%M-%S") + ".jpg"
                            blob = bucket.blob(f"{id}/{filename}")
                            cv2.imwrite(filename, img)
                            blob.upload_from_filename(filename)
                            print(f"{current_emotion} frame saved! {filename}")
                            make_photo = 0
                        # elif current_emotion == 'surprise':
                        #     surprise_count += 1
                        #     filename = "surprise_frame" + datetime.now().strftime("%Y-%m-%d_%H-%M-%S") + ".jpg"
                        #     blob = bucket.blob(f"{id}/{filename}")
                        #     cv2.imwrite(filename, img)
                        #     blob.upload_from_filename(filename)
                        #     print(f"Surprise frame saved! {filename}")
                        #     make_photo = 0
                        if make_photo == 0:
                            ref = db.reference(f"Students/{id}")
                            studentInfo['total_attendance'] += 1
                            ref.child('total_attendance').set(studentInfo['total_attendance'])
                            ref.child('last_attendance_time').set(datetime.now().strftime("%Y-%m-%d %H:%M:%S"))
                    except ValueError as e:
                        # If no face is detected, skip emotion analysis
                        pass

                if make_photo == 0:
                    counter += 1

                if counter >= 20:
                    counter = 0
                    modeType = 0
                    studentInfo = []
                    imgStudent = []
                    imgBackground[44:44 + 633, 808:808 + 414] = imgModeList[modeType]

    else:
        modeType = 0
        counter = 0
    #cv2.imshow('Webcam', img)
    cv2.imshow('Face attendance', imgBackground)
    cv2.waitKey(1)