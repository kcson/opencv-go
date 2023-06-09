import sys
import numpy as np
import cv2 as cv

tm = cv.TickMeter()
src = cv.imread('../imgs/lenna.bmp')
if src is None:
    print('image load fail !!')
    sys.exit()

face_classifier = cv.CascadeClassifier('haarcascade_frontalface_alt2.xml')
eye_classifier = cv.CascadeClassifier('haarcascade_eye.xml')
if face_classifier.empty() or eye_classifier.empty():
    print('xml load fail!!')
    sys.exit()

tm.start()
faces = face_classifier.detectMultiScale(src, scaleFactor=1.2, minSize=(100, 100))
for (x1, y1, w1, h1) in faces:
    cv.rectangle(src, (x1, y1, w1, h1), (255, 0, 255), 2)

    faceROI = src[y1:y1 + h1 // 2, x1:x1 + w1]
    eyes = eye_classifier.detectMultiScale(faceROI)
    for (x2, y2, w2, h2) in eyes:
        center = (x2 + w2 // 2, y2 + h2 // 2)
        cv.circle(faceROI, center, w2 // 2, (255, 0, 0), 2, cv.LINE_AA)
tm.stop()
print(tm.getTimeMilli())
cv.imshow('src', src)
cv.waitKey()
cv.destroyAllWindows()
