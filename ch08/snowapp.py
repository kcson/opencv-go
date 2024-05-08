import sys

import cv2 as cv
import numpy as np


def overlay(img, glasses, pos):
    sx = pos[0]
    ex = pos[0] + glasses.shape[1]
    sy = pos[1]
    ey = pos[1] + glasses.shape[0]

    if sx < 0 or sy < 0 or ex > img.shape[1] or ey > img.shape[0]:
        return

    img1 = img[sy:ey, sx:ex]
    img2 = glasses[:, :, 0:3]
    alpha = 1. - (glasses[:, :, 3] / 255.)

    img1[..., 0] = (img1[..., 0] * alpha + img2[..., 0] * (1. - alpha)).astype(np.uint8)
    img1[..., 1] = (img1[..., 1] * alpha + img2[..., 1] * (1. - alpha)).astype(np.uint8)
    img1[..., 2] = (img1[..., 2] * alpha + img2[..., 2] * (1. - alpha)).astype(np.uint8)


cap = cv.VideoCapture(0)
if not cap.isOpened():
    print('Cannot open camera')
    sys.exit()

w = int(cap.get(cv.CAP_PROP_FRAME_WIDTH))
h = int(cap.get(cv.CAP_PROP_FRAME_HEIGHT))

fourcc = cv.VideoWriter_fourcc(*'DIVX')
out = cv.VideoWriter('output.avi', fourcc, 30.0, (w, h))

face_classifier = cv.CascadeClassifier('haarcascade_frontalface_alt2.xml')
eye_classifier = cv.CascadeClassifier('haarcascade_eye.xml')
if face_classifier.empty() or eye_classifier.empty():
    print('No face detected')
    sys.exit()

glasses = cv.imread('glasses.png', cv.IMREAD_UNCHANGED)
if glasses is None:
    print('No glasses detected')
    sys.exit()

ew, eh = glasses.shape[:2]
ex1, ey1 = 240, 300
ex2, ey2 = 660, 300

while True:
    ret, frame = cap.read()
    if not ret:
        break

    faces = face_classifier.detectMultiScale(frame, scaleFactor=1.2, minSize=(100, 100), maxSize=(400, 400))
    for (x, y, w, h) in faces:
        faceROI = frame[y:y + h // 2, x:x + w]
        eyes = eye_classifier.detectMultiScale(faceROI)
        if len(eyes) != 2:
            continue

        # 두 개의 눈 중앙 위치를 (x1, y1), (x2, y2) 좌표로 저장
        x1 = x + eyes[0][0] + (eyes[0][2] // 2)
        y1 = y + eyes[0][1] + (eyes[0][3] // 2)
        x2 = x + eyes[1][0] + (eyes[1][2] // 2)
        y2 = y + eyes[1][1] + (eyes[1][3] // 2)

        if x1 > x2:
            x1, y1, x2, y2 = x2, y2, x1, y1

        fx = (x2 - x1) / (ex2 - ex1)
        glasses2 = cv.resize(glasses, (0, 0), fx=fx, fy=fx, interpolation=cv.INTER_AREA)

        pos = (x1 - int(ex1 * fx), y1 - int(ey1 * fx))

        overlay(frame, glasses2, pos)

    out.write(frame)
    cv.imshow('frame', frame)

    if cv.waitKey(1) == 27:
        break

cap.release()
out.release()
cv.destroyAllWindows()
