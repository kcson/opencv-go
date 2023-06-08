import sys
import numpy as np
import cv2 as cv

tm = cv.TickMeter()
src = cv.imread('../imgs/lenna.bmp')
if src is None:
    print('image load fail !!')
    sys.exit()

print(src.shape)
cv.imshow('src', src)
cv.waitKey()
cv.destroyAllWindows()
