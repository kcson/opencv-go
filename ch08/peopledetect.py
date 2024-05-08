import sys
import random

import cv2 as cv

cap = cv.VideoCapture('vtest.avi')
if not cap.isOpened():
    print('Cannot open camera')
    sys.exit()

hog = cv.HOGDescriptor()
hog.setSVMDetector(cv.HOGDescriptor_getDefaultPeopleDetector())

while True:
    ret, frame = cap.read()
    if not ret:
        break

    detected, _ = hog.detectMultiScale(frame)
    for (x, y, w, h) in detected:
        c = (random.randint(0, 255), random.randint(0, 255), random.randint(0, 255))
        cv.rectangle(frame, (x, y), (x+w, y+h), c, 3)

    cv.imshow('frame', frame)
    if cv.waitKey(10) == 27:
        break

cv.destroyAllWindows()