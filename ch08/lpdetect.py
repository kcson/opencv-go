import sys
import numpy as np
import cv2 as cv

tm = cv.TickMeter()
src = cv.imread('../imgs/vehicle.jpeg')
# src = cv.imread('car.jpeg')
if src is None:
    print('image load fail !!')
    sys.exit()

lp_classifier = cv.CascadeClassifier('haarcascade_russian_plate_number.xml')
if lp_classifier.empty():
    print('xml load fail!!')
    sys.exit()

tm.start()
lps = lp_classifier.detectMultiScale(src, scaleFactor=1.2)
for (x1, y1, w1, h1) in lps:
    cv.rectangle(src, (x1, y1, w1, h1), (255, 0, 255), 2)
tm.stop()
print(tm.getTimeMilli())
cv.imshow('src', src)
cv.waitKey()
cv.destroyAllWindows()
