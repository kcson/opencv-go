import sys
import numpy as np
import cv2 as cv

obj = cv.imread('../imgs/spades.png', cv.IMREAD_GRAYSCALE)
src = cv.imread('../imgs/symbols.png', cv.IMREAD_GRAYSCALE)

if obj is None or src is None:
    print('img read fail!!!')
    sys.exit()

_, obj_bin = cv.threshold(obj, 128, 255, cv.THRESH_BINARY_INV)
obj_contours, _ = cv.findContours(obj_bin, cv.RETR_EXTERNAL, cv.CHAIN_APPROX_NONE)
obj_pts = obj_contours[0]

_, src_bin = cv.threshold(src, 128, 255, cv.THRESH_BINARY_INV)
contours, _ = cv.findContours(src_bin, cv.RETR_EXTERNAL, cv.CHAIN_APPROX_NONE)

dst = cv.cvtColor(src, cv.COLOR_GRAY2BGR)

for pts in contours:
    if cv.contourArea(pts) < 1000:
        continue

    rc = cv.boundingRect(pts)
    cv.rectangle(dst, rc, (255, 0, 0), 1)

    dist = cv.matchShapes(obj_pts, pts, cv.CONTOURS_MATCH_I3, 0)
    cv.putText(dst, str(round(dist, 4)), (rc[0], rc[1] - 3), cv.FONT_HERSHEY_SIMPLEX, 0.5, (255, 0, 0), 1, cv.LINE_AA)

    if dist < 0.1:
        cv.rectangle(dst, rc, (0, 0, 255), 2)

cv.imshow('obj', obj)
cv.imshow('dst', dst)
cv.waitKey(0)
