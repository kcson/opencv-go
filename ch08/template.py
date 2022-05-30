import sys
import cv2 as cv
import numpy as np

src = cv.imread('../imgs/circuit.bmp', cv.IMREAD_GRAYSCALE)
templ = cv.imread('../imgs/crystal.bmp', cv.IMREAD_GRAYSCALE)

if src is None or templ is None:
    print('image read fail!!')
    sys.exit()

noise = np.zeros(src.shape, np.int32)
cv.randn(noise, 50, 10)
src = cv.add(src, noise, dtype=cv.CV_8UC3)

res = cv.matchTemplate(src, templ, cv.TM_CCOEFF_NORMED)
res_norm = cv.normalize(res, None, 0, 255, cv.NORM_MINMAX, cv.CV_8U)

_, maxv, _, maxloc = cv.minMaxLoc(res)

th, tw = templ.shape
dst = cv.cvtColor(src, cv.COLOR_GRAY2BGR)
cv.rectangle(dst, maxloc, (maxloc[0] + tw, maxloc[1] + th), (0, 0, 255), 2)

cv.imshow('res_norm', res_norm)
cv.imshow('dst', dst)
cv.waitKey()
cv.destroyAllWindows()
