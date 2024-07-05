import cv2
import face_recognition
import pickle
import os
import requests
import firebase_admin
from firebase_admin import credentials, storage

# Initialize Firebase Admin SDK
cred = credentials.Certificate("bushgardens-e7339-firebase-adminsdk-42jdi-c6d0b389d6.json")
firebase_admin.initialize_app(cred, {
    'databaseURL': 'https://bushgardens-e7339-default-rtdb.asia-southeast1.firebasedatabase.app/',
    'storageBucket': 'bushgardens-e7339.appspot.com'
})
# Define the folder in Firebase Storage
firebase_folder = "Images/"

# Define the local folder to save the images
local_folder = "Images"

# Create the local folder if it doesn't exist
if not os.path.exists(local_folder):
    os.makedirs(local_folder)

# Get the bucket
bucket = storage.bucket()

# List all files in the specified folder
blobs = bucket.list_blobs(prefix=firebase_folder)

# Initialize lists to store student IDs and image file names
studentIds = []
image_files = []

# Download each image file
for blob in blobs:
    # Skip folder names
    if not blob.name.endswith('/'):
        local_filename = os.path.join(local_folder, os.path.basename(blob.name))
        print(f"Downloading {blob.name} to {local_filename}")
        try:
            # Generate the download URL
            url = blob.generate_signed_url(version="v4", expiration=600)  # URL valid for 10 minutes
            # Download the file using requests
            response = requests.get(url)
            response.raise_for_status()  # Check if the request was successful
            with open(local_filename, 'wb') as f:
                f.write(response.content)
            student_id = os.path.splitext(os.path.basename(blob.name))[0]
            studentIds.append(student_id)
            image_files.append(local_filename)
        except Exception as e:
            print(f"Failed to download {blob.name}: {e}")
            continue

print("All images have been downloaded.")
print("Student IDs:", studentIds)
print("Image files:", image_files)

# Function to find encodings for a list of images
def findEncodings(imagesList):
    encodeList = []
    for img_path in imagesList:
        img = cv2.imread(img_path)
        img = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
        encode = face_recognition.face_encodings(img)
        if encode:  # Check if encoding is found
            encodeList.append(encode[0])
        else:
            print(f"Encoding not found for image {img_path}")

    return encodeList

print("Encoding Started ...")
encodeListKnown = findEncodings(image_files)
encodeListKnownWithIds = [encodeListKnown, studentIds]
print("Encoding Complete")

# Save the encodings and student IDs to a file
with open("EncodeFile.p", 'wb') as file:
    pickle.dump(encodeListKnownWithIds, file)

print("File Saved")
