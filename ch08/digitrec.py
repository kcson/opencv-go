import sys
import numpy as np
import cv2 as cv


def load_digits():
    None


src = cv.imread("../imgs/digits_print.bmp")
if src is None:
    print("image load fail!!")
    sys.exit()

cv.imshow('src', src)

cv.waitKey()
cv.destroyAllWindows()
